package commands

import (
	"encoding/xml"
)

// GetResults represents a get_results command request.
type GetResults struct {
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
	XMLName    xml.Name       `xml:"get_results_response"`
	Status     string         `xml:"status,attr"`
	StatusText string         `xml:"status_text,attr"`
	Results    []Result       `xml:"result"`
	Filters    ResultsFilters `xml:"filters"`
	Sort       ResultsSort    `xml:"sort"`
	Meta       ResultsMeta    `xml:"results"`
	Count      *ResultsCount  `xml:"result_count,omitempty"`
}

// Result represents a single result element in the get_results response.
type Result struct {
	ID               string           `xml:"id,attr"`
	Name             string           `xml:"name"`
	Owner            ResultOwner      `xml:"owner"`
	Comment          string           `xml:"comment"`
	CreationTime     string           `xml:"creation_time"`
	ModificationTime string           `xml:"modification_time"`
	UserTags         *ResultUserTags  `xml:"user_tags,omitempty"`
	Report           *ResultReportRef `xml:"report,omitempty"`
	Task             *ResultTaskRef   `xml:"task,omitempty"`
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

// ResultOwner represents the owner of a result.
type ResultOwner struct {
	Name string `xml:"name"`
}

// ResultUserTags represents user tags attached to a result.
type ResultUserTags struct {
	Count int         `xml:"count"`
	Tag   []ResultTag `xml:"tag"`
}

// ResultTag represents a single user tag.
type ResultTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// ResultReportRef represents a reference to a report.
type ResultReportRef struct {
	ID string `xml:"id,attr"`
}

// ResultTaskRef represents a reference to a task.
type ResultTaskRef struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// ResultHost represents a host element in a result.
type ResultHost struct {
	Value    string          `xml:",chardata"`
	Asset    ResultHostAsset `xml:"asset"`
	Hostname string          `xml:"hostname"`
}

// ResultHostAsset represents an asset attached to a host.
type ResultHostAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

// ResultNVT represents a network vulnerability test (NVT) in a result.
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

// ResultNVTSeverities represents severity information for an NVT.
type ResultNVTSeverities struct {
	Score float32 `xml:"score,attr"`
}

// ResultNVTEPSS represents EPSS information for an NVT.
type ResultNVTEPSS struct {
	MaxSeverity ResultNVTEPSSInfo `xml:"max_severity"`
	MaxEPSS     ResultNVTEPSSInfo `xml:"max_epss"`
}

// ResultNVTEPSSInfo represents detailed EPSS info for an NVT.
type ResultNVTEPSSInfo struct {
	Score      float32          `xml:"score"`
	Percentile float32          `xml:"percentile"`
	CVE        ResultNVTEPSSCVE `xml:"cve"`
}

// ResultNVTEPSSCVE represents a CVE reference in EPSS info.
type ResultNVTEPSSCVE struct {
	ID       string   `xml:"id,attr"`
	Severity *float32 `xml:"severity,omitempty"`
}

// ResultNVTRefs represents references for an NVT.
type ResultNVTRefs struct {
	Ref []ResultNVTRef `xml:"ref"`
}

// ResultNVTRef represents a single reference for an NVT.
type ResultNVTRef struct {
	ID   string `xml:"id,attr"`
	Type string `xml:"type,attr"`
}

// ResultQOD represents Quality of Detection information.
type ResultQOD struct {
	Value int    `xml:"value"`
	Type  string `xml:"type"`
}

// ResultDelta represents a delta element in a result.
type ResultDelta struct {
	Value     string               `xml:",chardata"`
	Result    *Result              `xml:"result"`
	Diff      string               `xml:"diff"`
	Notes     ResultDeltaNotes     `xml:"notes"`
	Overrides ResultDeltaOverrides `xml:"overrides"`
}

// ResultDeltaNotes represents notes in a delta.
type ResultDeltaNotes struct {
	Note []ResultNote `xml:"note"`
}

// ResultNote represents a note attached to a result or delta.
type ResultNote struct {
	ID               string             `xml:"id,attr"`
	Permissions      ResultNotePerms    `xml:"permissions"`
	Owner            ResultNoteOwner    `xml:"owner"`
	NVT              ResultNoteNVT      `xml:"nvt"`
	Text             ResultNoteText     `xml:"text"`
	CreationTime     string             `xml:"creation_time"`
	ModificationTime string             `xml:"modification_time"`
	Writable         string             `xml:"writable"`
	InUse            string             `xml:"in_use"`
	Active           string             `xml:"active"`
	Orphan           string             `xml:"orphan"`
	UserTags         ResultNoteUserTags `xml:"user_tags"`
	Hosts            string             `xml:"hosts"`
	Port             string             `xml:"port"`
	Severity         string             `xml:"severity"`
	Task             ResultNoteTask     `xml:"task"`
	EndTime          string             `xml:"end_time"`
	Result           ResultNoteResult   `xml:"result"`
}

// ResultNotePerms represents permissions for a note.
type ResultNotePerms struct {
	Permission []ResultNotePerm `xml:"permission"`
}

// ResultNotePerm represents a single permission for a note.
type ResultNotePerm struct {
	Name string `xml:"name"`
}

// ResultNoteOwner represents the owner of a note.
type ResultNoteOwner struct {
	Name string `xml:"name"`
}

// ResultNoteNVT represents an NVT in a note.
type ResultNoteNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
	Type string `xml:"type"`
}

