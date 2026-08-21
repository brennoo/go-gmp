package codegen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/brennoo/go-gmp/internal/schema"
)

// targetPackage is the package the generated files are written for. Support
// types named in TypeOverrides are spelled "commands.X" there so the table
// reads the same as the import a caller would write; inside the target package
// itself the qualifier has to go, and stripping it here is the one place that
// knows about it. The compile-check test re-adds it through type aliases.
const targetPackage = "commands"

// sharedFileSuffix names the file that holds nested types used by more than
// one command in the family. Emitting them once, rather than per command, is
// the whole point of the structural dedup pass: without it "nvt" would be
// declared by both create_note.go and modify_note.go and neither would build.
const sharedFileSuffix = "_shared.go"

// Generate renders Go source for every command in cmds.
//
// It returns a filename -> gofmt'd source map. The output is a pure function
// of (cmds, ov, pf, ro): no clock, no filesystem, no map iteration order leaks
// into it, and nothing here reads vendor/XML. That is what lets the golden
// test in generate_test.go run in CI with only the committed JSON snapshot.
//
// family is the manifest's declared family name (e.g. "nvt"), authoritative
// over the family noun this package would otherwise infer from the commands'
// own names: that inference assumes every command in the family shares one
// noun stem ("create_note"/"get_notes" -> "note"), which fails for a family
// like get_nvts/get_nvt_families/get_preferences/get_info whose commands
// share no common stem at all. An empty family falls back to the inferred
// noun, which every caller that has no manifest in scope (this package's own
// tests) passes.
//
// co is CardinalityOverrides, RequiredOverride's mirror image for schema Card
// rather than requiredness.
//
// ec is EnumConstantNames, keyed by literal wire value rather than by
// command; ecConsumed accumulates which of its entries this call resolved an
// alt value through. Unlike every other override table's own xConsumed map,
// ecConsumed is caller-owned and shared across every family in one
// GenerateAll run rather than reset per call: EnumConstantNames' entries are
// not command-scoped, so a family that never happens to need a given
// <relation> operator is not evidence the entry is unused corpus-wide — only
// GenerateAll, which sees every family, can decide that.
// validateCommandShapes fails closed on a snapshot too malformed to generate
// from at all, before any naming or emission work starts.
func validateCommandShapes(cmds []schema.Command) error {
	for i := range cmds {
		if cmds[i].Name == "" {
			return fmt.Errorf("codegen: command %d has no name", i)
		}
		if cmds[i].Request == nil || cmds[i].Response == nil {
			return fmt.Errorf("codegen: %s: snapshot is missing a request or response shape", cmds[i].Name)
		}
	}
	return nil
}

func Generate(cmds []schema.Command, ov map[string]Override, pf map[string]PointerOverride, ro map[string]RequiredOverride, fo map[string][]FieldOverride, fs map[string][]FieldSuppression, cw map[string]CommandWarning, da map[DefectKey]DefectAllowlistEntry, no map[string]NameOverride, tov map[string][]FieldTypeOverride, family string, co map[string]CardinalityOverride, ec map[string]EnumConstantName, ecConsumed map[string]bool) (map[string]string, error) {
	if len(cmds) == 0 {
		return nil, fmt.Errorf("codegen: no commands to generate")
	}
	g := &gen{
		ov: ov, pf: pf, ro: ro, fo: fo, fs: fs, cw: cw, da: da, no: no, tov: tov, co: co, ec: ec,
		byHash:      map[string]*structType{},
		taken:       map[string]bool{},
		byEnumName:  map[string]*enumType{},
		needsErrors: map[*strings.Builder]bool{},
		foConsumed:  map[string]bool{},
		fsConsumed:  map[string]bool{},
		cwConsumed:  map[string]bool{},
		daConsumed:  map[DefectKey]bool{},
		pfConsumed:  map[string]bool{},
		roConsumed:  map[string]bool{},
		noConsumed:  map[string]bool{},
		tovConsumed: map[string]bool{},
		coConsumed:  map[string]bool{},
		ecConsumed:  ecConsumed,
	}

	// Sorting decouples the output from the caller's slice order, so a
	// reordered snapshot list cannot silently renumber disambiguated names.
	g.cmds = append(g.cmds, cmds...)
	sort.Slice(g.cmds, func(i, j int) bool { return g.cmds[i].Name < g.cmds[j].Name })

	if err := validateCommandShapes(g.cmds); err != nil {
		return nil, err
	}
	if err := g.checkProvenance(); err != nil {
		return nil, err
	}
	if err := g.validateDefectAllowlist(); err != nil {
		return nil, err
	}
	if err := g.checkDefects(); err != nil {
		return nil, err
	}

	// familyNoun(g.cmds) is authoritative: it is what mergeFamilyFiles'
	// naming-convention check assumes Generate derives, so a caller-supplied
	// family can only fill in where derivation fails (a heterogeneous group
	// like nvt's get_info/get_nvt_families/get_nvts/get_preferences, whose
	// commands share no single trailing noun) — never override a name it did
	// derive, or a family declared under the wrong name would silently adopt
	// it instead of being caught.
	g.familyNoun = familyNoun(g.cmds)
	if g.familyNoun == "" {
		g.familyNoun = family
	}
	g.family = goName(g.familyNoun)
	if err := g.validateOverrideTables(); err != nil {
		return nil, err
	}
	if err := g.collect(); err != nil {
		return nil, err
	}
	if err := g.name(); err != nil {
		return nil, err
	}
	if err := g.nameEnums(); err != nil {
		return nil, err
	}
	out, err := g.emit()
	if err != nil {
		return nil, err
	}
	if err := g.checkTablesConsumed(); err != nil {
		return nil, err
	}
	return out, nil
}

// validateOverrideTables checks the shape of FieldOverrides, the
// FieldSuppressions/RelocatedFrom relocation link, and CommandWarnings.
// Grouped into one call so that adding a new override table costs Generate's
// cyclomatic complexity one branch total, not one per table it validates.
func (g *gen) validateOverrideTables() error {
	if err := g.validateFieldOverrides(); err != nil {
		return err
	}
	if err := validateFieldRelocations(g.fo, g.fs); err != nil {
		return err
	}
	if err := g.validateCommandWarnings(); err != nil {
		return err
	}
	if err := g.validateNameOverrides(); err != nil {
		return err
	}
	if err := g.validateFieldTypeOverrides(); err != nil {
		return err
	}
	if err := g.validatePointerOverrides(); err != nil {
		return err
	}
	if err := g.validateRequiredOverrides(); err != nil {
		return err
	}
	if err := g.validateCardinalityOverrides(); err != nil {
		return err
	}
	return g.validateEnumConstantNames()
}

// checkTablesConsumed fails the build on any stale override-table entry left
// after emission — one whose command was in scope for this call but never
// matched a generated field, command, or defect. Grouped for the same reason
// validateOverrideTables is.
func (g *gen) checkTablesConsumed() error {
	if err := g.checkFieldOverridesConsumed(); err != nil {
		return err
	}
	if err := g.checkFieldSuppressionsConsumed(); err != nil {
		return err
	}
	if err := g.checkCommandWarningsConsumed(); err != nil {
		return err
	}
	if err := g.checkDefectAllowlistConsumed(); err != nil {
		return err
	}
	if err := g.checkPointerOverridesConsumed(); err != nil {
		return err
	}
	if err := g.checkRequiredOverridesConsumed(); err != nil {
		return err
	}
	if err := g.checkCardinalityOverridesConsumed(); err != nil {
		return err
	}
	if err := g.checkNameOverridesConsumed(); err != nil {
		return err
	}
	return g.checkFieldTypeOverridesConsumed()
}

// sha256Pattern matches hex.EncodeToString's SHA-256 output: 64 lowercase
// hex digits.
var sha256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

var validExtractionMethods = map[schema.ExtractionMethod]bool{
	schema.ExtractionFile:   true,
	schema.ExtractionServed: true,
}

// checkProvenance fails the build when any command's schema provenance is
// absent or malformed, or when commands in this call describe different
// schemas — T-16a's fail-closed gate. It does not check ExtractedAt against the
// current time: Generate's output must stay a pure function of its inputs.
func (g *gen) checkProvenance() error {
	for i := range g.cmds {
		c := &g.cmds[i]
		p := c.Provenance
		switch {
		case p.Method == "":
			return fmt.Errorf("codegen: %s: schema provenance is missing Method; re-run Stage A's -update", c.Name)
		case !validExtractionMethods[p.Method]:
			return fmt.Errorf("codegen: %s: schema provenance has unknown Method %q", c.Name, p.Method)
		case strings.TrimSpace(p.Source) == "":
			return fmt.Errorf("codegen: %s: schema provenance is missing Source; re-run Stage A's -update", c.Name)
		case p.SourceSHA256 == "":
			return fmt.Errorf("codegen: %s: schema provenance is missing SourceSHA256; re-run Stage A's -update", c.Name)
		case !sha256Pattern.MatchString(p.SourceSHA256):
			return fmt.Errorf("codegen: %s: schema provenance SourceSHA256 %q is not 64 lowercase hex characters", c.Name, p.SourceSHA256)
		case p.ParsedSHA256 == "":
			return fmt.Errorf("codegen: %s: schema provenance is missing ParsedSHA256; re-run Stage A's -update", c.Name)
		case !sha256Pattern.MatchString(p.ParsedSHA256):
			return fmt.Errorf("codegen: %s: schema provenance ParsedSHA256 %q is not 64 lowercase hex characters", c.Name, p.ParsedSHA256)
		case strings.TrimSpace(p.DeclaredVersion) == "":
			return fmt.Errorf("codegen: %s: schema provenance is missing DeclaredVersion; re-run Stage A's -update", c.Name)
		case p.ExtractedAt.IsZero():
			return fmt.Errorf("codegen: %s: schema provenance is missing ExtractedAt; re-run Stage A's -update", c.Name)
		case p.ExtractedAt.Location() != time.UTC:
			return fmt.Errorf("codegen: %s: schema provenance ExtractedAt %v is not UTC", c.Name, p.ExtractedAt)
		}
	}
	first := g.cmds[0].Provenance
	for i := 1; i < len(g.cmds); i++ {
		if !sameSchemaSource(g.cmds[i].Provenance, first) {
			return fmt.Errorf("codegen: mixed schema provenance: %s and %s describe different schemas (%+v vs %+v)",
				g.cmds[0].Name, g.cmds[i].Name, first, g.cmds[i].Provenance)
		}
	}
	return nil
}

// sameSchemaSource compares everything but ExtractedAt, which says when a
// snapshot was written rather than what it was written from. Adding one family
// re-extracts only its own snapshots, so a batch routinely mixes timestamps
// while SourceSHA256 and ParsedSHA256 — the fields that actually identify the
// schema — agree.
func sameSchemaSource(a, b schema.Provenance) bool {
	a.ExtractedAt, b.ExtractedAt = time.Time{}, time.Time{}
	return a == b
}

// validDefectClasses are schema.parse.go's nine Defect* constants, checked
// by name so a typo'd Class in DefectAllowlist fails as malformed metadata
// rather than surfacing indirectly as a stale entry.
var validDefectClasses = map[string]bool{
	schema.DefectTypeForT:      true,
	schema.DefectDanglingRef:   true,
	schema.DefectDuplicateEle:  true,
	schema.DefectDuplicateRef:  true,
	schema.DefectEmptyPattern:  true,
	schema.DefectNoPattern:     true,
	schema.DefectUndefinedType: true,
	schema.DefectUnsupportedR:  true,
	schema.DefectRecursion:     true,
}

// validateDefectAllowlist checks the shape of every DefectAllowlist entry in
// g.da, independent of which command is in scope for this call: a missing
// Reason is a defect in the table itself, the same reasoning
// validateFieldOverrides/validateCommandWarnings apply to their own tables.
func (g *gen) validateDefectAllowlist() error {
	keys := make([]DefectKey, 0, len(g.da))
	for k := range g.da {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Command != keys[j].Command {
			return keys[i].Command < keys[j].Command
		}
		if keys[i].Class != keys[j].Class {
			return keys[i].Class < keys[j].Class
		}
		return keys[i].Path < keys[j].Path
	})
	for _, k := range keys {
		switch {
		case k.Command == "":
			return fmt.Errorf("codegen: DefectAllowlist[%+v]: Command is empty", k)
		case k.Path == "":
			return fmt.Errorf("codegen: DefectAllowlist[%+v]: Path is empty", k)
		case !validDefectClasses[k.Class]:
			return fmt.Errorf("codegen: DefectAllowlist[%+v]: Class is not a known schema.Defect* class", k)
		case g.da[k].Reason == "":
			return fmt.Errorf("codegen: DefectAllowlist[%+v]: Reason is empty", k)
		}
	}
	return nil
}

// checkDefects fails the build on any schema.Defect (schema.Command.Defects)
// carried by a command in this call that DefectAllowlist does not admit —
// T-15's fail-closed gate: an empty-pattern, no-pattern, recursion-cut or
// dangling-e node used to generate a silently truncated struct instead of
// refusing to generate at all. It also records which allowlist keys matched,
// so checkDefectAllowlistConsumed can flag one that no longer does.
func (g *gen) checkDefects() error {
	var bad []string
	for i := range g.cmds {
		c := &g.cmds[i]
		for _, d := range c.Defects {
			key := DefectKey{Command: c.Name, Class: d.Class, Path: d.Path}
			if _, ok := g.da[key]; !ok {
				bad = append(bad, fmt.Sprintf("%s: %s at %s (%s)", c.Name, d.Class, d.Path, d.Detail))
				continue
			}
			g.daConsumed[key] = true
		}
	}
	if len(bad) == 0 {
		return nil
	}
	sort.Strings(bad)
	return fmt.Errorf("codegen: schema defect(s) not in DefectAllowlist, refusing to generate: %s", strings.Join(bad, "; "))
}

