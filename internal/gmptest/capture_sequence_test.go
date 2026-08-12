//go:build integration

package gmptest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/brennoo/go-gmp/commands"
)

// The capture sequence. Ordered by dependency, and every step is allowed to
// fail: a step that cannot run records why and the ones after it record that
// their dependency is missing.

const (
	// captureComment goes on every created resource so anything left behind by
	// a run that died before its teardowns is obvious in the UI as well as by
	// name.
	captureComment = "Created by go-gmp fixture capture. Safe to delete."

	// readOnlyFilter is small enough that the response stays a fixture and
	// large enough that <filters>, <sort> and the *_count envelope are all
	// exercised. rows=1 would not distinguish a page from a single item.
	readOnlyFilter = "rows=2 first=1"

	// wrongPassword is deliberately not derived from the real one, so the
	// failure fixture cannot accidentally contain a prefix of it.
	wrongPassword = "go-gmp-fixture-deliberately-wrong-password"

	// nonexistentUUID is well-formed and cannot exist: gvmd's own resources are
	// random v4 UUIDs.
	nonexistentUUID = "00000000-0000-0000-0000-000000000000"

	// rowsPerPageSetting is gvmd's fixed UUID for the "Rows Per Page" user
	// setting. Chosen because it is per-user, cosmetic and reversible, which
	// rules out modifying anything the rest of the run depends on.
	rowsPerPageSetting = "5f5a8712-8017-11e1-8556-406186ea4fc5"
)

// captureState carries the identifiers one step creates and a later one needs.
type captureState struct {
	portListID   string
	credentialID string
	certCredID   string
	targetID     string
	filterID     string
	tagID        string
	groupID      string
	roleID       string
	userID       string
	permissionID string
	scheduleID   string
	alertID      string
	tlsCertID    string
	configID     string
	scannerID    string
	taskID       string
	reportID     string
	resultID     string
	nvtOID       string

	certPEM string
	keyPEM  string
}

func (r *captureRun) run() {
	st := &captureState{}

	r.capturePreAuth()
	r.captureReadOnly()
	r.captureLifecycle(st)
	r.captureContainerReport(st)
	r.captureAnnotations(st)
	r.captureDeliberateFailures()
	r.captureExpectedFailures(st)
}

// --- 1. Pre-auth -------------------------------------------------------------

// capturePreAuth records the two things that can only be seen before
// authenticating, plus the authenticate exchange itself in both directions.
//
// Each fixture gets its own connection. gvmd closes the session after refusing
// a command in the pre-authentication state, so reusing one would turn every
// fixture after the first into a transport error.
func (r *captureRun) capturePreAuth() {
	if conn, ok := r.dial("get_version", "pre-auth: get_version"); ok {
		defer conn.Close()
		r.sendOn(conn, "get_version", "pre-auth: get_version on a fresh unauthenticated connection",
			&commands.GetVersion{})

		// The one fixture whose request genuinely contains the harness
		// password. Its redaction spec is the reason the whole redactor exists.
		r.sendOn(conn, "authenticate",
			"pre-auth: authenticate with the harness credentials, password redacted",
			&commands.Authenticate{Credentials: &commands.AuthenticateCredentials{
				Username: *liveUser, Password: *livePass,
			}})
	}

	if conn, ok := r.dial("authenticate", "pre-auth: authenticate failure"); ok {
		defer conn.Close()
		r.sendOn(conn, "authenticate", "pre-auth: authenticate with a deliberately wrong password",
			&commands.Authenticate{Credentials: &commands.AuthenticateCredentials{
				Username: *liveUser, Password: wrongPassword,
			}})
	}

	if conn, ok := r.dial("help", "pre-auth: help before authenticate"); ok {
		defer conn.Close()
		r.sendOn(conn, "help", "pre-auth: help before authenticate, which gvmd should refuse",
			&commands.Help{})
	}
}

// --- 2. Authenticated read-only ---------------------------------------------

// readOnlyPresets supplies the attributes gvmd requires before it will answer
// at all. Without them the bare pass would capture four 400s instead of four
// responses; get_info with no type is captured deliberately later, as a
// missing-required-attribute fixture.
var readOnlyPresets = map[string]map[string]string{
	"get_aggregates":     {"Type": "task"},
	"get_assets":         {"Type": "host"},
	"get_info":           {"Type": "cve"},
	"get_resource_names": {"Type": "task"},
}

func (r *captureRun) captureReadOnly() {
	r.send("help", "read-only: help in gvmd's default format", &commands.Help{})
	r.send("help", `read-only: help format="XML", the whole GMP schema`, &commands.Help{Format: "XML"})
	r.send("describe_auth",
		"read-only: describe_auth, auth_conf_setting values redacted", &commands.DescribeAuth{})

	for _, name := range getCommands() {
		r.captureGet(name)
	}
}

// captureGet sends one get_* twice: bare, and with a paging filter.
//
// The filtered variant is the only way to see the <filters>/<sort>/<*_count>
// envelope, which is where the known role_count and agent_count gaps live. A
// request type with no Filter field cannot produce it, and that absence is
// recorded rather than skipped silently.
func (r *captureRun) captureGet(name string) {
	bare := Registry[name].NewReq()
	applyPreset(name, bare)
	r.send(name, "read-only: "+name+" bare", bare)

	filtered := Registry[name].NewReq()
	applyPreset(name, filtered)
	desc := fmt.Sprintf("read-only: %s with filter=%q", name, readOnlyFilter)
	if !setStringField(filtered, "Filter", readOnlyFilter) {
		r.skip(name, desc, "the request type has no settable string Filter field, so this command "+
			"cannot express pagination at all")
		return
	}
	r.send(name, desc, filtered)
}

