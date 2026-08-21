package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// This file holds the small file-system/build helpers shared by
// generate_test.go's TestGeneratedCodeCompiles (proves the generated source
// builds on its own, in a throwaway module) and validate.go's GeneratedGoRoots
// (proves what it reflects to, in a throwaway package inside this module).
// Both need to write files to disk, both need the GMPBool/GMPTime alias shim,
// and only one of them needs a go.mod: factoring the overlap here is what
// keeps them from drifting into two slightly different copies of the same
// bookkeeping.

// supportShim re-exports, under their unqualified names, every support type
// the override table points the generated code at. It is derived from
// TypeOverrides rather than hardcoded so a new override cannot silently make
// a compile or reflection check untestable.
func supportShim() string {
	seen := map[string]bool{}
	var aliases []string
	for _, o := range TypeOverrides {
		short, ok := strings.CutPrefix(o.GoType, targetPackage+".")
		if !ok || seen[short] {
			continue
		}
		seen[short] = true
		aliases = append(aliases, "type "+short+" = "+targetPackage+"."+short)
	}
	sort.Strings(aliases)
	return "package " + targetPackage + "\n\nimport \"github.com/brennoo/go-gmp/" + targetPackage + "\"\n\n" +
		strings.Join(aliases, "\n") + "\n"
}

var goDirectiveRE = regexp.MustCompile(`(?m)^go (\S+)$`)

// goDirective reads the "go X.Y" line out of a go.mod, so a scratch module can
// declare the same language version as the real one instead of a hardcoded
// guess that silently drifts.
func goDirective(gomod string) (string, error) {
	data, err := os.ReadFile(gomod) //nolint:gosec // the module's own go.mod
	if err != nil {
		return "", err
	}
	m := goDirectiveRE.FindSubmatch(data)
	if m == nil {
		return "", noGoDirectiveError(gomod)
	}
	return string(m[1]), nil
}

type noGoDirectiveError string

func (e noGoDirectiveError) Error() string { return "no go directive in " + string(e) }

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}

// writeScratchPackage materialises Generate's output plus the support shim
// into dir, so the out-of-process reflect helper can import it as an ordinary
// package. Split out of GeneratedGoRoots to keep that function under the
// cyclomatic-complexity limit .golangci.yml sets.
func writeScratchPackage(dir string, files map[string]string) error {
	for name, src := range files {
		if err := writeFile(filepath.Join(dir, name), src); err != nil {
			return fmt.Errorf("codegen: writing %s: %w", name, err)
		}
	}
	if err := writeFile(filepath.Join(dir, "support.go"), supportShim()); err != nil {
		return fmt.Errorf("codegen: writing support.go: %w", err)
	}
	return nil
}
