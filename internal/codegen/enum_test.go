package codegen

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// --- named type usage, presence rules (real corpus) -------------------------

// TestEnumFieldsUseNamedType pins T-19's core promise across several real
// families at once: an Alts-bearing field generates as a named string type
// with an exported const block, not a plain string, and the struct field
// itself references that type.
func TestEnumFieldsUseNamedType(t *testing.T) {
	got := generateOrFail(t)

	cases := []struct {
		file string
		want []string
	}{
		{"authenticate.go", []string{
			"Role       AuthRole `xml:\"role\"`",
			"type AuthRole string",
			"AuthRoleUser       AuthRole = \"User\"",
			"AuthRoleSuperAdmin AuthRole = \"Super Admin\"",
		}},
		{"describe_auth.go", []string{
			"TimeStatus        AuthTimeStatus `xml:\"time_status\"`",
			"type AuthTimeStatus string",
			"AuthTimeStatusExpired  AuthTimeStatus = \"expired\"",
		}},
		{"get_filters.go", []string{
			"Relation FilterRelation `xml:\"relation\"`",
			"type FilterRelation string",
			"FilterRelationEq       FilterRelation = \"=\"",
			"FilterRelationContains FilterRelation = \":\"",
			"type FilterOrder string",
			"FilterOrderAscending  FilterOrder = \"ascending\"",
		}},
		{"get_nvts.go", []string{
			"Category         NVTCategory     `xml:\"category,omitempty\"`",
			"type NVTCategory string",
			"NVTCategoryN0  NVTCategory = \"0\"",
			"NVTCategoryN11 NVTCategory = \"11\"",
		}},
		{"task_shared.go", []string{
			"type TaskUsageType string",
			"TaskUsageTypeScan  TaskUsageType = \"scan\"",
			"TaskUsageTypeAudit TaskUsageType = \"audit\"",
		}},
		{"credential_shared.go", []string{
			"type CredentialAlgorithm string",
			"CredentialAlgorithmAes CredentialAlgorithm = \"aes\"",
			"type CredentialType string",
			"CredentialTypeSNMP  CredentialType = \"snmp\"",
		}},
		{"user_shared.go", []string{
			"type UserSource string",
			"UserSourceLdapConnect   UserSource = \"ldap_connect\"",
		}},
	}
	for _, c := range cases {
		src, ok := got[c.file]
		if !ok {
			t.Fatalf("%s was not generated", c.file)
		}
		for _, want := range c.want {
			if !strings.Contains(src, want) {
				t.Errorf("%s: want substring %q, not found", c.file, want)
			}
		}
	}
}

// TestEnumOptionalAndRequiredPresenceRules confirms T-19 does not touch
// presence semantics: an optional Alts field stays a non-pointer with
// omitempty (zero value omitted, same as any other optional scalar), and a
// required Alts field keeps no omitempty (its zero value is still emitted).
func TestEnumOptionalAndRequiredPresenceRules(t *testing.T) {
	got := generateOrFail(t)

	optional := []struct{ file, want string }{
		{"get_nvts.go", "Category         NVTCategory     `xml:\"category,omitempty\"`"},
		{"create_task.go", "UsageType       TaskUsageType    `xml:\"usage_type,omitempty\"`"},
		{"user_shared.go", "Source UserSource `xml:\"source,omitempty\"`"},
		{"create_credential.go", "Type          CredentialType          `xml:\"type,omitempty\"`"},
	}
	for _, c := range optional {
		src, ok := got[c.file]
		if !ok {
			t.Fatalf("%s was not generated", c.file)
		}
		if !strings.Contains(src, c.want) {
			t.Errorf("%s: want optional enum field %q, not found", c.file, c.want)
		}
	}

	required := []struct{ file, want string }{
		{"authenticate.go", "Role       AuthRole `xml:\"role\"`"},
		{"describe_auth.go", "TimeStatus        AuthTimeStatus `xml:\"time_status\"`"},
		{"credential_shared.go", "Algorithm CredentialAlgorithm `xml:\"algorithm\"`"},
	}
	for _, c := range required {
		src, ok := got[c.file]
		if !ok {
			t.Fatalf("%s was not generated", c.file)
		}
		if !strings.Contains(src, c.want) {
			t.Errorf("%s: want required enum field %q, not found", c.file, c.want)
		}
		if strings.Contains(src, c.want+"omitempty") {
			t.Errorf("%s: required enum field %q must not carry omitempty", c.file, c.want)
		}
	}
}

// --- registerEnum / sameAltSet fail-closed checks ----------------------------

func newEnumTestGen() *gen {
	return &gen{
		byEnumName: map[string]*enumType{},
		ec:         EnumConstantNames,
		ecConsumed: map[string]bool{},
	}
}