// checkDefectAllowlistConsumed fails the build on a stale or mistyped
// DefectAllowlist entry, mirroring checkFieldOverridesConsumed's scoping:
// only a command that was actually part of this call can make "the defect it
// names never appeared" mean anything.
func (g *gen) checkDefectAllowlistConsumed() error {
	present := make(map[string]bool, len(g.cmds))
	for i := range g.cmds {
		present[g.cmds[i].Name] = true
	}
	var stale []string
	for key := range g.da {
		if present[key.Command] && !g.daConsumed[key] {
			stale = append(stale, fmt.Sprintf("%s:%s:%s", key.Command, key.Class, key.Path))
		}
	}
	if len(stale) == 0 {
		return nil
	}
	sort.Strings(stale)
	return fmt.Errorf("codegen: DefectAllowlist entry(ies) never matched a generated defect: %s", strings.Join(stale, ", "))
}

// validateCommandWarnings checks the shape of every CommandWarnings entry in
// g.cw, independent of which command is in scope for this call, the same
// reasoning validateFieldOverrides applies to g.fo.
func (g *gen) validateCommandWarnings() error {
	keys := make([]string, 0, len(g.cw))
	for k := range g.cw {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		w := g.cw[key]
		if w.Text == "" {
			return fmt.Errorf("codegen: CommandWarnings[%q]: Text is empty", key)
		}
		if w.Reason == "" {
			return fmt.Errorf("codegen: CommandWarnings[%q]: Reason is empty", key)
		}
	}
	return nil
}

// checkCommandWarningsConsumed fails the build on a stale or mistyped
// CommandWarnings entry, mirroring checkFieldOverridesConsumed's scoping:
// only a command that was actually part of this call can make "never
// emitted" mean anything. A "<command>_response" key is the response-side
// counterpart emitResponse attaches (see emitResponse), so it is present
// alongside the request-side "<command>" key.
func (g *gen) checkCommandWarningsConsumed() error {
	present := make(map[string]bool, len(g.cmds)*2)
	for i := range g.cmds {
		present[g.cmds[i].Name] = true
		present[g.cmds[i].Name+"_response"] = true
	}
	var stale []string
	for command := range g.cw {
		if present[command] && !g.cwConsumed[command] {
			stale = append(stale, command)
		}
	}
	if len(stale) == 0 {
		return nil
	}
	sort.Strings(stale)
	return fmt.Errorf("codegen: CommandWarnings entry(ies) never attached to a generated command: %s", strings.Join(stale, ", "))
}

// validatePointerOverrides checks the shape of every PointerOverrides entry in
// g.pf, independent of which command is in scope for this call, the same
// reasoning validateRequiredOverrides applies to g.ro: a malformed key or a
// missing Reason is a defect in the table itself.
func (g *gen) validatePointerOverrides() error {
	keys := make([]string, 0, len(g.pf))
	for k := range g.pf {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		command, field, ok := strings.Cut(key, ".")
		if !ok || command == "" || field == "" {
			return fmt.Errorf(`codegen: PointerOverrides[%q]: key must be "<command>.<field>"`, key)
		}
		if g.pf[key].Reason == "" {
			return fmt.Errorf("codegen: PointerOverrides[%q]: Reason is empty", key)
		}
	}
	return nil
}

// checkPointerOverridesConsumed fails the build on a stale or mistyped
// PointerOverrides entry: one whose command is among g.cmds (so this call was
// in a position to apply it) but which never matched a generated optional
// top-level field — a typo'd field name, an entry aimed at a required or
// nested field (fieldOf only ever consults g.pf from its CardOptional/
// InChoice, command != "" case), or a schema correction that removed the
// field the override existed to work around. Scoped by command like
// checkRequiredOverridesConsumed, for the same GenerateAll/CheckOutput
// per-family reason.
func (g *gen) checkPointerOverridesConsumed() error {
	return checkOverridesConsumed(g.cmds, g.pf, g.pfConsumed, "PointerOverrides", "optional field")
}

// pointerOverrideForced reports whether PointerOverrides names this optional
// field, recording the match in g.pfConsumed as a side effect. Factored out
// of fieldOf so the lookup-and-record pair costs that function's cyclomatic
// complexity nothing beyond the plain boolean it already had.
func (g *gen) pointerOverrideForced(command, name string) bool {
	key := command + "." + name
	_, forced := g.pf[key]
	if forced {
		g.pfConsumed[key] = true
	}
	return forced
}

// validateRequiredOverrides checks the shape of every RequiredOverrides entry
// in g.ro, independent of which command is in scope for this call, the same
// reasoning validateFieldOverrides applies to g.fo: a malformed key or a
// missing Reason is a defect in the table itself, not something that could
// ever become valid depending on which family is being generated.
func (g *gen) validateRequiredOverrides() error {
	keys := make([]string, 0, len(g.ro))
	for k := range g.ro {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		command, field, ok := strings.Cut(key, ".")
		if !ok || command == "" || field == "" {
			return fmt.Errorf(`codegen: RequiredOverrides[%q]: key must be "<command>.<field>"`, key)
		}
		if g.ro[key].Reason == "" {
			return fmt.Errorf("codegen: RequiredOverrides[%q]: Reason is empty", key)
		}
	}
	return nil
}

// checkRequiredOverridesConsumed fails the build on a stale or mistyped
// RequiredOverrides entry: one whose command is among g.cmds (so this call
// was in a position to apply it) but which never matched a generated
// required field — a typo'd field name, an entry aimed at an optional or
// nested field (fieldOf only ever consults g.ro from its CardRequired,
// command != "" case), or a schema correction that removed the requiredness
// the override existed to work around. Scoped by command like
// checkFieldOverridesConsumed, for the same GenerateAll/CheckOutput
// per-family reason.
func (g *gen) checkRequiredOverridesConsumed() error {
	return checkOverridesConsumed(g.cmds, g.ro, g.roConsumed, "RequiredOverrides", "required field")
}

// checkOverridesConsumed is the shared staleness check behind
// checkRequiredOverridesConsumed and checkCardinalityOverridesConsumed: an
// override table entry whose command is in scope but was never consumed is a
// typo, a misaimed entry, or a schema correction that outdated it.
func checkOverridesConsumed[V any](cmds []schema.Command, table map[string]V, consumed map[string]bool, tableName, kind string) error {
	present := make(map[string]bool, len(cmds))
	for i := range cmds {
		present[cmds[i].Name] = true
	}
	var stale []string
	for key := range table {
		command, _, _ := strings.Cut(key, ".")
		if present[command] && !consumed[key] {
			stale = append(stale, key)
		}
	}
	if len(stale) == 0 {
		return nil
	}
	sort.Strings(stale)
	return fmt.Errorf("codegen: %s entry(ies) never matched a generated %s: %s", tableName, kind, strings.Join(stale, ", "))
}

// validateCardinalityOverrides checks the shape of every CardinalityOverrides
// entry in g.co, mirroring validateRequiredOverrides: these are defects in
// the table itself, independent of which command is in scope for this call.
func (g *gen) validateCardinalityOverrides() error {
	keys := make([]string, 0, len(g.co))
	for k := range g.co {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		command, field, ok := strings.Cut(key, ".")
		if !ok || command == "" || field == "" {
			return fmt.Errorf(`codegen: CardinalityOverrides[%q]: key must be "<command>.<field>"`, key)
		}
		if g.co[key].Reason == "" {
			return fmt.Errorf("codegen: CardinalityOverrides[%q]: Reason is empty", key)
		}
	}
	return nil
}

// checkCardinalityOverridesConsumed fails the build on a stale or mistyped
// CardinalityOverrides entry: one whose command is among g.cmds but which
// never matched a generated schema-singular top-level field. Mirrors
// checkRequiredOverridesConsumed.
func (g *gen) checkCardinalityOverridesConsumed() error {
	return checkOverridesConsumed(g.cmds, g.co, g.coConsumed, "CardinalityOverrides", "schema-singular field")
}

// validateFieldOverrides checks the shape of every FieldOverrides entry in
// g.fo, independent of which command is in scope for this call: these are
// defects in the table itself (a typo'd GoType, a missing Reason, two
// chardata entries fighting over one struct), not something that could ever
// become valid depending on which family is being generated. Keys are sorted
// first so a bad entry is reported the same way regardless of Go's random map
// iteration order.
func (g *gen) validateFieldOverrides() error {
	keys := make([]string, 0, len(g.fo))
	for k := range g.fo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		hasChardata := false
		type nameAttr struct {
			name string
			attr bool
		}
		seen := map[nameAttr]bool{}

		for i, e := range g.fo[key] {
			if e.Reason == "" {
				return fmt.Errorf("codegen: FieldOverrides[%q][%d]: Reason is empty", key, i)
			}
			if e.GoType == "" {
				return fmt.Errorf("codegen: FieldOverrides[%q][%d]: GoType is empty", key, i)
			}
			if _, err := parser.ParseExpr(e.GoType); err != nil {
				return fmt.Errorf("codegen: FieldOverrides[%q][%d]: GoType %q is not a valid Go type expression: %w", key, i, e.GoType, err)
			}
			if e.Chardata {
				if e.Name != "" || e.Attr || e.OmitEmpty {
					return fmt.Errorf("codegen: FieldOverrides[%q][%d]: Chardata is mutually exclusive with Name, Attr and OmitEmpty", key, i)
				}
				if hasChardata {
					return fmt.Errorf("codegen: FieldOverrides[%q]: more than one chardata override", key)
				}
				hasChardata = true
				continue
			}
			if e.Name == "" {
				return fmt.Errorf("codegen: FieldOverrides[%q][%d]: Name is empty on a non-chardata override", key, i)
			}
			na := nameAttr{name: e.Name, attr: e.Attr}
			if seen[na] {
				return fmt.Errorf("codegen: FieldOverrides[%q]: duplicate wire name %q (attr=%v)", key, e.Name, e.Attr)
			}
			seen[na] = true
		}
	}
	return nil
}

// validateFieldRelocations checks the shape of every FieldSuppressions entry
// and every FieldOverride.RelocatedFrom reference, independent of which
// command is in scope for this call, the same way validateFieldOverrides
// checks FieldOverrides. It enforces the bidirectional link a relocation pair
// depends on: a FieldSuppressions entry's RelocatedTo must resolve to a
// FieldOverrides entry whose own RelocatedFrom points back at it, and vice
// versa, so deleting either half without the other fails every future
// Generate call rather than silently leaving a suppressed field with nowhere
// to reappear.
func validateFieldRelocations(fo map[string][]FieldOverride, fs map[string][]FieldSuppression) error {
	fsKeys := make([]string, 0, len(fs))
	for k := range fs {
		fsKeys = append(fsKeys, k)
	}
	sort.Strings(fsKeys)
	for _, key := range fsKeys {
		if err := validateFieldSuppressionRelocations(key, fs[key], fo); err != nil {
			return err
		}
	}

	foKeys := make([]string, 0, len(fo))
	for k := range fo {
		foKeys = append(foKeys, k)
	}
	sort.Strings(foKeys)
	for _, key := range foKeys {
		if err := validateFieldOverrideRelocations(key, fo[key], fs); err != nil {
			return err
		}
	}
	return nil
}

// validateFieldSuppressionRelocations checks the FieldSuppressions entries
// under key: at most one Chardata suppression, mutually exclusive with
// Name/Attr, and every relocation has a matching FieldOverrides counterpart.
func validateFieldSuppressionRelocations(key string, entries []FieldSuppression, fo map[string][]FieldOverride) error {
	hasChardata := false
	for i, s := range entries {
		if s.Chardata {
			if s.Name != "" || s.Attr {
				return fmt.Errorf("codegen: FieldSuppressions[%q][%d]: Chardata is mutually exclusive with Name and Attr", key, i)
			}
			if hasChardata {
				return fmt.Errorf("codegen: FieldSuppressions[%q]: more than one chardata suppression", key)
			}
			hasChardata = true
		} else if s.Name == "" {
			return fmt.Errorf("codegen: FieldSuppressions[%q][%d]: Name is empty on a non-chardata suppression", key, i)
		}
		if s.Reason == "" {
			return fmt.Errorf("codegen: FieldSuppressions[%q][%d]: Reason is empty", key, i)
		}
		if s.RelocatedTo == "" {
			return fmt.Errorf("codegen: FieldSuppressions[%q][%d]: RelocatedTo is empty", key, i)
		}
		if !hasRelocationCounterpart(fo[s.RelocatedTo], s.Name, s.Attr, s.Chardata, key) {
			return fmt.Errorf(
				"codegen: FieldSuppressions[%q][%d]: RelocatedTo %q has no FieldOverrides entry with "+
					"Name %q, Attr %v, Chardata %v and RelocatedFrom %q",
				key, i, s.RelocatedTo, s.Name, s.Attr, s.Chardata, key)
		}
	}
	return nil
}

// validateFieldOverrideRelocations checks that every FieldOverrides entry
// under key with a RelocatedFrom has a matching FieldSuppressions counterpart.
func validateFieldOverrideRelocations(key string, entries []FieldOverride, fs map[string][]FieldSuppression) error {
	for i, e := range entries {
		if e.RelocatedFrom == "" {
			continue
		}
		if !hasSuppressionCounterpart(fs[e.RelocatedFrom], e.Name, e.Attr, e.Chardata, key) {
			return fmt.Errorf(
				"codegen: FieldOverrides[%q][%d]: RelocatedFrom %q has no FieldSuppressions entry with "+
					"Name %q, Attr %v, Chardata %v and RelocatedTo %q",
				key, i, e.RelocatedFrom, e.Name, e.Attr, e.Chardata, key)
		}
	}
	return nil
}

