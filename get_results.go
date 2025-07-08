package gmp

import (
	"encoding/xml"
	"time"
)

// GetResultsCommand represents a get_results command request.
type GetResultsCommand struct {
	XMLName          xml.Name `xml:"get_results"`
	ResultID         string   `xml:"result_id,attr,omitempty"`
	Filter           string   `xml:"filter,attr,omitempty"`
	FiltID           string   `xml:"filt_id,attr,omitempty"`
	TaskID           string   `xml:"task_id,attr,omitempty"`
	NotesDetails     bool     `xml:"notes_details,attr,omitempty"`
	OverridesDetails bool     `xml:"overrides_details,attr,omitempty"`
	Details          bool     `xml:"details,attr,omitempty"`
	GetCounts        bool     `xml:"get_counts,attr,omitempty"`
}

// GetResultsResponse represents a get_results command response.
type GetResultsResponse struct {
	XMLName     xml.Name          `xml:"get_results_response"`
	Status      string            `xml:"status,attr"`
	StatusText  string            `xml:"status_text,attr"`
	Results     []Result          `xml:"result"`
	Filters     GetResultsFilters `xml:"filters"`
	Sort        GetResultsSort    `xml:"sort"`
	ResultsMeta GetResultsMeta    `xml:"results"`
	ResultCount *GetResultsCount  `xml:"result_count,omitempty"`
}

type Result struct {
	ID               string           `xml:"id,attr"`
	Name             string           `xml:"name"`
	Owner            ResultOwner      `xml:"owner"`
	Comment          string           `xml:"comment"`
	CreationTime     time.Time        `xml:"creation_time"`
	ModificationTime time.Time        `xml:"modification_time"`
	UserTags         *ResultUserTags  `xml:"user_tags,omitempty"`
	Report           *ResultReport    `xml:"report,omitempty"`
	Task             *ResultTask      `xml:"task,omitempty"`
	Host             ResultHost       `xml:"host"`
	Port             string           `xml:"port"`
	Path             *string          `xml:"path,omitempty"`
	NVT              ResultNVT        `xml:"nvt"`
	ScanNVTVersion   string           `xml:"scan_nvt_version"`
	Threat           string           `xml:"threat"`
	Severity         float32          `xml:"severity"`
	QOD              ResultQOD        `xml:"qod"`
	OriginalThreat   *string          `xml:"original_threat,omitempty"`
	OriginalSeverity *float32         `xml:"original_severity,omitempty"`
	Compliance       string           `xml:"compliance"`
	Description      string           `xml:"description"`
	Delta            *ResultDelta     `xml:"delta,omitempty"`
	Detection        *ResultDetection `xml:"detection,omitempty"`
	Notes            *ResultNotes     `xml:"notes,omitempty"`
	Overrides        *ResultOverrides `xml:"overrides,omitempty"`
	Tickets          *ResultTickets   `xml:"tickets,omitempty"`
}

type ResultOwner struct {
	Name string `xml:"name"`
}

type ResultUserTags struct {
	Count int                 `xml:"count"`
	Tag   []ResultUserTagsTag `xml:"tag"`
}

type ResultUserTagsTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

type ResultReport struct {
	ID string `xml:"id,attr"`
}

type ResultTask struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type ResultHost struct {
	Value    string          `xml:",chardata"`
	Asset    ResultHostAsset `xml:"asset"`
	Hostname string          `xml:"hostname"`
}

type ResultHostAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

type ResultNVTSeverities struct {
	Score float32 `xml:"score,attr"`
}

type ResultNVT struct {
	OID        string              `xml:"oid,attr"`
	Name       string              `xml:"name"`
	Type       string              `xml:"type"`
	Family     string              `xml:"family"`
	CVSSBase   float32             `xml:"cvss_base"`
	Severities ResultNVTSeverities `xml:"severities"`
	CPE        string              `xml:"cpe"`
	Tags       string              `xml:"tags"`
	EPSS       *ResultNVTEPSS      `xml:"epss,omitempty"`
	Refs       ResultNVTRefs       `xml:"refs"`
}

