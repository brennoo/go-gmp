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
	"encoding/xml"
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
	// captureComment goes on every created resource so anything left behind
	// by a died run is obvious in the UI as well as by name.
	captureComment = "Created by go-gmp fixture capture. Safe to delete."

	// readOnlyFilter is small enough to stay a fixture, large enough to
	// exercise <filters>, <sort> and the *_count envelope.
	readOnlyFilter = "rows=2 first=1"

	// wrongPassword is not derived from the real one, so the failure fixture
	// cannot accidentally contain a prefix of it.
	wrongPassword = "go-gmp-fixture-deliberately-wrong-password"

	// nonexistentUUID is well-formed and cannot exist: gvmd's own resources
	// are random v4 UUIDs.
	nonexistentUUID = "00000000-0000-0000-0000-000000000000"

	// rowsPerPageSetting is gvmd's fixed UUID for the "Rows Per Page" user
	// setting: per-user, cosmetic and reversible.
	rowsPerPageSetting = "5f5a8712-8017-11e1-8556-406186ea4fc5"
)

// captureState carries the identifiers one step creates and a later one needs.
type captureState struct {
	portListID     string
	credentialID   string
	certCredID     string
	targetID       string
	filterID       string
	tagID          string
	groupID        string
	roleID         string
	userID         string
	permissionID   string
	scheduleID     string
	alertID        string
	tlsCertID      string
	configID       string
	reportFormatID string
	reportConfigID string
	scannerID      string
	taskID         string
	reportID       string
	resultID       string
	nvtOID         string

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
// Each fixture gets its own connection, since gvmd closes the session after
// refusing a command pre-authentication, and reusing one would turn every
// fixture after the first into a transport error.
func (r *captureRun) capturePreAuth() {
	if conn, ok := r.dial("get_version", "pre-auth: get_version"); ok {
		defer conn.Close()
		r.sendOn(conn, "get_version", "pre-auth: get_version on a fresh unauthenticated connection",
			&commands.GetVersion{})

		// The one fixture whose request genuinely contains the harness
		// password — the reason the whole redactor exists.
		r.sendOn(conn, "authenticate",
			"pre-auth: authenticate with the harness credentials, password redacted",
			&commands.Authenticate{Credentials: commands.AuthCredentials{
				Username: *liveUser, Password: *livePass,
			}})
	}

	if conn, ok := r.dial("authenticate", "pre-auth: authenticate failure"); ok {
		defer conn.Close()
		r.sendOn(conn, "authenticate", "pre-auth: authenticate with a deliberately wrong password",
			&commands.Authenticate{Credentials: commands.AuthCredentials{
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
// at all. get_info with no type is captured separately, deliberately, as a
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

// captureGet sends one get_* twice: bare, and with a paging filter. The
// filtered variant is the only way to see the <filters>/<sort>/<*_count>
// envelope; a request type with no Filter field can't produce it, and that
// absence is recorded rather than skipped silently.
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

// getByID re-reads a single resource by its identifier. The field name is
// resolved by reflection so a request type missing its id attribute (get_notes
// and get_roles have no fields at all) becomes a recorded skip rather than a
// compile-time obstacle to capturing everything else.
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
// identifier. Teardown goes on t.Cleanup (LIFO), so registering in creation
// order gives reverse-dependency-order teardown for free.
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
	r.lifecycleReportFormat(st)
	r.lifecycleReportConfig(st)
	r.lifecycleScanner(st)
	r.lifecycleTask(st)
	r.lifecycleSetting()
	r.lifecycleAuth()
}

// makeKeypair generates the throwaway certificate and private key used by the
// client-certificate credential, the TLS certificate resource and the
// scanner. Generating rather than embedding is what makes the redaction test
// real: the private key is a genuine secret registered with the sweep, so a
// spec that fails to remove it from a request aborts the run.
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

	// The PEM body, not the whole armoured block: a sweep term containing
	// newlines would miss a re-wrapped copy.
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

	resp, _ := r.send("create_port_range", "lifecycle: create_port_range",
		&commands.CreatePortRange{
			PortList: commands.PortRangePortList{ID: st.portListID},
			Start:    "20", End: "22", Type: "TCP",
		})
	if id := responseID(resp); id != "" {
		r.send("delete_port_range", "lifecycle: delete_port_range",
			&commands.DeletePortRange{PortRangeID: id})
	}
}

func (r *captureRun) lifecycleCredentials(st *captureState) {
	// The password is generated here and registered with the sweep, so a
	// spec that misses it aborts the run.
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
				Comment:      ptr(captureComment + " (modified)"),
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
			Key:         &commands.CredentialKey{Private: st.keyPEM},
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
	// Mixed case and punctuation in case gvmd enforces a password policy.
	secret := "Gg1!" + hex.EncodeToString(b[:])
	r.red.AddSecret(label, secret)
	return secret
}

func (r *captureRun) lifecycleTarget(st *captureState) {
	req := &commands.CreateTarget{
		Name:       r.resourceName("target"),
		Comment:    captureComment,
		Hosts:      "127.0.0.1",
		AliveTests: &commands.TargetAliveTests{AliveTest: []string{"Consider Alive"}},
	}
	if st.portListID != "" {
		req.PortList = &commands.TargetPortList{ID: st.portListID}
	}
	if st.credentialID != "" {
		req.SSHCredential = &commands.TargetSSHCredential{ID: st.credentialID, Port: "22"}
	}

	st.targetID = r.created("create_target",
		"lifecycle: create_target using the created port list and credential", req,
		func(id string) {
			r.send("delete_target", "teardown: delete_target ultimate",
				&commands.DeleteTarget{TargetID: id, Ultimate: true})
		})
	if st.targetID == "" {
		return
	}
	r.getByID("get_targets", "TargetID", st.targetID,
		"lifecycle: get_targets for the target just created")
	modifiedComment := captureComment + " (modified)"
	r.send("modify_target", "lifecycle: modify_target comment",
		&commands.ModifyTarget{TargetID: st.targetID, Comment: &modifiedComment})
}

// lifecycleFilter also produces the restore fixture: the one resource cheap
// enough to delete non-ultimately and pull back out of the trashcan.
func (r *captureRun) lifecycleFilter(st *captureState) {
	st.filterID = r.created("create_filter", "lifecycle: create_filter",
		&commands.CreateFilter{
			Name:    r.resourceName("filter"),
			Comment: captureComment,
			Term:    "rows=2 first=1 sort=name",
			Type:    "result",
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

	// The rawTag fallback exercises gvmd's tolerance of empty <type>/<name>/
	// <trash> children inside <resource>, as recorded historically in
	// testdata/live/create_tag/001_request.xml; TagResource itself only
	// carries ID today.
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
		&commands.ModifyTag{TagID: st.tagID, Value: ptr("fixture-modified")})
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
				Modify: commands.Bool(true),
				Value:  r.generatedSecret("created user's password, modified"),
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
			Subject: commands.PermissionSubject{ID: st.userID, Type: "user"},
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
			// delete_tls_certificate is unmodelled in commands/ and Registry, so
			// it can't be a fixture, but it must still be sent to avoid leaking
			// the resource.
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
	source := firstID(resp, "config")
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

func (r *captureRun) lifecycleReportFormat(st *captureState) {
	resp, ok := r.send("get_report_formats", "lifecycle: get_report_formats to choose a format to copy",
		&commands.GetReportFormats{Filter: readOnlyFilter})
	if !ok {
		return
	}
	source := firstID(resp, "report_format")
	if source == "" {
		r.skip("create_report_format", "lifecycle: create_report_format as a copy",
			"get_report_formats returned no report format to copy; the feed import owner is probably unset "+
				"(see hack/gvmd-up.sh step 5)")
		return
	}

	st.reportFormatID = r.created("create_report_format",
		"lifecycle: create_report_format copying an existing format",
		&commands.CreateReportFormat{Name: r.resourceName("report-format"), Copy: source},
		func(id string) {
			r.send("delete_report_format", "teardown: delete_report_format ultimate",
				&commands.DeleteReportFormat{ReportFormatID: id, Ultimate: true})
		})
	if st.reportFormatID == "" {
		return
	}
	r.getByID("get_report_formats", "ReportFormatID", st.reportFormatID,
		"lifecycle: get_report_formats for the copied report format")
	r.send("modify_report_format", "lifecycle: modify_report_format name",
		&commands.ModifyReportFormat{
			ReportFormatID: st.reportFormatID,
			Name:           r.resourceName("report-format") + "-modified",
		})
}

// lifecycleReportConfig skips create_report_config when no configurable
// report format exists — see docs/review-2026-08.md's delete_X sibling batch
// section for why this stack never has one.
func (r *captureRun) lifecycleReportConfig(st *captureState) {
	resp, ok := r.send("get_report_formats", "lifecycle: get_report_formats to find a configurable format",
		&commands.GetReportFormats{})
	if !ok {
		return
	}
	formatID := firstConfigurableReportFormatID(resp)
	if formatID == "" {
		r.skip("create_report_config", "lifecycle: create_report_config",
			"no configurable report format exists on this stack (command_create_report_config.xml requires "+
				"one); the default feed install's built-in formats are all non-configurable")
		return
	}

	st.reportConfigID = r.created("create_report_config", "lifecycle: create_report_config",
		&commands.CreateReportConfig{
			Name:         r.resourceName("report-config"),
			Comment:      captureComment,
			ReportFormat: commands.ReportConfigReportFormat{ID: formatID},
		},
		func(id string) {
			r.send("delete_report_config", "teardown: delete_report_config ultimate",
				&commands.DeleteReportConfig{ReportConfigID: id, Ultimate: true})
		})
	if st.reportConfigID == "" {
		return
	}
	r.getByID("get_report_configs", "ReportConfigID", st.reportConfigID,
		"lifecycle: get_report_configs for the report config just created")
	r.send("modify_report_config", "lifecycle: modify_report_config comment",
		&commands.ModifyReportConfig{ReportConfigID: st.reportConfigID, Comment: captureComment + " (modified)"})
}

func (r *captureRun) lifecycleScanner(st *captureState) {
	if st.certCredID == "" || st.certPEM == "" {
		r.missingDep("create_scanner", "lifecycle: create_scanner", "the client-certificate credential")
		return
	}

	resp, _ := r.send("create_scanner", "lifecycle: create_scanner, ca_pub redacted",
		&commands.CreateScanner{
			Name:       r.resourceName("scanner"),
			Comment:    captureComment,
			Host:       "127.0.0.1",
			Port:       "9391",
			Type:       "2",
			CAPub:      st.certPEM,
			Credential: commands.ScannerCredential{ID: st.certCredID},
		})

	st.scannerID = r.acceptCreated("create_scanner", "lifecycle: create_scanner", resp, func(id string) {
		r.send("delete_scanner", "teardown: delete_scanner ultimate",
			&commands.DeleteScanner{ScannerID: id, Ultimate: true})
	})
	if st.scannerID == "" {
		return
	}
	r.getByID("get_scanners", "ScannerID", st.scannerID,
		"lifecycle: get_scanners for the scanner just created, ca_pub redacted")
	comment := captureComment + " (modified)"
	r.send("modify_scanner", "lifecycle: modify_scanner comment",
		&commands.ModifyScanner{
			ScannerID: st.scannerID,
			Host:      "127.0.0.1",
			Port:      "9391",
			Type:      "2",
			Comment:   &comment,
		})
}

// lifecycleTask creates a container task (no target), the only kind that can
// exist without a scanner and the only way to get non-empty get_reports/
// get_results/get_aggregates on a stack where OpenVAS has never run. Which
// spelling of "no target" gvmd wants is open, so both python-gvm's
// <target id="0"/> and an empty-string form are tried as fallbacks.
func (r *captureRun) lifecycleTask(st *captureState) {
	name := r.resourceName("container-task")
	rawWithZero := fmt.Sprintf(
		`<create_task><name>%s</name><comment>%s</comment><target id="0"/></create_task>`,
		name, captureComment)
	rawWithEmpty := fmt.Sprintf(
		`<create_task><name>%s</name><comment>%s</comment><target id=""/></create_task>`,
		name, captureComment)

	resp := r.sendWithFallbacks("create_task",
		"lifecycle: create_task, typed request (expected to fail: Scanner has no omitempty, "+
			"so it always sends an empty scanner id gvmd now rejects)",
		&commands.CreateTask{
			Name:    name,
			Comment: captureComment,
			Target:  commands.TaskTarget{ID: "0"},
		},
		[]rawAttempt{
			{desc: `lifecycle: create_task container form with <target id="0"/>`, xml: rawWithZero},
			{desc: `lifecycle: create_task container form with <target id=""/>`, xml: rawWithEmpty},
		})

	st.taskID = r.acceptCreated("create_task", "lifecycle: create_task", resp, func(id string) {
		r.send("delete_task", "teardown: delete_task ultimate",
			&commands.DeleteTask{TaskID: id, Ultimate: true})
	})
	if st.taskID == "" {
		return
	}
	r.getByID("get_tasks", "TaskID", st.taskID, "lifecycle: get_tasks for the container task")
	comment := captureComment + " (modified)"
	r.send("modify_task", "lifecycle: modify_task comment",
		&commands.ModifyTask{TaskID: st.taskID, Comment: &comment})
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

// lifecycleAuth captures modify_auth without touching gvmd's auth config: it
// writes method:file's current auth_conf_setting values straight back, then
// re-reads and fails hard if anything changed. Byte-identical request/
// response XML only proves the library didn't lose data in transit, not that
// gvmd left the config alone, hence the second read.
func (r *captureRun) lifecycleAuth() {
	before, ok := r.describeAuthFileSettings("lifecycle: describe_auth before modify_auth")
	if !ok {
		return
	}
	settings := make([]commands.AuthGroupAuthConfSetting, len(before))
	for i, s := range before {
		settings[i] = commands.AuthGroupAuthConfSetting{Key: s.Key, Value: s.Value}
	}
	if _, ok := r.send("modify_auth", "lifecycle: modify_auth, writes method:file's own current settings back unchanged",
		&commands.ModifyAuth{Group: commands.ModifyAuthGroup{Name: "method:file", AuthConfSetting: settings}}); !ok {
		return
	}
	after, ok := r.describeAuthFileSettings("lifecycle: describe_auth after modify_auth, verifying the no-op")
	if !ok {
		return
	}
	if changed := diffAuthSettingKeys(before, after); len(changed) > 0 {
		// Values are security-sensitive and must never reach a test log; only
		// the differing keys are reported.
		r.t.Fatalf("modify_auth changed method:file's stored settings instead of writing them back unchanged "+
			"(keys with a different value or presence: %v) — gvmd's auth config may now be in a bad state", changed)
	}
}

// describeAuthFileSettings calls describe_auth and returns method:file's
// auth_conf_setting list.
func (r *captureRun) describeAuthFileSettings(desc string) ([]commands.AuthConfSetting, bool) {
	resp, ok := r.send("describe_auth", desc, &commands.DescribeAuth{})
	if !ok {
		r.missingDep("modify_auth", "lifecycle: modify_auth", "describe_auth")
		return nil, false
	}
	var parsed commands.DescribeAuthResponse
	if err := xml.Unmarshal([]byte(resp), &parsed); err != nil {
		r.fail("modify_auth", "lifecycle: modify_auth", "parsing describe_auth response: "+err.Error())
		return nil, false
	}
	for i := range parsed.Group {
		if parsed.Group[i].Name == "method:file" {
			return parsed.Group[i].AuthConfSetting, true
		}
	}
	r.missingDep("modify_auth", "lifecycle: modify_auth", "method:file group in describe_auth")
	return nil, false
}

// diffAuthSettingKeys compares two method:file auth_conf_setting lists and
// returns the keys whose set of values differs. A plain map[string]string
// would collapse before=[enable=true, enable=false] and after=[enable=false]
// into an equal-looking single value each, hiding a real change; comparing
// the full multiset per key (order-normalized, duplicates kept) catches that.
func diffAuthSettingKeys(before, after []commands.AuthConfSetting) []string {
	beforeVals := authSettingMultimap(before)
	afterVals := authSettingMultimap(after)

	keys := make(map[string]struct{}, len(beforeVals)+len(afterVals))
	for k := range beforeVals {
		keys[k] = struct{}{}
	}
	for k := range afterVals {
		keys[k] = struct{}{}
	}

	var changed []string
	for k := range keys {
		if !equalValueMultisets(beforeVals[k], afterVals[k]) {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}

// authSettingMultimap groups setting values by key, preserving duplicates.
func authSettingMultimap(settings []commands.AuthConfSetting) map[string][]string {
	m := make(map[string][]string, len(settings))
	for _, s := range settings {
		m[s.Key] = append(m[s.Key], s.Value)
	}
	return m
}

// equalValueMultisets reports whether a and b hold the same values the same
// number of times each, ignoring order.
func equalValueMultisets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	as, bs := append([]string(nil), a...), append([]string(nil), b...)
	sort.Strings(as)
	sort.Strings(bs)
	for i := range as {
		if as[i] != bs[i] {
			return false
		}
	}
	return true
}

// --- 4. The container-task trick ---------------------------------------------

// captureContainerReport uploads a synthetic report into the container task,
// the only way to get non-empty get_reports/get_results/get_aggregates on a
// stack where no scan has ever run and no scanner exists.
func (r *captureRun) captureContainerReport(st *captureState) {
	if st.taskID == "" {
		r.missingDep("create_report", "container task: create_report", "the container task")
		return
	}

	// Without format_id on the outer <report>, gvmd never parses the nested
	// <results>. The built-in "XML" format's id is per-install, not the fixed
	// UUID GMP.xml.in's example uses, so it must be resolved by name.
	formatsResp, ok := r.send("get_report_formats",
		"container task: get_report_formats to find the native XML format",
		&commands.GetReportFormats{})
	if !ok {
		return
	}
	formatID := reportFormatIDByName(formatsResp, "XML")
	if formatID == "" {
		r.skip("create_report", "container task: create_report",
			`get_report_formats returned no format named "XML"; the feed import owner is probably unset `+
				"(see hack/gvmd-up.sh step 5)")
		return
	}

	now := time.Now().UTC()
	payload := fmt.Sprintf(`<create_report><report format_id="%s"><report>`+
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
		formatID, now.Add(-time.Minute).Format(time.RFC3339), fixtureNVTOID, now.Format(time.RFC3339), st.taskID)

	resp, ok := r.sendRaw("create_report", "container task: create_report uploading a synthetic report", payload)
	if !ok {
		return
	}
	st.nvtOID = fixtureNVTOID
	st.reportID = r.acceptCreated("create_report", "container task: create_report", resp, func(id string) {
		// This 404s intermittently: gvmd can discard a container task's
		// synthetic report on its own once in_assets=1 finishes importing it,
		// independent of delete_asset's ordering — see docs/review-2026-08.md.
		r.send("delete_asset", "teardown: delete_asset for the imported report's assets",
			&commands.DeleteAsset{ReportID: id})
		r.send("delete_report", "teardown: delete_report, occasionally already gone",
			&commands.DeleteReport{ReportID: id})
	})
	if st.reportID == "" {
		return
	}

	r.getByID("get_reports", "ReportID", st.reportID, "container task: get_reports for the imported report")
	r.send("get_results", "container task: get_results for the container task",
		&commands.GetResults{TaskID: st.taskID, Details: commands.Bool(true)})
	r.send("get_aggregates", "container task: get_aggregates over results",
		&commands.GetAggregates{Type: "result", GroupColumn: "severity"})
	r.send("get_aggregates", "container task: get_aggregates with sort/data_column/text_column",
		&commands.GetAggregates{
			Type:        "result",
			GroupColumn: "severity",
			Sort:        []commands.AggregateSort{{Field: "severity", Order: "descending"}},
			DataColumn2: []string{"severity"},
			TextColumn:  []string{"name"},
		})

	if results, ok := r.send("get_results", "container task: get_results to pick a result id",
		&commands.GetResults{TaskID: st.taskID}); ok {
		st.resultID = firstID(results, "result")
	}
}

// fixtureNVTOID is a real OID from the Greenbone family so the imported result
// looks like something gvmd would have produced. Nothing depends on the plugin
// existing; a container import is not validated against the NVT feed.
const fixtureNVTOID = "1.3.6.1.4.1.25623.1.0.90022"

// firstID won't do here: <report_format> also carries <owner><name>, and
// encoding/xml's "name" tag binds only the format's own direct child.
func reportFormatIDByName(doc, name string) string {
	var parsed struct {
		Formats []struct {
			ID   string `xml:"id,attr"`
			Name string `xml:"name"`
		} `xml:"report_format"`
	}
	if err := xml.Unmarshal([]byte(doc), &parsed); err != nil {
		return ""
	}
	for _, f := range parsed.Formats {
		if f.Name == name {
			return f.ID
		}
	}
	return ""
}

func firstConfigurableReportFormatID(doc string) string {
	var parsed struct {
		Formats []struct {
			ID           string           `xml:"id,attr"`
			Configurable commands.GMPBool `xml:"configurable"`
		} `xml:"report_format"`
	}
	if err := xml.Unmarshal([]byte(doc), &parsed); err != nil {
		return ""
	}
	for _, f := range parsed.Formats {
		if bool(f.Configurable) {
			return f.ID
		}
	}
	return ""
}

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
			Comment:    captureComment,
			Result:     commands.TicketResult{ID: st.resultID},
			AssignedTo: commands.TicketAssignedTo{User: commands.TicketUser{ID: st.userID}},
			OpenNote:   "go-gmp fixture ticket.",
		},
		func(id string) {
			r.send("delete_ticket", "teardown: delete_ticket ultimate",
				&commands.DeleteTicket{TicketID: id, Ultimate: true})
		})
	if id == "" {
		return
	}
	r.getByID("get_tickets", "TicketID", id, "annotations: get_tickets for the ticket just created")
	modifiedComment := captureComment + " (modified)"
	r.send("modify_ticket", "annotations: modify_ticket comment",
		&commands.ModifyTicket{TicketID: id, Comment: &modifiedComment})
}

// --- 6. Deliberate failures --------------------------------------------------

// captureDeliberateFailures records what gvmd says when it says no. These are
// the fixtures that settle T-7b (does gvmd ever answer under a root element
// other than <command>_response?); the unauthenticated case, filed under
// get_tasks, is the strongest candidate since gvmd refuses everything but
// get_version pre-auth and is believed to answer <gmp_response>.
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

	// Malformed input gets its own connection: it may get no answer at all,
	// and the connection is discarded rather than reused.
	if conn, ok := r.dial("get_tasks", "failure: malformed XML"); ok {
		defer conn.Close()
		r.exchange(conn, "get_tasks", "failure: malformed XML with a mismatched end tag",
			`<get_tasks><unclosed></get_tasks>`)
	}

	// An entirely unknown command is the cleanest probe for T-7b and can't be
	// a fixture, so its root element is reported in the summary instead.
	if conn, ok := r.dial("get_version", "failure: unknown command probe"); ok {
		defer conn.Close()
		r.unfiled(conn, "unknown top-level command <no_such_command/>", `<no_such_command/>`, "")
	}
}

// --- 7. Expected-to-fail commands --------------------------------------------

// captureExpectedFailures records commands that cannot succeed on this stack
// (Enterprise-only features, anything needing a scanner to have run): an
// error response is still a response, and covers a path nothing else does.
func (r *captureRun) captureExpectedFailures(st *captureState) {
	r.send("sync_config", "expected failure: sync_config with no feed sync configured",
		&commands.SyncConfig{})
	r.send("modify_license", "expected failure: modify_license without an Enterprise appliance",
		&commands.ModifyLicense{})
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
	id := firstID(resp, "report_format")
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
