package commands

import "encoding/xml"

// GetFeeds represents a get_feeds command request.
type GetFeeds struct {
	XMLName xml.Name `xml:"get_feeds"`
	Type    string   `xml:"type,attr,omitempty"`
}

// GetFeedsResponse represents a get_feeds command response.
type GetFeedsResponse struct {
	XMLName             xml.Name `xml:"get_feeds_response"`
	Status              string   `xml:"status,attr"`
	StatusText          string   `xml:"status_text,attr"`
	FeedOwnerSet        bool     `xml:"feed_owner_set,omitempty"`
	FeedRolesSet        bool     `xml:"feed_roles_set,omitempty"`
	FeedResourcesAccess bool     `xml:"feed_resources_access,omitempty"`
	Feeds               []Feed   `xml:"feed,omitempty"`
}

// Feed represents a <feed> element in the get_feeds response.
type Feed struct {
	Type             string                `xml:"type,omitempty"`
	Name             string                `xml:"name,omitempty"`
	Version          string                `xml:"version,omitempty"`
	Description      string                `xml:"description,omitempty"`
	SyncNotAvailable *FeedSyncNotAvailable `xml:"sync_not_available,omitempty"`
	CurrentlySyncing *FeedCurrentlySyncing `xml:"currently_syncing,omitempty"`
}

// FeedSyncNotAvailable represents the <sync_not_available> element in a feed.
type FeedSyncNotAvailable struct {
	Error string `xml:"error"`
}

// FeedCurrentlySyncing represents the <currently_syncing> element in a feed.
type FeedCurrentlySyncing struct {
	Timestamp string `xml:"timestamp"`
}
