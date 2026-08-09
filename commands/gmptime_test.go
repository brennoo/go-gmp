package commands

import (
	"testing"
	"time"
)

func TestGMPTimeUnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    time.Time
		wantErr bool
	}{
		{
			name: "rfc3339",
			in:   "2011-11-08T19:57:06+02:00",
			want: time.Date(2011, 11, 8, 19, 57, 6, 0, time.FixedZone("", 2*60*60)),
		},
		{
			name: "rfc3339 with fractional seconds",
			in:   "2011-11-08T19:57:06.123456789+02:00",
			want: time.Date(2011, 11, 8, 19, 57, 6, 123456789, time.FixedZone("", 2*60*60)),
		},
		{
			name: "empty is zero value, not an error",
			in:   "",
			want: time.Time{},
		},
		{
			name: "whitespace-only is zero value, not an error",
			in:   "\n            ",
			want: time.Time{},
		},
		{
			name: "leading and trailing whitespace is trimmed",
			in:   "\n            2019-08-13T12:28:44+02:00",
			want: time.Date(2019, 8, 13, 12, 28, 44, 0, time.FixedZone("", 2*60*60)),
		},
		{
			name: "NVD feed style with milliseconds and no colon in the offset",
			in:   "2012-10-26T13:18:00.000+0000",
			want: time.Date(2012, 10, 26, 13, 18, 0, 0, time.FixedZone("", 0)),
		},
		{
			name: "NVD feed style without milliseconds",
			in:   "2012-10-26T13:18:00+0000",
			want: time.Date(2012, 10, 26, 13, 18, 0, 0, time.FixedZone("", 0)),
		},
		{
			name: "C ctime format",
			in:   "Mon Feb  1 19:11:20 2010",
			want: time.Date(2010, 2, 1, 19, 11, 20, 0, time.UTC),
		},
		{
			name:    "unparseable garbage is an error",
			in:      "not a time",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got GMPTime
			err := got.UnmarshalText([]byte(tt.in))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("UnmarshalText(%q) = nil error, want an error", tt.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalText(%q) = %v, want nil error", tt.in, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("UnmarshalText(%q) = %v, want %v", tt.in, got.Time, tt.want)
			}
		})
	}
}

func TestGMPTimeMarshalTextIsPromoted(t *testing.T) {
	gt := GMPTime{Time: time.Date(2011, 11, 8, 19, 57, 6, 0, time.UTC)}
	got, err := gt.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}
	want := "2011-11-08T19:57:06Z"
	if string(got) != want {
		t.Errorf("MarshalText() = %q, want %q", got, want)
	}
}
