package schema

import (
	"fmt"
	"sort"
	"strings"
)

// Report renders the conformance report. The report is the deliverable, so it
// leads with the totals a reviewer needs to judge scale and only then goes
// per-command.
type Report struct {
	SchemaDir string
	Results   []Result
	// Skipped names commands present in the schema with no Go type, and Go
	// types with no schema file. Both are drift, and both belong in the totals.
	SchemaOnly []string
	GoOnly     []string
}

// String renders the whole report. Output is deterministic: every list is
// sorted and nothing carries a timestamp or a path outside the repository.
func (r Report) String() string {
	var b strings.Builder
	r.writeHeader(&b)
	r.writeTotals(&b)
	r.writeDefects(&b)
	r.writeWorst(&b)
	r.writeConfirmed(&b)
	r.writePerCommand(&b)
	return b.String()
}

func (r Report) writeHeader(b *strings.Builder) {
	b.WriteString(`# GMP schema conformance report

Regenerate with: go test ./internal/schema -run TestConformanceReport -update

What this is: the GMP schema DSL, read directly, diffed against the
hand-written structs in commands/. The schema revision compared here is the one
vendored at ` + r.SchemaDir + `, which matches commands/ one-for-one, so a
difference found here is a transcription error and not version drift.

This is a *schema* oracle. internal/gmptest is a *runtime* oracle: it
round-trips captured documents and records what they lose. Findings marked
[fidelity] were reached by both, from different evidence, by tools that share
no code — those are the highest-confidence class here.

Categories:
  MISSING       the schema declares it, no Go field carries it
  EXTRA         a Go field with nothing behind it in the schema
  CARDINALITY   repeated-vs-slice, or an optional whose absence Go cannot express
  KIND          an attribute modelled as an element, or the reverse
  OMITEMPTY     omitempty where the schema says required, or absent where optional
  TYPE          the Go type cannot hold the schema's declared lexical space
  SCHEMA-DEFECT the schema is wrong; nothing is asked of commands/

Findings marked [suppressed] are counted but not charged to commands/. They are
almost all request-side over-constraint: the schema marks create_report's
report/@id, config_id, format_id, extension and content_type all required,
though real payloads only ever populate one submission path, so response-side
optionality is trusted and request-side is not.

`)
}

func (r Report) writeTotals(b *strings.Builder) {
	var (
		byCat        = map[string]int{}
		bySup        = map[string]int{}
		bySide       = map[string]int{}
		conformant   int
		suppressed   int
		confirmed    int
		actionable   int
		withFindings int
	)
	for _, res := range r.Results {
		if res.Conformant() {
			conformant++
		} else {
			withFindings++
		}
		for _, f := range res.Findings {
			byCat[f.Category]++
			switch {
			case f.Suppressed:
				bySup[f.Category]++
				suppressed++
			case f.Category != CatDefect:
				actionable++
				bySide[f.Side]++
			}
			if f.Confirmed {
				confirmed++
			}
		}
	}

	b.WriteString("## Totals\n\n")
	fmt.Fprintf(b, "commands compared:            %d\n", len(r.Results))
	fmt.Fprintf(b, "fully conformant:             %d\n", conformant)
	fmt.Fprintf(b, "with actionable findings:     %d\n", withFindings)
	fmt.Fprintf(b, "actionable findings:          %d  (request %d, response %d)\n",
		actionable, bySide["request"], bySide["response"])
	fmt.Fprintf(b, "suppressed findings:          %d\n", suppressed)
	fmt.Fprintf(b, "confirmed by fidelity corpus: %d\n", confirmed)
	fmt.Fprintf(b, "in the schema, not in Go:     %d%s\n", len(r.SchemaOnly), inline(r.SchemaOnly))
	fmt.Fprintf(b, "in Go, not in the schema:     %d%s\n\n", len(r.GoOnly), inline(r.GoOnly))

	b.WriteString("### Findings by category\n\n")
	fmt.Fprintf(b, "  %-14s %10s %10s\n", "", "actionable", "suppressed")
	for _, c := range Categories {
		fmt.Fprintf(b, "  %-14s %10d %10d\n", c, byCat[c]-bySup[c], bySup[c])
	}
	b.WriteString("\n")
}

func inline(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return "  (" + strings.Join(ss, ", ") + ")"
}

