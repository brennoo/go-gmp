package codegen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/internal/schema"
)

// This file is the pre-swap validation step: before commands/create_note.go,
// modify_note.go and get_notes.go are replaced by Generate's output, this is
// what proves the replacement is a match-or-improve on schema conformance,
// using the same schema.Comparer/schema.ModelOf machinery
// internal/schema/conformance_test.go already trusts for the 122 hand-written
// commands.
//
// # Design choice: how the "after" side gets a real reflect.Type
//
// schema.ModelOf takes a Go value and reflects it; it cannot work from source
// text or from the schema.Node tree Generate consumed; that would just be
// comparing the schema to itself. The generated code has to actually be
// compiled, and reflected, somewhere.
//
// The complication is that the generated types collide by name with the
// hand-written ones (both declare `type Note struct`, `type CreateNote
// struct`, ...) and both are still needed side by side: the hand-written
// ones for "before", the generated ones for "after". They cannot both be
// package commands in the same build.
//
// TestGeneratedCodeCompiles (generate_test.go) solves the same collision for
// its own purpose — proving the generated code builds independently of
// commands/ — by copying it into a throwaway module (its own go.mod, a
// replace directive back to this repo) with the package clause rewritten.
// That proves compilation, but a throwaway *module* cannot be reflected from
// here: reflect.TypeOf needs the type compiled into this process's binary,
// and a separate module is a separate `go build`, not something this test
// binary can link against without cgo-style plugin loading (fragile,
// platform-specific, and unnecessary for a prototype covering three
// commands).
//
// So GeneratedGoRoots takes a different shortcut, using the fact that the
// generated code already lives inside this module rather than needing to
// pretend it doesn't:
//
//  1. Write Generate's output, unmodified, into a throwaway *directory*
//     (not a throwaway module) under internal/codegen/testdata — still part
//     of module github.com/brennoo/go-gmp, so it can import
//     internal/schema (an internal package, so a foreign module could not
//     import it even with a replace directive: Go's internal-import rule
//     is keyed on the importing package's own import path, not on disk
//     layout).
//  2. Write a tiny generated `package main` next to it that imports that
//     directory, calls schema.ModelOf on each command's request/response
//     zero value — the exact function TestGoModelMatchesRegistry already
//     trusts for the hand-written side — and prints the resulting
//     schema.GoRoot values as JSON.
//  3. Run it with `go run` and decode the JSON.
//
// This is a leaner version of the plan's option C (a subprocess "reflect
// server" that re-derives ModelOf's output): because everything stays
// inside one module, no go.mod/replace scaffolding is needed at all, only a
// scratch directory under testdata (which `go build ./...`/`go vet ./...`
// already skip by name, and which is removed before this function returns).
// It is a real, out-of-process reflection of the real generated source,
// not a self-comparison of the schema against itself (the risk flagged for
// option A) and not a second copy of ModelOf's logic (the risk of a
// go/types-based static approach).
const modulePath = "github.com/brennoo/go-gmp"

// ErrNoGoToolchain is returned by GeneratedGoRoots when no "go" binary is on
// PATH. Callers that can skip (tests) should check for it with errors.Is.
var ErrNoGoToolchain = errors.New("codegen: no go toolchain on PATH")

// CommandRoots is one command's reflected request/response shape.
type CommandRoots struct {
	Request  schema.GoRoot
	Response schema.GoRoot
}