func hasRelocationCounterpart(entries []FieldOverride, name string, attr, chardata bool, relocatedFrom string) bool {
	for _, e := range entries {
		if e.Name == name && e.Attr == attr && e.Chardata == chardata && e.RelocatedFrom == relocatedFrom {
			return true
		}
	}
	return false
}

func hasSuppressionCounterpart(entries []FieldSuppression, name string, attr, chardata bool, relocatedTo string) bool {
	for _, s := range entries {
		if s.Name == name && s.Attr == attr && s.Chardata == chardata && s.RelocatedTo == relocatedTo {
			return true
		}
	}
	return false
}

// validateOverrideKey checks that key has the "<command>:/<path>" shape
// fieldOverrideKey produces. Without this, a malformed key (missing colon,
// empty command, empty path) never matches a generated struct and is
// silently ignored rather than caught as the stale entry it is: staleness
// checks key their command by cutting on ":", so a key with no ":" at all
// yields a command that matches no real command and is treated as
// out-of-scope, not stale.
func validateOverrideKey(key string) error {
	command, rest, ok := strings.Cut(key, ":/")
	if !ok || command == "" || rest == "" {
		return fmt.Errorf(`key must be "<command>:/<path>"`)
	}
	return nil
}

// validateNameOverrides checks the shape of every NameOverrides entry in
// g.no, independent of which command is in scope for this call, the same
// reasoning validateFieldOverrides applies to g.fo.
func (g *gen) validateNameOverrides() error {
	keys := make([]string, 0, len(g.no))
	for k := range g.no {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if err := validateOverrideKey(key); err != nil {
			return fmt.Errorf("codegen: NameOverrides[%q]: %w", key, err)
		}
		e := g.no[key]
		if e.Reason == "" {
			return fmt.Errorf("codegen: NameOverrides[%q]: Reason is empty", key)
		}
		if !token.IsIdentifier(e.GoName) || !ast.IsExported(e.GoName) {
			return fmt.Errorf("codegen: NameOverrides[%q]: GoName %q is not a valid exported Go identifier", key, e.GoName)
		}
	}
	return nil
}

// validateEnumConstantNames checks the shape of every EnumConstantNames
// entry in g.ec, independent of which command is in scope for this call: the
// table is keyed by literal wire value, not by command, so there is no
// per-call subset of it to validate.
func (g *gen) validateEnumConstantNames() error {
	keys := make([]string, 0, len(g.ec))
	for k := range g.ec {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if key == "" {
			return fmt.Errorf("codegen: EnumConstantNames: key (the wire alt value) is empty")
		}
		e := g.ec[key]
		if e.Reason == "" {
			return fmt.Errorf("codegen: EnumConstantNames[%q]: Reason is empty", key)
		}
		if !isIdentifierFragment(e.Suffix) {
			return fmt.Errorf("codegen: EnumConstantNames[%q]: Suffix %q is not a valid Go identifier fragment", key, e.Suffix)
		}
	}
	return nil
}

// validateFieldTypeOverrides checks the shape of every FieldTypeOverrides
// entry in g.tov, independent of which command is in scope for this call,
// the same reasoning validateFieldOverrides applies to g.fo. Unlike an
// invalid FieldOverrides.GoType, which only ever adds a field, an invalid
// FieldTypeOverride.GoType replaces the type of a field the schema already
// declares — validating it as a parseable Go type expression here means a
// bad entry fails at Generate time instead of surfacing as a compile error
// in generated code with Generate having already reported success.
func (g *gen) validateFieldTypeOverrides() error {
	keys := make([]string, 0, len(g.tov))
	for k := range g.tov {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	type nameAttr struct {
		name string
		attr bool
	}

	for _, key := range keys {
		if err := validateOverrideKey(key); err != nil {
			return fmt.Errorf("codegen: FieldTypeOverrides[%q]: %w", key, err)
		}
		seen := map[nameAttr]bool{}
		for i, e := range g.tov[key] {
			if e.Name == "" {
				return fmt.Errorf("codegen: FieldTypeOverrides[%q][%d]: Name is empty", key, i)
			}
			if e.Reason == "" {
				return fmt.Errorf("codegen: FieldTypeOverrides[%q][%d]: Reason is empty", key, i)
			}
			if e.GoType == "" {
				return fmt.Errorf("codegen: FieldTypeOverrides[%q][%d]: GoType is empty", key, i)
			}
			if _, err := parser.ParseExpr(e.GoType); err != nil {
				return fmt.Errorf("codegen: FieldTypeOverrides[%q][%d]: GoType %q is not a valid Go type expression: %w", key, i, e.GoType, err)
			}
			if e.Pointer && strings.HasPrefix(e.GoType, "*") {
				return fmt.Errorf("codegen: FieldTypeOverrides[%q][%d]: Pointer is set but GoType %q is already a pointer", key, i, e.GoType)
			}
			na := nameAttr{name: e.Name, attr: e.Attr}
			if seen[na] {
				return fmt.Errorf("codegen: FieldTypeOverrides[%q]: duplicate wire name %q (attr=%v)", key, e.Name, e.Attr)
			}
			seen[na] = true
		}
	}
	return nil
}

// checkFieldOverridesConsumed fails the build on a stale or mistyped
// FieldOverrides entry: one whose command is among g.cmds (so this call was
// in a position to apply it) but which never matched a node emitStruct
// looked up. An entry whose command isn't part of this call at all is not
// checked here — GenerateAll and CheckOutput both call Generate once per
// family, often with only a subset of FieldOverrides' commands in scope (a
// review-only family's commands, for instance, never reach CheckOutput's
// Adopted()-filtered calls), so "unconsumed" only means something once the
// relevant command was actually part of the run.
func (g *gen) checkFieldOverridesConsumed() error {
	stale := staleTableKeys(g.cmds, g.fo, g.foConsumed)
	if len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("codegen: FieldOverrides entry(ies) never matched a generated field: %s", strings.Join(stale, ", "))
}

// checkFieldSuppressionsConsumed fails the build on a stale or mistyped
// FieldSuppressions entry, using the same in-scope-command test as
// checkFieldOverridesConsumed.
func (g *gen) checkFieldSuppressionsConsumed() error {
	stale := staleTableKeys(g.cmds, g.fs, g.fsConsumed)
	if len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("codegen: FieldSuppressions entry(ies) never matched a generated field: %s", strings.Join(stale, ", "))
}

// checkFieldTypeOverridesConsumed fails the build on a stale or mistyped
// FieldTypeOverrides entry, using the same in-scope-command test as
// checkFieldOverridesConsumed.
func (g *gen) checkFieldTypeOverridesConsumed() error {
	stale := staleTableKeys(g.cmds, g.tov, g.tovConsumed)
	if len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("codegen: FieldTypeOverrides entry(ies) never matched a generated field: %s", strings.Join(stale, ", "))
}

// checkNameOverridesConsumed fails the build on a stale or mistyped
// NameOverrides entry, using the same in-scope-command test as
// checkFieldOverridesConsumed.
func (g *gen) checkNameOverridesConsumed() error {
	stale := staleTableKeys(g.cmds, g.no, g.noConsumed)
	if len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("codegen: NameOverrides entry(ies) never matched a generated type: %s", strings.Join(stale, ", "))
}

// staleTableKeys returns the keys of a "<command>:/..." keyed table (either
// FieldOverrides or FieldSuppressions — the value type itself is never
// inspected) that are in scope for this run — the command is among cmds, so
// this call was in a position to apply the entry — but never marked
// consumed, sorted for a deterministic error message. Shared by
// checkFieldOverridesConsumed and checkFieldSuppressionsConsumed, which
// differ only in which table and consumed-set they check.
func staleTableKeys[V any](cmds []schema.Command, table map[string]V, consumed map[string]bool) []string {
	present := make(map[string]bool, len(cmds))
	for i := range cmds {
		present[cmds[i].Name] = true
	}
	var stale []string
	for key := range table {
		command, _, _ := strings.Cut(key, ":")
		if present[command] && !consumed[key] {
			stale = append(stale, key)
		}
	}
	sort.Strings(stale)
	return stale
}

// GenerateAll renders every family in m, one Generate call per family, and
// merges the results into a single filename -> source map.
//
// Per-family rather than one global call is a deliberate choice. Generate's
// structural dedup, its type-name disambiguation ladder and its shared-file
// decision are all scoped to the commands it is handed; giving it exactly one
// family's commands is what guarantees that adding an unrelated family to the
// manifest can never rename a type or move a struct in an existing one. A
// single global call would have to reproduce that isolation internally, and
// nothing in the output needs to be shared across families to make it work.
//
// Stated as the rule it is: structurally-equal shapes are deduplicated
// family-locally and never globally. Two families that happen to describe the
// same <nvt> shape each get their own type for it, in their own shared file,
// under their own family-qualified name. The duplication is the price of the
// property that matters more — a family's generated names, files and bytes
// depend on that family's commands and nothing else — because the alternative
// makes every future manifest addition a potential rename of already-committed,
// already-called types across commands/, client/ and internal/gmptest.
// Two tests hold the line: TestGenerateAllIsPerFamily (a synthetic twin family
// whose shapes are identical to note's, the adversarial case) and
// TestRealFamiliesGenerateIndependently (the real note and override families,
// each generated alone and compared byte-for-byte against its slice of the
// combined run).
//
// cmds may contain commands in any order and must cover every command the
// manifest names; a command not claimed by any family is an error, since it
// would otherwise be silently dropped.
func GenerateAll(m Manifest, cmds []schema.Command, ov map[string]Override, pf map[string]PointerOverride, ro map[string]RequiredOverride, fo map[string][]FieldOverride, fs map[string][]FieldSuppression, cw map[string]CommandWarning, da map[DefectKey]DefectAllowlistEntry, no map[string]NameOverride, tov map[string][]FieldTypeOverride, co map[string]CardinalityOverride, ec map[string]EnumConstantName) (map[string]string, error) {
	byName := make(map[string]schema.Command, len(cmds))
	for _, c := range cmds {
		byName[c.Name] = c
	}

	out := map[string]string{}
	from := map[string]string{} // filename -> the family that emitted it
	ecConsumed := map[string]bool{}
	for _, f := range m.Families {
		fam := make([]schema.Command, 0, len(f.Commands))
		for _, name := range f.Commands {
			c, ok := byName[name]
			if !ok {
				return nil, fmt.Errorf("codegen: family %q: no snapshot loaded for %s", f.Name, name)
			}
			fam = append(fam, c)
		}
		files, err := Generate(fam, ov, pf, ro, fo, fs, cw, da, no, tov, f.Name, co, ec, ecConsumed)
		if err != nil {
			return nil, fmt.Errorf("codegen: family %q: %w", f.Name, err)
		}
		if err := mergeFamilyFiles(out, from, f, files); err != nil {
			return nil, err
		}
	}

	if err := checkEnumConstantNamesConsumed(ec, ecConsumed); err != nil {
		return nil, err
	}
	return out, nil
}

// checkEnumConstantNamesConsumed fails closed on an EnumConstantNames entry
// no family's generation used: a corpus-wide table (see overrides.go) can go
// stale if the schema stops using a value like "=" anywhere, and nothing else
// would ever notice. Shared by GenerateAll and CheckOutput so gengmp -check
// enforces the same guarantee as the real write path.
func checkEnumConstantNamesConsumed(ec map[string]EnumConstantName, ecConsumed map[string]bool) error {
	var stale []string
	for key := range ec {
		if !ecConsumed[key] {
			stale = append(stale, key)
		}
	}
	if len(stale) > 0 {
		sort.Strings(stale)
		return fmt.Errorf("codegen: EnumConstantNames entries never used by any alt value across the "+
			"corpus: %q", stale)
	}
	return nil
}

// mergeFamilyFiles folds one family's Generate output into out.
//
// Both checks it makes should be unreachable by construction — a command
// belongs to one family (LoadManifest enforces it) and the shared file is named
// after the family — but "should be unreachable" is exactly the property worth
// asserting: silently overwriting one family's file with another's would
// produce a plausible-looking tree that is missing types.
func mergeFamilyFiles(out, from map[string]string, f Family, files map[string]string) error {
	allowed := map[string]bool{f.SharedFile(): true}
	for _, name := range f.Commands {
		allowed[CommandFile(name)] = true
	}
	for name, src := range files {
		if !allowed[name] {
			return fmt.Errorf("codegen: family %q emitted %s, which the manifest's naming convention does not "+
				"account for (expected %s or one of the family's command files); is the family's name wrong?",
				f.Name, name, f.SharedFile())
		}
		if other, ok := from[name]; ok {
			return fmt.Errorf("codegen: families %q and %q both emit %s", other, f.Name, name)
		}
		from[name] = f.Name
		out[name] = src
	}
	return nil
}

// --- collection -------------------------------------------------------------

// structType is one generated Go struct: a nested element shape, plus every
// place in the family that shape was reached from.
type structType struct {
	node *schema.Node // the first occurrence; all others are structurally equal
	uses []use

	goName string
	file   string
}

// use is one path a shape was reached by, e.g. get_notes /note/result/nvt.
type use struct {
	command string
	path    []string // ancestor element names, outermost first, excluding the node
}

// enumType is one generated named string type for an <alts>-constrained
// field: every place in the family an Alts-bearing scalar leaf named node.Name
// was reached, plus the ordered constants derived from node.Alts.
//
// Grouped by wire name rather than by structural shape (unlike structType):
// a scalar leaf has no children for typeHash to fold Alts into, so
// registerEnum is what typeHash's identical-shape rule becomes for a leaf.
type enumType struct {
	name string       // the wire element/attribute name, e.g. "condition"
	node *schema.Node // the first occurrence; every other must share its Alts
	uses []use

	goName string
	file   string
	consts []enumConst // in node.Alts order
}

// enumConst is one exported constant an enumType declares: goName is already
// qualified by the enum type's own name (e.g. "AlertConditionAlways"), not
// just the suffix, so emitEnum can write it directly.
type enumConst struct {
	value  string
	goName string
}

type gen struct {
	cmds   []schema.Command
	ov     map[string]Override
	pf     map[string]PointerOverride
	ro     map[string]RequiredOverride
	fo     map[string][]FieldOverride
	fs     map[string][]FieldSuppression
	cw     map[string]CommandWarning
	da     map[DefectKey]DefectAllowlistEntry
	no     map[string]NameOverride
	tov    map[string][]FieldTypeOverride
	co     map[string]CardinalityOverride
	ec     map[string]EnumConstantName
	family string
	// familyNoun is family's pre-goName snake_case form, used only by fileFor
	// for the shared-file basename (see familyNoun's doc comment).
	familyNoun string

	byHash map[string]*structType
	types  []*structType   // first-occurrence order
	taken  map[string]bool // every Go type name already handed out

	// byEnumName groups Alts-bearing scalar leaves by wire name within this
	// family: two occurrences of the same name must agree on their alt set
	// (see registerEnum), the scalar-leaf counterpart to typeHash folding
	// Alts into a struct child's content hash.
	byEnumName map[string]*enumType
	enums      []*enumType // first-occurrence order

	// needsErrors marks, by output builder, a file whose generated collision
	// codec (emitCollisionUnmarshal) reports an unexpected extra occurrence by
	// returning an error rather than silently skipping it — see
	// emitCollisionMethods' CardRequired guard.
	needsErrors map[*strings.Builder]bool

	// foConsumed records which FieldOverrides keys matched a node during this
	// run, so checkFieldOverridesConsumed can flag a stale one.
	foConsumed map[string]bool

	// fsConsumed records which FieldSuppressions keys matched a node during
	// this run, so checkFieldSuppressionsConsumed can flag a stale one.
	fsConsumed map[string]bool

	// cwConsumed records which CommandWarnings keys were attached to a
	// generated command during this run, so checkCommandWarningsConsumed can
	// flag a stale one.
	cwConsumed map[string]bool

	// daConsumed records which DefectAllowlist keys matched a schema.Defect
	// during this run, so checkDefectAllowlistConsumed can flag a stale one.
	daConsumed map[DefectKey]bool

	// pfConsumed records which PointerOverrides keys matched an optional
	// field during this run, so checkPointerOverridesConsumed can flag a
	// stale or mistyped one.
	pfConsumed map[string]bool

	// roConsumed records which RequiredOverrides keys matched a required
	// top-level field during this run, so checkRequiredOverridesConsumed can
	// flag a stale or mistyped one.
	roConsumed map[string]bool

	// noConsumed records which NameOverrides keys matched a structType during
	// this run, so checkNameOverridesConsumed can flag a stale one.
	noConsumed map[string]bool

	// tovConsumed records which FieldTypeOverrides keys matched a node during
	// this run, so checkFieldTypeOverridesConsumed can flag a stale one.
	tovConsumed map[string]bool

	// coConsumed records which CardinalityOverrides keys matched a
	// schema-singular top-level field during this run, so
	// checkCardinalityOverridesConsumed can flag a stale or mistyped one.
	coConsumed map[string]bool

	// ecConsumed records which EnumConstantNames keys resolved an alt value
	// during this run. Caller-owned and shared across every Generate call in
	// one GenerateAll run — see Generate's doc comment.
	ecConsumed map[string]bool
}

// fieldOverrideKey builds a FieldOverrides key from a structType's use: the
// command that reached it, the ancestor path (outermost first, excluding the
// node itself), and the node's own name.
func fieldOverrideKey(command string, path []string, name string) string {
	full := make([]string, 0, len(path)+1)
	full = append(full, path...)
	full = append(full, name)
	return command + ":/" + strings.Join(full, "/")
}

// collect walks every command's request and response, recording one entry per
// distinct nested struct shape and one per distinct Alts-bearing wire name.
func (g *gen) collect() error {
	for i := range g.cmds {
		c := &g.cmds[i]
		if err := g.walk(c.Name, c.Request, nil); err != nil {
			return err
		}
		if err := g.walk(c.Name, c.Response, nil); err != nil {
			return err
		}
	}
	return nil
}

func (g *gen) walk(command string, n *schema.Node, path []string) error {
	for _, kid := range n.Children {
		if len(kid.Alts) > 0 {
			// Checked before the Children guard below: a "dual-role" node
			// (get_report_configs.type, get_report_formats.trust/type) has
			// both its own Alts (its chardata text is enum-constrained) and
			// Children (attributes), so it needs both a struct and an enum
			// type registered for it.
			if err := g.registerEnum(command, kid, path); err != nil {
				return err
			}
		}
		if len(kid.Children) == 0 {
			continue // a scalar leaf needs no struct type of its own
		}
		h := typeHash(kid)
		st, ok := g.byHash[h]
		if !ok {
			st = &structType{node: kid}
			g.byHash[h] = st
			g.types = append(g.types, st)
		}
		st.uses = append(st.uses, use{command: command, path: append([]string(nil), path...)})
		// A fresh slice per branch: sharing the backing array would let a
		// later sibling overwrite an ancestor path already recorded above.
		down := make([]string, len(path), len(path)+1)
		copy(down, path)
		if err := g.walk(command, kid, append(down, kid.Name)); err != nil {
			return err
		}
	}
	return nil
}

// registerEnum records one Alts-bearing leaf occurrence, grouped by wire name
// within this family. A second occurrence of the same name must declare the
// same alt set (order-independent) — the counterpart, for a leaf typeHash
// never visits, of typeHash already folding Alts into a struct child's
// content hash. Per T-19's scope, this is a hard generation error rather
// than a family-specific alias: no real disagreement exists in the corpus
// today, and the user's guidance is not to build an override mechanism to
// speculatively paper over one that might appear later.
func (g *gen) registerEnum(command string, n *schema.Node, path []string) error {
	if len(n.Alts) == 0 {
		return fmt.Errorf("codegen: %s: enum node %q reached with empty Alts", command, n.Name)
	}
	seen := make(map[string]bool, len(n.Alts))
	for _, a := range n.Alts {
		if a == "" {
			return fmt.Errorf("codegen: %s: enum %q declares an empty alt value", command, n.Name)
		}
		if seen[a] {
			return fmt.Errorf("codegen: %s: enum %q declares duplicate alt value %q", command, n.Name, a)
		}
		seen[a] = true
	}

	et, ok := g.byEnumName[n.Name]
	if !ok {
		et = &enumType{name: n.Name, node: n}
		g.byEnumName[n.Name] = et
		g.enums = append(g.enums, et)
	} else if !sameAltSet(et.node.Alts, n.Alts) {
		return fmt.Errorf(
			"codegen: %s: enum %q declares alt set %v, but %s already established %v for the same "+
				"name in this family; structurally-shared enum fields must agree on their alt set "+
				"(see docs/review-2026-08.md T-19)",
			command, n.Name, n.Alts, ownerOfUses(et.uses), et.node.Alts)
	}
	et.uses = append(et.uses, use{command: command, path: append([]string(nil), path...)})
	return nil
}

// sameAltSet reports whether a and b declare the same alt values, ignoring
// order: two occurrences of one wire name are allowed to declare their alts
// in different schema-document order and still be the same enum.
func sameAltSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sa := append([]string(nil), a...)
	sb := append([]string(nil), b...)
	sort.Strings(sa)
	sort.Strings(sb)
	for i := range sa {
		if sa[i] != sb[i] {
			return false
		}
	}
	return true
}

