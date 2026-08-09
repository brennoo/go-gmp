package commands

import (
	"strings"
	"time"
)

// GMPTime is a time.Time that tolerates the timestamp shapes gvmd is
// observed to actually emit, beyond strict RFC 3339: an empty element
// (treated as the zero value, not an error), a leading/trailing whitespace
// wrapped value, the Java/NVD feed style with a non-colon UTC offset (e.g.
// vulnerability feed metadata timestamps), and the C ctime() format that
// some scan-engine-generated report timestamps use despite the schema
// labelling the field iso_time.
type GMPTime struct {
	time.Time
}

var gmpTimeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.000-0700", // NVD/Java feed style, e.g. NVT/CVE metadata
	"2006-01-02T15:04:05-0700",
	time.ANSIC, // C ctime(), e.g. some report timestamps
}

// UnmarshalText implements encoding.TextUnmarshaler, which encoding/xml uses
// for both elements and attributes. It shadows time.Time's own
// UnmarshalText (RFC 3339 only) promoted by embedding.
func (t *GMPTime) UnmarshalText(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "" {
		*t = GMPTime{}
		return nil
	}
	var firstErr error
	for _, layout := range gmpTimeLayouts {
		if parsed, err := time.Parse(layout, s); err == nil {
			*t = GMPTime{parsed}
			return nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
