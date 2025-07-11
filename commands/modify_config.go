package commands

import "encoding/xml"

// ModifyConfig represents a modify_config command request.
type ModifyConfig struct {
	XMLName         xml.Name                     `xml:"modify_config"`
	ConfigID        string                       `xml:"config_id,attr,omitempty"`
	Name            string                       `xml:"name,omitempty"`
	Comment         string                       `xml:"comment,omitempty"`
	Preference      []ConfigPreference           `xml:"preference,omitempty"`
	FamilySelection *ModifyConfigFamilySelection `xml:"family_selection,omitempty"`
	NVTSelection    []ModifyConfigNVTSelection   `xml:"nvt_selection,omitempty"`
}

// ModifyConfigResponse represents a modify_config command response.
type ModifyConfigResponse struct {
	XMLName    xml.Name `xml:"modify_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

type ModifyConfigFamilySelection struct {
	Growing bool                 `xml:"growing,omitempty"`
	Family  []ModifyConfigFamily `xml:"family,omitempty"`
}

type ModifyConfigFamily struct {
	All     bool   `xml:"all,omitempty"`
	Growing bool   `xml:"growing,omitempty"`
	Name    string `xml:"name,omitempty"`
}

type ModifyConfigNVTSelection struct {
	Family string            `xml:"family,omitempty"`
	NVT    []ModifyConfigNVT `xml:"nvt,omitempty"`
}

type ModifyConfigNVT struct {
	OID string `xml:"oid,attr,omitempty"`
}