// BeforeGoRoots reflects the commands/ types as they stand right now, for
// every family the manifest names. Unlike the generated side, this needs no
// subprocess: commands is a normal, non-internal package, and internal/codegen
// already depends on it for exactly this purpose, so schema.ModelOf can be
// called on it directly, in-process, the same way conformance_test.go calls it
// on gmptest.Registry's entries.
//
// "Before" means different things per family, deliberately. For an adopted
// family (note, override, role, group, tag) these types are already the
// generator's output, so the comparison is a drift detector on the committed
// files. For a review-only family (none currently) these would still be the
// hand-written originals, so the comparison would be the pre-swap gate note,
// override, role, group, and tag all went through.
//
// The map is written out by hand rather than derived from the manifest because
// reflection needs the concrete types: adding a family here is the one place
// that cannot be manifest-driven, and Validate fails loudly on a command it
// has no entry for rather than silently skipping it.
func BeforeGoRoots() map[string]CommandRoots {
	return map[string]CommandRoots{
		"create_note": {
			Request:  schema.ModelOf(&commands.CreateNote{}),
			Response: schema.ModelOf(&commands.CreateNoteResponse{}),
		},
		"modify_note": {
			Request:  schema.ModelOf(&commands.ModifyNote{}),
			Response: schema.ModelOf(&commands.ModifyNoteResponse{}),
		},
		"get_notes": {
			Request:  schema.ModelOf(&commands.GetNotes{}),
			Response: schema.ModelOf(&commands.GetNotesResponse{}),
		},
		"create_override": {
			Request:  schema.ModelOf(&commands.CreateOverride{}),
			Response: schema.ModelOf(&commands.CreateOverrideResponse{}),
		},
		"modify_override": {
			Request:  schema.ModelOf(&commands.ModifyOverride{}),
			Response: schema.ModelOf(&commands.ModifyOverrideResponse{}),
		},
		"get_overrides": {
			Request:  schema.ModelOf(&commands.GetOverrides{}),
			Response: schema.ModelOf(&commands.GetOverridesResponse{}),
		},
		"create_role": {
			Request:  schema.ModelOf(&commands.CreateRole{}),
			Response: schema.ModelOf(&commands.CreateRoleResponse{}),
		},
		"modify_role": {
			Request:  schema.ModelOf(&commands.ModifyRole{}),
			Response: schema.ModelOf(&commands.ModifyRoleResponse{}),
		},
		"get_roles": {
			Request:  schema.ModelOf(&commands.GetRoles{}),
			Response: schema.ModelOf(&commands.GetRolesResponse{}),
		},
		"create_group": {
			Request:  schema.ModelOf(&commands.CreateGroup{}),
			Response: schema.ModelOf(&commands.CreateGroupResponse{}),
		},
		"modify_group": {
			Request:  schema.ModelOf(&commands.ModifyGroup{}),
			Response: schema.ModelOf(&commands.ModifyGroupResponse{}),
		},
		"get_groups": {
			Request:  schema.ModelOf(&commands.GetGroups{}),
			Response: schema.ModelOf(&commands.GetGroupsResponse{}),
		},
		"create_tag": {
			Request:  schema.ModelOf(&commands.CreateTag{}),
			Response: schema.ModelOf(&commands.CreateTagResponse{}),
		},
		"modify_tag": {
			Request:  schema.ModelOf(&commands.ModifyTag{}),
			Response: schema.ModelOf(&commands.ModifyTagResponse{}),
		},
		"get_tags": {
			Request:  schema.ModelOf(&commands.GetTags{}),
			Response: schema.ModelOf(&commands.GetTagsResponse{}),
		},
	}
}

// packageDir is internal/codegen's own directory, resolved once from the
// running binary's source location rather than from the process's current
// working directory. Working-directory-relative paths ("../..", "testdata")
// only resolve correctly when something happens to run from this package's
// directory (true for `go test ./internal/codegen`, not guaranteed for a
// caller elsewhere in the tree), so anything that needs to find this
// package's testdata directory from outside a test binary needs this instead.
func packageDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("codegen: runtime.Caller could not resolve this file's path")
	}
	return filepath.Dir(file), nil
}