func applyPreset(name string, req any) {
	for field, value := range readOnlyPresets[name] {
		setStringField(req, field, value)
	}
}

func getCommands() []string {
	var out []string
	for name := range Registry {
		if strings.HasPrefix(name, "get_") {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// getByID re-reads a single resource by its identifier.
//
// The field name is resolved by reflection so that a request type missing its
// id attribute — get_notes and get_roles have no fields whatsoever — becomes a
// recorded skip instead of a compile-time obstacle to capturing everything else.
func (r *captureRun) getByID(command, idField, id, desc string) {
	req := Registry[command].NewReq()
	if !setStringField(req, idField, id) {
		r.skip(command, desc, "the request type has no settable string "+idField+
			" field, so a single resource cannot be requested")
		return
	}
	r.send(command, desc, req)
}

// created runs a create_* command, registers its teardown and returns the new
// identifier.
//
// The teardown goes on t.Cleanup, which is LIFO, so registering as we create in
// dependency order gives teardown in reverse dependency order for free.
func (r *captureRun) created(command, desc string, req any, teardown func(id string)) string {
	resp, ok := r.send(command, desc, req)
	if !ok {
		return ""
	}
	return r.acceptCreated(command, desc, resp, teardown)
}

func (r *captureRun) acceptCreated(command, desc, resp string, teardown func(id string)) string {
	id := responseID(resp)
	if id == "" {
		_, status, statusText := describeResponse(resp)
		r.skip(command, desc, fmt.Sprintf(
			"no id in the response (status %s %q), so everything downstream of it is skipped",
			status, statusText))
		return ""
	}
	if teardown != nil {
		r.t.Cleanup(func() { teardown(id) })
	}
	return id
}

func (r *captureRun) missingDep(command, desc, dep string) {
	r.skip(command, desc, "skipped: "+dep+" was not created")
}

// --- 3. Create / modify / delete cycles --------------------------------------

func (r *captureRun) captureLifecycle(st *captureState) {
	r.makeKeypair(st)

	r.lifecyclePortList(st)
	r.lifecycleCredentials(st)
	r.lifecycleTarget(st)
	r.lifecycleFilter(st)
	r.lifecycleTag(st)
	r.lifecycleGroup(st)
	r.lifecycleRole(st)
	r.lifecycleUser(st)
	r.lifecyclePermission(st)
	r.lifecycleSchedule(st)
	r.lifecycleAlert(st)
	r.lifecycleTLSCertificate(st)
	r.lifecycleConfig(st)
	r.lifecycleScanner(st)
	r.lifecycleTask(st)
	r.lifecycleSetting()
}

// makeKeypair generates the throwaway certificate and private key used by the
// client-certificate credential, the TLS certificate resource and the scanner.
//
// Generating rather than embedding is what makes the redaction test real: the
// private key is a genuine secret that only exists inside this process, it is
// registered with the sweep, and if a spec fails to remove it from a request
// the run aborts. An embedded fixture key would be a secret nobody cares about
// and the sweep would be theatre.
func (r *captureRun) makeKeypair(st *captureState) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		r.fail("create_credential", "setup: generate a throwaway keypair", err.Error())
		return
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		r.fail("create_credential", "setup: generate a certificate serial", err.Error())
		return
	}
	tmpl := x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: r.resourceName("cert")},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		r.fail("create_credential", "setup: self-sign a throwaway certificate", err.Error())
		return
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		r.fail("create_credential", "setup: marshal the throwaway private key", err.Error())
		return
	}

	st.certPEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	st.keyPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))

	// The PEM body, not the whole armoured block: the BEGIN/END lines are not
	// secret and a sweep term containing newlines would miss a re-wrapped copy.
	r.red.AddSecret("generated private key", pemBody(st.keyPEM))
}

// pemBody returns the longest base64 line of a PEM block, which is the part
// worth sweeping for.
func pemBody(p string) string {
	longest := ""
	for _, line := range strings.Split(p, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "-----") && len(line) > len(longest) {
			longest = line
		}
	}
	return longest
}

func (r *captureRun) lifecyclePortList(st *captureState) {
	st.portListID = r.created("create_port_list", "lifecycle: create_port_list",
		&commands.CreatePortList{
			Name:      r.resourceName("port-list"),
			Comment:   captureComment,
			PortRange: "T:1-10,U:1-10",
		},
		func(id string) {
			r.send("delete_port_list", `teardown: delete_port_list ultimate`,
				&commands.DeletePortList{PortListID: id, Ultimate: true})
		})
	if st.portListID == "" {
		return
	}

	r.getByID("get_port_lists", "PortListID", st.portListID,
		"lifecycle: get_port_lists for the port list just created")
	r.send("modify_port_list", "lifecycle: modify_port_list comment",
		&commands.ModifyPortList{PortListID: st.portListID, Comment: captureComment + " (modified)"})

	// CreatePortRange models <port_list> as an element whose text is the id,
	// while GMP wants <port_list id="..."/>. The typed request is sent first so
	// the rejection is on the record, then the shape gvmd actually accepts.
	resp := r.sendWithFallbacks("create_port_range", "lifecycle: create_port_range",
		&commands.CreatePortRange{PortListID: st.portListID, Start: "20", End: "22", Type: "TCP"},
		[]rawAttempt{{
			desc: `lifecycle: create_port_range with <port_list id="…"/> instead of the struct's text form`,
			xml: fmt.Sprintf(`<create_port_range><port_list id="%s"/><start>20</start>`+
				`<end>22</end><type>TCP</type></create_port_range>`, st.portListID),
		}})
	if id := responseID(resp); id != "" {
		r.send("delete_port_range", "lifecycle: delete_port_range",
			&commands.DeletePortRange{PortRangeID: id})
	}
}