// typeHash identifies a generated struct by the shape it would have.
//
// The node's own name is part of the hash: <result id=""/> and <task id=""/>
// are the same shape but different concepts, and collapsing them would produce
// a type no reader could name. The node's own Card and Kind are deliberately
// *not* part of it: those describe the field that holds the struct, not the
// struct, which is why create_note's required <nvt> and modify_note's optional
// <nvt> share one type.
func typeHash(n *schema.Node) string {
	var b strings.Builder
	writeHash(&b, n)
	return b.String()
}

func writeHash(b *strings.Builder, n *schema.Node) {
	b.WriteString(n.Name)
	if n.Chardata {
		// A node with mixed content gets a Value field the same shape's
		// pure-element twin does not; without this marker two commands whose
		// children are identical but which disagree on bare chardata would
		// collapse into one struct and one of them would be wrong.
		b.WriteString("#text:")
		b.WriteString(n.Type)
	}
	b.WriteByte('{')
	for _, c := range n.Children {
		fmt.Fprintf(b, "%s|%s|%s|%s|%s|", c.Name, c.Kind, c.Card, c.Type, strings.Join(c.Alts, "\x1f"))
		writeHash(b, c)
		b.WriteByte(';')
	}
	b.WriteByte('}')
}

// --- naming -----------------------------------------------------------------

// familyNoun derives the noun the commands are about: create_note, modify_note
// and get_notes are all about a "note". Dropping the leading verb and the
// plural is enough for GMP, whose command names are uniformly <verb>_<noun>.
// If the commands disagree there is no family, and every nested type falls
// back to being named after its own command. Returned in the schema's own
// snake_case spelling (e.g. "report_config") rather than goName's Go-cased
// form — fileFor needs this form, and lowercasing a multi-word Go name back
// down (strings.ToLower(g.family)) loses the word boundary goName introduced,
// turning "ReportConfig" into "reportconfig" instead of "report_config".
func familyNoun(cmds []schema.Command) string {
	noun := ""
	for _, c := range cmds {
		parts := strings.SplitN(c.Name, "_", 2)
		if len(parts) != 2 || parts[1] == "" {
			return ""
		}
		n := strings.TrimSuffix(parts[1], "s")
		if noun == "" {
			noun = n
		} else if noun != n {
			return ""
		}
	}
	return noun
}

// name assigns every collected shape a unique Go type name, in
// first-occurrence order so the assignment is stable.
func (g *gen) name() error {
	for i := range g.cmds {
		g.taken[goName(g.cmds[i].Name)] = true
		g.taken[goName(g.cmds[i].Name)+"Response"] = true
	}
	for _, st := range g.types {
		goName, err := g.pick(st)
		if err != nil {
			return err
		}
		st.goName = goName
		g.taken[st.goName] = true
		st.file = g.fileFor(st)
	}
	return nil
}

// pick consults NameOverrides first (see its doc comment for why a name
// sometimes cannot be derived from this family's own commands alone), then
// walks a ladder of increasingly qualified candidates and takes the first one
// free: the bare element name, then the element name prefixed by one more
// ancestor at a time, then the same qualified by the command, then a counter.
// Short names win when nothing contends for them, and the collisions that do
// happen (get_notes has three structurally different <nvt>s) are resolved with
// the context that actually distinguishes them.
func (g *gen) pick(st *structType) (string, error) {
	override, ok, err := g.nameOverrideFor(st)
	if err != nil {
		return "", err
	}
	if ok {
		if g.taken[override] {
			return "", fmt.Errorf("codegen: NameOverrides assigns %q, which is already taken by another generated type", override)
		}
		return override, nil
	}

	first := st.uses[0]
	for _, c := range g.candidates(st, first) {
		if !g.taken[c] {
			return c, nil
		}
	}
	base := g.qualify(goName(st.node.Name))
	for i := 2; ; i++ {
		if c := fmt.Sprintf("%s%d", base, i); !g.taken[c] {
			return c, nil
		}
	}
}

// nameEnums assigns every collected enum type a unique Go type name and
// derives its constants. It runs strictly after name(), so g.taken already
// reflects every struct name: this is what resolves the 3 "dual-role" nodes
// (get_report_configs.type, get_report_formats.trust/type — see registerEnum)
// for free, the same way any other struct-vs-struct name collision is
// resolved, since the struct half of a dual-role node claims the bare name
// first and the enum ladder here falls through to the next qualified
// candidate.
//
// Enum types do not consult NameOverrides: that table is struct-scoped
// (fieldOverrideKey-keyed) and has no entries aimed at an enum today; if a
// real, reviewed collision appears, it gets its own override then rather
// than widening NameOverrides' meaning now.
func (g *gen) nameEnums() error {
	for _, et := range g.enums {
		et.goName = g.pickEnum(et)
		g.taken[et.goName] = true
		et.file = g.fileForEnum(et)
		if err := g.buildEnumConsts(et); err != nil {
			return err
		}
	}
	return nil
}

// pickEnum is pick's counterpart for an enumType: the same qualified-name
// ladder, without a NameOverrides consultation (see nameEnums). Unlike pick,
// it cannot fail: there's no override lookup to reject, and the fallback
// loop always terminates on some suffixed candidate.
func (g *gen) pickEnum(et *enumType) string {
	first := et.uses[0]
	for _, c := range g.candidateLadder(goName(et.name), first) {
		if !g.taken[c] {
			return c
		}
	}
	base := g.qualify(goName(et.name))
	for i := 2; ; i++ {
		if c := fmt.Sprintf("%s%d", base, i); !g.taken[c] {
			return c
		}
	}
}

// buildEnumConsts derives one exported constant per alt value, in the order
// node.Alts declares them, and fails closed if two different alt values
// collapse onto the same constant name: silently keeping only one would drop
// a legitimate wire value from the generated API.
func (g *gen) buildEnumConsts(et *enumType) error {
	byName := make(map[string]string, len(et.node.Alts))
	for _, v := range et.node.Alts {
		suffix, err := g.altSuffix(v)
		if err != nil {
			return fmt.Errorf("codegen: enum %s: %w", et.goName, err)
		}
		name := et.goName + suffix
		if other, ok := byName[name]; ok {
			return fmt.Errorf(
				"codegen: enum %s: alt values %q and %q both produce constant name %q",
				et.goName, other, v, name)
		}
		byName[name] = v
		et.consts = append(et.consts, enumConst{value: v, goName: name})
	}
	return nil
}

