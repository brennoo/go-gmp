package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// lockEntry is the subset of the mirror's lock file entry schemadrift
// actually reads. The mirror itself (and its full lock schema) now lives in
// github.com/brennoo/gvm-mirror; go-gmp keeps only a vendored copy of its
// lock file at hack/image-mirror.lock.json, and this narrower struct is
// self-documenting about the dependency: encoding/json ignores every other
// field.
type lockEntry struct {
	SourceRepository  string `json:"source_repository"`
	SourceChannel     string `json:"source_channel"`
	DestinationDigest string `json:"destination_digest"`
}

type lock struct {
	Entries []lockEntry `json:"entries"`
}

// readLock reads and parses the lock file at path. There is deliberately no
// corresponding write: this is a read-only view of a vendored copy that
// belongs to github.com/brennoo/gvm-mirror, not something schemadrift owns.
func readLock(path string) (lock, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is an operator-supplied CLI flag, not untrusted input
	if err != nil {
		return lock{}, err
	}
	var l lock
	if err := json.Unmarshal(data, &l); err != nil {
		return lock{}, fmt.Errorf("parsing lock file %s: %w", path, err)
	}
	return l, nil
}
