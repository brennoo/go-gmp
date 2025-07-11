package commands

import "encoding/xml"

// GetTags represents a get_tags command request.
type GetTags struct {
	XMLName xml.Name `xml:"get_tags"`
}

// GetTagsResponse represents a get_tags command response.
type GetTagsResponse struct {
	XMLName    xml.Name `xml:"get_tags_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Tags       []Tag    `xml:"tag,omitempty"`
}

// Tag represents a <tag> element in the get_tags response.
type Tag struct {
	ID               string        `xml:"id,attr"`
	Name             string        `xml:"name"`
	Comment          string        `xml:"comment,omitempty"`
	CreationTime     string        `xml:"creation_time,omitempty"`     // ISO 8601 string
	ModificationTime string        `xml:"modification_time,omitempty"` // ISO 8601 string
	Writable         string        `xml:"writable,omitempty"`          // "0"/"1"
	InUse            string        `xml:"in_use,omitempty"`            // "0"/"1"
	Value            string        `xml:"value,omitempty"`
	Active           string        `xml:"active,omitempty"` // "0"/"1"
	Resources        []TagResource `xml:"resource,omitempty"`
}

// TagResource represents a <resource> element for a tag.
type TagResource struct {
	ID    string `xml:"id,attr"`
	Type  string `xml:"type"`
	Name  string `xml:"name"`
	Trash string `xml:"trash"` // "0"/"1"
}