func (r *captureRun) lifecycleCredentials(st *captureState) {
	// A username/password credential. The password is generated here and
	// registered with the sweep, so a spec that misses it aborts the run.
	password := r.generatedSecret("credential password")
	st.credentialID = r.created("create_credential",
		"lifecycle: create_credential type=up, password redacted",
		&commands.CreateCredential{
			Name:     r.resourceName("credential-up"),
			Comment:  captureComment,
			Type:     "up",
			Login:    "go-gmp-fixture",
			Password: password,
		},
		func(id string) {
			r.send("delete_credential", "teardown: delete_credential ultimate",
				&commands.DeleteCredential{CredentialID: id, Ultimate: true})
		})

	if st.credentialID != "" {
		r.getByID("get_credentials", "CredentialID", st.credentialID,
			"lifecycle: get_credentials for the username/password credential")
		r.send("modify_credential", "lifecycle: modify_credential, new password redacted",
			&commands.ModifyCredential{
				CredentialID: st.credentialID,
				Comment:      captureComment + " (modified)",
				Password:     r.generatedSecret("credential password, modified"),
			})
	}

	if st.certPEM == "" {
		r.missingDep("create_credential",
			"lifecycle: create_credential type=cc", "the throwaway keypair")
		return
	}
	st.certCredID = r.created("create_credential",
		"lifecycle: create_credential type=cc with a generated certificate and private key",
		&commands.CreateCredential{
			Name:        r.resourceName("credential-cc"),
			Comment:     captureComment,
			Type:        "cc",
			Certificate: st.certPEM,
			Key:         &commands.CreateCredentialKey{Private: st.keyPEM},
		},
		func(id string) {
			r.send("delete_credential", "teardown: delete_credential ultimate",
				&commands.DeleteCredential{CredentialID: id, Ultimate: true})
		})
	if st.certCredID != "" {
		r.getByID("get_credentials", "CredentialID", st.certCredID,
			"lifecycle: get_credentials for the client-certificate credential")
	}
}

// generatedSecret returns a fresh random value and registers it with the sweep.
func (r *captureRun) generatedSecret(label string) string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		r.t.Fatalf("Failed to generate a secret for %s: %v", label, err)
	}
	// Mixed case and punctuation because gvmd may enforce a password policy,
	// and a rejected create_credential would cost the whole branch.
	secret := "Gg1!" + hex.EncodeToString(b[:])
	r.red.AddSecret(label, secret)
	return secret
}

func (r *captureRun) lifecycleTarget(st *captureState) {
	req := &commands.CreateTarget{
		Name:       r.resourceName("target"),
		Comment:    captureComment,
		Hosts:      "127.0.0.1",
		AliveTests: "Consider Alive",
	}
	if st.portListID != "" {
		req.PortList = &commands.TargetPortList{ID: st.portListID}
	}
	if st.credentialID != "" {
		req.SSHCredential = &commands.TargetCredential{ID: st.credentialID, Port: "22"}
	}

	st.targetID = r.created("create_target",
		"lifecycle: create_target using the created port list and credential", req,
		func(id string) {
			r.send("delete_target", "teardown: delete_target ultimate",
				&commands.DeleteTarget{TargetID: id, Ultimate: commands.Bool(true)})
		})
	if st.targetID == "" {
		return
	}
	r.getByID("get_targets", "TargetID", st.targetID,
		"lifecycle: get_targets for the target just created")
	r.send("modify_target", "lifecycle: modify_target comment",
		&commands.ModifyTarget{TargetID: st.targetID, Comment: captureComment + " (modified)"})
}

// lifecycleFilter also produces the restore fixture: it is the one resource
// cheap enough to delete non-ultimately and pull back out of the trashcan,
// which is the only way to capture <restore> against something real.
func (r *captureRun) lifecycleFilter(st *captureState) {
	st.filterID = r.created("create_filter", "lifecycle: create_filter",
		&commands.CreateFilter{
			Name:    r.resourceName("filter"),
			Comment: captureComment,
			Term:    "rows=2 first=1 sort=name",
			Type:    "task",
		},
		func(id string) {
			r.send("delete_filter", "teardown: delete_filter ultimate",
				&commands.DeleteFilter{FilterID: id, Ultimate: true})
		})
	if st.filterID == "" {
		return
	}

	r.getByID("get_filters", "FilterID", st.filterID,
		"lifecycle: get_filters for the filter just created")
	modifiedTerm := "rows=3 first=1 sort=name"
	r.send("modify_filter", "lifecycle: modify_filter term",
		&commands.ModifyFilter{FilterID: st.filterID, Term: &modifiedTerm})

	r.send("delete_filter", "lifecycle: delete_filter into the trashcan, not ultimate",
		&commands.DeleteFilter{FilterID: st.filterID, Ultimate: false})
	r.send("restore", "lifecycle: restore the filter out of the trashcan",
		&commands.Restore{ID: st.filterID})
}