// ResultNoteText represents the text of a note.
type ResultNoteText struct {
	Value   string `xml:",chardata"`
	Excerpt string `xml:"excerpt,attr"`
}

// ResultNoteUserTags represents user tags in a note.
type ResultNoteUserTags struct {
	Count int             `xml:"count"`
	Tag   []ResultNoteTag `xml:"tag"`
}

// ResultNoteTag represents a single tag in a note.
type ResultNoteTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// ResultNoteTask represents a task in a note.
type ResultNoteTask struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Trash string `xml:"trash"`
}

// ResultNoteResult represents a result reference in a note.
type ResultNoteResult struct {
	ID          string         `xml:"id,attr"`
	Host        ResultNoteHost `xml:"host"`
	Port        string         `xml:"port"`
	NVT         ResultNoteNVT  `xml:"nvt"`
	Severity    string         `xml:"severity"`
	Threat      string         `xml:"threat"`
	QOD         ResultNoteQOD  `xml:"qod"`
	Description string         `xml:"description"`
}

// ResultNoteHost represents a host in a note result.
type ResultNoteHost struct {
	Value string          `xml:",chardata"`
	Asset ResultNoteAsset `xml:"asset"`
}

// ResultNoteAsset represents an asset in a note host.
type ResultNoteAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

// ResultNoteQOD represents QOD in a note result.
type ResultNoteQOD struct {
	Value int    `xml:"value"`
	Type  string `xml:"type"`
}

// ResultDeltaOverrides represents overrides in a delta.
type ResultDeltaOverrides struct {
	Override []ResultOverride `xml:"override"`
}

// ResultOverride represents an override in a result or delta.
type ResultOverride struct {
	ID               string                 `xml:"id,attr"`
	Permissions      ResultOverridePerms    `xml:"permissions"`
	Owner            ResultOverrideOwner    `xml:"owner"`
	NVT              ResultOverrideNVT      `xml:"nvt"`
	CreationTime     string                 `xml:"creation_time"`
	ModificationTime string                 `xml:"modification_time"`
	Writable         string                 `xml:"writable"`
	InUse            string                 `xml:"in_use"`
	Active           string                 `xml:"active"`
	Text             ResultOverrideText     `xml:"text"`
	Threat           string                 `xml:"threat"`
	Severity         string                 `xml:"severity"`
	NewThreat        string                 `xml:"new_threat"`
	NewSeverity      string                 `xml:"new_severity"`
	Orphan           string                 `xml:"orphan"`
	UserTags         ResultOverrideUserTags `xml:"user_tags"`
	Hosts            string                 `xml:"hosts"`
	Port             string                 `xml:"port"`
	Task             ResultOverrideTask     `xml:"task"`
	EndTime          string                 `xml:"end_time"`
	Result           ResultOverrideResult   `xml:"result"`
}

