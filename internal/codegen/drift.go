package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brennoo/go-gmp/internal/schema"
)

// Drift is one shape difference between a live-served command and its
// committed testdata/schema snapshot.
type Drift struct {
	Path   string
	Kind   string // "added", "removed", or "retyped"
	Detail string // populated for "retyped" only
}

// FamilyDrift collects the Drift entries found in one adopted command.
type FamilyDrift struct {
	Command string
	Nodes   []Drift
}

// DriftReport is what changed between a live gvmd's served schema and the
// committed testdata/schema/*.json snapshots — see docs/review-2026-08.md, T-17.
type DriftReport struct {
	// Gained is every command the live server offers that manifest.json does
	// not know about at all.
	Gained []string
	// Removed is every manifest.json command the live server no longer offers.
	Removed []string
	// Changed is every adopted command present on both sides whose request or
	// response shape differs from its committed snapshot, ranked by node
	// count descending so the biggest drift sorts first.
	Changed []FamilyDrift
}

// DetectDrift compares live (freshly extracted from a live gvmd) against the
// committed snapshots under snapshotDir, scoped by m. Gained/Removed cover
// every command m knows about; node-level shape diffing (Changed) is scoped
// to m.Adopted() only, since a review-only family has no committed generated
// behavior to protect. Defects are not compared: they describe upstream
// schema malformation, not the request/response shape this exists to track.
func DetectDrift(live []schema.Command, m Manifest, snapshotDir string) (DriftReport, error) {
	byName := make(map[string]schema.Command, len(live))
	for _, c := range live {
		byName[c.Name] = c
	}

	var report DriftReport
	for name := range byName {
		if _, ok := m.FamilyOf(name); !ok {
			report.Gained = append(report.Gained, name)
		}
	}
	sort.Strings(report.Gained)

	for _, name := range m.AllCommands() {
		if _, ok := byName[name]; !ok {
			report.Removed = append(report.Removed, name)
		}
	}
	sort.Strings(report.Removed)

	for _, name := range m.Adopted().AllCommands() {
		liveCmd, ok := byName[name]
		if !ok {
			continue // already recorded in Removed
		}
		snapshots, err := LoadSnapshots(snapshotDir, []string{name})
		if err != nil {
			return DriftReport{}, fmt.Errorf("loading committed snapshot for %s: %w", name, err)
		}
		committed := snapshots[0]

		var nodes []Drift
		nodes = append(nodes, diffTree(rootPath(committed.Request, liveCmd.Request), committed.Request, liveCmd.Request)...)
		nodes = append(nodes, diffTree(rootPath(committed.Response, liveCmd.Response), committed.Response, liveCmd.Response)...)
		if len(nodes) > 0 {
			report.Changed = append(report.Changed, FamilyDrift{Command: name, Nodes: nodes})
		}
	}
	sort.Slice(report.Changed, func(i, j int) bool {
		if len(report.Changed[i].Nodes) != len(report.Changed[j].Nodes) {
			return len(report.Changed[i].Nodes) > len(report.Changed[j].Nodes)
		}
		return report.Changed[i].Command < report.Changed[j].Command
	})

	return report, nil
}

// rootPath names a root node the same way internal/schema/conform.go's
// Compare does, so a Drift.Path can be cross-checked against a fidelity
// finding's Path by string equality. Prefers want, since got may be nil (the
// whole request or response side vanished upstream).
func rootPath(want, got *schema.Node) string {
	if want != nil {
		return "/" + want.Name
	}
	if got != nil {
		return "/" + got.Name
	}
	return ""
}