func (r *captureRun) lifecycleTag(st *captureState) {
	if st.targetID == "" {
		r.missingDep("create_tag", "lifecycle: create_tag on the created target", "the target")
		return
	}

	rawTag := fmt.Sprintf(`<create_tag><name>%s</name><comment>%s</comment><value>fixture</value>`+
		`<resources><type>target</type><resource id="%s"/></resources></create_tag>`,
		r.resourceName("tag"), captureComment, st.targetID)

	// TagResource only carries ID today. It used to also have non-omitempty
	// <type>, <name> and <trash> children, and gvmd's tolerance of the
	// resulting empty elements inside <resource> is exactly what got recorded
	// verbatim in testdata/live/create_tag/001_request.xml — that capture
	// predates the narrower type and is historical client output, not a
	// requirement this call still exercises.
	resp := r.sendWithFallbacks("create_tag", "lifecycle: create_tag on the created target",
		&commands.CreateTag{
			Name:    r.resourceName("tag"),
			Comment: captureComment,
			Value:   "fixture",
			Resources: commands.TagResources{
				Type:     "target",
				Resource: []commands.TagResource{{ID: st.targetID}},
			},
		},
		[]rawAttempt{{desc: "lifecycle: create_tag with a bare <resource id=…/>", xml: rawTag}})

	st.tagID = r.acceptCreated("create_tag", "lifecycle: create_tag", resp, func(id string) {
		r.send("delete_tag", "teardown: delete_tag ultimate",
			&commands.DeleteTag{TagID: id, Ultimate: true})
	})
	if st.tagID == "" {
		return
	}
	r.getByID("get_tags", "TagID", st.tagID, "lifecycle: get_tags for the tag just created")
	r.send("modify_tag", "lifecycle: modify_tag value",
		&commands.ModifyTag{TagID: st.tagID, Value: "fixture-modified"})
}

func (r *captureRun) lifecycleGroup(st *captureState) {
	st.groupID = r.created("create_group", "lifecycle: create_group",
		&commands.CreateGroup{Name: r.resourceName("group"), Comment: captureComment},
		func(id string) {
			r.send("delete_group", "teardown: delete_group ultimate",
				&commands.DeleteGroup{GroupID: id, Ultimate: true})
		})
	if st.groupID == "" {
		return
	}
	r.getByID("get_groups", "GroupID", st.groupID, "lifecycle: get_groups for the group just created")
	r.send("modify_group", "lifecycle: modify_group comment",
		&commands.ModifyGroup{GroupID: st.groupID, Comment: captureComment + " (modified)"})
}

func (r *captureRun) lifecycleRole(st *captureState) {
	st.roleID = r.created("create_role", "lifecycle: create_role",
		&commands.CreateRole{Name: r.resourceName("role"), Comment: captureComment},
		func(id string) {
			r.send("delete_role", "teardown: delete_role ultimate",
				&commands.DeleteRole{RoleID: id, Ultimate: true})
		})
	if st.roleID == "" {
		return
	}
	r.getByID("get_roles", "RoleID", st.roleID, "lifecycle: get_roles for the role just created")
	r.send("modify_role", "lifecycle: modify_role comment",
		&commands.ModifyRole{RoleID: st.roleID, Comment: captureComment + " (modified)"})
}

func (r *captureRun) lifecycleUser(st *captureState) {
	st.userID = r.created("create_user", "lifecycle: create_user, password redacted",
		&commands.CreateUser{
			Name:     r.resourceName("user"),
			Comment:  captureComment,
			Password: r.generatedSecret("created user's password"),
		},
		func(id string) {
			r.send("delete_user", "teardown: delete_user", &commands.DeleteUser{UserID: id})
		})
	if st.userID == "" {
		return
	}
	r.getByID("get_users", "UserID", st.userID, "lifecycle: get_users for the user just created")
	r.send("modify_user", "lifecycle: modify_user, new password redacted",
		&commands.ModifyUser{
			UserID: st.userID,
			Name:   r.resourceName("user"),
			Password: &commands.UserPassword{
				Modify: "1",
				Text:   r.generatedSecret("created user's password, modified"),
			},
		})
}

func (r *captureRun) lifecyclePermission(st *captureState) {
	if st.userID == "" {
		r.missingDep("create_permission", "lifecycle: create_permission", "the user")
		return
	}
	st.permissionID = r.created("create_permission",
		"lifecycle: create_permission granting get_tasks to the created user",
		&commands.CreatePermission{
			Name:    "get_tasks",
			Comment: captureComment,
			Subject: &commands.PermissionSubject{ID: st.userID, Type: "user"},
		},
		func(id string) {
			r.send("delete_permission", "teardown: delete_permission ultimate",
				&commands.DeletePermission{PermissionID: id, Ultimate: true})
		})
	if st.permissionID == "" {
		return
	}
	r.getByID("get_permissions", "PermissionID", st.permissionID,
		"lifecycle: get_permissions for the permission just created")
	r.send("modify_permission", "lifecycle: modify_permission comment",
		&commands.ModifyPermission{PermissionID: st.permissionID, Comment: captureComment + " (modified)"})
}

