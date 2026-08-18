package gmptest

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// The redactor is the one part of the live-capture path that can be tested with
// no server, so it is tested hard. Everything here runs against hand-written
// XML: if a test needed a captured fixture it would be testing the capture, and
// the capture is what these tests exist to make safe.

func TestRedactPathsRewritesTextValues(t *testing.T) {
	doc := []byte(`<authenticate><credentials><username>admin</username><password>hunter2</password></credentials></authenticate>`)

	got, applied, err := RedactPaths(doc, []string{"/authenticate/credentials/password/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	if strings.Contains(string(got), "hunter2") {
		t.Errorf("The password survived redaction:\n%s", got)
	}
	if !strings.Contains(string(got), RedactionPlaceholder) {
		t.Errorf("No placeholder was written:\n%s", got)
	}
	if !strings.Contains(string(got), "<username>admin</username>") {
		t.Errorf("A path that was not in the spec was rewritten:\n%s", got)
	}
	if len(applied) != 1 || applied[0] != "/authenticate/credentials/password/#text" {
		t.Errorf("applied = %v, want the one spec path", applied)
	}
}

// TestRedactPathsIsByteSurgical pins the property that makes live fixtures
// usable as evidence: everything outside the redacted span is untouched,
// including self-closing tags, single-quoted attributes, comments and
// indentation, none of which survive a decode/re-encode round trip.
func TestRedactPathsIsByteSurgical(t *testing.T) {
	doc := []byte("<?xml version='1.0'?>\n" +
		"<!-- a comment -->\n" +
		"<create_credential>\n" +
		"  <name>x</name>\n" +
		"  <empty/>\n" +
		"  <odd  spaced = 'yes' />\n" +
		"  <password>s3cret</password>\n" +
		"</create_credential>\n")

	got, _, err := RedactPaths(doc, []string{"/create_credential/password/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}

	want := strings.Replace(string(doc), "s3cret", RedactionPlaceholder, 1)
	if string(got) != want {
		t.Errorf("Redaction changed bytes it should not have.\n--- want\n%s\n--- got\n%s", want, got)
	}
}

func TestRedactPathsRewritesAttributeValues(t *testing.T) {
	doc := []byte(`<describe_auth_response status="200"><group name="method:ldap_connect">` +
		`<auth_conf_setting key="ldaphost" value="ldap.example"/>` +
		`<auth_conf_setting key='bindpw' value = 'topsecret' />` +
		`</group></describe_auth_response>`)

	got, applied, err := RedactPaths(doc, []string{"/describe_auth_response/group/auth_conf_setting/@value"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	if strings.Contains(string(got), "topsecret") {
		t.Errorf("The bind password survived redaction:\n%s", got)
	}
	if strings.Contains(string(got), "ldap.example") {
		t.Errorf("The spec names every @value, so the host should have gone too:\n%s", got)
	}
	if !strings.Contains(string(got), `key="ldaphost"`) || !strings.Contains(string(got), `key='bindpw'`) {
		t.Errorf("The key attributes should be untouched, quoting included:\n%s", got)
	}
	if strings.Count(string(got), RedactionPlaceholder) != 2 {
		t.Errorf("Expected both occurrences to be redacted:\n%s", got)
	}
	if len(applied) != 1 {
		t.Errorf("applied = %v, want one distinct path", applied)
	}
}

func TestRedactPathsHandlesEveryOccurrenceOfARepeatedPath(t *testing.T) {
	doc := []byte(`<get_credentials_response>` +
		`<credential id="a"><package format="key">AAAA</package></credential>` +
		`<credential id="b"><package format="key">BBBB</package></credential>` +
		`</get_credentials_response>`)

	got, _, err := RedactPaths(doc, []string{"/get_credentials_response/credential/package/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	for _, secret := range []string{"AAAA", "BBBB"} {
		if strings.Contains(string(got), secret) {
			t.Errorf("%s survived redaction:\n%s", secret, got)
		}
	}
	if strings.Count(string(got), RedactionPlaceholder) != 2 {
		t.Errorf("Expected two placeholders:\n%s", got)
	}
}

// TestRedactPathsRedactsTLSCertificateData proves the committed
// get_tls_certificates spec actually reaches the certificate bytes: this is
// the one response element the whole family exists to protect, and it only
// appears in the wire document when IncludeCertificateData was requested.
func TestRedactPathsRedactsTLSCertificateData(t *testing.T) {
	const fakeDER = "MIIBSynthEXAMPLEnotarealcertificateFAKEbase64QQQQ=="
	doc := []byte(`<get_tls_certificates_response status="200" status_text="OK">` +
		`<tls_certificate id="a"><name>fixture</name>` +
		`<certificate format="DER">` + fakeDER + `</certificate>` +
		`</tls_certificate></get_tls_certificates_response>`)

	got, applied, err := RedactPaths(doc,
		redactionSpecs["get_tls_certificates"])
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	if strings.Contains(string(got), fakeDER) {
		t.Errorf("The certificate bytes survived redaction:\n%s", got)
	}
	if !strings.Contains(string(got), RedactionPlaceholder) {
		t.Errorf("No placeholder was written:\n%s", got)
	}
	if !strings.Contains(string(got), `format="DER"`) {
		t.Errorf("The sibling format attribute should be untouched:\n%s", got)
	}
	if len(applied) != 1 || applied[0] != "/get_tls_certificates_response/tls_certificate/certificate/#text" {
		t.Errorf("applied = %v, want the one spec path", applied)
	}
}

// TestRedactPathsDistinguishesDepth guards against a prefix or suffix match
// standing in for a real path comparison: only the fully qualified path is a
// hit, so a <password> somewhere else in the tree is left alone.
func TestRedactPathsDistinguishesDepth(t *testing.T) {
	doc := []byte(`<a><password>keep</password><b><password>drop</password></b></a>`)

	got, _, err := RedactPaths(doc, []string{"/a/b/password/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	if !strings.Contains(string(got), ">keep<") {
		t.Errorf("The shallower password should have survived:\n%s", got)
	}
	if strings.Contains(string(got), ">drop<") {
		t.Errorf("The specified password should have gone:\n%s", got)
	}
}

// TestRedactPathsMergesAdjacentTextTokens covers CDATA, which the decoder
// returns as its own CharData token. Without span merging the output would
// carry two placeholders where the document had one value.
func TestRedactPathsMergesAdjacentTextTokens(t *testing.T) {
	doc := []byte(`<a><private><![CDATA[--BEGIN--]]>tail</private></a>`)

	got, _, err := RedactPaths(doc, []string{"/a/private/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	if strings.Contains(string(got), "BEGIN") || strings.Contains(string(got), "tail") {
		t.Errorf("Part of the value survived:\n%s", got)
	}
	if n := strings.Count(string(got), RedactionPlaceholder); n != 1 {
		t.Errorf("Expected one merged placeholder, got %d:\n%s", n, got)
	}
}

// TestRedactPathsCoversEscapedText proves the raw byte span is replaced rather
// than the decoded value: an entity-escaped secret must not be left behind in
// its escaped form.
func TestRedactPathsCoversEscapedText(t *testing.T) {
	doc := []byte(`<a><password>zzz&amp;qqq&#x21;jjj</password></a>`)

	got, _, err := RedactPaths(doc, []string{"/a/password/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	for _, fragment := range []string{"zzz", "&amp;", "&#x21;", "qqq", "jjj"} {
		if strings.Contains(string(got), fragment) {
			t.Errorf("The escaped fragment %q survived:\n%s", fragment, got)
		}
	}
}

func TestRedactPathsWithNoMatchReturnsTheInput(t *testing.T) {
	doc := []byte(`<get_tasks_response status="200"><task id="1"/></get_tasks_response>`)

	got, applied, err := RedactPaths(doc, []string{"/get_tasks_response/task/password/#text"})
	if err != nil {
		t.Fatalf("RedactPaths returned an error: %v", err)
	}
	if string(got) != string(doc) {
		t.Errorf("A non-matching spec changed the document:\n%s", got)
	}
	if applied != nil {
		t.Errorf("applied = %v, want nil so the caller can trust byte_faithful", applied)
	}
}

func TestRedactPathsRejectsAMalformedDocument(t *testing.T) {
	_, _, err := RedactPaths([]byte(`<a><b></a>`), []string{"/a/b/#text"})
	if err == nil {
		t.Fatal("Expected an error for a document that is not well-formed")
	}
}

func TestValidateRedactionPath(t *testing.T) {
	valid := []string{
		"/a/#text",
		"/a/b/c/#text",
		"/a/@id",
		"/{urn:x}a/@{urn:x}id",
	}
	for _, p := range valid {
		if err := ValidateRedactionPath(p); err != nil {
			t.Errorf("ValidateRedactionPath(%q) = %v, want nil", p, err)
		}
	}

	invalid := map[string]string{
		"a/#text":     "missing leading slash",
		"/a/":         "trailing slash",
		"/a//b/#text": "empty segment",
		"/a/b/#text ": "whitespace",
		"/a/b":        "names an element, which has no value of its own",
		"/a/@":        "attribute with no name",
	}
	for p, why := range invalid {
		if err := ValidateRedactionPath(p); err == nil {
			t.Errorf("ValidateRedactionPath(%q) = nil, want an error (%s)", p, why)
		}
	}
}

// TestRedactorRefusesASecretBearingCommandWithNoSpec drives the deny-by-default
// rule with injected maps. It cannot be driven through NewRedactor, because
// TestRedactionSpecsAreConsistent asserts the real maps never reach this state
// — which is the whole design: the runtime check exists for the window between
// someone adding a command to the set and someone running the tests.
func TestRedactorRefusesASecretBearingCommandWithNoSpec(t *testing.T) {
	r := newRedactor(map[string][]string{}, map[string]bool{"create_credential": true})

	if _, err := r.PathsFor("create_credential"); err == nil {
		t.Fatal("PathsFor returned nil for a secret-bearing command with no spec")
	}
	if _, _, err := r.Redact("create_credential", []byte(`<create_credential/>`)); err == nil {
		t.Fatal("Redact wrote through a secret-bearing command with no spec")
	}
	if _, err := r.PathsFor("get_version"); err != nil {
		t.Fatalf("PathsFor(get_version) = %v, want nil: an unlisted command has nothing to redact", err)
	}
}

// TestRedactionSpecsAreConsistent keeps the two maps and the registry from
// drifting apart. Each clause is a distinct way the corpus could leak: a spec
// for a command nobody marked secret-bearing means the set is stale; a
// secret-bearing command with no spec means the tool will refuse to capture it;
// a path rooted at the wrong element silently matches nothing.
func TestRedactionSpecsAreConsistent(t *testing.T) {
	for command, paths := range redactionSpecs {
		if !secretBearingCommands[command] {
			t.Errorf("%s has a redaction spec but is not in secretBearingCommands", command)
		}
		if _, ok := Registry[command]; !ok {
			t.Errorf("%s has a redaction spec but is not a known GMP command", command)
		}
		if len(paths) == 0 {
			t.Errorf("%s has an empty redaction spec; remove it or fill it in", command)
		}
		for _, p := range paths {
			if err := ValidateRedactionPath(p); err != nil {
				t.Errorf("%s: %v", command, err)
				continue
			}
			root := strings.Split(strings.TrimPrefix(p, "/"), "/")[0]
			if root != command && root != command+"_response" {
				t.Errorf("%s: path %q is rooted at %q, which is neither the request nor the "+
					"response element, so it can never match", command, p, root)
			}
		}
	}

	for command := range secretBearingCommands {
		if _, ok := Registry[command]; !ok {
			t.Errorf("%s is listed as secret-bearing but is not a known GMP command", command)
		}
		if _, ok := redactionSpecs[command]; !ok {
			t.Errorf("%s is listed as secret-bearing but has no redaction spec, so capture "+
				"will refuse to write its fixture", command)
		}
	}
}

func TestSweepSeparatesSecretsFromIdentifiers(t *testing.T) {
	r := NewRedactor()
	r.AddSecret("-gmp.pass", "hunter2")
	r.AddIdentifier("-gmp.user", "admin")

	fatal, notices := r.Sweep([]byte(`<owner><name>admin</name></owner>`))
	if len(fatal) != 0 {
		t.Errorf("fatal = %v, want none: a username in <owner> is not a leak", fatal)
	}
	if len(notices) != 1 || notices[0] != "-gmp.user" {
		t.Errorf("notices = %v, want the username reported", notices)
	}

	fatal, _ = r.Sweep([]byte(`<password>hunter2</password>`))
	if len(fatal) != 1 || fatal[0] != "-gmp.pass" {
		t.Errorf("fatal = %v, want the password flagged", fatal)
	}
	for _, label := range fatal {
		if strings.Contains(label, "hunter2") {
			t.Errorf("The sweep reported the secret itself (%q); labels only", label)
		}
	}
}

func TestSweepIgnoresEmptyAndDuplicateTerms(t *testing.T) {
	r := NewRedactor()
	r.AddSecret("empty", "")
	r.AddSecret("dup", "x")
	r.AddSecret("dup", "x")

	if got := len(r.Terms()); got != 1 {
		t.Errorf("Terms() has %d entries, want 1: an empty value would match every document", got)
	}
}

// TestCommittedSpecsRedactRealRequests marshals the actual commands.* types the
// capture tool sends and checks the committed specs against the bytes that come
// out. A spec written against a guessed element name would pass every test
// above and fail here.
func TestCommittedSpecsRedactRealRequests(t *testing.T) {
	r := NewRedactor()

	cases := []struct {
		command string
		req     any
		secrets []string
		keep    []string
	}{
		{
			command: "authenticate",
			req: &commands.Authenticate{
				Credentials: commands.AuthCredentials{Username: "admin", Password: "PW-AUTH"},
			},
			secrets: []string{"PW-AUTH"},
			keep:    []string{"admin"},
		},
		{
			command: "create_credential",
			req: &commands.CreateCredential{
				Name:      "c",
				Login:     "root",
				Password:  "PW-CRED",
				Community: "PW-COMMUNITY",
				Key: &commands.CredentialKey{
					Phrase:  "PW-PHRASE",
					Private: "PW-PRIVATE",
					Public:  "ssh-ed25519 PUBLIC",
				},
				Privacy: &commands.CredentialPrivacy{Algorithm: "aes", Password: "PW-PRIVACY"},
			},
			secrets: []string{"PW-CRED", "PW-COMMUNITY", "PW-PHRASE", "PW-PRIVATE", "PW-PRIVACY"},
			keep:    []string{"root", "ssh-ed25519 PUBLIC", "aes"},
		},
		{
			command: "modify_credential",
			req: &commands.ModifyCredential{
				CredentialID: "id",
				Password:     "PW-CRED",
				Community:    "PW-COMMUNITY",
				Key:          &commands.CredentialKey{Phrase: "PW-PHRASE", Private: "PW-PRIVATE"},
				Privacy:      &commands.CredentialPrivacy{Password: "PW-PRIVACY"},
			},
			secrets: []string{"PW-CRED", "PW-COMMUNITY", "PW-PHRASE", "PW-PRIVATE", "PW-PRIVACY"},
		},
		{
			command: "create_user",
			req:     &commands.CreateUser{Name: "u", Password: "PW-USER"},
			secrets: []string{"PW-USER"},
			keep:    []string{"<name>u</name>"},
		},
		{
			command: "modify_user",
			req: &commands.ModifyUser{
				UserID:   "id",
				Name:     "u",
				Password: &commands.UserPassword{Modify: commands.Bool(true), Value: "PW-USER"},
			},
			secrets: []string{"PW-USER"},
			keep:    []string{`modify="1"`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.command, func(t *testing.T) {
			doc, err := xml.Marshal(tc.req)
			if err != nil {
				t.Fatalf("Failed to marshal the request: %v", err)
			}
			for _, s := range tc.secrets {
				if !strings.Contains(string(doc), s) {
					t.Fatalf("The test value %q is not in the marshalled request, so this case "+
						"proves nothing:\n%s", s, doc)
				}
			}

			got, applied, err := r.Redact(tc.command, doc)
			if err != nil {
				t.Fatalf("Redact returned an error: %v", err)
			}
			for _, s := range tc.secrets {
				if strings.Contains(string(got), s) {
					t.Errorf("%s survived redaction:\n%s", s, got)
				}
			}
			for _, k := range tc.keep {
				if !strings.Contains(string(got), k) {
					t.Errorf("Non-secret content %q was rewritten:\n%s", k, got)
				}
			}
			if len(applied) == 0 {
				t.Errorf("No path was applied, so the spec matched nothing:\n%s", got)
			}
		})
	}
}

// TestSweepBacksUpTheSpecs is the end-to-end shape of the capture guard: redact
// with the committed spec, then sweep for the literal secret. Nothing may be
// written between those two steps.
func TestSweepBacksUpTheSpecs(t *testing.T) {
	r := NewRedactor()
	r.AddSecret("-gmp.pass", "PW-AUTH")

	doc, err := xml.Marshal(&commands.Authenticate{
		Credentials: commands.AuthCredentials{Username: "admin", Password: "PW-AUTH"},
	})
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if fatal, _ := r.Sweep(doc); len(fatal) == 0 {
		t.Fatal("The sweep missed an unredacted password, so it cannot be a backstop")
	}

	redacted, _, err := r.Redact("authenticate", doc)
	if err != nil {
		t.Fatalf("Redact returned an error: %v", err)
	}
	if fatal, _ := r.Sweep(redacted); len(fatal) != 0 {
		t.Errorf("The sweep still finds %v after redaction:\n%s", fatal, redacted)
	}
}