// diffTree recursively diffs want against got, reporting an added/removed
// leaf for a nil side and, for a node present on both, a retyped entry (if
// its Type or Card changed) followed by its children's own diffs.
func diffTree(path string, want, got *schema.Node) []Drift {
	switch {
	case want == nil && got == nil:
		return nil
	case want == nil:
		return []Drift{{Path: path, Kind: "added"}}
	case got == nil:
		return []Drift{{Path: path, Kind: "removed"}}
	}

	var out []Drift
	if detail := retypeDetail(want, got); detail != "" {
		out = append(out, Drift{Path: path, Kind: "retyped", Detail: detail})
	}
	out = append(out, diffChildren(path, want.Children, got.Children)...)
	return out
}

// diffChildren pairs want's and got's children by (kind, name) and, within a
// shared key, by occurrence order — get_filters's schema legitimately
// declares two same-named <filters> children (see internal/schema/conform.go's
// walk, which pairs the same way), so a plain name-keyed map would conflate
// them instead of comparing each to its own counterpart.
func diffChildren(parentPath string, want, got []*schema.Node) []Drift {
	wantByKey := groupByKey(want)
	gotByKey := groupByKey(got)

	var out []Drift
	for _, key := range unionKeys(want, got) {
		ws, gs := wantByKey[key], gotByKey[key]
		n := len(ws)
		if len(gs) > n {
			n = len(gs)
		}
		for i := 0; i < n; i++ {
			var w, g *schema.Node
			if i < len(ws) {
				w = ws[i]
			}
			if i < len(gs) {
				g = gs[i]
			}
			var path string
			if w != nil {
				path = w.Path(parentPath)
			} else {
				path = g.Path(parentPath)
			}
			out = append(out, diffTree(path, w, g)...)
		}
	}
	return out
}

func groupByKey(nodes []*schema.Node) map[string][]*schema.Node {
	out := make(map[string][]*schema.Node, len(nodes))
	for _, n := range nodes {
		k := nodeKey(n)
		out[k] = append(out[k], n)
	}
	return out
}

func unionKeys(want, got []*schema.Node) []string {
	seen := make(map[string]bool)
	var keys []string
	for _, n := range want {
		if k := nodeKey(n); !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	for _, n := range got {
		if k := nodeKey(n); !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func nodeKey(n *schema.Node) string { return n.Kind.String() + ":" + n.Name }

func retypeDetail(want, got *schema.Node) string {
	var parts []string
	if want.Type != got.Type {
		parts = append(parts, fmt.Sprintf("type: %s -> %s", want.Type, got.Type))
	}
	if want.Card != got.Card {
		parts = append(parts, fmt.Sprintf("card: %s -> %s", want.Card, got.Card))
	}
	return strings.Join(parts, "; ")
}

// String renders report for a human, ranking Gained/Removed alphabetically
// and Changed by how much drifted (see DetectDrift) so the output prioritizes
// work instead of just saying "something changed".
func (r DriftReport) String() string {
	if len(r.Gained) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0 {
		return "No schema drift detected.\n"
	}

	var b strings.Builder
	if len(r.Gained) > 0 {
		fmt.Fprintf(&b, "Gained (upstream offers, manifest.json does not know about):\n")
		for _, name := range r.Gained {
			fmt.Fprintf(&b, "  + %s\n", name)
		}
	}
	if len(r.Removed) > 0 {
		fmt.Fprintf(&b, "Removed (manifest.json expects, upstream no longer offers):\n")
		for _, name := range r.Removed {
			fmt.Fprintf(&b, "  - %s\n", name)
		}
	}
	if len(r.Changed) > 0 {
		fmt.Fprintf(&b, "Changed (adopted families, shape differs from the committed snapshot):\n")
		for _, fd := range r.Changed {
			fmt.Fprintf(&b, "  %s (%d node(s)):\n", fd.Command, len(fd.Nodes))
			for _, d := range fd.Nodes {
				if d.Detail != "" {
					fmt.Fprintf(&b, "    %s %s: %s\n", d.Kind, d.Path, d.Detail)
				} else {
					fmt.Fprintf(&b, "    %s %s\n", d.Kind, d.Path)
				}
			}
		}
	}
	return b.String()
}
