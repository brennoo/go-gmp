package commands

import (
	"encoding/xml"
	"time"
)

// GetTags represents a get_tags command request.
type GetTags struct {
	XMLName xml.Name `xml:"get_tags"`
	TagID   string   `xml:"tag_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
}

// GetTagsResponse represents a get_tags command response.
type GetTagsResponse struct {
	XMLName    xml.Name `xml:"get_tags_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Tags       []Tag    `xml:"tag,omitempty"`
	Truncated  string   `xml:"truncated,omitempty"`
}

// Tag represents a <tag> element in the get_tags response.
type Tag struct {
	ID               string          `xml:"id,attr"`
	Name             string          `xml:"name"`
	Comment          string          `xml:"comment,omitempty"`
	CreationTime     time.Time       `xml:"creation_time,omitempty"`
	ModificationTime time.Time       `xml:"modification_time,omitempty"`
	Writable         bool            `xml:"writable,omitempty"`
	InUse            bool            `xml:"in_use,omitempty"`
	Value            string          `xml:"value,omitempty"`
	Active           bool            `xml:"active,omitempty"`
	Orphan           bool            `xml:"orphan,omitempty"`
	Owner            *TagOwner       `xml:"owner,omitempty"`
	Permissions      []TagPermission `xml:"permissions>permission,omitempty"`
	Resources        []TagResource   `xml:"resource,omitempty"`
}

// TagOwner represents the owner of a tag.
type TagOwner struct {
	Name string `xml:"name"`
}

// TagPermission represents a permission on a tag.
type TagPermission struct {
	Name string `xml:"name"`
}

// TagResource represents a <resource> element for a tag.
type TagResource struct {
	ID    string `xml:"id,attr"`
	Type  string `xml:"type"`
	Name  string `xml:"name"`
	Trash bool   `xml:"trash"`
}