// GeneratedGoRoots runs Generate(cmds, ov) — one family's commands, the unit
// Generate is scoped to, see GenerateAll — and reflects its output the way
// BeforeGoRoots reflects the hand-written types, via the out-of-process
// mechanism this file's package doc explains. The temporary directory it
// creates under testdata/ is removed before it returns, on both the success
// and the error path.
func GeneratedGoRoots(cmds []schema.Command, ov map[string]Override) (map[string]CommandRoots, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", ErrNoGoToolchain, err)
	}

	files, err := Generate(cmds, ov)
	if err != nil {
		return nil, fmt.Errorf("codegen: generate: %w", err)
	}

	pkgDir, err := packageDir()
	if err != nil {
		return nil, err
	}
	moduleRoot := filepath.Dir(filepath.Dir(pkgDir)) // .../internal/codegen -> .../internal -> module root

	scratch, err := os.MkdirTemp(filepath.Join(pkgDir, "testdata"), "genreflect-")
	if err != nil {
		return nil, fmt.Errorf("codegen: creating scratch dir: %w", err)
	}
	defer os.RemoveAll(scratch) //nolint:errcheck // best-effort cleanup of a throwaway dir

	genDir := filepath.Join(scratch, "gen")
	mainDir := filepath.Join(scratch, "main")
	if err := os.MkdirAll(genDir, 0o750); err != nil {
		return nil, fmt.Errorf("codegen: creating scratch gen dir: %w", err)
	}
	if err := os.MkdirAll(mainDir, 0o750); err != nil {
		return nil, fmt.Errorf("codegen: creating scratch main dir: %w", err)
	}

	// Unlike TestGeneratedCodeCompiles, the generated files' package clause
	// (package commands) does not need rewriting: genDir is a distinct
	// import path from the real commands/, so there is no declaration
	// collision, only a same-named package living at a different path — Go
	// does not require a package's name to match its directory or import
	// path.
	if err := writeScratchPackage(genDir, files); err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(moduleRoot, scratch)
	if err != nil {
		return nil, fmt.Errorf("codegen: relativizing scratch dir: %w", err)
	}
	genImportPath := modulePath + "/" + filepath.ToSlash(filepath.Join(rel, "gen"))

	src, err := reflectHelperSource(cmds, genImportPath)
	if err != nil {
		return nil, err
	}
	if err := writeFile(filepath.Join(mainDir, "main.go"), src); err != nil {
		return nil, fmt.Errorf("codegen: writing the reflect helper: %w", err)
	}

	pattern := "./" + filepath.ToSlash(filepath.Join(rel, "main"))
	cmd := exec.CommandContext(context.Background(), goBin, "run", pattern) //nolint:gosec // goBin is from exec.LookPath, pattern is built from a temp dir this function created
	cmd.Dir = moduleRoot
	// GOWORK=off keeps a developer's workspace file from redirecting module
	// resolution; GOFLAGS is cleared so an inherited -mod flag cannot apply.
	// Both mirror TestGeneratedCodeCompiles.
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("codegen: reflect helper failed: %w\n%s", err, stderr.String())
	}

	var decoded map[string]struct {
		Request  schema.GoRoot `json:"request"`
		Response schema.GoRoot `json:"response"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		return nil, fmt.Errorf("codegen: decoding the reflect helper's output: %w\n%s", err, stdout.String())
	}

	out := make(map[string]CommandRoots, len(decoded))
	for name, r := range decoded {
		out[name] = CommandRoots{Request: r.Request, Response: r.Response}
	}
	return out, nil
}

// reflectHelperSource generates the tiny `package main` GeneratedGoRoots runs:
// it imports the just-written generated package under genImportPath, calls
// schema.ModelOf on every command's request and response zero value — the
// same function, not a reimplementation of it — and JSON-encodes the result
// to stdout.
func reflectHelperSource(cmds []schema.Command, genImportPath string) (string, error) {
	var b strings.Builder
	b.WriteString("// Code generated by internal/codegen for validation. DO NOT EDIT.\n\n")
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n\n")
	fmt.Fprintf(&b, "\tgen %q\n", genImportPath)
	b.WriteString("\t\"github.com/brennoo/go-gmp/internal/schema\"\n")
	b.WriteString(")\n\n")
	b.WriteString("type roots struct {\n")
	b.WriteString("\tRequest  schema.GoRoot `json:\"request\"`\n")
	b.WriteString("\tResponse schema.GoRoot `json:\"response\"`\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() {\n")
	b.WriteString("\tout := map[string]roots{}\n")
	for i := range cmds {
		name := cmds[i].Name
		reqType := goName(name)
		respType := reqType + "Response"
		fmt.Fprintf(&b, "\tout[%q] = roots{Request: schema.ModelOf(&gen.%s{}), Response: schema.ModelOf(&gen.%s{})}\n",
			name, reqType, respType)
	}
	b.WriteString("\tif err := json.NewEncoder(os.Stdout).Encode(out); err != nil {\n")
	b.WriteString("\t\tpanic(err)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return "", fmt.Errorf("codegen: reflect helper source does not parse: %w\n%s", err, b.String())
	}
	return string(formatted), nil
}

// Validate builds the before/after schema-conformance comparison for cmds —
// normally the 3-command note family LoadSnapshots returns — as one
// schema.Result per command against the current hand-written commands/
// types and one against what Generate(cmds, ov) would produce. Both sides
// are produced by schema.Comparer.Compare, the same function
// TestConformanceReport uses for all 122 commands: nothing here
// reimplements its diff logic.
//
// dropped is the fidelity corpus's per-command dropped-path set
// (schema.LoadDroppedPaths); passing nil just means no [fidelity]
// cross-reference, not an error.
func Validate(cmds []schema.Command, ov map[string]Override, dropped map[string]map[string]bool) (before, after []schema.Result, err error) {
	beforeRoots := BeforeGoRoots()
	afterRoots, err := GeneratedGoRoots(cmds, ov)
	if err != nil {
		return nil, nil, err
	}

	cmp := &schema.Comparer{Dropped: dropped}
	for i := range cmds {
		cmd := cmds[i]
		br, ok := beforeRoots[cmd.Name]
		if !ok {
			return nil, nil, fmt.Errorf("codegen: BeforeGoRoots has no hand-written model for %s", cmd.Name)
		}
		ar, ok := afterRoots[cmd.Name]
		if !ok {
			return nil, nil, fmt.Errorf("codegen: the reflect helper reported no model for %s", cmd.Name)
		}
		before = append(before, cmp.Compare(cmd, br.Request, br.Response))
		after = append(after, cmp.Compare(cmd, ar.Request, ar.Response))
	}
	return before, after, nil
}
