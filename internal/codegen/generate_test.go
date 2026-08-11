package codegen

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// goldenDir holds the generator's committed output. These are .go files under
// testdata, so the go tool ignores them: they are golden *bytes*, and the only
// thing that ever compiles them is TestGeneratedCodeCompiles below.
const goldenDir = "testdata/golden"

// TestGenerate is the test CI actually runs for T-9. It needs only the
// committed JSON snapshot, never vendor/XML, so stale generated output shows
// up as an ordinary test failure on any checkout.
func TestGenerate(t *testing.T) {
	got := generateOrFail(t)

	if *update {
		if err := os.RemoveAll(goldenDir); err != nil {
			t.Fatalf("Failed to clear %s: %v", goldenDir, err)
		}
		if err := os.MkdirAll(goldenDir, 0o750); err != nil {
			t.Fatalf("Failed to create %s: %v", goldenDir, err)
		}
		for name, src := range got {
			if err := os.WriteFile(filepath.Join(goldenDir, name), []byte(src), 0o600); err != nil {
				t.Fatalf("Failed to write %s: %v", name, err)
			}
		}
		return
	}

	entries, err := os.ReadDir(goldenDir)
	if err != nil {
		t.Fatalf("No committed goldens (%v).\nRun with -update to create them.", err)
	}
	committed := map[string]bool{}
	for _, e := range entries {
		committed[e.Name()] = true
		want, err := os.ReadFile(filepath.Join(goldenDir, e.Name())) //nolint:gosec // a fixed testdata path
		if err != nil {
			t.Fatalf("Failed to read %s: %v", e.Name(), err)
		}
		src, ok := got[e.Name()]
		if !ok {
			t.Errorf("%s is committed but the generator no longer emits it. Re-run with -update.", e.Name())
			continue
		}
		if src != string(want) {
			t.Errorf("%s changed. Re-run with -update and review the diff.\n"+
				"committed %d bytes, now %d bytes\n%s", e.Name(), len(want), len(src), firstDiff(string(want), src))
		}
	}
	for name := range got {
		if !committed[name] {
			t.Errorf("The generator emits %s but it is not committed. Re-run with -update.", name)
		}
	}
}

// TestGenerateIsDeterministic guards the golden's only real precondition:
// nothing in Generate may depend on map iteration order.
func TestGenerateIsDeterministic(t *testing.T) {
	a, b := generateOrFail(t), generateOrFail(t)
	if len(a) != len(b) {
		t.Fatalf("Two runs produced %d and %d files", len(a), len(b))
	}
	for name, src := range a {
		if b[name] != src {
			t.Errorf("%s differs between two runs of Generate", name)
		}
	}
}

// TestGeneratedCodeCompiles is acceptance criterion 2: the generated source
// has to build without the hand-written commands/create_note.go,
// modify_note.go or get_notes.go, or the generator could be quietly leaning on
// types it did not emit.
//
// It covers every family the manifest declares, review-only ones included, and
// compiles them all into one package — which is also the only check that two
// families never declare the same type name. Family-local dedup means each
// family names its shapes without consulting the others, so "note and override
// both describe an <nvt>" is resolved by the family prefix (NoteNVT vs
// OverrideNVT) rather than by a shared type; if that prefixing ever failed to
// separate two families, this test is where the redeclaration surfaces.
//
// Those three files still exist on disk, so the check cannot simply build the
// output as package commands: every type would be redeclared. Instead the
// files are copied into a throwaway module with only their package clause
// rewritten — that one line is the entire transformation — and the support
// types the override table targets (commands.GMPBool, commands.GMPTime) are
// re-supplied as aliases to the real commands package. Everything the
// generated code needs beyond the standard library therefore still comes from
// commands, and nothing else does.
func TestGeneratedCodeCompiles(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	// Every family's command files have to be in the build, or "it compiles"
	// would be a claim about whichever subset happened to be generated.
	for _, f := range manifestOrFail(t).Families {
		for _, c := range f.Commands {
			if _, ok := got[CommandFile(c)]; !ok {
				t.Fatalf("family %q: %s was not generated, so this test would not compile it", f.Name, CommandFile(c))
			}
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}

	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	// supportShim() re-exports GMPBool/GMPTime under package commands (see its
	// doc comment); the package clause is rewritten to genproto below along
	// with every generated file, so it lands in the same throwaway package.
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))

	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}

	// context.Background rather than t.Context: this module still declares
	// go 1.23, where testing.T has no Context method.
	cmd := exec.CommandContext(context.Background(), goBin, "build", "./...")
	cmd.Dir = dir
	// GOWORK=off keeps a developer's workspace file from redirecting the
	// replace; GOFLAGS is cleared so an inherited -mod flag cannot apply.
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Generated code does not compile on its own: %v\n%s", err, out)
	}
}

// mustWriteFile is writeFile (scratch.go) with the *testing.T ceremony a test
// wants but a reusable helper shouldn't carry.
func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := writeFile(path, content); err != nil {
		t.Fatalf("Failed to write %s: %v", path, err)
	}
}

// generateOrFail renders every family the manifest declares, review-only ones
// included: the goldens are where a review-only family's output lives, so the
// golden test has to see it.
func generateOrFail(t *testing.T) map[string]string {
	t.Helper()
	return generateManifest(t, manifestOrFail(t))
}

// generateAdoptedOrFail renders only the families whose output is committed
// under commands/ — what the CLI writes and what CheckOutput inspects.
func generateAdoptedOrFail(t *testing.T) map[string]string {
	t.Helper()
	return generateManifest(t, manifestOrFail(t).Adopted())
}

func generateManifest(t *testing.T, m Manifest) map[string]string {
	t.Helper()
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	out, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}
	return out
}

// manifestOrFail loads the committed manifest, which every test in this
// package sources its command list from so none of them carries a second
// hardcoded copy.
func manifestOrFail(t *testing.T) Manifest {
	t.Helper()
	m, err := LoadManifest(manifestFile)
	if err != nil {
		t.Fatalf("Failed to load the manifest: %v", err)
	}
	return m
}

// firstDiff points at the first differing line, so a failure names the change
// instead of dumping thirteen kilobytes of struct definitions.
func firstDiff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) && i < len(g); i++ {
		if w[i] != g[i] {
			return "first difference at line " + strconv.Itoa(i+1) + ":\n  committed: " + w[i] + "\n  now:       " + g[i]
		}
	}
	return "one file is a prefix of the other"
}