// altSuffix derives an enum constant's name suffix from its wire value.
// goName handles ordinary words (snake_case, kebab-case, space-separated),
// but two shapes of value defeat it:
//
//   - A value that is purely numeric (GMP uses these for a handful of
//     integer-coded fields, e.g. NVT category codes) has no letters for
//     goName to capitalize, so it round-trips unchanged — a legal Go
//     identifier fragment everywhere except its required first rune.
//     Prefixing "N" fixes this mechanically for any such value, with no
//     per-value table: unlike a punctuation symbol, a digit string has no
//     protocol-wide meaning to look up, so there's nothing an override table
//     could add. Checked against the raw value, not goName's output: this
//     must not fire on a value goName legitimately transforms into something
//     else.
//   - A value that is nothing but punctuation — GMP's <relation> filter
//     operators (=, :, ~, >, <) chief among them — gives goName no letters
//     to capitalize and no separator to split on. For ':' specifically,
//     goName treats the character itself as a word separator, leaving no
//     word at all; goName's empty-input fallback then silently returns the
//     generic name "Field", which would be wrong here, not merely invalid,
//     so it must never reach that fallback for a value with no word content.
//     A relation operator's meaning is intrinsic to the symbol wherever it
//     appears, so EnumConstantNames is the reviewed, protocol-wide fallback
//     for this whole case.
func (g *gen) altSuffix(value string) (string, error) {
	if isAllDigits(value) {
		return "N" + value, nil
	}
	if len(strings.FieldsFunc(value, isNameSeparator)) > 0 {
		if s := goName(value); isIdentifierFragment(s) {
			return s, nil
		}
	}
	if o, ok := g.ec[value]; ok {
		if g.ecConsumed != nil {
			g.ecConsumed[value] = true
		}
		return o.Suffix, nil
	}
	return "", fmt.Errorf(
		"alt value %q does not produce a valid Go identifier and has no EnumConstantNames entry", value)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// isIdentifierFragment reports whether s is safe to use, on its own or
// concatenated after another identifier, as a piece of a Go identifier:
// every rune is a letter or underscore, and any digit is not the first rune.
func isIdentifierFragment(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		switch {
		case r == '_' || unicode.IsLetter(r):
		case i > 0 && unicode.IsDigit(r):
		default:
			return false
		}
	}
	return true
}

// nameOverrideFor looks up st in NameOverrides by every path it was reached
// from, mirroring fieldOverridesFor's per-use agreement discipline: every use
// — including one with no matching g.no entry, which resolves to "no
// opinion" rather than being skipped — must agree, or an override meant for
// one use could silently leak onto a structurally-deduplicated sibling use
// that had none.
func (g *gen) nameOverrideFor(st *structType) (name string, ok bool, err error) {
	seen := map[string]bool{}
	haveFirst := false
	var firstKey string
	var firstName string
	var firstOK bool

	for _, u := range st.uses {
		key := fieldOverrideKey(u.command, u.path, st.node.Name)
		if seen[key] {
			continue
		}
		seen[key] = true

		var resName string
		var resOK bool
		if entry, matched := g.no[key]; matched {
			g.noConsumed[key] = true
			resName, resOK = entry.GoName, true
		}

		if !haveFirst {
			haveFirst, firstKey, firstName, firstOK = true, key, resName, resOK
			continue
		}
		if resName != firstName || resOK != firstOK {
			return "", false, fmt.Errorf(
				"codegen: %s: NameOverrides for %q and %q disagree; a structurally-deduplicated type "+
					"must resolve to the same name from every path that reaches it",
				goName(st.node.Name), firstKey, key)
		}
	}
	return firstName, firstOK, nil
}

func (g *gen) candidates(st *structType, first use) []string {
	return g.candidateLadder(goName(st.node.Name), first)
}

// candidateLadder builds the shared candidate sequence pick and pickEnum both
// walk: the bare (family-qualified) name, then the name prefixed by one more
// ancestor at a time, then the same qualified by the command that first
// reached it. Factored out of candidates rather than folding enumType into
// candidates itself, to keep the already-golden-tested struct naming path
// untouched by the enum feature.
func (g *gen) candidateLadder(self string, first use) []string {
	var out []string
	add := func(s string) {
		for _, seen := range out {
			if seen == s {
				return
			}
		}
		out = append(out, s)
	}

	// Family-qualified, from the bare name outward through the ancestors.
	add(g.qualify(self))
	for i := len(first.path) - 1; i >= 0; i-- {
		add(g.qualify(goName(strings.Join(first.path[i:], "_")) + self))
	}
	// Then the same ladder qualified by the command that first reached it.
	cmd := goName(first.command)
	add(cmd + self)
	for i := len(first.path) - 1; i >= 0; i-- {
		add(cmd + goName(strings.Join(first.path[i:], "_")) + self)
	}
	return out
}

// qualify prefixes a name with the family noun, unless it already starts with
// it: <note_count> is NoteCount, not NoteNoteCount.
func (g *gen) qualify(s string) string {
	if g.family == "" || strings.HasPrefix(s, g.family) {
		return s
	}
	return g.family + s
}

// fileFor puts a shape used by one command in that command's file, and a shape
// reached from more than one command in the family's shared file.
func (g *gen) fileFor(st *structType) string {
	return fileForUses(st.uses, g.sharedFile())
}

// fileForEnum mirrors fileFor for an enumType.
func (g *gen) fileForEnum(et *enumType) string {
	return fileForUses(et.uses, g.sharedFile())
}

// sharedFile names the family's shared file, falling back to "commands" for
// the handful of families whose commands share no common noun (familyNoun's
// doc comment explains why).
func (g *gen) sharedFile() string {
	base := g.familyNoun
	if base == "" {
		base = "commands"
	}
	return base + sharedFileSuffix
}

// fileForUses puts a shape used by exactly one command in that command's
// file, and a shape reached from more than one command in shared.
func fileForUses(uses []use, shared string) string {
	first := uses[0].command
	for _, u := range uses[1:] {
		if u.command != first {
			return shared
		}
	}
	return first + ".go"
}

// --- Go type mapping ---------------------------------------------------------

// scalarTypes maps schema type names onto Go types for the leaves the
// mechanical rule can decide on its own. Anything absent is a string: the DSL's
// remaining leaf types (uuid, oid, name, status, threat, compliance_status,
// ...) are all xsd:token or xsd:string patterns, and so is an unnamed leaf.
//
// severity is here rather than in TypeOverrides because it is not an override
// of anything: its declared pattern only ever matches a decimal number, so
// float64 *is* the mechanical reading of it.
var scalarTypes = map[string]string{
	"text":     "string",
	"integer":  "int",
	"decimal":  "float64",
	"severity": "float64",
}

// pointerScalars are the Go types whose zero value is a legal, meaningful wire
// value, so an optional field of that type has to be a pointer for "absent"
// and "zero" to stay distinguishable. This is T-5's rule. A string is not one
// of them: in this protocol an empty string and an absent element mean the
// same thing, which is why the hand-written code uses `string` + omitempty.
// GMPTime belongs here for a related but distinct reason: its zero value
// (year 1) is a legal-looking timestamp, and encoding/xml's isEmptyValue
// never treats a struct as empty, so omitempty on a non-pointer GMPTime is a
// silent no-op that fabricates a zero timestamp on re-encode instead of
// omitting the field.
var pointerScalars = map[string]bool{
	"int": true, "float64": true, "GMPBool": true, "GMPAtoiBool": true, "GMPTime": true,
}

// goTypeOf resolves a node's Go type name, unqualified for targetPackage.
func (g *gen) goTypeOf(n *schema.Node) (string, error) {
	if len(n.Children) > 0 {
		return g.byHash[typeHash(n)].goName, nil
	}
	return g.scalarGoType(n)
}

// scalarGoType is goTypeOf's leaf-resolution rule, factored out so
// emitChardataOverridable can share it directly: a chardata node's own Type
// is always a lexical leaf type (e.g. "text", "base64"), never a nested
// shape, so it has no use for goTypeOf's structural branch. It takes the
// whole node, not just n.Type, so it can also route an Alts-bearing node
// (including a "dual-role" struct node whose own chardata text is itself
// enum-constrained — see registerEnum) to its collected enum type.
func (g *gen) scalarGoType(n *schema.Node) (string, error) {
	if len(n.Alts) > 0 {
		et, ok := g.byEnumName[n.Name]
		if !ok {
			// Unreachable by construction: collect's walk() visits every node
			// emit() later visits, so an Alts-bearing node reached here was
			// already registered. Kept as an error, not a panic or a silent
			// "string" fallback, per this codebase's fail-closed convention.
			return "", fmt.Errorf("codegen: %s: enum type was never registered during collection", n.Name)
		}
		return et.goName, nil
	}
	t := schema.Resolve(n.Type)
	if o, ok := g.ov[t]; ok {
		return strings.TrimPrefix(o.GoType, targetPackage+"."), nil
	}
	if s, ok := scalarTypes[t]; ok {
		return s, nil
	}
	return "string", nil
}

type field struct {
	name   string
	goType string
	tag    string
	node   *schema.Node // the node this field was built from; nil is never valid

	// wireTag is the xml struct tag this field would have carried had it not
	// lost a wire-name collision (see emitFields), preserved so writeField can
	// still record its wire name — under a distinct tag key, since `go vet`'s
	// structtag analyzer rejects two fields sharing one xml tag value outright
	// — for internal/schema's gomodel.go to recover by reflection.
	wireTag string

	// overrideReason is non-empty when this field was synthesized from
	// FieldOverrides rather than the schema; writeField emits it as a
	// gmpoverride tag for internal/schema/gomodel.go to recover, and
	// conform.go's comparer uses it to suppress the resulting EXTRA finding.
	overrideReason string
}

// fieldOf applies the DSL -> Go cardinality mapping. command is the request's
// or response's own command name, used only to look up PointerOverrides and
// RequiredOverrides; it is "" when n belongs to a nested struct (emitStruct),
// since a shape shared by more than one command has no single owner for
// either table to key on — see PointerOverrides' doc comment.
func (g *gen) fieldOf(command string, n *schema.Node, used map[string]bool, typeOv map[fieldKey]typeOverride) (field, error) {
	base, err := g.goTypeOf(n)
	if err != nil {
		return field{}, err
	}
	forcedPointer := false
	if t, ok := typeOv[fieldKey{name: n.Name, attr: n.Kind == schema.KindAttribute}]; ok {
		base = t.goType
		forcedPointer = t.pointer
	}
	goType, opts := base, ""
	var overrideReason string

	switch {
	case n.Card == schema.CardRepeated:
		goType = "[]" + base
		opts = ",omitempty"
	case n.Card == schema.CardOptional || n.InChoice:
		// InChoice means this field was reached through an <or>: the schema
		// marks each branch required, but at most one branch is ever present,
		// so its cardinality is advisory (see schema.Node.InChoice) and it
		// needs the same omitempty treatment as a genuinely optional field.
		forced := g.pointerOverrideForced(command, n.Name)
		if len(n.Children) > 0 || pointerScalars[base] || forced {
			goType = "*" + base
		}
		opts = ",omitempty"
	case n.Card == schema.CardRequired:
		var err error
		goType, opts, overrideReason, err = g.requiredFieldType(command, n, base, forcedPointer)
		if err != nil {
			return field{}, err
		}
	}

	overrideReason, err = g.withChardataRelocationReason(command, n, overrideReason)
	if err != nil {
		return field{}, err
	}
	if n.PatchReason != "" {
		if overrideReason != "" {
			return field{}, fmt.Errorf(
				"codegen: %s.%s: has both a SchemaPatch reason and another override/relocation reason; "+
					"ambiguous which gmpoverride tag to emit", command, n.Name)
		}
		overrideReason = n.PatchReason
	}

	tag := n.Name
	if n.Kind == schema.KindAttribute {
		tag += ",attr"
	}
	tag += opts

	return field{name: uniqueName(goName(n.Name), used), goType: goType, tag: tag, node: n, overrideReason: overrideReason}, nil
}

// requiredFieldType computes the Go type/tag-options/override-reason for a
// field outside any <or> that the schema declares exactly-one. Present by
// definition: no pointer, no omitempty — unless RequiredOverrides names this
// field as an exception live evidence disproved that default for. A plain
// string leaf just gets omitempty (resume_task/start_task's report_id). A
// structured field additionally needs a pointer to stay distinguishable from
// its zero value on re-encode — see T-23's count-wrapper entries, where the
// wrapper struct's own zero value (all-zero Value/Filtered/Page) is otherwise
// indistinguishable from "genuinely absent," fabricating an element the wire
// never sent.
//
// forcedPointer is FieldTypeOverrides[...].Pointer for this field: the same
// structured-pointer exception RequiredOverrides applies via the "command"
// case above, reached here instead because the field belongs to a nested,
// structurally-deduplicated shape with no single owning command for
// RequiredOverrides' "<command>.<field>" key to name (see get_results/
// get_reports/create_report/get_notes/get_overrides' shared NVT struct).
func (g *gen) requiredFieldType(command string, n *schema.Node, base string, forcedPointer bool) (goType, opts, overrideReason string, err error) {
	goType = base
	key := command + "." + n.Name
	if _, ok := g.ro[key]; ok {
		if len(n.Children) > 0 {
			goType = "*" + base
		} else if base != "string" {
			return "", "", "", fmt.Errorf(
				"codegen: RequiredOverrides[%q]: only evidenced for plain string or structured fields, "+
					"got Go type %q; no live evidence yet justifies overriding requiredness for a "+
					"non-string scalar field", key, base)
		}
		g.roConsumed[key] = true
		opts = ",omitempty"
	}
	if forcedPointer {
		if len(n.Children) == 0 {
			return "", "", "", fmt.Errorf(
				"codegen: FieldTypeOverrides[%q]: Pointer is only evidenced for structured fields, "+
					"got a scalar leaf (Go type %q); a scalar's own zero value is either always meaningful "+
					"(pointerScalars) or always equivalent to absence, so no live evidence yet justifies "+
					"this exception for one", key, base)
		}
		goType = "*" + base
		opts = ",omitempty"
	}
	// CardinalityOverrides names this field as one the schema declares
	// exactly-one but live evidence proves repeats on the wire (see
	// get_nvts.nvt): unlike RequiredOverrides, this is not restricted to
	// string leaves, since the whole point is forcing a structured shape
	// into a slice.
	if ov, ok := g.co[key]; ok {
		goType = "[]" + base
		opts = ",omitempty"
		overrideReason = ov.Reason
		g.coConsumed[key] = true
	}
	return goType, opts, overrideReason, nil
}

