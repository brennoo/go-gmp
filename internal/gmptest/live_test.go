//go:build integration

package gmptest

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"

	gmp "github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/commands"
)

// liveClient is the authenticated client the readiness ladder in live.go
// leaves behind, shared by every test in this package. Package-level rather
// than threaded through *testing.T because the ladder can take the better
// part of an hour on a first boot, and each test paying for it is not an option.
var liveClient gmp.Client

// liveConn is the same authenticated transport liveClient sits on, exposed
// separately because capture_test.go needs RawXML: a fixture must be the
// bytes gvmd sent, and a response the library can't decode is precisely the
// most valuable thing the corpus can record.
//
// TestMain owns its lifetime and closes it after m.Run; nothing else may.
var liveConn gmp.Connection

var (
	liveOnce sync.Once
	errLive  error
)

// TestMain parses flags and closes the shared connection.
//
// flag.Parse must be called explicitly: a custom TestMain opts out of the
// parsing testing.Main would otherwise do, and the -gmp.* flags in live.go
// must be parsed before WaitReady reads them.
//
// The readiness ladder deliberately does not run here — it used to, which
// made merely *selecting* a test pay for it, including `-run TestCaptureLive`
// with no -capture. Moving it behind requireLive makes the cost proportional
// to what's actually run, while keeping "exactly once per binary" (sync.Once)
// and "fail loudly" (t.Fatalf on first use).
func TestMain(m *testing.M) {
	flag.Parse()

	// os.Exit does not run deferred calls, so the connection is closed
	// explicitly around m.Run rather than via defer.
	code := m.Run()
	if liveConn != nil {
		_ = liveConn.Close()
	}
	os.Exit(code)
}

// requireLive returns the shared authenticated client, running the readiness
// ladder on first use and failing the calling test if gvmd never gets there.
func requireLive(t *testing.T) gmp.Client {
	t.Helper()

	liveOnce.Do(func() {
		progress := func(format string, args ...any) {
			fmt.Fprintf(os.Stderr, "gmptest: "+format+"\n", args...)
		}
		cli, conn, err := WaitReady(context.Background(), progress)
		if err != nil {
			errLive = err
			return
		}
		liveClient, liveConn = cli, conn
	})

	if errLive != nil {
		t.Fatalf("Live gvmd never became ready: %v", errLive)
	}
	return liveClient
}

// TestGetVersionSmoke is deliberately minimal: it only proves the ladder in
// live.go left behind a working, authenticated client.
func TestGetVersionSmoke(t *testing.T) {
	cli := requireLive(t)

	resp, err := cli.GetVersion(context.Background(), &commands.GetVersion{})
	if err != nil {
		t.Fatalf("GetVersion: %v", err)
	}
	if resp.Version == "" {
		t.Fatal("GetVersion returned an empty version")
	}
}