func (r *captureRun) lifecycleSchedule(st *captureState) {
	// A fixed, already-past one-shot event: nothing this schedule is attached
	// to can ever be triggered by it, which matters because the run also
	// creates a task.
	ical := "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//go-gmp//fixture//EN\r\n" +
		"BEGIN:VEVENT\r\nDTSTAMP:20200101T000000Z\r\nDTSTART:20200101T000000Z\r\n" +
		"DURATION:PT1H\r\nUID:go-gmp-fixture@example.invalid\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"

	st.scheduleID = r.created("create_schedule", "lifecycle: create_schedule with a past one-shot event",
		&commands.CreateSchedule{
			Name:      r.resourceName("schedule"),
			Comment:   captureComment,
			ICalendar: ical,
			Timezone:  "UTC",
		},
		func(id string) {
			r.send("delete_schedule", "teardown: delete_schedule ultimate",
				&commands.DeleteSchedule{ScheduleID: id, Ultimate: true})
		})
	if st.scheduleID == "" {
		return
	}
	r.getByID("get_schedules", "ScheduleID", st.scheduleID,
		"lifecycle: get_schedules for the schedule just created")
	r.send("modify_schedule", "lifecycle: modify_schedule comment",
		&commands.ModifySchedule{ScheduleID: st.scheduleID, Comment: captureComment + " (modified)"})
}

func (r *captureRun) lifecycleAlert(st *captureState) {
	// Syslog rather than Email: it needs no SMTP configuration and test_alert
	// on it later cannot send anything anywhere.
	st.alertID = r.created("create_alert", "lifecycle: create_alert with a Syslog method",
		&commands.CreateAlert{
			Name:      r.resourceName("alert"),
			Comment:   captureComment,
			Event:     commands.AlertEvent{Value: "Task run status changed", Data: []commands.AlertData{{Name: "status", Value: "Done"}}},
			Condition: commands.AlertCondition{Value: "Always"},
			Method:    commands.AlertMethod{Value: "Syslog"},
			Filter:    commands.AlertFilter{ID: st.filterID},
		},
		func(id string) {
			r.send("delete_alert", "teardown: delete_alert ultimate",
				&commands.DeleteAlert{AlertID: id, Ultimate: true})
		})
	if st.alertID == "" {
		return
	}
	r.getByID("get_alerts", "AlertID", st.alertID, "lifecycle: get_alerts for the alert just created")
	modifiedComment := captureComment + " (modified)"
	r.send("modify_alert", "lifecycle: modify_alert comment",
		&commands.ModifyAlert{AlertID: st.alertID, Comment: &modifiedComment})
}

func (r *captureRun) lifecycleTLSCertificate(st *captureState) {
	if st.certPEM == "" {
		r.missingDep("create_tls_certificate", "lifecycle: create_tls_certificate", "the throwaway keypair")
		return
	}
	st.tlsCertID = r.created("create_tls_certificate",
		"lifecycle: create_tls_certificate from the generated certificate, redacted",
		&commands.CreateTLSCertificate{
			Name:        r.resourceName("tls-certificate"),
			Comment:     captureComment,
			Certificate: base64.StdEncoding.EncodeToString([]byte(st.certPEM)),
		},
		func(id string) {
			// GMP has delete_tls_certificate; commands/ does not model it, and
			// neither does Registry, so it cannot be a fixture — there is no
			// directory to file it under. It is still sent, unfiled, because a
			// teardown that silently skipped would leak the resource.
			r.unfiled(r.conn, "teardown: delete_tls_certificate (unmodelled command, not filed)",
				fmt.Sprintf(`<delete_tls_certificate tls_certificate_id="%s" ultimate="1"/>`, id),
				"delete_tls_certificate_response")
		})
	if st.tlsCertID == "" {
		return
	}
	r.getByID("get_tls_certificates", "TLSCertificateID", st.tlsCertID,
		"lifecycle: get_tls_certificates for the certificate just created")
	r.send("get_tls_certificates", "lifecycle: get_tls_certificates with include_certificate_data, certificate redacted",
		&commands.GetTLSCertificates{TLSCertificateID: st.tlsCertID, IncludeCertificateData: commands.Bool(true)})
	modifiedComment := captureComment + " (modified)"
	r.send("modify_tls_certificate", "lifecycle: modify_tls_certificate comment",
		&commands.ModifyTLSCertificate{TLSCertificateID: st.tlsCertID, Comment: &modifiedComment})
}

func (r *captureRun) lifecycleConfig(st *captureState) {
	// Copy an existing scan config rather than build one: a config created from
	// nothing has no NVT selection and gvmd rejects attaching it to a task.
	resp, ok := r.send("get_configs", "lifecycle: get_configs to choose a config to copy",
		&commands.GetConfigs{Filter: readOnlyFilter})
	if !ok {
		return
	}
	source := firstAttr(resp, "config", "id")
	if source == "" {
		r.skip("create_config", "lifecycle: create_config as a copy",
			"get_configs returned no config to copy; the feed import owner is probably unset "+
				"(see hack/gvmd-up.sh step 5)")
		return
	}

	st.configID = r.created("create_config", "lifecycle: create_config copying an existing config",
		&commands.CreateConfig{Name: r.resourceName("config"), Comment: captureComment, Copy: source},
		func(id string) {
			r.send("delete_config", "teardown: delete_config ultimate",
				&commands.DeleteConfig{ConfigID: id, Ultimate: true})
		})
	if st.configID == "" {
		return
	}
	r.getByID("get_configs", "ConfigID", st.configID, "lifecycle: get_configs for the copied config")
	r.send("modify_config", "lifecycle: modify_config comment",
		&commands.ModifyConfig{ConfigID: st.configID, Comment: captureComment + " (modified)"})
}

