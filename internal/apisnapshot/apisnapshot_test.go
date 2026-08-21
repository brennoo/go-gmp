// Package apisnapshot renders the module's public API — every exported
// declaration in every publicly importable package — into a single
// deterministic, committed golden. See TestPublicAPISnapshot.
package apisnapshot

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Declared locally rather than reused from internal/gmptest: importing
// gmptest for one bool would pull the whole fidelity corpus package in.
var update = flag.Bool("update", false, "rewrite the committed public API golden")

const goldenPath = "testdata/PUBLIC_API.txt"

// pkg is one publicly importable package this module ships. The list below
// is fixed, not discovered by walking the tree, so internal/ and examples/
// restructuring never silently joins or leaves this surface.
type pkg struct {
	dir        string // relative to the repo root
	importPath string
}

var targetPackages = []pkg{
	{dir: ".", importPath: "github.com/brennoo/go-gmp"},
	{dir: "client", importPath: "github.com/brennoo/go-gmp/client"},
	{dir: "commands", importPath: "github.com/brennoo/go-gmp/commands"},
	{dir: "commands/filtering", importPath: "github.com/brennoo/go-gmp/commands/filtering"},
	{dir: "commands/pagination", importPath: "github.com/brennoo/go-gmp/commands/pagination"},
	{dir: "connections", importPath: "github.com/brennoo/go-gmp/connections"},
}

// TestPublicAPISnapshot builds the public API golden and checks it against
// the committed one. Regenerate with:
// go test ./internal/apisnapshot -run TestPublicAPISnapshot -update
//
// Any diff here is a reviewed API decision, never a mechanical fix-the-build
// step: this covers every exported declaration this module has promised callers.
func TestPublicAPISnapshot(t *testing.T) {
	got, err := buildSnapshot()
	if err != nil {
		t.Fatalf("Failed to build the public API snapshot: %v", err)
	}

	if *update {
		if err := os.WriteFile(goldenPath, []byte(got), 0o600); err != nil {
			t.Fatalf("Failed to write %s: %v", goldenPath, err)
		}
		return
	}

	want, err := os.ReadFile(goldenPath) //nolint:gosec // a fixed testdata path
	if err != nil {
		t.Fatalf("No committed public API snapshot (%v).\nRun with -update to create it.", err)
	}
	if got != string(want) {
		t.Fatalf("%s changed. Re-run with -update and review the diff.\n"+
			"committed %d bytes, now %d bytes\n%s",
			goldenPath, len(want), len(got), firstDiff(string(want), got))
	}
}

// TestPublicAPISnapshotIsDeterministic guards the golden's only real
// precondition: building it twice from the same on-disk source must produce
// byte-identical output.
func TestPublicAPISnapshotIsDeterministic(t *testing.T) {
	a, err := buildSnapshot()
	if err != nil {
		t.Fatalf("Failed to build the public API snapshot: %v", err)
	}
	b, err := buildSnapshot()
	if err != nil {
		t.Fatalf("Failed to build the public API snapshot: %v", err)
	}
	if a != b {
		t.Fatalf("Building the public API snapshot twice from the same input produced different output:\n%s", firstDiff(a, b))
	}
}

// firstDiff points at the first differing line, so a failure names the
// change instead of dumping hundreds of kilobytes of declarations.
func firstDiff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) && i < len(g); i++ {
		if w[i] != g[i] {
			return "first difference at line " + strconv.Itoa(i+1) + ":\n  committed: " + w[i] + "\n  now:       " + g[i]
		}
	}
	return "one file is a prefix of the other"
}

// entry is one exported declaration, ready to render.
type entry struct {
	sortKey string
	node    ast.Node
}

// buildSnapshot parses every target package's non-test .go files directly
// from the repo tree (never vendor/XML, nothing gitignored) and renders one
// section per package.
func buildSnapshot() (string, error) {
	var b strings.Builder
	b.WriteString(header)

	total := 0
	sections := make([]string, len(targetPackages))
	for i, p := range targetPackages {
		entries, err := collectPackage(p.dir)
		if err != nil {
			return "", fmt.Errorf("%s: %w", p.importPath, err)
		}
		total += len(entries)
		sections[i] = renderSection(p.importPath, entries)
	}

	fmt.Fprintf(&b, "total packages: %d\ntotal exported declarations: %d\n\n", len(targetPackages), total)
	for _, s := range sections {
		b.WriteString(s)
	}
	return b.String(), nil
}

const header = `# Public API snapshot

Regenerate with: go test ./internal/apisnapshot -run TestPublicAPISnapshot -update

T-27: every exported declaration (types with full field lists and tags,
func/method signatures, consts with values, vars) across every publicly
importable package this module ships — github.com/brennoo/go-gmp and its
client, commands, commands/filtering, commands/pagination and connections
subpackages. Any diff here is a reviewed API decision, never a mechanical
fix-the-build step.

Deliberately excluded: internal/... and examples/... (not importable by
consumers), all _test.go files, and doc comments (churn without API meaning
— struct tags are kept, since xml/gmpwire/gmpoverride tags are wire-behavior
API here, not documentation).

client's section is one line by design: client.New returns the gmp.Client
interface, and every method on its unexported *client implementation is only
reachable through that interface, whose signature this snapshot already
captures in the root package's Client section. A method added to *client
without a matching Client entry is a compile error, not something this
golden needs to also catch.

`