func (r Report) writeDefects(b *strings.Builder) {
	classes := map[string]int{}
	cmds := map[string]map[string]bool{}
	for _, res := range r.Results {
		for _, f := range res.Findings {
			if f.Category != CatDefect {
				continue
			}
			c := strings.SplitN(f.Detail, ":", 2)[0]
			classes[c]++
			if cmds[c] == nil {
				cmds[c] = map[string]bool{}
			}
			cmds[c][res.Command] = true
		}
	}

	b.WriteString("### Schema defects by class\n\n")
	b.WriteString(`How much of the schema is unusable as written. None of these are
transcription errors; they are upstream bugs the parser has to survive.

These counts are lower than the ones in docs/review-2026-08.md T-9, which were
taken against upstream main (155 commands). Two reasons, both expected:

  - This revision predates the commands that carry most of the damage. The 68
    dangling <e> references upstream are concentrated in get_scan_report,
    get_agent_groups, get_report_tls_certificates and get_integration_configs,
    none of which exist here.
  - Only the shape is walked. A <type> naming UUID inside a <filter_keywords>
    <column> describes filter-term syntax, not an element or an attribute, so
    it is out of scope: 28 of the 31 UUID uses are of that kind.

`)
	keys := sortedKeys(classes)
	for _, k := range keys {
		fmt.Fprintf(b, "  %-16s %5d sites in %d commands\n", k, classes[k], len(cmds[k]))
	}
	if len(keys) == 0 {
		b.WriteString("  (none)\n")
	}
	b.WriteString("\n")
}

// writeWorst ranks commands by actionable findings, which is the closest thing
// to "how much of this matters" the report can state without judging each
// finding by hand.
func (r Report) writeWorst(b *strings.Builder) {
	type row struct {
		name string
		n    int
		cats map[string]int
	}
	rows := make([]row, 0, len(r.Results))
	for _, res := range r.Results {
		x := row{name: res.Command, cats: map[string]int{}}
		for _, f := range res.Findings {
			if f.Category == CatDefect || f.Suppressed {
				continue
			}
			x.n++
			x.cats[f.Category]++
		}
		if x.n > 0 {
			rows = append(rows, x)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].n != rows[j].n {
			return rows[i].n > rows[j].n
		}
		return rows[i].name < rows[j].name
	})

	b.WriteString("### Commands by actionable findings\n\n")
	const top = 20
	for i, x := range rows {
		if i == top {
			fmt.Fprintf(b, "  ... and %d more with fewer\n", len(rows)-top)
			break
		}
		fmt.Fprintf(b, "  %-26s %4d   %s\n", x.name, x.n, summarise(x.cats))
	}
	b.WriteString("\n")
}

func (r Report) writeConfirmed(b *strings.Builder) {
	var lines []string
	for _, res := range r.Results {
		for _, f := range res.Findings {
			if f.Confirmed {
				lines = append(lines, fmt.Sprintf("  %-24s %s", res.Command, f.Path))
			}
		}
	}
	sort.Strings(lines)

	b.WriteString("### Confirmed by both oracles\n\n")
	b.WriteString(`Paths the schema says exist and that a captured document was also seen to
lose. Fix these first: neither oracle is guessing.

Response side only: the fidelity corpus keeps per-path DROPPED lists for
responses and only aggregate counts for requests, so a request-side finding can
never carry this mark. That is a gap in the cross-reference, not evidence that
requests are in better shape — get_notes models none of its seven request
attributes, and one of its own captured request fixtures uses three of them.

`)
	if len(lines) == 0 {
		b.WriteString("  (none)\n\n")
		return
	}
	b.WriteString(strings.Join(lines, "\n"))
	b.WriteString("\n\n")
}

func (r Report) writePerCommand(b *strings.Builder) {
	b.WriteString("## Per command\n\n")
	for _, res := range r.Results {
		counts := map[string]int{}
		for _, f := range res.Findings {
			counts[f.Category]++
		}
		fmt.Fprintf(b, "### %s  (%s)\n\n", res.Command, res.File)
		if len(res.Findings) == 0 {
			b.WriteString("conformant\n\n")
			continue
		}
		fmt.Fprintf(b, "%s\n\n", summarise(counts))
		for _, f := range res.Findings {
			fmt.Fprintf(b, "  %-13s %-8s %s\n", f.Category, f.Side, f.Path)
			fmt.Fprintf(b, "                %s%s\n", f.Detail, tags(f))
		}
		b.WriteString("\n")
	}
}

func summarise(counts map[string]int) string {
	var parts []string
	for _, c := range Categories {
		if counts[c] > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", c, counts[c]))
		}
	}
	return strings.Join(parts, "  ")
}

func tags(f Finding) string {
	var out string
	if f.Confirmed {
		out += "  [fidelity]"
	}
	if f.Suppressed {
		out += "  [suppressed: " + f.Reason + "]"
	}
	return out
}

func sortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