func (r *captureRun) lifecycleScanner(st *captureState) {
	if st.certCredID == "" || st.certPEM == "" {
		r.missingDep("create_scanner", "lifecycle: create_scanner", "the client-certificate credential")
		return
	}

	rawScanner := fmt.Sprintf(`<create_scanner><name>%s</name><comment>%s</comment>`+
		`<host>127.0.0.1</host><port>9391</port><type>2</type><ca_pub>%s</ca_pub>`+
		`<credential id="%s"/></create_scanner>`,
		r.resourceName("scanner"), captureComment, xmlEscape(st.certPEM), st.certCredID)

	// ScannerCredential's <name>, <login> and <trash> are not omitempty, so the
	// typed request cannot express the bare <credential id="…"/> GMP documents.
	resp := r.sendWithFallbacks("create_scanner", "lifecycle: create_scanner, ca_pub redacted",
		&commands.CreateScanner{
			Name:       r.resourceName("scanner"),
			Comment:    captureComment,
			Host:       "127.0.0.1",
			Port:       9391,
			Type:       "2",
			CAPub:      st.certPEM,
			Credential: &commands.ScannerCredential{ID: st.certCredID},
		},
		[]rawAttempt{{
			desc: `lifecycle: create_scanner with a bare <credential id="…"/>, ca_pub redacted`,
			xml:  rawScanner,
		}})

	st.scannerID = r.acceptCreated("create_scanner", "lifecycle: create_scanner", resp, func(id string) {
		r.send("delete_scanner", "teardown: delete_scanner ultimate",
			&commands.DeleteScanner{ScannerID: id, Ultimate: true})
	})
	if st.scannerID == "" {
		return
	}
	r.getByID("get_scanners", "ScannerID", st.scannerID,
		"lifecycle: get_scanners for the scanner just created, ca_pub redacted")
	r.send("modify_scanner", "lifecycle: modify_scanner comment",
		&commands.ModifyScanner{ScannerID: st.scannerID, Comment: captureComment + " (modified)"})
}

// lifecycleTask creates a container task: a task with no target, which is the
// only kind that can exist without a scanner and the only way to get non-empty
// get_reports, get_results and get_aggregates out of a stack where OpenVAS has
// never run.
//
// Which spelling of "no target" gvmd wants is the open question. python-gvm's
// create_container_task sends <target id="0"/>; the empty-string form is tried
// as a fallback. Both attempts are captured, so whichever one works, the
// other's rejection is documented.
func (r *captureRun) lifecycleTask(st *captureState) {
	name := r.resourceName("container-task")
	rawWithZero := fmt.Sprintf(
		`<create_task><name>%s</name><comment>%s</comment><target id="0"/></create_task>`,
		name, captureComment)
	rawWithEmpty := fmt.Sprintf(
		`<create_task><name>%s</name><comment>%s</comment><target id=""/></create_task>`,
		name, captureComment)

	resp := r.sendWithFallbacks("create_task", "lifecycle: create_task as a container task",
		&commands.CreateTask{
			Name:    name,
			Comment: captureComment,
			Target:  &commands.TaskTarget{ID: "0"},
		},
		[]rawAttempt{
			{desc: `lifecycle: create_task container form with <target id="0"/>`, xml: rawWithZero},
			{desc: `lifecycle: create_task container form with <target id=""/>`, xml: rawWithEmpty},
		})

	st.taskID = r.acceptCreated("create_task", "lifecycle: create_task", resp, func(id string) {
		r.send("delete_task", "teardown: delete_task ultimate",
			&commands.DeleteTask{TaskID: id, Ultimate: commands.Bool(true)})
	})
	if st.taskID == "" {
		return
	}
	r.getByID("get_tasks", "TaskID", st.taskID, "lifecycle: get_tasks for the container task")
	r.send("modify_task", "lifecycle: modify_task comment",
		&commands.ModifyTask{TaskID: st.taskID, Comment: captureComment + " (modified)"})
}

// lifecycleSetting modifies one cosmetic per-user setting. Rows Per Page is
// chosen because nothing else in the run reads it, and gvmd wants the value
// base64-encoded.
func (r *captureRun) lifecycleSetting() {
	r.getByID("get_settings", "SettingID", rowsPerPageSetting,
		"lifecycle: get_settings for the Rows Per Page setting")

	value := base64.StdEncoding.EncodeToString([]byte("11"))
	raw := fmt.Sprintf(`<modify_setting setting_id="%s"><value>%s</value></modify_setting>`,
		rowsPerPageSetting, value)

	// ModifySetting's <name> is not omitempty, so the typed request sends an
	// empty one alongside setting_id.
	r.sendWithFallbacks("modify_setting", "lifecycle: modify_setting Rows Per Page",
		&commands.ModifySetting{SettingID: rowsPerPageSetting, Value: value},
		[]rawAttempt{{desc: "lifecycle: modify_setting with no empty <name> element", xml: raw}})
}

// --- 4. The container-task trick ---------------------------------------------

