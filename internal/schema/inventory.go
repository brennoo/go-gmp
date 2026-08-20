package schema

import (
	"fmt"
	"sort"
	"strings"
)

// ScalarRow is one schema-optional scalar leaf: a field with no children on
// either side, where the schema declares CardOptional. T-24's audit exists to
// answer, for every such field, whether its legal wire domain includes the Go
// zero value and, if so, whether the generated type is a pointer.
type ScalarRow struct {
	Command    string
	Side       string // "request" or "response"
	Path       string
	SchemaType string
	Alts       []string // non-empty for <alts>-constrained (enum-like) fields
	GoType     string
	Pointer    bool
	OmitEmpty  bool
}

// Inventory is the full corpus of optional scalar leaves, built by wiring
// Comparer.OnOptionalScalar into a full pass over every adopted command.
type Inventory struct {
	Rows []ScalarRow
}

// String renders the inventory deterministically: sorted by command, then
// side, then path, with nothing environment-dependent in the output.
func (inv Inventory) String() string {
	rows := append([]ScalarRow(nil), inv.Rows...)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Command != rows[j].Command {
			return rows[i].Command < rows[j].Command
		}
		if rows[i].Side != rows[j].Side {
			return rows[i].Side < rows[j].Side
		}
		return rows[i].Path < rows[j].Path
	})

	var b strings.Builder
	b.WriteString(`# Optional scalar field inventory

Regenerate with: go test ./internal/schema -run TestOptionalScalarInventory -update

T-24: every schema-optional field with no children on either side (int,
decimal, GMPBool/GMPAtoiBool, plain string, including <alts>-constrained
"enum-like" strings). Pointer/OmitEmpty are the generated Go field's actual
shape. A row with Alts populated and Pointer=false is a candidate for T-19's
enum type-design work; T-24 makes no type decision for those fields.

`)
	fmt.Fprintf(&b, "total optional scalar fields: %d\n\n", len(rows))
	for _, r := range rows {
		fmt.Fprintf(&b, "%-28s %-9s %-40s %-10s", r.Command, r.Side, r.Path, r.SchemaType)
		if len(r.Alts) > 0 {
			fmt.Fprintf(&b, " alts=[%s]", strings.Join(r.Alts, ","))
		}
		fmt.Fprintf(&b, " -> %s pointer=%t omitempty=%t\n", r.GoType, r.Pointer, r.OmitEmpty)
	}
	return b.String()
}