func TestEnumRegistrationRejectsEmptyAlts(t *testing.T) {
	g := newEnumTestGen()
	n := &schema.Node{Name: "status", Alts: nil}
	if err := g.registerEnum("get_things", n, nil); err == nil {
		t.Fatal("expected an error for a node reached with empty Alts, got nil")
	}
}

func TestEnumRegistrationRejectsDuplicateAltValue(t *testing.T) {
	g := newEnumTestGen()
	n := &schema.Node{Name: "status", Alts: []string{"on", "off", "on"}}
	err := g.registerEnum("get_things", n, nil)
	if err == nil {
		t.Fatal("expected an error for a duplicate alt value, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate alt value") {
		t.Errorf("error %q does not mention the duplicate alt value", err)
	}
}

func TestEnumRegistrationAgreesAcrossReorderedOccurrences(t *testing.T) {
	g := newEnumTestGen()
	first := &schema.Node{Name: "status", Alts: []string{"on", "off"}}
	if err := g.registerEnum("create_thing", first, nil); err != nil {
		t.Fatalf("first occurrence: unexpected error: %v", err)
	}
	second := &schema.Node{Name: "status", Alts: []string{"off", "on"}}
	if err := g.registerEnum("modify_thing", second, nil); err != nil {
		t.Fatalf("second occurrence (reordered alt set): unexpected error: %v", err)
	}
	if len(g.enums) != 1 {
		t.Fatalf("want one deduped enumType, got %d", len(g.enums))
	}
	if len(g.byEnumName["status"].uses) != 2 {
		t.Fatalf("want both occurrences recorded as uses, got %d", len(g.byEnumName["status"].uses))
	}
}

func TestEnumRegistrationRejectsDisagreeingAltSets(t *testing.T) {
	g := newEnumTestGen()
	first := &schema.Node{Name: "status", Alts: []string{"on", "off"}}
	if err := g.registerEnum("create_thing", first, nil); err != nil {
		t.Fatalf("first occurrence: unexpected error: %v", err)
	}
	second := &schema.Node{Name: "status", Alts: []string{"on", "off", "unknown"}}
	err := g.registerEnum("modify_thing", second, nil)
	if err == nil {
		t.Fatal("expected an error for a disagreeing alt set, got nil")
	}
	if !strings.Contains(err.Error(), "must agree on their alt set") {
		t.Errorf("error %q does not explain the disagreement", err)
	}
}

// --- altSuffix: punctuation, digits, and the goName(":") fallback guard -----

func TestAltSuffixHandlesOrdinaryWords(t *testing.T) {
	g := newEnumTestGen()
	cases := map[string]string{
		"always":       "Always",
		"ldap_connect": "LdapConnect",
		"Super Admin":  "SuperAdmin",
	}
	for value, want := range cases {
		got, err := g.altSuffix(value)
		if err != nil {
			t.Errorf("altSuffix(%q): unexpected error: %v", value, err)
			continue
		}
		if got != want {
			t.Errorf("altSuffix(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestAltSuffixPrefixesAllDigitValues(t *testing.T) {
	g := newEnumTestGen()
	cases := map[string]string{"0": "N0", "7": "N7", "11": "N11"}
	for value, want := range cases {
		got, err := g.altSuffix(value)
		if err != nil {
			t.Errorf("altSuffix(%q): unexpected error: %v", value, err)
			continue
		}
		if got != want {
			t.Errorf("altSuffix(%q) = %q, want %q", value, got, want)
		}
	}
}

// TestAltSuffixRoutesPunctuationThroughEnumConstantNames covers the relation
// operator table end to end, including ":" specifically: goName treats ':'
// itself as a word separator, so splitting it produces no word at all and
// goName's empty-input fallback would otherwise silently return the
// misleading name "Field". altSuffix must never reach that fallback for a
// punctuation-only value.
func TestAltSuffixRoutesPunctuationThroughEnumConstantNames(t *testing.T) {
	g := newEnumTestGen()
	cases := map[string]string{
		"=": "Eq",
		":": "Contains",
		"~": "Regex",
		">": "GT",
		"<": "LT",
	}
	for value, want := range cases {
		got, err := g.altSuffix(value)
		if err != nil {
			t.Errorf("altSuffix(%q): unexpected error: %v", value, err)
			continue
		}
		if got != want {
			t.Errorf("altSuffix(%q) = %q, want %q", value, got, want)
		}
		if !g.ecConsumed[value] {
			t.Errorf("altSuffix(%q) did not mark %q as consumed in EnumConstantNames", value, value)
		}
	}
}

func TestAltSuffixFailsClosedWithoutAnOverride(t *testing.T) {
	g := newEnumTestGen()
	if _, err := g.altSuffix("@"); err == nil {
		t.Fatal("expected an error for a punctuation value with no EnumConstantNames entry, got nil")
	}
}

// --- buildEnumConsts: colliding constant names ------------------------------

func TestBuildEnumConstsRejectsCollidingConstantNames(t *testing.T) {
	g := newEnumTestGen()
	et := &enumType{
		goName: "Status",
		node:   &schema.Node{Alts: []string{"on", "On"}},
	}
	err := g.buildEnumConsts(et)
	if err == nil {
		t.Fatal("expected an error for two alt values colliding on one constant name, got nil")
	}
	if !strings.Contains(err.Error(), "both produce constant name") {
		t.Errorf("error %q does not explain the collision", err)
	}
}

func TestBuildEnumConstsOrdersConstantsBySchemaOrder(t *testing.T) {
	g := newEnumTestGen()
	et := &enumType{
		goName: "Status",
		node:   &schema.Node{Alts: []string{"expired", "inactive", "unknown", "valid"}},
	}
	if err := g.buildEnumConsts(et); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"StatusExpired", "StatusInactive", "StatusUnknown", "StatusValid"}
	if len(et.consts) != len(want) {
		t.Fatalf("got %d consts, want %d", len(et.consts), len(want))
	}
	for i, c := range et.consts {
		if c.goName != want[i] {
			t.Errorf("const %d: got %q, want %q", i, c.goName, want[i])
		}
		if c.value != et.node.Alts[i] {
			t.Errorf("const %d: value %q does not match schema order %q", i, c.value, et.node.Alts[i])
		}
	}
}

// --- runtime round trip: presence rules and open-enum unknown values --------

// enumScratchModule writes the full generated tree into a throwaway Go
// module, the same way credentialScratchModule/ticketScratchModule do, for
// the families TestEnumGeneratedRoundTrips exercises.
func enumScratchModule(t *testing.T, dir string) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	for _, name := range []string{
		CommandFile("authenticate"), CommandFile("get_nvts"), CommandFile("get_filters"),
	} {
		if _, ok := got[name]; !ok {
			t.Fatalf("%s was not generated, so this test would not compile it", name)
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
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
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	return goBin
}

// TestEnumGeneratedRoundTrips proves T-19's runtime contract on real
// generated types: a required enum field's zero value is still emitted (no
// omitempty), an optional one is dropped, a known constant round-trips as
// its literal wire value, an unrecognized server value still decodes rather
// than being rejected (the enum is open, not validated), and a punctuation-
// derived constant (EnumConstantNames) round-trips as its symbol.
func TestEnumGeneratedRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := enumScratchModule(t, dir)
	mustWriteFile(t, filepath.Join(dir, "main.go"), enumRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"required-enum-zero-value-present: <authenticate_response status=\"\" status_text=\"\"><role></role><timezone></timezone></authenticate_response>\n" +
		"optional-enum-zero-value-omitted: <GetNVTsNVT oid=\"1.2.3\"><name>n</name></GetNVTsNVT>\n" +
		"optional-enum-known-value: <GetNVTsNVT oid=\"1.2.3\"><name>n</name><category>0</category></GetNVTsNVT>\n" +
		"optional-enum-unknown-value-roundtrip: <GetNVTsNVT oid=\"1.2.3\"><name>n</name><category>99</category></GetNVTsNVT> -> Category=\"99\"\n" +
		"punctuation-constant-relation-eq: <FilterKeyword><column>name</column><relation>=</relation><value>x</value></FilterKeyword>\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

const enumRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func printXML(label string, v any) {
	b, err := xml.Marshal(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, label, err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", label, b)
}

func main() {
	printXML("required-enum-zero-value-present", &genproto.AuthenticateResponse{})

	printXML("optional-enum-zero-value-omitted", &genproto.GetNVTsNVT{OID: "1.2.3", Name: "n"})
	printXML("optional-enum-known-value", &genproto.GetNVTsNVT{OID: "1.2.3", Name: "n", Category: genproto.NVTCategoryN0})

	b, err := xml.Marshal(&genproto.GetNVTsNVT{OID: "1.2.3", Name: "n", Category: genproto.NVTCategory("99")})
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshal unknown category:", err)
		os.Exit(1)
	}
	var back genproto.GetNVTsNVT
	if err := xml.Unmarshal(b, &back); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal unknown category:", err)
		os.Exit(1)
	}
	fmt.Printf("optional-enum-unknown-value-roundtrip: %s -> Category=%q\n", b, back.Category)

	printXML("punctuation-constant-relation-eq", &genproto.FilterKeyword{Column: "name", Relation: genproto.FilterRelationEq, Value: "x"})
}
`