// captureContainerReport uploads a synthetic report into the container task.
//
// This is the unlock for the whole reporting half of the protocol: with a
// report in place, get_reports, get_results and get_aggregates all return
// content on a stack where no scan has ever run and no scanner exists. The
// payload shape follows testdata/upstream/create_report/001_request.xml, which
// nests <report> inside <report>.
func (r *captureRun) captureContainerReport(st *captureState) {
	if st.taskID == "" {
		r.missingDep("create_report", "container task: create_report", "the container task")
		return
	}

	now := time.Now().UTC()
	payload := fmt.Sprintf(`<create_report><report><report>`+
		`<scan_start>%s</scan_start>`+
		`<results start="1" max="-1"><result>`+
		`<name>go-gmp fixture result</name>`+
		`<description>Synthetic result uploaded by the go-gmp fixture capture harness.</description>`+
		`<host>127.0.0.1<hostname>localhost</hostname></host>`+
		`<port>general/tcp</port>`+
		`<nvt oid="%s"><type>nvt</type><name>go-gmp fixture NVT</name>`+
		`<family>General</family><cvss_base>5.0</cvss_base></nvt>`+
		`<threat>Medium</threat><severity>5.0</severity>`+
		`<qod><value>80</value><type>remote_banner</type></qod>`+
		`</result></results>`+
		`<scan_end>%s</scan_end>`+
		`</report></report><task id="%s"/><in_assets>1</in_assets></create_report>`,
		now.Add(-time.Minute).Format(time.RFC3339), fixtureNVTOID, now.Format(time.RFC3339), st.taskID)

	resp, ok := r.sendRaw("create_report", "container task: create_report uploading a synthetic report", payload)
	if !ok {
		return
	}
	st.nvtOID = fixtureNVTOID
	st.reportID = r.acceptCreated("create_report", "container task: create_report", resp, func(id string) {
		// in_assets=1 means the import created host and OS assets. They are not
		// owned by the task, so deleting the task does not remove them.
		r.send("delete_asset", "teardown: delete_asset for the imported report's assets",
			&commands.DeleteAsset{ReportID: id})
		r.send("delete_report", "teardown: delete_report", &commands.DeleteReport{ReportID: id})
	})
	if st.reportID == "" {
		return
	}

	r.getByID("get_reports", "ReportID", st.reportID, "container task: get_reports for the imported report")
	r.send("get_results", "container task: get_results for the container task",
		&commands.GetResults{TaskID: st.taskID, Details: commands.Bool(true)})
	r.send("get_aggregates", "container task: get_aggregates over results",
		&commands.GetAggregates{Type: "result", GroupColumn: "severity"})

	if results, ok := r.send("get_results", "container task: get_results to pick a result id",
		&commands.GetResults{TaskID: st.taskID}); ok {
		st.resultID = firstAttr(results, "result", "id")
	}
}

// fixtureNVTOID is a real OID from the Greenbone family so the imported result
// looks like something gvmd would have produced. Nothing depends on the plugin
// existing; a container import is not validated against the NVT feed.
const fixtureNVTOID = "1.3.6.1.4.1.25623.1.0.90022"

// --- 5. Annotations that need a result --------------------------------------

func (r *captureRun) captureAnnotations(st *captureState) {
	r.annotationNote(st)
	r.annotationOverride(st)
	r.annotationTicket(st)
}

func (r *captureRun) annotationNote(st *captureState) {
	if st.nvtOID == "" {
		r.missingDep("create_note", "annotations: create_note", "the imported report")
		return
	}
	id := r.created("create_note", "annotations: create_note on the imported result's NVT",
		&commands.CreateNote{
			Text: "go-gmp fixture note. " + captureComment,
			NVT:  commands.NoteNVT{OID: st.nvtOID},
		},
		func(id string) {
			r.send("delete_note", "teardown: delete_note ultimate",
				&commands.DeleteNote{NoteID: id, Ultimate: true})
		})
	if id == "" {
		return
	}
	r.getByID("get_notes", "NoteID", id, "annotations: get_notes for the note just created")
	r.send("modify_note", "annotations: modify_note text",
		&commands.ModifyNote{NoteID: id, Text: "go-gmp fixture note, modified."})
}

func (r *captureRun) annotationOverride(st *captureState) {
	if st.nvtOID == "" {
		r.missingDep("create_override", "annotations: create_override", "the imported report")
		return
	}
	id := r.created("create_override", "annotations: create_override on the imported result's NVT",
		&commands.CreateOverride{
			Text:        "go-gmp fixture override. " + captureComment,
			NVT:         commands.OverrideNVT{OID: st.nvtOID},
			NewSeverity: commands.Float64(2.0),
		},
		func(id string) {
			r.send("delete_override", "teardown: delete_override ultimate",
				&commands.DeleteOverride{OverrideID: id, Ultimate: true})
		})
	if id == "" {
		return
	}
	r.getByID("get_overrides", "OverrideID", id,
		"annotations: get_overrides for the override just created")
	r.send("modify_override", "annotations: modify_override text",
		&commands.ModifyOverride{OverrideID: id, Text: "go-gmp fixture override, modified."})
}

func (r *captureRun) annotationTicket(st *captureState) {
	if st.resultID == "" || st.userID == "" {
		r.missingDep("create_ticket", "annotations: create_ticket", "a result and the created user")
		return
	}
	id := r.created("create_ticket", "annotations: create_ticket against the imported result",
		&commands.CreateTicket{
			Comment:  captureComment,
			Result:   commands.TicketResult{ID: st.resultID},
			Assigned: commands.TicketAssignedTo{User: &commands.TicketUser{ID: st.userID}},
			OpenNote: "go-gmp fixture ticket.",
		},
		func(id string) {
			r.send("delete_ticket", "teardown: delete_ticket ultimate",
				&commands.DeleteTicket{TicketID: id, Ultimate: true})
		})
	if id == "" {
		return
	}
	r.getByID("get_tickets", "TicketID", id, "annotations: get_tickets for the ticket just created")
	r.send("modify_ticket", "annotations: modify_ticket comment",
		&commands.ModifyTicket{TicketID: id, Comment: captureComment + " (modified)"})
}

// --- 6. Deliberate failures --------------------------------------------------

