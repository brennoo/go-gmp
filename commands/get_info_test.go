package commands

import (
	"encoding/xml"
	"strings"
	"testing"
)

// T-26: command_get_info.xml's <pattern> declares exactly one <info> shape
// (the CVE/CPE/NVT/etc. data item), but every live 200 response and both
// upstream fixtures also send a second, disjoint <info start="..." max="..."/>
// pagination wrapper sharing the same element name. The two are told apart by
// attribute signature (@id vs @start/@max), the same duplicate-wire collision
// codec get_filters' two <filters> definitions already use.

func TestGetInfoResponseItemsOnlyRoundTrips(t *testing.T) {
	in := `<get_info_response status="200" status_text="OK">` +
		`<info id="CVE-1999-0001"><name>CVE-1999-0001</name></info>` +
		`<info id="CVE-1999-0002"><name>CVE-1999-0002</name></info>` +
		`</get_info_response>`
	var resp GetInfoResponse
	if err := xml.Unmarshal([]byte(in), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(resp.Info) != 2 || resp.Info[0].ID != "CVE-1999-0001" || resp.Info[1].ID != "CVE-1999-0002" {
		t.Fatalf("Info = %+v, want 2 data items in order", resp.Info)
	}
	if resp.Info2 != nil {
		t.Fatalf("Info2 = %+v, want nil since no wrapper was sent", resp.Info2)
	}

	out, err := xml.Marshal(&resp)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := strings.Count(string(out), "<info "); got != 2 {
		t.Fatalf("Marshal = %s, want 2 <info> elements, got %d", out, got)
	}
}

func TestGetInfoResponseWrapperOnlyRoundTrips(t *testing.T) {
	in := `<get_info_response status="400" status_text="No type specified."><info start="1" max="1000"></info></get_info_response>`
	var resp GetInfoResponse
	if err := xml.Unmarshal([]byte(in), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(resp.Info) != 0 {
		t.Fatalf("Info = %+v, want no data items", resp.Info)
	}
	if resp.Info2 == nil || resp.Info2.Start != 1 || resp.Info2.Max != 1000 {
		t.Fatalf("Info2 = %+v, want Start=1 Max=1000", resp.Info2)
	}

	out, err := xml.Marshal(&resp)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(out), `<info start="1" max="1000">`) {
		t.Fatalf("Marshal = %s, want the wrapper preserved", out)
	}
}

// TestGetInfoResponseMixedRoundTrips is the live-fixture shape: N data items
// plus the pagination wrapper in the same response, byte-faithfully
// preserved at the path level — the exact case that used to DROP @start/@max
// and fabricate a 12th data item out of the wrapper (11/001_fidelity.txt's
// 11 -> 12 DUPLICATED finding).
func TestGetInfoResponseMixedRoundTrips(t *testing.T) {
	in := `<get_info_response status="200" status_text="OK">` +
		`<info id="CVE-1999-0001"><name>CVE-1999-0001</name></info>` +
		`<info id="CVE-1999-0002"><name>CVE-1999-0002</name></info>` +
		`<info start="1" max="1000"></info>` +
		`</get_info_response>`
	var resp GetInfoResponse
	if err := xml.Unmarshal([]byte(in), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(resp.Info) != 2 {
		t.Fatalf("Info = %+v, want 2 data items", resp.Info)
	}
	if resp.Info2 == nil || resp.Info2.Start != 1 || resp.Info2.Max != 1000 {
		t.Fatalf("Info2 = %+v, want Start=1 Max=1000", resp.Info2)
	}

	out, err := xml.Marshal(&resp)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := strings.Count(string(out), "<info"); got != 3 {
		t.Fatalf("Marshal = %s, want 3 <info...> elements (2 items + 1 wrapper), got %d", out, got)
	}
	if !strings.Contains(string(out), `<info start="1" max="1000">`) {
		t.Fatalf("Marshal = %s, want the wrapper preserved alongside the items", out)
	}
}

// TestGetInfoResponseAmbiguousInfoFailsClosed pins the discriminator's
// fail-closed contract: an <info> matching neither the item signature (@id)
// nor the wrapper signature (@start/@max) must error, not silently guess.
func TestGetInfoResponseAmbiguousInfoFailsClosed(t *testing.T) {
	in := `<get_info_response status="200" status_text="OK"><info></info></get_info_response>`
	var resp GetInfoResponse
	err := xml.Unmarshal([]byte(in), &resp)
	if err == nil {
		t.Fatalf("Unmarshal succeeded with resp=%+v, want an error for an <info> matching neither known variant", resp)
	}
	if !strings.Contains(err.Error(), "unrecognized") {
		t.Fatalf("err = %v, want an \"unrecognized\" classification error", err)
	}
}

// TestGetInfoResponseExtraWrapperFailsClosed pins the other half of the
// fail-closed contract: a second occurrence of the singular wrapper variant
// is a real defect (unlike a second data item, which is the normal repeated
// case), and must error rather than silently overwrite Info2.
func TestGetInfoResponseExtraWrapperFailsClosed(t *testing.T) {
	in := `<get_info_response status="200" status_text="OK">` +
		`<info start="1" max="1000"></info>` +
		`<info start="2" max="1000"></info>` +
		`</get_info_response>`
	var resp GetInfoResponse
	err := xml.Unmarshal([]byte(in), &resp)
	if err == nil {
		t.Fatalf("Unmarshal succeeded with resp=%+v, want an error for a second wrapper occurrence", resp)
	}
	if !strings.Contains(err.Error(), "unexpected extra") {
		t.Fatalf("err = %v, want an \"unexpected extra\" error", err)
	}
}