// ResultOverridePerms represents permissions for an override.
type ResultOverridePerms struct {
	Permission []ResultOverridePerm `xml:"permission"`
}

// ResultOverridePerm represents a single permission for an override.
type ResultOverridePerm struct {
	Name string `xml:"name"`
}

// ResultOverrideOwner represents the owner of an override.
type ResultOverrideOwner struct {
	Name string `xml:"name"`
}

// ResultOverrideNVT represents an NVT in an override.
type ResultOverrideNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
	Type string `xml:"type"`
}

// ResultOverrideText represents the text of an override.
type ResultOverrideText struct {
	Text    string `xml:",chardata"`
	Excerpt string `xml:"excerpt,attr"`
}

// ResultOverrideUserTags represents user tags in an override.
type ResultOverrideUserTags struct {
	Count int                 `xml:"count"`
	Tag   []ResultOverrideTag `xml:"tag"`
}

// ResultOverrideTag represents a single tag in an override.
type ResultOverrideTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// ResultOverrideTask represents a task in an override.
type ResultOverrideTask struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Trash string `xml:"trash"`
}

// ResultOverrideResult represents a result reference in an override.
type ResultOverrideResult struct {
	ID          string             `xml:"id,attr"`
	Host        ResultOverrideHost `xml:"host"`
	Port        string             `xml:"port"`
	NVT         ResultOverrideNVT  `xml:"nvt"`
	Threat      string             `xml:"threat"`
	Severity    string             `xml:"severity"`
	QOD         ResultOverrideQOD  `xml:"qod"`
	Description string             `xml:"description"`
}

// ResultOverrideHost represents a host in an override result.
type ResultOverrideHost struct {
	Value string              `xml:",chardata"`
	Asset ResultOverrideAsset `xml:"asset"`
}

// ResultOverrideAsset represents an asset in an override host.
type ResultOverrideAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

// ResultOverrideQOD represents QOD in an override result.
type ResultOverrideQOD struct {
	Value int    `xml:"value"`
	Type  string `xml:"type"`
}

// ResultDetection represents detection information in a result.
type ResultDetection struct {
	Result ResultDetectionResult `xml:"result"`
}

// ResultDetectionResult represents a detection result.
type ResultDetectionResult struct {
	ID      string                 `xml:"id,attr"`
	Details ResultDetectionDetails `xml:"details"`
}

// ResultDetectionDetails represents details in a detection result.
type ResultDetectionDetails struct {
	Detail []ResultDetectionDetail `xml:"detail"`
}

// ResultDetectionDetail represents a single detail in detection.
type ResultDetectionDetail struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// ResultNotes represents notes in a result.
type ResultNotes struct {
	Note []ResultNote `xml:"note"`
}

// ResultOverrides represents overrides in a result.
type ResultOverrides struct {
	Override []ResultOverride `xml:"override"`
}

// ResultTickets represents tickets in a result.
type ResultTickets struct {
	Ticket []ResultTicketRef `xml:"ticket"`
}

// ResultTicketRef represents a ticket reference in a result.
type ResultTicketRef struct {
	ID string `xml:"id,attr"`
}

// ResultsFilters represents filters in a get_results response.
type ResultsFilters struct {
	ID       string                 `xml:"id,attr"`
	Term     string                 `xml:"term"`
	Name     *string                `xml:"name,omitempty"`
	Keywords []ResultsFilterKeyword `xml:"keywords>keyword"`
}

// ResultsFilterKeyword represents a filter keyword in filters.
type ResultsFilterKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// ResultsSort represents sorting information in a get_results response.
type ResultsSort struct {
	Value string           `xml:",chardata"`
	Field ResultsSortField `xml:"field"`
}

// ResultsSortField represents a sort field in sorting information.
type ResultsSortField struct {
	Order string `xml:"order"`
}

// ResultsMeta represents meta information (pagination) in a get_results response.
type ResultsMeta struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// ResultsCount represents the result count in a get_results response.
type ResultsCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