// captureDeliberateFailures records what gvmd says when it says no.
//
// These are the fixtures that settle T-7b: whether gvmd ever answers under a
// root element other than <command>_response. The unauthenticated case is the
// strongest candidate — gvmd refuses everything but get_version before
// authentication, and the refusal is believed to come back as <gmp_response>,
// which is filed under get_tasks because that is the command that was sent.
func (r *captureRun) captureDeliberateFailures() {
	if conn, ok := r.dial("get_tasks", "failure: unauthenticated get_tasks"); ok {
		defer conn.Close()
		r.sendOn(conn, "get_tasks", "failure: get_tasks on an unauthenticated connection",
			&commands.GetTasks{})
	}

	r.send("get_tasks", "failure: get_tasks for a well-formed but nonexistent UUID",
		&commands.GetTasks{TaskID: nonexistentUUID})

	r.sendRaw("get_info", "failure: get_info with no type attribute, which GMP requires",
		`<get_info/>`)
	r.sendRaw("create_target", "failure: create_target with neither hosts nor asset_hosts",
		`<create_target><name>go-gmp-fixture-missing-hosts</name></create_target>`)

	// Malformed input gets its own connection and may get no answer at all,
	// which is itself a result: the exchange times out, the transport error is
	// recorded, and the connection is discarded rather than reused.
	if conn, ok := r.dial("get_tasks", "failure: malformed XML"); ok {
		defer conn.Close()
		r.exchange(conn, "get_tasks", "failure: malformed XML with a mismatched end tag",
			`<get_tasks><unclosed></get_tasks>`)
	}

	// An entirely unknown command is the cleanest probe for T-7b, and the only
	// one that cannot be a fixture: there is no registered command directory to
	// file it under. Its root element is reported in the summary instead.
	if conn, ok := r.dial("get_version", "failure: unknown command probe"); ok {
		defer conn.Close()
		r.unfiled(conn, "unknown top-level command <no_such_command/>", `<no_such_command/>`, "")
	}
}

// --- 7. Expected-to-fail commands --------------------------------------------

// captureExpectedFailures records the commands that cannot succeed on this
// stack: Enterprise-only features, and anything that needs a scanner to have
// actually run.
//
// They are captured because an error response is still a response, and the
// error path is the part of these types nothing else in the corpus covers.
func (r *captureRun) captureExpectedFailures(st *captureState) {
	r.send("sync_config", "expected failure: sync_config with no feed sync configured",
		&commands.SyncConfig{})
	r.send("modify_license", "expected failure: modify_license without an Enterprise appliance",
		&commands.ModifyLicense{})
	r.send("modify_agents", "expected failure: modify_agents with no agents", &commands.ModifyAgents{})
	r.send("delete_agents", "expected failure: delete_agents with no agents", &commands.DeleteAgents{})
	r.sendRaw("run_wizard", "expected failure: run_wizard for a wizard that does not exist",
		`<run_wizard><name>go-gmp-fixture-no-such-wizard</name></run_wizard>`)

	r.expectedTaskFailures(st)

	if st.alertID != "" {
		r.send("test_alert", "expected failure: test_alert on the Syslog alert",
			&commands.TestAlert{AlertID: st.alertID})
	} else {
		r.missingDep("test_alert", "expected failure: test_alert", "the alert")
	}

	if st.scannerID != "" {
		r.send("verify_scanner", "expected failure: verify_scanner against a scanner with nothing behind it",
			&commands.VerifyScanner{ScannerID: st.scannerID})
	} else {
		r.send("verify_scanner", "expected failure: verify_scanner for a nonexistent scanner",
			&commands.VerifyScanner{ScannerID: nonexistentUUID})
	}

	r.verifyAReportFormat()
}

func (r *captureRun) expectedTaskFailures(st *captureState) {
	if st.taskID == "" {
		r.missingDep("start_task", "expected failure: task commands on a container task", "the container task")
		return
	}
	r.send("start_task", "expected failure: start_task on a container task, which has no target",
		&commands.StartTask{TaskID: st.taskID})
	r.send("stop_task", "expected failure: stop_task on a task that is not running",
		&commands.StopTask{TaskID: st.taskID})
	r.send("resume_task", "expected failure: resume_task on a task that never ran",
		&commands.ResumeTask{TaskID: st.taskID})
	r.send("move_task", "expected failure: move_task to a nonexistent slave",
		&commands.MoveTask{TaskID: st.taskID, SlaveID: nonexistentUUID})
}

func (r *captureRun) verifyAReportFormat() {
	resp, ok := r.send("get_report_formats", "expected failure: get_report_formats to pick one to verify",
		&commands.GetReportFormats{})
	if !ok {
		return
	}
	id := firstAttr(resp, "report_format", "id")
	if id == "" {
		r.skip("verify_report_format", "expected failure: verify_report_format",
			"get_report_formats returned none; the feed import owner is probably unset")
		return
	}
	r.send("verify_report_format", "verify_report_format against a feed report format",
		&commands.VerifyReportFormat{ReportFormatID: id})
}

// captureEmptyTrashcan runs after every teardown, because it destroys whatever
// the non-ultimate delete in lifecycleFilter and any failed ultimate delete
// left behind. Registered in TestCaptureLive, not called from run().
func (r *captureRun) captureEmptyTrashcan() {
	r.send("empty_trashcan", "teardown: empty_trashcan after every delete", &commands.EmptyTrashcan{})
}

func xmlEscape(s string) string {
	rep := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return rep.Replace(s)
}