type ResultNVTEPSS struct {
	MaxSeverity ResultNVTEPSSInfo `xml:"max_severity"`
	MaxEPSS     ResultNVTEPSSInfo `xml:"max_epss"`
}

type ResultNVTEPSSInfo struct {
	Score      float32          `xml:"score"`
	Percentile float32          `xml:"percentile"`
	CVE        ResultNVTEPSSCVE `xml:"cve"`
}

type ResultNVTEPSSCVE struct {
	ID       string   `xml:"id,attr"`
	Severity *float32 `xml:"severity,omitempty"`
}

type ResultNVTRefs struct {
	Ref []ResultNVTRefsRef `xml:"ref"`
}

type ResultNVTRefsRef struct {
	ID   string `xml:"id,attr"`
	Type string `xml:"type,attr"`
}

type ResultQOD struct {
	Value int    `xml:"value"`
	Type  string `xml:"type"`
}

type ResultDelta struct {
	Value     string               `xml:",chardata"`
	Result    *Result              `xml:"result"`
	Diff      string               `xml:"diff"`
	Notes     ResultDeltaNotes     `xml:"notes"`
	Overrides ResultDeltaOverrides `xml:"overrides"`
}

type ResultDeltaNotes struct {
	Note []Note `xml:"note"`
}

type Note struct {
	ID               string          `xml:"id,attr"`
	Permissions      NotePermissions `xml:"permissions"`
	Owner            NoteOwner       `xml:"owner"`
	NVT              NoteNVT         `xml:"nvt"`
	Text             NoteText        `xml:"text"`
	CreationTime     time.Time       `xml:"creation_time"`
	ModificationTime time.Time       `xml:"modification_time"`
	Writable         bool            `xml:"writable"`
	InUse            bool            `xml:"in_use"`
	Active           bool            `xml:"active"`
	Orphan           bool            `xml:"orphan"`
	UserTags         NoteUserTags    `xml:"user_tags"`
	Hosts            string          `xml:"hosts"`
	Port             string          `xml:"port"`
	Severity         string          `xml:"severity"`
	Task             NoteTask        `xml:"task"`
	EndTime          string          `xml:"end_time"`
	Result           NoteResult      `xml:"result"`
}

type NotePermissions struct {
	Permission []NotePermissionsPermission `xml:"permission"`
}

type NotePermissionsPermission struct {
	Name string `xml:"name"`
}

type NoteOwner struct {
	Name string `xml:"name"`
}

type NoteNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
	Type string `xml:"type"`
}

type NoteText struct {
	Value   string `xml:",chardata"`
	Excerpt bool   `xml:"excerpt,attr"`
}

type NoteUserTags struct {
	Count int               `xml:"count"`
	Tag   []NoteUserTagsTag `xml:"tag"`
}

type NoteUserTagsTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

type NoteTask struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Trash bool   `xml:"trash"`
}

type NoteResult struct {
	ID          string         `xml:"id,attr"`
	Host        NoteResultHost `xml:"host"`
	Port        string         `xml:"port"`
	NVT         NoteResultNVT  `xml:"nvt"`
	Severity    string         `xml:"severity"`
	Threat      string         `xml:"threat"`
	QOD         NoteResultQOD  `xml:"qod"`
	Description string         `xml:"description"`
}

type NoteResultHost struct {
	Value string              `xml:",chardata"`
	Asset NoteResultHostAsset `xml:"asset"`
}

type NoteResultHostAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

type NoteResultNVT struct {
	OID      string `xml:"oid,attr"`
	Name     string `xml:"name"`
	Type     string `xml:"type"`
	CVSSBase string `xml:"cvss_base"`
	CVE      string `xml:"cve"`
	BID      int    `xml:"bid"`
}

type NoteResultQOD struct {
	Value int    `xml:"value"`
	Type  string `xml:"type"`
}

type ResultDeltaOverrides struct {
	Override []Override `xml:"override"`
}

