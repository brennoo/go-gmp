// Command gengmp materializes internal/codegen's generated Go source into a
// directory, or (with -check) reports how the already-committed source differs
// from what would be generated now.
//
// It is a thin wrapper: load the generation manifest, load the committed
// Stage-A schema snapshots it names, call codegen, write or compare what comes
// back. It never reads vendor/XML/commands itself — that AGPL-licensed
// extraction step already happened when the snapshots under -schema were
// committed (internal/codegen/snapshot_test.go -update), and keeping this CLI
// on the far side of that boundary is what T-9's licensing isolation
// (acceptance criterion 9) requires: nothing downstream of the snapshot ever
// needs to know vendor/XML existed.
//
// Writing is a manual, human-reviewed step, not something CI runs to write
// code: see hack/generate.sh. -check is the read-only counterpart CI can run —
// it touches neither the output directory nor the git index.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/brennoo/go-gmp/internal/codegen"
)

// checkOutDir is where -check looks when -out is not given. Unlike the write
// path, defaulting is safe here: nothing is written, and commands/ is the one
// real generated tree this repo has.
const checkOutDir = "./commands"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "gengmp:", err)
		os.Exit(1)
	}
}

type options struct {
	manifest string
	schema   string
	out      string
	check    bool
}

func parse(args []string) (options, error) {
	var o options
	fs := flag.NewFlagSet("gengmp", flag.ContinueOnError)
	fs.StringVar(&o.manifest, "manifest", "internal/codegen/manifest.json",
		"the generation manifest: the single declaration of which command families are generated")
	fs.StringVar(&o.schema, "schema", "internal/codegen/testdata/schema",
		"directory of committed Stage-A schema snapshots (internal/codegen/snapshot.go); never vendor/XML directly")
	fs.StringVar(&o.out, "out", "",
		"directory to write generated Go files into. Must be passed explicitly — there is no destructive default, "+
			"so a forgotten -out fails loudly instead of writing into the current directory. "+
			"With -check it names the existing directory to inspect instead, and defaults to "+checkOutDir+".")
	fs.BoolVar(&o.check, "check", false,
		"report how the generated files already in -out differ from what would be generated now, and exit non-zero "+
			"if they do. Writes nothing, anywhere.")
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	if fs.NArg() > 0 {
		return options{}, fmt.Errorf("unexpected argument %q; gengmp takes flags only", fs.Arg(0))
	}
	switch {
	case o.check && o.out == "":
		o.out = checkOutDir
	case !o.check && o.out == "":
		return options{}, fmt.Errorf("-out is required; e.g. -out ./commands (see hack/generate.sh)")
	}
	return o, nil
}

func run(args []string) error {
	o, err := parse(args)
	if err != nil {
		return err
	}

	manifest, err := codegen.LoadManifest(o.manifest)
	if err != nil {
		return err
	}
	if o.check {
		return check(manifest, o)
	}
	return write(manifest, o)
}

// check is the read-only mode. It never creates, modifies or deletes anything
// under -out and never shells out to git, so it behaves identically on a clean
// tree, a dirty tree and a checkout without a .git directory.
func check(manifest codegen.Manifest, o options) error {
	info, err := os.Stat(o.out)
	if err != nil {
		return fmt.Errorf("-check needs an existing directory to inspect: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("-out %s is not a directory", o.out)
	}

	report, err := codegen.CheckOutput(manifest, o.schema, o.out, codegen.TypeOverrides, codegen.PointerOverrides, codegen.FieldOverrides, codegen.CommandWarnings, codegen.DefectAllowlist)
	if err != nil {
		return err
	}
	if report.OK() {
		// Adopted(), matching CheckOutput's own scope: a review-only family's
		// files are not in o.out and were never inspected, so counting them
		// here would overstate what the clean report covers.
		fmt.Printf("%s is up to date with the manifest and schema (%d file(s) accounted for)\n",
			o.out, len(manifest.Adopted().ExpectedOutputFiles(o.out)))
		if skipped := manifest.ReviewOnlyFamilies(); len(skipped) > 0 {
			fmt.Printf("not checked, review-only: %s\n", strings.Join(skipped, ", "))
		}
		return nil
	}
	fmt.Print(report.String())
	return fmt.Errorf("%s is out of date; re-run hack/generate.sh and review the diff", o.out)
}

// write materialises the generated tree. Like check, it works from the
// manifest's adopted families only: a review-only family's output is still
// under review as a golden, and writing it here would overwrite the
// hand-written files it is being reviewed against.
func write(manifest codegen.Manifest, o options) error {
	if skipped := manifest.ReviewOnlyFamilies(); len(skipped) > 0 {
		fmt.Printf("skipping review-only families (not written to %s): %s\n", o.out, strings.Join(skipped, ", "))
	}
	manifest = manifest.Adopted()

	cmds, err := codegen.LoadSnapshots(o.schema, manifest.AllCommands())
	if err != nil {
		return fmt.Errorf("loading snapshots from %s: %w", o.schema, err)
	}

	files, err := codegen.GenerateAll(manifest, cmds, codegen.TypeOverrides, codegen.PointerOverrides, codegen.FieldOverrides, codegen.CommandWarnings, codegen.DefectAllowlist)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if err := os.MkdirAll(o.out, 0o750); err != nil {
		return fmt.Errorf("creating %s: %w", o.out, err)
	}

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	total := 0
	for _, name := range names {
		src := files[name]
		path := filepath.Join(o.out, name)
		if err := os.WriteFile(path, []byte(src), 0o600); err != nil { //nolint:gosec // -out is an operator-supplied destination, by design
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("wrote %-40s %6d bytes\n", path, len(src))
		total += len(src)
	}
	fmt.Printf("%d file(s), %d bytes total, written to %s\n", len(names), total, o.out)
	return nil
}
