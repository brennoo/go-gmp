package commands

import "encoding/xml"

// GetFeedsCommand represents a get_feeds command request.
type GetFeedsCommand struct {
	XMLName xml.Name `xml:"get_feeds"`
	Type    string   `xml:"type,attr,omitempty"`
}

// GetFeedsResponse represents a get_feeds command response.
type GetFeedsResponse struct {
	XMLName             xml.Name `xml:"get_feeds_response"`
	Status              string   `xml:"status,attr"`
	StatusText          string   `xml:"status_text,attr"`
	FeedOwnerSet        string   `xml:"feed_owner_set,omitempty"`
	FeedRolesSet        string   `xml:"feed_roles_set,omitempty"`
	FeedResourcesAccess string   `xml:"feed_resources_access,omitempty"`
	Feeds               []Feed   `xml:"feed,omitempty"`
}

// Feed represents a feed entry in the get_feeds response.
type Feed struct {
	Type             string            `xml:"type,omitempty"`
	Name             string            `xml:"name,omitempty"`
	Version          string            `xml:"version,omitempty"`
	Description      string            `xml:"description,omitempty"`
	SyncNotAvailable *SyncNotAvailable `xml:"sync_not_available,omitempty"`
	CurrentlySyncing *CurrentlySyncing `xml:"currently_syncing,omitempty"`
}

type SyncNotAvailable struct {
	Error string `xml:"error,omitempty"`
}

type CurrentlySyncing struct {
	Timestamp string `xml:"timestamp,omitempty"`
}