// withChardataRelocationReason returns existing unchanged unless n's struct
// type has its own chardata relocated away (FieldSuppressions[...].Chardata),
// which leaves nothing on that struct to carry the fact — see emitStruct — so
// the field pointing to it carries the suppression's Reason instead, for
// internal/schema's comparer to recover by reflection and suppress the
// resulting source-side MISSING.
func (g *gen) withChardataRelocationReason(command string, n *schema.Node, existing string) (string, error) {
	if len(n.Children) == 0 {
		return existing, nil
	}
	st, ok := g.byHash[typeHash(n)]
	if !ok {
		return existing, nil
	}
	_, chardataSuppression, err := g.fieldSuppressionsFor(st)
	if err != nil {
		return "", err
	}
	if chardataSuppression == nil {
		return existing, nil
	}
	if existing != "" {
		return "", fmt.Errorf(
			"codegen: %s.%s: has both a CardinalityOverride reason and a chardata-relocation "+
				"suppression reason; ambiguous which gmpoverride tag to emit", command, n.Name)
	}
	return chardataSuppression.Reason, nil
}

// emitChardata declares the field for a mixed-content node's bare text, if it
// has any. It is called before emitFields so the chardata field reads first,
// matching the hand-written override types this generator is meant to match
// (e.g. commands.OverrideText's Text field precedes its Excerpt attribute).
func (g *gen) emitChardata(b *strings.Builder, n *schema.Node, used map[string]bool) error {
	return g.emitChardataOverridable(b, n, used, nil)
}

// emitChardataOverridable is emitChardata plus a FieldOverrides chardata
// entry for the same node, if any. The schema wins when both are set: a node
// the schema itself declares Chardata for has no need of the override, and
// double-declaring the field would collide. override is only ever non-nil
// from emitStruct, since FieldOverrides addresses nested struct shapes by
// path, and a command's own request/response envelope has none.
func (g *gen) emitChardataOverridable(b *strings.Builder, n *schema.Node, used map[string]bool, override *FieldOverride) error {
	switch {
	case n.Chardata:
		name := uniqueName("Value", used)
		t, err := g.scalarGoType(n)
		if err != nil {
			return err
		}
		fmt.Fprintf(b, "%s %s `xml:\",chardata\"`\n", name, t)
	case override != nil:
		name := uniqueName("Value", used)
		fmt.Fprintf(b, "%s %s `xml:\",chardata\" gmpoverride:%q`\n", name, override.GoType, override.Reason)
	}
	return nil
}

func uniqueName(name string, used map[string]bool) string {
	out := name
	for i := 2; used[out]; i++ {
		out = fmt.Sprintf("%s%d", name, i)
	}
	used[out] = true
	return out
}

// --- emission ---------------------------------------------------------------

func (g *gen) emit() (map[string]string, error) {
	files := map[string]*strings.Builder{}
	needsXML := map[string]bool{}
	get := func(name string) *strings.Builder {
		b, ok := files[name]
		if !ok {
			b = &strings.Builder{}
			files[name] = b
		}
		return b
	}

	for i := range g.cmds {
		c := &g.cmds[i]
		file := c.Name + ".go"
		needsXML[file] = true
		b := get(file)
		if err := g.emitRequest(b, c); err != nil {
			return nil, err
		}
		if err := g.emitResponse(b, c); err != nil {
			return nil, err
		}
	}
	for _, st := range g.types {
		if err := g.emitStruct(get(st.file), st); err != nil {
			return nil, err
		}
	}
	for _, et := range g.enums {
		emitEnum(get(et.file), et)
	}

	out := make(map[string]string, len(files))
	for name, b := range files {
		var src strings.Builder
		src.WriteString(GeneratedMarker + "\n\n")
		src.WriteString("package " + targetPackage + "\n\n")
		var imports []string
		if needsXML[name] {
			imports = append(imports, "encoding/xml")
		}
		if g.needsErrors[b] {
			imports = append(imports, "errors")
		}
		switch len(imports) {
		case 0:
		case 1:
			fmt.Fprintf(&src, "import %q\n\n", imports[0])
		default:
			sort.Strings(imports)
			src.WriteString("import (\n")
			for _, im := range imports {
				fmt.Fprintf(&src, "%q\n", im)
			}
			src.WriteString(")\n\n")
		}
		src.WriteString(b.String())

		formatted, err := format.Source([]byte(src.String()))
		if err != nil {
			return nil, fmt.Errorf("codegen: %s: generated source does not parse: %w", name, err)
		}
		out[name] = string(formatted)
	}
	return out, nil
}

func (g *gen) emitRequest(b *strings.Builder, c *schema.Command) error {
	name := goName(c.Name)
	fmt.Fprintf(b, "// %s is a %s command request. See the GMP protocol docs for details.\n", name, c.Name)
	if w, ok := g.cw[c.Name]; ok {
		g.cwConsumed[c.Name] = true
		b.WriteString("//\n")
		for _, line := range strings.Split(w.Text, "\n") {
			fmt.Fprintf(b, "// %s\n", line)
		}
	}
	fmt.Fprintf(b, "type %s struct {\n", name)
	used := map[string]bool{"XMLName": true}
	fmt.Fprintf(b, "XMLName xml.Name `xml:%q`\n", c.Name)
	if err := g.emitChardata(b, c.Request, used); err != nil {
		return err
	}
	extra, _, err := g.fieldOverridesFor(&structType{node: c.Request, uses: []use{{command: c.Name}}})
	if err != nil {
		return err
	}
	typeOv, err := g.fieldTypeOverridesFor(&structType{node: c.Request, uses: []use{{command: c.Name}}})
	if err != nil {
		return err
	}
	fields, collided, err := g.emitFields(c.Request.Children, used, nil, typeOv, c.Name)
	if err != nil {
		return err
	}
	for _, ov := range extra {
		fields = append(fields, g.overrideField(ov, used))
	}
	for _, f := range fields {
		writeField(b, f)
	}
	b.WriteString("}\n\n")

	if collided {
		if err := g.emitCollisionMethods(b, name, true, false, c.Request.Chardata, fields, c.Name); err != nil {
			return fmt.Errorf("codegen: %s: request: %w", c.Name, err)
		}
	}
	return nil
}

func (g *gen) emitResponse(b *strings.Builder, c *schema.Command) error {
	name := goName(c.Name) + "Response"
	root := c.Response.Name
	if root == "" {
		root = c.Name + "_response"
	}
	fmt.Fprintf(b, "// %s is a %s command response. See the GMP protocol docs for details.\n", name, c.Name)
	if w, ok := g.cw[c.Name+"_response"]; ok {
		g.cwConsumed[c.Name+"_response"] = true
		b.WriteString("//\n")
		for _, line := range strings.Split(w.Text, "\n") {
			fmt.Fprintf(b, "// %s\n", line)
		}
	}
	fmt.Fprintf(b, "type %s struct {\n", name)
	used := map[string]bool{"XMLName": true, "Status": true, "StatusText": true}
	fmt.Fprintf(b, "XMLName xml.Name `xml:%q`\n", root)
	// status/status_text are gvmd's universal response envelope, declared by
	// every command's schema and modelled identically by all 122 hand-written
	// response types. They are emitted first, from here rather than from the
	// schema children, so their order and spelling cannot drift per command.
	b.WriteString("Status string `xml:\"status,attr\"`\n")
	b.WriteString("StatusText string `xml:\"status_text,attr\"`\n")
	if err := g.emitChardata(b, c.Response, used); err != nil {
		return err
	}
	extra, _, err := g.fieldOverridesFor(&structType{node: c.Response, uses: []use{{command: c.Name}}})
	if err != nil {
		return err
	}
	typeOv, err := g.fieldTypeOverridesFor(&structType{node: c.Response, uses: []use{{command: c.Name}}})
	if err != nil {
		return err
	}
	fields, collided, err := g.emitFields(c.Response.Children, used, envelopeAttrs, typeOv, c.Name)
	if err != nil {
		return err
	}
	for _, ov := range extra {
		fields = append(fields, g.overrideField(ov, used))
	}
	for _, f := range fields {
		writeField(b, f)
	}
	b.WriteString("}\n\n")

	fmt.Fprintf(b, "// GetStatus returns the status and status text from the response.\n")
	fmt.Fprintf(b, "func (r *%s) GetStatus() (string, string) {\n\treturn r.Status, r.StatusText\n}\n\n", name)

	if collided {
		if err := g.emitCollisionMethods(b, name, true, true, c.Response.Chardata, fields, root); err != nil {
			return fmt.Errorf("codegen: %s: response: %w", c.Name, err)
		}
	}
	return nil
}

// fieldKey identifies a schema-declared child by wire name and
// attribute-vs-element kind, for use as an emitFields skip-set key.
type fieldKey struct {
	name string
	attr bool
}

// typeOverride is fieldTypeOverridesFor's per-field resolution: the retyped
// Go type, and whether FieldTypeOverride.Pointer additionally forces it to a
// pointer with omitempty despite the schema marking the field required — see
// FieldTypeOverride's doc comment.
type typeOverride struct {
	goType  string
	pointer bool
}

// envelopeAttrs are the response attributes emitRequest/emitResponse already
// declare by hand, and must not be declared a second time from the schema.
var envelopeAttrs = map[fieldKey]bool{
	{name: "status", attr: true}:      true,
	{name: "status_text", attr: true}: true,
}

func (g *gen) emitStruct(b *strings.Builder, st *structType) error {
	fmt.Fprintf(b, "// %s is the <%s> element of %s. See the GMP protocol docs for details.\n",
		st.goName, st.node.Name, g.ownerOf(st))
	fmt.Fprintf(b, "type %s struct {\n", st.goName)
	used := map[string]bool{}

	extra, chardataOverride, err := g.fieldOverridesFor(st)
	if err != nil {
		return err
	}
	suppressed, chardataSuppression, err := g.fieldSuppressionsFor(st)
	if err != nil {
		return err
	}
	typeOv, err := g.fieldTypeOverridesFor(st)
	if err != nil {
		return err
	}
	if chardataSuppression == nil {
		if err := g.emitChardataOverridable(b, st.node, used, chardataOverride); err != nil {
			return err
		}
	}
	fields, collided, err := g.emitFields(st.node.Children, used, suppressed, typeOv, "")
	if err != nil {
		return err
	}
	for _, ov := range extra {
		fields = append(fields, g.overrideField(ov, used))
	}
	for _, f := range fields {
		writeField(b, f)
	}
	b.WriteString("}\n\n")

	chardata := (st.node.Chardata || chardataOverride != nil) && chardataSuppression == nil
	if collided {
		if err := g.emitCollisionMethods(b, st.goName, false, false, chardata, fields, ""); err != nil {
			return fmt.Errorf("codegen: %s: %w", st.goName, err)
		}
	}
	return nil
}

// emitEnum declares the named string type and its constants for one
// <alts>-constrained wire name. Open, not closed: the underlying type is a
// plain string, so decoding an unrecognized server value neither fails nor
// loses the value.
func emitEnum(b *strings.Builder, et *enumType) {
	fmt.Fprintf(b, "// %s is the <%s> value of %s. See the GMP protocol docs for details.\n",
		et.goName, et.name, ownerOfUses(et.uses))
	fmt.Fprintf(b, "type %s string\n\n", et.goName)
	b.WriteString("const (\n")
	for _, c := range et.consts {
		fmt.Fprintf(b, "%s %s = %q\n", c.goName, et.goName, c.value)
	}
	b.WriteString(")\n\n")
}

// overrideResult is one key's resolved FieldOverrides for a structType: the
// chardata entry, if any, and the ordinary element/attribute entries.
type overrideResult struct {
	extra    []FieldOverride
	chardata *FieldOverride
}

// fieldOverridesFor looks up st in FieldOverrides by every path it was
// reached from. FieldOverrides' doc comment says a shape keyed this way must
// be reached by exactly one command, so every use has to resolve to the same
// overrideResult — including a use with no matching g.fo entry at all, which
// resolves to the empty overrideResult rather than being skipped: skipping it
// is exactly what let an override on one use silently leak onto a structurally
// -deduplicated sibling use that had none. Any disagreement between two uses'
// results is an error instead of a silent pick.
func (g *gen) fieldOverridesFor(st *structType) (extra []FieldOverride, chardataOverride *FieldOverride, err error) {
	seen := map[string]bool{}
	haveFirst := false
	var firstKey string
	var first overrideResult

	for _, u := range st.uses {
		key := fieldOverrideKey(u.command, u.path, st.node.Name)
		if seen[key] {
			continue
		}
		seen[key] = true

		var res overrideResult
		if entries, ok := g.fo[key]; ok {
			g.foConsumed[key] = true
			for _, e := range entries {
				e := e
				if e.Chardata {
					res.chardata = &e
					continue
				}
				if schemaChildCollision(st.node, e) {
					return nil, nil, fmt.Errorf(
						"codegen: FieldOverrides[%q]: override %q duplicates an existing schema-declared child of <%s>",
						key, e.Name, st.node.Name)
				}
				res.extra = append(res.extra, e)
			}
		}

		if !haveFirst {
			haveFirst, firstKey, first = true, key, res
			extra, chardataOverride = res.extra, res.chardata
			continue
		}
		if !overrideResultEqual(&first, &res) {
			return nil, nil, fmt.Errorf(
				"codegen: %s: FieldOverrides for %q and %q disagree; a structurally-deduplicated type "+
					"must resolve to the same override set from every path that reaches it",
				st.goName, firstKey, key)
		}
	}
	return extra, chardataOverride, nil
}