// collectPackage parses dir's non-test .go files and returns every exported
// top-level declaration, sorted for a stable, diff-friendly order.
func collectPackage(dir string) ([]entry, error) {
	full := filepath.Join("../..", dir)
	files, err := os.ReadDir(full)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	fset := token.NewFileSet()
	var entries []entry
	for _, name := range names {
		path := filepath.Join(full, name)
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		entries = append(entries, exportedDecls(f)...)
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].sortKey < entries[j].sortKey })
	return entries, nil
}

// exportedDecls extracts every exported top-level declaration from one
// parsed file, splitting a multi-spec type decl into one entry per type (so
// an unrelated sibling type in the same `type (...)` block never masks a
// change) while keeping a const/var GenDecl whole (iota position depends on
// spec order within the block, so splitting would misrepresent it).
func exportedDecls(f *ast.File) []entry {
	var out []entry
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for _, spec := range d.Specs {
					ts := spec.(*ast.TypeSpec)
					if ts.Name.IsExported() {
						out = append(out, entry{
							sortKey: ts.Name.Name,
							node:    &ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{ts}},
						})
					}
				}
			case token.CONST, token.VAR:
				if anyExportedSpec(d.Specs) {
					out = append(out, entry{sortKey: firstName(d.Specs), node: d})
				}
			}
		case *ast.FuncDecl:
			if e, ok := exportedFunc(d); ok {
				out = append(out, e)
			}
		}
	}
	return out
}

func anyExportedSpec(specs []ast.Spec) bool {
	for _, spec := range specs {
		vs := spec.(*ast.ValueSpec)
		for _, name := range vs.Names {
			if name.IsExported() {
				return true
			}
		}
	}
	return false
}

// firstName sorts a const/var block by its first declared name, exported or
// not: iota-based blocks conventionally lead with an exported name anyway,
// and this only needs to be a stable, deterministic key, not a filter.
func firstName(specs []ast.Spec) string {
	vs := specs[0].(*ast.ValueSpec)
	return vs.Names[0].Name
}

// exportedFunc reports whether a func/method declaration belongs in the
// snapshot: a plain func iff its own name is exported, a method iff both its
// receiver's base type name and the method name are exported. Body is
// cleared before returning so the entry renders as a signature only.
func exportedFunc(d *ast.FuncDecl) (entry, bool) {
	if d.Recv == nil {
		if !d.Name.IsExported() {
			return entry{}, false
		}
		clone := *d
		clone.Body = nil
		return entry{sortKey: d.Name.Name, node: &clone}, true
	}

	recvType := d.Recv.List[0].Type
	base := receiverBaseIdent(recvType)
	if base == nil || !base.IsExported() || !d.Name.IsExported() {
		return entry{}, false
	}
	clone := *d
	clone.Body = nil
	return entry{sortKey: base.Name + "." + d.Name.Name, node: &clone}, true
}

// receiverBaseIdent unwraps a receiver type expression down to its base
// identifier. A generic receiver like *Iterator[T] parses as
// StarExpr{X: IndexExpr{X: Ident}} — the pointer wraps the type-parameter
// index, not the reverse — so the star is stripped first.
func receiverBaseIdent(expr ast.Expr) *ast.Ident {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	switch e := expr.(type) {
	case *ast.Ident:
		return e
	case *ast.IndexExpr:
		id, _ := e.X.(*ast.Ident)
		return id
	case *ast.IndexListExpr:
		id, _ := e.X.(*ast.Ident)
		return id
	}
	return nil
}

// renderSection prints one package's entries under its import-path heading.
func renderSection(importPath string, entries []entry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n", importPath)
	fset := token.NewFileSet()
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	for _, e := range entries {
		var buf strings.Builder
		if err := cfg.Fprint(&buf, fset, e.node); err != nil {
			fmt.Fprintf(&b, "// ERROR rendering %s: %v\n\n", e.sortKey, err)
			continue
		}
		b.WriteString(collapseBlankLines(buf.String()))
		b.WriteString("\n\n")
	}
	return b.String()
}

var blankRunRE = regexp.MustCompile(`\n[ \t]*\n(\n[ \t]*)+`)

// collapseBlankLines removes go/printer's dependence on the parsed source's
// original line spacing: dropping ParseComments hides comment text but a
// comment's blank-line gap still leaks through fset's line bookkeeping.
// Adding or deleting a doc comment must never change this golden.
func collapseBlankLines(s string) string {
	s = blankRunRE.ReplaceAllString(s, "\n\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.Join(lines, "\n")
}
