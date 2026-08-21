package gmptest

import (
	"fmt"
	"sort"
)

// RequestCase is one curated, populated request value registered for the
// request golden suite.
type RequestCase struct {
	Command string // GMP command name, e.g. "get_tasks". Must be a key in Registry.
	Name    string // Case name, unique within Command, e.g. "with_filter".
	Req     any    // A pointer to a populated commands.* request value.
}

// requestCases accumulates every RegisterRequestCase call made from cases_*.go
// init functions.
var requestCases []RequestCase

// requestCaseKeys guards against duplicate command/name pairs across files.
var requestCaseKeys = map[string]bool{}

// RegisterRequestCase adds a curated request value to the request golden
// suite. Call it from an init() in a cases_<group>.go file.
//
// It panics rather than returning an error: these calls run during package
// initialization, before any test has started, so there is no test.T to
// report a failure to. A typo in command or a copy/pasted name must stop the
// whole test binary from starting, not silently drop a case from the suite.
func RegisterRequestCase(command, name string, req any) {
	if _, ok := Registry[command]; !ok {
		panic(fmt.Sprintf("gmptest: RegisterRequestCase(%q, %q): %q is not a known GMP command; check Registry in registry.go", command, name, command))
	}
	if req == nil {
		panic(fmt.Sprintf("gmptest: RegisterRequestCase(%q, %q): req is nil", command, name))
	}

	key := command + "/" + name
	if requestCaseKeys[key] {
		panic(fmt.Sprintf("gmptest: RegisterRequestCase(%q, %q): duplicate case", command, name))
	}
	requestCaseKeys[key] = true

	requestCases = append(requestCases, RequestCase{Command: command, Name: name, Req: req})
}

// RequestCases returns every registered case, sorted by command then name so
// that test iteration order does not depend on init() ordering across files
// (which Go does not guarantee beyond "within a file, in declaration order").
func RequestCases() []RequestCase {
	out := make([]RequestCase, len(requestCases))
	copy(out, requestCases)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Command != out[j].Command {
			return out[i].Command < out[j].Command
		}
		return out[i].Name < out[j].Name
	})
	return out
}