// fieldSuppressionsFor looks up st in FieldSuppressions by every path it was
// reached from, mirroring fieldOverridesFor's per-use agreement discipline:
// every use — including one with no matching g.fs entry, which resolves to
// the empty set rather than being skipped — must resolve to the same
// suppression set, or a suppression on one use could leak onto a
// structurally-deduplicated sibling use that had none. The chardata
// suppression is tracked separately from the (name, attr)-keyed map since a
// chardata entry has neither; its Reason is returned (rather than a bare
// bool) so callers emitting the relocation's source-side field can carry it
// forward as that field's own gmpoverride tag, for internal/schema's
// comparer to recover by reflection.
func (g *gen) fieldSuppressionsFor(st *structType) (map[fieldKey]bool, *FieldSuppression, error) {
	seen := map[string]bool{}
	haveFirst := false
	var firstKey string
	var first map[fieldKey]bool
	var firstChardata *FieldSuppression

	for _, u := range st.uses {
		key := fieldOverrideKey(u.command, u.path, st.node.Name)
		if seen[key] {
			continue
		}
		seen[key] = true

		res := map[fieldKey]bool{}
		var chardataSuppression *FieldSuppression
		if entries, ok := g.fs[key]; ok {
			g.fsConsumed[key] = true
			for _, e := range entries {
				e := e
				if e.Chardata {
					chardataSuppression = &e
					continue
				}
				res[fieldKey{name: e.Name, attr: e.Attr}] = true
			}
		}

		if !haveFirst {
			haveFirst, firstKey, first, firstChardata = true, key, res, chardataSuppression
			continue
		}
		if !fieldKeySetEqual(first, res) || !fieldSuppressionEqual(firstChardata, chardataSuppression) {
			return nil, nil, fmt.Errorf(
				"codegen: %s: FieldSuppressions for %q and %q disagree; a structurally-deduplicated type "+
					"must resolve to the same suppression set from every path that reaches it",
				st.goName, firstKey, key)
		}
	}
	return first, firstChardata, nil
}

// fieldSuppressionEqual reports whether two chardata suppression results
// from fieldSuppressionsFor agree: both absent, or both present with the
// same relocation target and reason.
func fieldSuppressionEqual(a, b *FieldSuppression) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	return a == nil || (a.RelocatedTo == b.RelocatedTo && a.Reason == b.Reason)
}