type Override struct {
	ID               string              `xml:"id,attr"`
	Permissions      OverridePermissions `xml:"permissions"`
	Owner            OverrideOwner       `xml:"owner"`
	NVT              OverrideNVT         `xml:"nvt"`
	CreationTime     time.Time           `xml:"creation_time"`
	ModificationTime time.Time           `xml:"modification_time"`
	Writable         bool                `xml:"writable"`
	InUse            bool                `xml:"in_use"`
	Active           bool                `xml:"active"`
	Text             OverrideText        `xml:"text"`
	Threat           string              `xml:"threat"`
	Severity         string              `xml:"severity"`
	NewThreat        string              `xml:"new_threat"`
	NewSeverity      string              `xml:"new_severity"`
	Orphan           bool                `xml:"orphan"`
	UserTags         OverrideUserTags    `xml:"user_tags"`
	Hosts            string              `xml:"hosts"`
	Port             string              `xml:"port"`
	Task             OverrideTask        `xml:"task"`
	EndTime          string              `xml:"end_time"`
	Result           OverrideResult      `xml:"result"`
}

type OverridePermissions struct {
	Permission []OverridePermissionsPermission `xml:"permission"`
}

type OverridePermissionsPermission struct {
	Name string `xml:"name"`
}

type OverrideOwner struct {
	Name string `xml:"name"`
}

type OverrideNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
	Type string `xml:"type"`
}

type OverrideText struct {
	Text    string `xml:",chardata"`
	Excerpt bool   `xml:"excerpt,attr"`
}

type OverrideUserTags struct {
	Count int                   `xml:"count"`
	Tag   []OverrideUserTagsTag `xml:"tag"`
}

type OverrideUserTagsTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

type OverrideTask struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Trash bool   `xml:"trash"`
}

type OverrideResult struct {
	ID          string             `xml:"id,attr"`
	Host        OverrideResultHost `xml:"host"`
	Port        string             `xml:"port"`
	NVT         OverrideResultNVT  `xml:"nvt"`
	Threat      string             `xml:"threat"`
	Severity    string             `xml:"severity"`
	QOD         OverrideResultQOD  `xml:"qod"`
	Description string             `xml:"description"`
}

type OverrideResultHost struct {
	Value string                  `xml:",chardata"`
	Asset OverrideResultHostAsset `xml:"asset"`
}

type OverrideResultHostAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

type OverrideResultNVT struct {
	OID      string `xml:"oid,attr"`
	Name     string `xml:"name"`
	Type     string `xml:"type"`
	CVSSBase string `xml:"cvss_base"`
	CVE      string `xml:"cve"`
	BID      int    `xml:"bid"`
}

type OverrideResultQOD struct {
	Value int    `xml:"value"`
	Type  string `xml:"type"`
}

type ResultDetection struct {
	Result ResultDetectionResult `xml:"result"`
}

type ResultDetectionResult struct {
	ID      string                       `xml:"id,attr"`
	Details ResultDetectionResultDetails `xml:"details"`
}

type ResultDetectionResultDetails struct {
	Detail []ResultDetectionResultDetailsDetail `xml:"detail"`
}

type ResultDetectionResultDetailsDetail struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type ResultNotes struct {
	Note []Note `xml:"note"`
}

type ResultOverrides struct {
	Override []Override `xml:"override"`
}

type ResultTickets struct {
	Ticket []ResultTicketsTicket `xml:"ticket"`
}

type ResultTicketsTicket struct {
	ID string `xml:"id,attr"`
}

type GetResultsFilters struct {
	ID       string                     `xml:"id,attr"`
	Term     string                     `xml:"term"`
	Name     *string                    `xml:"name,omitempty"`
	Keywords []GetResultsFiltersKeyword `xml:"keywords>keyword"`
}

type GetResultsFiltersKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

type GetResultsSort struct {
	Value string              `xml:",chardata"`
	Field GetResultsSortField `xml:"field"`
}

type GetResultsSortField struct {
	Order string `xml:"order"`
}

type GetResultsMeta struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

type GetResultsCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
