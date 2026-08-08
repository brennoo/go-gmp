package gmptest

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// fixtureFile matches the four files that make up one fixture.
var fixtureFile = regexp.MustCompile(`^(\d{3})_(request\.xml|response\.xml|meta\.json|fidelity\.txt)$`)

// TestCorpusIsComplete walks the fixture tree directly rather than through
// LoadFixtures, which discovers from *_response.xml and would not notice
// anything orphaned alongside it.
//
// Without this, deleting a request document degrades silently: HasRequest goes
// false, the round-trip test skips that fixture, and the suite stays green
// while quietly covering less.
func TestCorpusIsComplete(t *testing.T) {
	var problems []string

	for _, prov := range []Provenance{Upstream, Live} {
		base := filepath.Join(TestdataDir, string(prov))
		if _, err := os.Stat(base); err != nil {
			continue
		}

		cmdDirs, err := os.ReadDir(base)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", base, err)
		}

		for _, cd := range cmdDirs {
			if !cd.IsDir() {
				// go.mod and MANIFEST.md are expected at the provenance root.
				switch cd.Name() {
				case "go.mod", "MANIFEST.md":
				default:
					problems = append(problems, base+"/"+cd.Name()+": unexpected file at the provenance root")
				}
				continue
			}
			if _, known := Registry[cd.Name()]; !known {
				problems = append(problems, base+"/"+cd.Name()+": directory does not name a known command")
			}
			problems = append(problems, checkCommandDir(t, filepath.Join(base, cd.Name()), prov)...)
		}
	}

	sort.Strings(problems)
	for _, p := range problems {
		t.Errorf("%s", p)
	}
}

// checkCommandDir verifies that every fixture in dir has exactly the file set
// it should, and that nothing else is present.
func checkCommandDir(t *testing.T, dir string, prov Provenance) []string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{dir + ": " + err.Error()}
	}

	var problems []string
	seen := map[string]map[string]bool{} // id -> kind -> present

	for _, e := range entries {
		m := fixtureFile.FindStringSubmatch(e.Name())
		if m == nil {
			problems = append(problems, dir+"/"+e.Name()+": not a fixture file")
			continue
		}
		if seen[m[1]] == nil {
			seen[m[1]] = map[string]bool{}
		}
		seen[m[1]][m[2]] = true
	}

	// A fixture is a response plus its provenance sidecar; a request is
	// required for upstream, where every example documents both directions.
	required := []string{"response.xml", "meta.json", "fidelity.txt"}
	if prov == Upstream {
		required = append(required, "request.xml")
	}

	for id, kinds := range seen {
		var missing []string
		for _, want := range required {
			if !kinds[want] {
				missing = append(missing, want)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			problems = append(problems, dir+"/"+id+": missing "+strings.Join(missing, ", "))
		}
	}
	return problems
}

// TestCorpusIsDiscoverable pins the two claims LoadFixtures makes: it finds
// everything on disk, and it reads provenance from the tree rather than a list.
func TestCorpusIsDiscoverable(t *testing.T) {
	fixtures, err := LoadFixtures(TestdataDir)
	if err != nil {
		t.Fatalf("Failed to load fixtures: %v", err)
	}
	if len(fixtures) == 0 {
		t.Skip("Skipping: no fixtures present")
	}

	onDisk, err := filepath.Glob(filepath.Join(TestdataDir, "*", "*", "*_response.xml"))
	if err != nil {
		t.Fatalf("Failed to glob: %v", err)
	}
	if len(fixtures) != len(onDisk) {
		t.Errorf("LoadFixtures found %d fixtures, but %d response documents are on disk", len(fixtures), len(onDisk))
	}

	for _, f := range fixtures {
		if f.Provenance == Upstream && f.Meta.ByteFaithful {
			t.Errorf("%s/%s claims byte_faithful, but upstream examples are documentation, not wire bytes", f.Command, f.ID)
		}
	}
}