// fieldKeySetEqual compares two fieldSuppressionsFor results for value
// equality, mirroring overrideResultEqual's role for fieldOverridesFor.
func fieldKeySetEqual(a, b map[fieldKey]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// fieldTypeOverridesFor looks up st in FieldTypeOverrides by every path it
// was reached from, mirroring fieldSuppressionsFor's per-use agreement
// discipline: every use must resolve to the same retype set, or a correction
// on one use could leak onto a structurally-deduplicated sibling use that had
// none.
func (g *gen) fieldTypeOverridesFor(st *structType) (map[fieldKey]typeOverride, error) {
	seen := map[string]bool{}
	haveFirst := false
	var firstKey string
	var first map[fieldKey]typeOverride

	for _, u := range st.uses {
		key := fieldOverrideKey(u.command, u.path, st.node.Name)
		if seen[key] {
			continue
		}
		seen[key] = true

		res := map[fieldKey]typeOverride{}
		if entries, ok := g.tov[key]; ok {
			g.tovConsumed[key] = true
			for _, e := range entries {
				res[fieldKey{name: e.Name, attr: e.Attr}] = typeOverride{goType: e.GoType, pointer: e.Pointer}
			}
		}

		if !haveFirst {
			haveFirst, firstKey, first = true, key, res
			continue
		}
		if !fieldTypeSetEqual(first, res) {
			return nil, fmt.Errorf(
				"codegen: %s: FieldTypeOverrides for %q and %q disagree; a structurally-deduplicated type "+
					"must resolve to the same retype set from every path that reaches it",
				st.goName, firstKey, key)
		}
	}
	return first, nil
}

// fieldTypeSetEqual compares two fieldTypeOverridesFor results for value
// equality, mirroring fieldKeySetEqual's role for fieldSuppressionsFor.
func fieldTypeSetEqual(a, b map[fieldKey]typeOverride) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// schemaChildCollision reports whether e would duplicate an already
// schema-declared child of node: same wire name, same attribute-vs-element
// kind. A FieldOverride exists to add a field the schema cannot produce, so
// one that shadows a real schema field is a defect in the table, not a
// legitimate addition.
func schemaChildCollision(node *schema.Node, e FieldOverride) bool {
	kind := schema.KindElement
	if e.Attr {
		kind = schema.KindAttribute
	}
	for _, c := range node.Children {
		if c.Name == e.Name && c.Kind == kind {
			return true
		}
	}
	return false
}

// overrideResultEqual compares two resolved overrideResults for value
// equality: chardata by value (nil only equals nil), extra as a multiset
// (content, not order — a single key's slice is ordered in practice, but nothing
// should rely on that to prove equality).
func overrideResultEqual(a, b *overrideResult) bool {
	switch {
	case a.chardata == nil && b.chardata != nil, a.chardata != nil && b.chardata == nil:
		return false
	case a.chardata != nil && b.chardata != nil && *a.chardata != *b.chardata:
		return false
	}
	if len(a.extra) != len(b.extra) {
		return false
	}
	used := make([]bool, len(b.extra))
	for _, x := range a.extra {
		found := false
		for j, y := range b.extra {
			if !used[j] && x == y {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// overrideField builds the Go field for one non-chardata FieldOverride, via a
// synthetic schema.Node so it participates in the same collision/attribute
// machinery a schema-declared field would (see emitCollisionMethods).
func (g *gen) overrideField(ov FieldOverride, used map[string]bool) field {
	kind := schema.KindElement
	if ov.Attr {
		kind = schema.KindAttribute
	}
	node := &schema.Node{Name: ov.Name, Kind: kind, Card: schema.CardRequired}

	tag := ov.Name
	if ov.Attr {
		tag += ",attr"
	}
	if ov.OmitEmpty {
		tag += ",omitempty"
	}

	return field{
		name: uniqueName(goName(ov.Name), used), goType: ov.GoType, tag: tag, node: node,
		overrideReason: ov.Reason,
	}
}

// ownerOf describes, for the doc comment, which commands reach a shape.
func (g *gen) ownerOf(st *structType) string {
	return ownerOfUses(st.uses)
}

// ownerOfUses is ownerOf's shared implementation, usable from both a
// structType and an enumType (and from registerEnum's error message, before
// any enumType even exists to hold the uses).
func ownerOfUses(uses []use) string {
	var names []string
	seen := map[string]bool{}
	for _, u := range uses {
		if !seen[u.command] {
			seen[u.command] = true
			names = append(names, u.command)
		}
	}
	sort.Strings(names)
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

// emitFields builds the field list for a struct's schema-derived children.
// It reports collided when two KindElement children share one wire name — the
// DSL's <id> escape hatch (see internal/schema's env.scope) lets that happen
// deliberately (get_filters's two <filters> definitions) — since the second
// and later occurrence cannot keep its own tag: encoding/xml's reflection-based
// field routing rejects two fields with an identical explicit tag outright.
// Attribute names cannot collide this way: a start tag cannot declare one
// attribute twice.
func (g *gen) emitFields(kids []*schema.Node, used map[string]bool, skip map[fieldKey]bool, typeOv map[fieldKey]typeOverride, command string) ([]field, bool, error) {
	var fields []field
	seen := map[string]bool{}
	collided := false
	for _, kid := range kids {
		if skip[fieldKey{name: kid.Name, attr: kid.Kind == schema.KindAttribute}] {
			continue
		}
		f, err := g.fieldOf(command, kid, used, typeOv)
		if err != nil {
			return nil, false, err
		}
		if kid.Kind == schema.KindElement {
			if seen[kid.Name] {
				f.wireTag = f.tag
				f.tag = "-"
				collided = true
			}
			seen[kid.Name] = true
		}
		fields = append(fields, f)
	}
	return fields, collided, nil
}

func writeField(b *strings.Builder, f field) {
	switch {
	case f.wireTag != "" && f.overrideReason != "":
		// A collision-group member (get_info's schema-patched wrapper
		// sibling) that also needs a suppression reason: gmpwire and
		// gmpoverride are independent tag keys — internal/schema/gomodel.go's
		// leafNode reads gmpoverride unconditionally, regardless of which
		// path (xml or gmpwire) supplied the field's name/opts — so both can
		// coexist on one field without either shadowing the other.
		fmt.Fprintf(b, "%s %s `xml:%q gmpwire:%q gmpoverride:%q`\n", f.name, f.goType, f.tag, f.wireTag, f.overrideReason)
	case f.wireTag != "":
		// gmpwire carries this field's real wire identity for
		// internal/schema/gomodel.go's addField to recover by reflection —
		// xml:"-" alone would make the field invisible to the structural
		// conformance oracle, reopening the MISSING/EXTRA gap this tag exists
		// to close.
		fmt.Fprintf(b, "%s %s `xml:%q gmpwire:%q`\n", f.name, f.goType, f.tag, f.wireTag)
	case f.overrideReason != "":
		fmt.Fprintf(b, "%s %s `xml:%q gmpoverride:%q`\n", f.name, f.goType, f.tag, f.overrideReason)
	default:
		fmt.Fprintf(b, "%s %s `xml:%q`\n", f.name, f.goType, f.tag)
	}
}

// --- duplicate-wire-tag codec generation -------------------------------------
//
// A struct built by emitFields can end up with two fields sharing one wire
// name (see emitFields' doc comment). encoding/xml refuses to marshal or
// unmarshal such a struct through its normal reflection-based path in either
// direction ("field X with tag Y conflicts with field Z"), because both
// directions build their field-routing table from struct tags. A custom
// UnmarshalXML/MarshalXML pair sidesteps that lookup entirely by driving the
// decoder/encoder directly.
//
// The decoder routes by attribute signature, not occurrence order: gvmd's
// schema does not promise wire order matches declaration order, so routing
// the Nth occurrence to the Nth field is unsound against a reordered
// document. Instead collisionDiscriminators computes, per member, the
// start-tag attribute name(s) that appear on it and no sibling (get_filters's
// two <filters> are disjoint: one carries id/term/keywords/name, the other
// start/max), and emitClassifyFunc emits a function that reads a decoded
// start tag's attributes and returns which member they identify — failing
// closed on zero matches or more than one, rather than guessing.
//
// Two shapes fail loudly rather than guess: mixed content combined with a
// collision (chardata and a hand-driven token loop don't agree on who owns
// the decoder's cursor), and a dynamic attribute field sharing a struct with
// a collision (attribute names can't collide, so nothing in the corpus needs
// this; out of scope until something does).
func (g *gen) emitCollisionMethods(b *strings.Builder, typeName string, hasXMLName, envelope, chardata bool, fields []field, rootName string) error {
	if chardata {
		return fmt.Errorf("mixed content combined with a wire-tag collision is not supported")
	}
	for _, f := range fields {
		if f.node.Kind == schema.KindAttribute {
			return fmt.Errorf("dynamic attribute field %s cannot coexist with a wire-tag collision", f.name)
		}
	}
	byName, order := collisionGroups(fields)
	disc := map[string][][]string{}
	for _, name := range order {
		fs := byName[name]
		if len(fs) < 2 {
			continue
		}
		d, err := collisionDiscriminators(typeName, name, fs)
		if err != nil {
			return err
		}
		disc[name] = d
	}
	g.emitCollisionUnmarshal(b, typeName, hasXMLName, envelope, fields, disc)
	g.emitCollisionMarshal(b, typeName, hasXMLName, envelope, fields, rootName)
	return nil
}

// collisionDiscriminators computes, for each member of a same-wire-name
// collision group, the start-tag attribute names that appear on that
// member's schema node and on no sibling's — the signature that proves a
// decoded start tag is this occurrence and not another. It errors when a
// member has none, since dispatch can't fail closed without one.
func collisionDiscriminators(typeName, name string, fs []field) ([][]string, error) {
	attrSets := make([]map[string]bool, len(fs))
	for i, f := range fs {
		set := map[string]bool{}
		for _, c := range f.node.Children {
			if c.Kind == schema.KindAttribute {
				set[c.Name] = true
			}
		}
		attrSets[i] = set
	}

	disc := make([][]string, len(fs))
	for i := range fs {
		var names []string
		for a := range attrSets[i] {
			exclusive := true
			for j := range fs {
				if j != i && attrSets[j][a] {
					exclusive = false
					break
				}
			}
			if exclusive {
				names = append(names, a)
			}
		}
		if len(names) == 0 {
			return nil, fmt.Errorf("codegen: %s: wire-tag collision on <%s>: member %d (%s) has no "+
				"attribute that uniquely identifies it among its siblings; attribute-signature dispatch "+
				"requires every colliding member to be distinguishable by a start-tag attribute no sibling shares",
				typeName, name, i, fs[i].name)
		}
		sort.Strings(names)
		disc[i] = names
	}
	return disc, nil
}

// collisionGroups orders a collided struct's fields by their shared wire
// name, first-occurrence order, so the emitted switch's case order matches
// the struct's declaration order.
func collisionGroups(fields []field) (byName map[string][]field, order []string) {
	byName = map[string][]field{}
	for _, f := range fields {
		name := f.node.Name
		if _, ok := byName[name]; !ok {
			order = append(order, name)
		}
		byName[name] = append(byName[name], f)
	}
	return byName, order
}

func (g *gen) emitCollisionUnmarshal(b *strings.Builder, typeName string, hasXMLName, envelope bool, fields []field, disc map[string][][]string) {
	byName, order := collisionGroups(fields)
	// emitCollisionMethods only calls this when collided is true, which
	// guarantees at least one name below has 2+ fields, so a classify
	// function's "unrecognized"/"ambiguous" errors and the "unexpected extra
	// occurrence" branch are always emitted, and this file always needs the
	// errors package.
	g.needsErrors[b] = true

	for _, name := range order {
		if len(byName[name]) > 1 {
			g.emitClassifyFunc(b, typeName, name, disc[name])
		}
	}

	fmt.Fprintf(b, "// UnmarshalXML decodes %s by hand: two of its fields share one wire tag,\n", typeName)
	fmt.Fprintf(b, "// which encoding/xml's reflection-based field routing cannot express. The\n")
	fmt.Fprintf(b, "// colliding members are told apart by attribute signature (see the\n")
	fmt.Fprintf(b, "// classify functions above), not by arrival order. See the generator's\n")
	fmt.Fprintf(b, "// emitCollisionMethods for why this is safe.\n")
	fmt.Fprintf(b, "func (r *%s) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {\n", typeName)
	// A decode target isn't necessarily zero-valued — decoding into an
	// already-populated value (deliberately, or by reuse) must not leave
	// stale data from whatever the receiver held before this call.
	fmt.Fprintf(b, "*r = %s{}\n", typeName)
	if hasXMLName {
		b.WriteString("r.XMLName = start.Name\n")
	}
	if envelope {
		b.WriteString("for _, a := range start.Attr {\n")
		b.WriteString("switch a.Name.Local {\n")
		b.WriteString("case \"status\":\nr.Status = a.Value\n")
		b.WriteString("case \"status_text\":\nr.StatusText = a.Value\n")
		b.WriteString("}\n}\n")
	}
	for _, name := range order {
		if len(byName[name]) > 1 {
			fmt.Fprintf(b, "seen%s := make([]bool, %d)\n", goName(name), len(byName[name]))
		}
	}
	b.WriteString("for {\n")
	b.WriteString("tok, err := d.Token()\n")
	b.WriteString("if err != nil {\nreturn err\n}\n")
	b.WriteString("switch t := tok.(type) {\n")
	b.WriteString("case xml.StartElement:\n")
	b.WriteString("switch t.Name.Local {\n")
	for _, name := range order {
		fs := byName[name]
		fmt.Fprintf(b, "case %q:\n", name)
		if len(fs) == 1 {
			b.WriteString(decodeInto("r."+fs[0].name, fs[0]))
			continue
		}
		fmt.Fprintf(b, "idx, err := %s(t.Attr)\n", classifyFuncName(typeName, name))
		b.WriteString("if err != nil {\nreturn err\n}\n")
		b.WriteString("switch idx {\n")
		for i, f := range fs {
			fmt.Fprintf(b, "case %d:\n", i)
			// A slice member (get_info's repeated data-item variant, alongside
			// its singular pagination-wrapper sibling) legitimately occurs more
			// than once; only a non-slice member's second occurrence is a real
			// defect worth failing closed on.
			if !strings.HasPrefix(f.goType, "[]") {
				fmt.Fprintf(b, "if seen%s[%d] {\n", goName(name), i)
				fmt.Fprintf(b, "return errors.New(%q)\n", fmt.Sprintf("%s: unexpected extra <%s> element", typeName, name))
				b.WriteString("}\n")
				fmt.Fprintf(b, "seen%s[%d] = true\n", goName(name), i)
			}
			b.WriteString(decodeInto("r."+f.name, f))
		}
		b.WriteString("}\n")
	}
	b.WriteString("default:\nif err := d.Skip(); err != nil {\nreturn err\n}\n")
	b.WriteString("}\n") // end switch t.Name.Local
	b.WriteString("case xml.EndElement:\n")
	b.WriteString("return nil\n")
	b.WriteString("}\n") // end switch tok.(type)
	b.WriteString("}\n") // end for
	b.WriteString("}\n\n")
}

// classifyFuncName names the package-level function emitClassifyFunc emits
// for one collision group, so emitCollisionUnmarshal's dispatch switch and
// the function definition agree on what to call.
func classifyFuncName(typeName, name string) string {
	return "classify" + typeName + goName(name)
}

// emitClassifyFunc emits a function that identifies which member of a
// same-wire-name collision group a decoded start tag's attributes belong to.
// disc[i] is member i's discriminator attribute name(s), from
// collisionDiscriminators; zero or multiple matches both fail closed.
func (g *gen) emitClassifyFunc(b *strings.Builder, typeName, name string, disc [][]string) {
	fnName := classifyFuncName(typeName, name)
	fmt.Fprintf(b, "// %s identifies which schema-declared <%s> variant a start tag\n", fnName, name)
	fmt.Fprintf(b, "// belongs to, from its attributes: %s's colliding members are distinguished\n", typeName)
	b.WriteString("// by disjoint attribute signatures, not by arrival order, so a reordered or\n")
	b.WriteString("// malformed document fails closed instead of being misrouted.\n")
	fmt.Fprintf(b, "func %s(attrs []xml.Attr) (int, error) {\n", fnName)
	b.WriteString("has := map[string]bool{}\n")
	b.WriteString("for _, a := range attrs {\nhas[a.Name.Local] = true\n}\n")
	b.WriteString("match := -1\n")
	for i, names := range disc {
		var conds []string
		for _, nm := range names {
			conds = append(conds, fmt.Sprintf("has[%q]", nm))
		}
		fmt.Fprintf(b, "if %s {\n", strings.Join(conds, " || "))
		fmt.Fprintf(b, "if match != -1 && match != %d {\nreturn -1, errors.New(%q)\n}\n",
			i, fmt.Sprintf("%s: ambiguous <%s> element: attributes match more than one known variant", typeName, name))
		fmt.Fprintf(b, "match = %d\n", i)
		b.WriteString("}\n")
	}
	fmt.Fprintf(b, "if match == -1 {\nreturn -1, errors.New(%q)\n}\n",
		fmt.Sprintf("%s: unrecognized <%s> element: no known variant's attributes are present", typeName, name))
	b.WriteString("return match, nil\n")
	b.WriteString("}\n\n")
}

// decodeInto emits the statement(s) that decode one child element, held in
// the loop variable t, into fieldExpr, matching the field's cardinality: a
// slice appends a freshly decoded element, a pointer allocates before
// decoding into it, and everything else decodes in place.
func decodeInto(fieldExpr string, f field) string {
	switch {
	case strings.HasPrefix(f.goType, "[]"):
		elem := strings.TrimPrefix(f.goType, "[]")
		return fmt.Sprintf(
			"var v %s\nif err := d.DecodeElement(&v, &t); err != nil {\nreturn err\n}\n%s = append(%s, v)\n",
			elem, fieldExpr, fieldExpr)
	case strings.HasPrefix(f.goType, "*"):
		elem := strings.TrimPrefix(f.goType, "*")
		return fmt.Sprintf(
			"%s = new(%s)\nif err := d.DecodeElement(%s, &t); err != nil {\nreturn err\n}\n",
			fieldExpr, elem, fieldExpr)
	default:
		return fmt.Sprintf(
			"if err := d.DecodeElement(&%s, &t); err != nil {\nreturn err\n}\n",
			fieldExpr)
	}
}

func (g *gen) emitCollisionMarshal(b *strings.Builder, typeName string, hasXMLName, envelope bool, fields []field, rootName string) {
	fmt.Fprintf(b, "// MarshalXML encodes %s by hand, mirroring UnmarshalXML: manual\n", typeName)
	fmt.Fprintf(b, "// EncodeElement calls bypass the same reflection-based tag routing that\n")
	fmt.Fprintf(b, "// cannot express two fields sharing one wire tag.\n")
	fmt.Fprintf(b, "func (r *%s) MarshalXML(e *xml.Encoder, start xml.StartElement) error {\n", typeName)
	if hasXMLName {
		b.WriteString("start.Name = r.XMLName\n")
		// xml.Marshal's own defaultStart falls back to the Go type name, not
		// this type's declared wire name, when there is no field-tag context
		// to read it from (true for any top-level Marshal call on a type with
		// a custom MarshalXML) — so a zero-value r.XMLName must not be copied
		// over that default verbatim, or the result has no name at all
		// ("xml: start tag with no name"). Falling back to the literal,
		// schema-known root name here matches what a plain XMLName struct tag
		// would have produced through the reflection-based path.
		fmt.Fprintf(b, "if start.Name.Local == \"\" {\nstart.Name.Local = %q\n}\n", rootName)
	}
	if envelope {
		b.WriteString("start.Attr = append(start.Attr[:0:0], ")
		b.WriteString(`xml.Attr{Name: xml.Name{Local: "status"}, Value: r.Status}, `)
		b.WriteString("xml.Attr{Name: xml.Name{Local: \"status_text\"}, Value: r.StatusText})\n")
	}
	b.WriteString("if err := e.EncodeToken(start); err != nil {\nreturn err\n}\n")
	for _, f := range fields {
		b.WriteString(encodeStmt("r."+f.name, f))
	}
	b.WriteString("if err := e.EncodeToken(start.End()); err != nil {\nreturn err\n}\n")
	b.WriteString("return nil\n")
	b.WriteString("}\n\n")
}

// encodeStmt emits the statement(s) that encode fieldExpr as a child element
// named after f.node's wire name, via a manual EncodeElement call — the part
// that legally produces two same-named sibling elements, since it bypasses
// the reflection-based tag lookup that would otherwise reject the second one.
func encodeStmt(fieldExpr string, f field) string {
	name := f.node.Name
	switch {
	case strings.HasPrefix(f.goType, "[]"):
		return fmt.Sprintf(
			"for _, v := range %s {\nif err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: %q}}); err != nil {\nreturn err\n}\n}\n",
			fieldExpr, name)
	case strings.HasPrefix(f.goType, "*"):
		return fmt.Sprintf(
			"if %s != nil {\nif err := e.EncodeElement(%s, xml.StartElement{Name: xml.Name{Local: %q}}); err != nil {\nreturn err\n}\n}\n",
			fieldExpr, fieldExpr, name)
	default:
		return fmt.Sprintf(
			"if err := e.EncodeElement(%s, xml.StartElement{Name: xml.Name{Local: %q}}); err != nil {\nreturn err\n}\n",
			fieldExpr, name)
	}
}

// --- identifiers -------------------------------------------------------------

// initialisms are the acronyms and other fixed-casing exceptions the
// hand-written commands/ package already spells this way — most are full-caps
// acronyms, but "icalendar" (RFC 5545's own spelling, not an acronym) is a
// mixed-case exception. Keeping the generator's spelling identical is what
// lets generated types drop in beside them without renaming every call site.
//
// Plurals of these (e.g. "nvts", "cves") are handled generically by goWord,
// not listed here — only add an entry for a word that isn't itself just
// initialism+"s".
var initialisms = map[string]string{
	"api":       "API",
	"ca":        "CA",
	"cpe":       "CPE",
	"cve":       "CVE",
	"cvss":      "CVSS",
	"dn":        "DN",
	"dns":       "DNS",
	"epss":      "EPSS",
	"esxi":      "ESXI",
	"hr":        "HR",
	"http":      "HTTP",
	"icalendar": "ICalendar",
	"icmp":      "ICMP",
	"id":        "ID",
	"ip":        "IP",
	"lsc":       "LSC",
	"md5":       "MD5",
	"nvt":       "NVT",
	"oci":       "OCI",
	"oid":       "OID",
	"os":        "OS",
	"qod":       "QOD",
	"sha256":    "SHA256",
	"smb":       "SMB",
	"snmp":      "SNMP",
	"ssh":       "SSH",
	"tcp":       "TCP",
	"tls":       "TLS",
	"udp":       "UDP",
	"url":       "URL",
	"uuid":      "UUID",
	"vt":        "VT",
	"xml":       "XML",
}

// goName turns a schema element or command name into an exported Go
// identifier: snake_case to CamelCase, with known acronyms capitalised.
func goName(s string) string {
	var b strings.Builder
	for _, part := range strings.FieldsFunc(s, isNameSeparator) {
		b.WriteString(goWord(part))
	}
	if b.Len() == 0 {
		return "Field"
	}
	return b.String()
}

func isNameSeparator(r rune) bool {
	return r == '_' || r == '-' || r == '.' || r == ':' || r == ' '
}

// goWord capitalises a single already-split word, expanding it via
// initialisms if it's a known acronym or the plural of one — GMP schema
// tokens routinely pluralize an initialism without a separator ("nvts",
// "cves", "cpes"), and goName's caller has already split on underscores, so
// this is the only place that can still see "nvt" hiding inside "nvts".
func goWord(part string) string {
	lower := strings.ToLower(part)
	if up, ok := initialisms[lower]; ok {
		return up
	}
	if base, ok := strings.CutSuffix(lower, "s"); ok {
		if up, ok := initialisms[base]; ok {
			return up + "s"
		}
	}
	return strings.ToUpper(lower[:1]) + lower[1:]
}
