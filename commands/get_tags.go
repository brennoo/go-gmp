package commands

import "encoding/xml"

// GetTags represents a get_tags command request.
type GetTags struct {
	XMLName xml.Name `xml:"get_tags"`
}

// GetTagsResponse represents a get_tags command response.
type GetTagsResponse struct {
	XMLName    xml.Name  `xml:"get_tags_response"`
	Status     string    `xml:"status,attr"`
	StatusText string    `xml:"status_text,attr"`
	Tags       []TagInfo `xml:"tag,omitempty"`
}

type TagInfo struct {
	ID               string        `xml:"id,attr"`
	Name             string        `xml:"name"`
	Comment          string        `xml:"comment,omitempty"`
	CreationTime     string        `xml:"creation_time,omitempty"`
	ModificationTime string        `xml:"modification_time,omitempty"`
	Writable         string        `xml:"writable,omitempty"`
	InUse            string        `xml:"in_use,omitempty"`
	Value            string        `xml:"value,omitempty"`
	Active           string        `xml:"active,omitempty"`
	Resources        []TagResource `xml:"resource,omitempty"`
}

type TagResource struct {
	ID    string `xml:"id,attr"`
	Type  string `xml:"type"`
	Name  string `xml:"name"`
	Trash string `xml:"trash"`
}
