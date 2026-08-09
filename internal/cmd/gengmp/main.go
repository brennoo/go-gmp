// Command gengmp materializes internal/codegen's generated Go source into a
// directory.
//
// It is a thin wrapper: load the committed Stage-A schema snapshots, call
// codegen.Generate, write what comes back. It never reads vendor/XML/commands
// itself — that AGPL-licensed extraction step already happened when the
// snapshots under -schema were committed (internal/codegen/snapshot_test.go
// -update), and keeping this CLI on the far side of that boundary is what
// T-9's licensing isolation (acceptance criterion 9) requires: nothing
// downstream of the snapshot ever needs to know vendor/XML existed.
//
// This is a manual, human-reviewed step, not something CI runs to write code:
// see hack/generate.sh.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/brennoo/go-gmp/internal/codegen"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "gengmp:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("gengmp", flag.ContinueOnError)
	schemaDir := fs.String("schema", "internal/codegen/testdata/schema",
		"directory of committed Stage-A schema snapshots (internal/codegen/snapshot.go); never vendor/XML directly")
	outDir := fs.String("out", ".",
		"directory to write generated Go files into. Must be passed explicitly — there is no destructive default, "+
			"so a forgotten -out fails loudly instead of writing into the current directory.")
	if err := fs.Parse(args); err != nil {
		return err
	}

	explicitOut := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "out" {
			explicitOut = true
		}
	})
	if !explicitOut {
		return fmt.Errorf("-out is required; e.g. -out ./commands (see hack/generate.sh)")
	}

	cmds, err := codegen.LoadSnapshots(*schemaDir)
	if err != nil {
		return fmt.Errorf("loading snapshots from %s: %w", *schemaDir, err)
	}

	files, err := codegen.Generate(cmds, codegen.TypeOverrides)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if err := os.MkdirAll(*outDir, 0o750); err != nil {
		return fmt.Errorf("creating %s: %w", *outDir, err)
	}

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	total := 0
	for _, name := range names {
		src := files[name]
		path := filepath.Join(*outDir, name)
		if err := os.WriteFile(path, []byte(src), 0o600); err != nil { //nolint:gosec // -out is an operator-supplied destination, by design
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("wrote %-40s %6d bytes\n", path, len(src))
		total += len(src)
	}
	fmt.Printf("%d file(s), %d bytes total, written to %s\n", len(names), total, *outDir)
	return nil
}
