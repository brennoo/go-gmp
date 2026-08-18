package gmptest

import (
	"encoding/xml"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// commandsDir is resolved from this file's own location rather than the
// process working directory. go test does chdir into the package directory, but
// this is a long-lived drift detector and should not depend on that.
var commandsDir = func() string {
	_, self, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("..", "..", "commands")
	}
	return filepath.Join(filepath.Dir(self), "..", "..", "commands")
}()

// xmlNameTags parses commands/ and returns every type that declares an
// XMLName field, keyed by its xml tag.
//
// This deliberately re-derives from source rather than reusing Registry: a
// completeness check that consults the thing it is checking proves nothing.
func xmlNameTags(t *testing.T) map[string]string {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, commandsDir, nil, 0)
	if err != nil {
		t.Fatalf("Failed to parse %s: %v", commandsDir, err)
	}

	tags := map[string]string{}
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			ast.Inspect(file, func(n ast.Node) bool {
				spec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}
				st, ok := spec.Type.(*ast.StructType)
				if !ok || st.Fields == nil {
					return true
				}
				for _, f := range st.Fields.List {
					if len(f.Names) != 1 || f.Names[0].Name != "XMLName" || f.Tag == nil {
						continue
					}
					raw, err := strconv.Unquote(f.Tag.Value)
					if err != nil {
						continue
					}
					tag := reflect.StructTag(raw).Get("xml")
					if tag = strings.Split(tag, ",")[0]; tag != "" {
						tags[tag] = spec.Name.Name
					}
				}
				return true
			})
		}
	}
	if len(tags) == 0 {
		t.Fatalf("Parsed %s but found no XMLName tags", commandsDir)
	}
	return tags
}

// rootCommands applies the rule that defines a root command: a tag is one iff a
// sibling type is tagged <tag>_response. No exception list — a type like
// VulnsFilters also declares XMLName and is excluded by the rule itself
// rather than by being named here.
func rootCommands(tags map[string]string) map[string][2]string {
	roots := map[string][2]string{}
	for tag, typ := range tags {
		if resp, ok := tags[tag+"_response"]; ok {
			roots[tag] = [2]string{typ, resp}
		}
	}
	return roots
}

func TestRegistryMatchesSource(t *testing.T) {
	if _, err := os.Stat(commandsDir); err != nil {
		t.Skipf("Skipping: %s is not readable (%v)", commandsDir, err)
	}

	roots := rootCommands(xmlNameTags(t))

	var missing, extra []string
	for tag := range roots {
		if _, ok := Registry[tag]; !ok {
			missing = append(missing, tag)
		}
	}
	for tag := range Registry {
		if _, ok := roots[tag]; !ok {
			extra = append(extra, tag)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)

	if len(missing) > 0 {
		t.Errorf("Commands in commands/ but not in Registry (%d): %v", len(missing), missing)
	}
	if len(extra) > 0 {
		t.Errorf("Commands in Registry but not in commands/ (%d): %v", len(extra), extra)
	}
}

// TestRegistryConstructorsMatchTheirTags catches a copy/paste in the table that
// the name comparison cannot see: an entry whose constructors return the wrong
// type. Registry["get_tasks"] returning a GetTargets would otherwise pass.
func TestRegistryConstructorsMatchTheirTags(t *testing.T) {
	for name, cmd := range Registry {
		if got := rootTag(t, cmd.NewReq()); got != name {
			t.Errorf("Registry[%q].NewReq builds a type tagged %q", name, got)
		}
		want := name + "_response"
		if got := rootTag(t, cmd.NewResp()); got != want {
			t.Errorf("Registry[%q].NewResp builds a type tagged %q, want %q", name, got, want)
		}
	}
}

// rootTag reports the xml tag on v's XMLName field.
func rootTag(t *testing.T, v any) string {
	t.Helper()
	rt := reflect.TypeOf(v).Elem()
	f, ok := rt.FieldByName("XMLName")
	if !ok {
		t.Fatalf("%s has no XMLName field", rt.Name())
	}
	if f.Type != reflect.TypeOf(xml.Name{}) {
		t.Fatalf("%s.XMLName is %s, not xml.Name", rt.Name(), f.Type)
	}
	return strings.Split(f.Tag.Get("xml"), ",")[0]
}

// TestRegistryCoversEveryCommandFile is a coarse cross-check against the file
// system, independent of both the AST walk and the table.
func TestRegistryCoversEveryCommandFile(t *testing.T) {
	if _, err := os.Stat(commandsDir); err != nil {
		t.Skipf("Skipping: %s is not readable (%v)", commandsDir, err)
	}
	if len(Registry) != len(rootCommands(xmlNameTags(t))) {
		t.Fatalf("Registry has %d commands", len(Registry))
	}
	t.Logf("Registry covers %d root commands", len(Registry))
}
