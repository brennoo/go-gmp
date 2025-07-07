package gmp

import "encoding/xml"

// CreateReportCommand represents a create_report command request.
type CreateReportCommand struct {
	XMLName  xml.Name          `xml:"create_report"`
	Report   *ReportWrapper    `xml:"report"`
	Task     *CreateReportTask `xml:"task,omitempty"`
	InAssets *bool             `xml:"in_assets,omitempty"`
}

// ReportWrapper represents the <report> element for create_report and get_reports.
type ReportWrapper struct {
	ID          string `xml:"id,attr,omitempty"`
	FormatID    string `xml:"format_id,attr,omitempty"`
	ConfigID    string `xml:"config_id,attr,omitempty"`
	Extension   string `xml:"extension,attr,omitempty"`
	ContentType string `xml:"content_type,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	Base64      string `xml:"base64,omitempty"`

	Owner            *ReportOwner     `xml:"owner,omitempty"`
	Name             string           `xml:"name,omitempty"`
	Comment          string           `xml:"comment,omitempty"`
	CreationTime     string           `xml:"creation_time,omitempty"`
	ModificationTime string           `xml:"modification_time,omitempty"`
	Writable         string           `xml:"writable,omitempty"`
	InUse            string           `xml:"in_use,omitempty"`
	Task             *ReportTask      `xml:"task,omitempty"`
	ReportFormat     *ReportFormatRef `xml:"report_format,omitempty"`
	ReportConfig     *ReportConfigRef `xml:"report_config,omitempty"`
	Report           *ReportContent   `xml:"report,omitempty"`
}

// ReportOwner represents the <owner> element.
type ReportOwner struct {
	Name string `xml:"name,omitempty"`
}

// ReportTask represents the <task> element inside <report>.
type ReportTask struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}

// ReportFormatRef represents the <report_format> element inside <report>.
type ReportFormatRef struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}

// ReportConfigRef represents the <report_config> element inside <report>.
type ReportConfigRef struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}

// ReportContent represents the nested <report> element (full report data, can be expanded as needed).
type ReportContent struct {
	ID              string                 `xml:"id,attr,omitempty"`
	Type            string                 `xml:"type,attr,omitempty"`
	GMP             *ReportGMP             `xml:"gmp,omitempty"`
	ReportFormat    *ReportFormatParams    `xml:"report_format,omitempty"`
	Sort            *ReportSort            `xml:"sort,omitempty"`
	Filters         *ReportFilters         `xml:"filters,omitempty"`
	Delta           *ReportDelta           `xml:"delta,omitempty"`
	Timezone        string                 `xml:"timezone,omitempty"`
	TimezoneAbbrev  string                 `xml:"timezone_abbrev,omitempty"`
	Permissions     *ReportPermissions     `xml:"permissions,omitempty"`
	UserTags        *ReportUserTags        `xml:"user_tags,omitempty"`
	ScanRunStatus   string                 `xml:"scan_run_status,omitempty"`
	ResultCount     *ReportResultCount     `xml:"result_count,omitempty"`
	Severity        *ReportSeverity        `xml:"severity,omitempty"`
	ComplianceCount *ReportComplianceCount `xml:"compliance_count,omitempty"`
	Compliance      *ReportCompliance      `xml:"compliance,omitempty"`
	Task            *ReportTaskDetail      `xml:"task,omitempty"`
	Ports           *ReportPorts           `xml:"ports,omitempty"`
	Results         *ReportResults         `xml:"results,omitempty"`
	Hosts           *ReportHosts           `xml:"hosts,omitempty"`
	ClosedCVEs      *ReportCount           `xml:"closed_cves,omitempty"`
	Vulns           *ReportCount           `xml:"vulns,omitempty"`
	OS              *ReportCount           `xml:"os,omitempty"`
	Apps            *ReportCount           `xml:"apps,omitempty"`
	SSLCerts        *ReportCount           `xml:"ssl_certs,omitempty"`
	Host            []ReportHost           `xml:"host,omitempty"`
	TLSCertificates *ReportTLSCertificates `xml:"tls_certificates,omitempty"`
	Timestamp       string                 `xml:"timestamp,omitempty"`
	ScanStart       string                 `xml:"scan_start,omitempty"`
	ScanEnd         string                 `xml:"scan_end,omitempty"`
	Errors          *ReportErrors          `xml:"errors,omitempty"`
}

// CreateReportTask represents the <task> element for create_report.
type CreateReportTask struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateReportResponse represents a create_report command response.
type CreateReportResponse struct {
	XMLName    xml.Name `xml:"create_report_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

type ReportGMP struct {
	Version string `xml:"version,omitempty"`
}

type ReportFormatParams struct {
	Params []ReportFormatParam `xml:"param,omitempty"`
}
type ReportFormatParam struct {
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

type ReportSort struct {
	Fields []ReportSortField `xml:"field,omitempty"`
}
type ReportSortField struct {
	Order string `xml:"order,omitempty"`
}

type ReportFilters struct {
	ID              string             `xml:"id,attr,omitempty"`
	Levels          string             `xml:"levels,omitempty"`
	Term            string             `xml:"term,omitempty"`
	Phrase          string             `xml:"phrase,omitempty"`
	Notes           string             `xml:"notes,omitempty"`
	Overrides       string             `xml:"overrides,omitempty"`
	ApplyOverrides  string             `xml:"apply_overrides,omitempty"`
	ResultHostsOnly string             `xml:"result_hosts_only,omitempty"`
	MinQod          string             `xml:"min_qod,omitempty"`
	Filter          []string           `xml:"filter,omitempty"`
	Delta           *ReportDeltaStates `xml:"delta,omitempty"`
}
type ReportDeltaStates struct {
	Changed string `xml:"changed,omitempty"`
	Gone    string `xml:"gone,omitempty"`
	New     string `xml:"new,omitempty"`
	Same    string `xml:"same,omitempty"`
}

type ReportDelta struct {
	Report *ReportContent `xml:"report,omitempty"`
}

type ReportPermissions struct {
	Permissions []ReportPermission `xml:"permission,omitempty"`
}
type ReportPermission struct {
	Name string `xml:"name,omitempty"`
}

type ReportUserTags struct {
	Count int         `xml:"count,omitempty"`
	Tags  []ReportTag `xml:"tag,omitempty"`
}
type ReportTag struct {
	ID      string `xml:"id,attr,omitempty"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

type ReportResultCount struct {
	Full     int `xml:"full,omitempty"`
	Filtered int `xml:"filtered,omitempty"`
}

type ReportSeverity struct {
	Full     string `xml:"full,omitempty"`
	Filtered string `xml:"filtered,omitempty"`
}

type ReportComplianceCount struct {
	Full       int `xml:"full,omitempty"`
	Filtered   int `xml:"filtered,omitempty"`
	Yes        int `xml:"yes,omitempty"`
	No         int `xml:"no,omitempty"`
	Incomplete int `xml:"incomplete,omitempty"`
	Undefined  int `xml:"undefined,omitempty"`
}

type ReportCompliance struct {
	Full     string `xml:"full,omitempty"`
	Filtered string `xml:"filtered,omitempty"`
}

type ReportTaskDetail struct {
	ID       string            `xml:"id,attr,omitempty"`
	Name     string            `xml:"name,omitempty"`
	Comment  string            `xml:"comment,omitempty"`
	Target   *ReportTaskTarget `xml:"target,omitempty"`
	Progress int               `xml:"progress,omitempty"`
	UserTags *ReportUserTags   `xml:"user_tags,omitempty"`
}
type ReportTaskTarget struct {
	ID      string `xml:"id,attr,omitempty"`
	Trash   string `xml:"trash,omitempty"`
	Name    string `xml:"name,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

type ReportPorts struct {
	Start int          `xml:"start,attr,omitempty"`
	Max   int          `xml:"max,attr,omitempty"`
	Ports []ReportPort `xml:"port,omitempty"`
}
type ReportPort struct {
	Host     string `xml:"host,omitempty"`
	Severity string `xml:"severity,omitempty"`
	Threat   string `xml:"threat,omitempty"`
}

type ReportResults struct {
	Start   int            `xml:"start,attr,omitempty"`
	Max     int            `xml:"max,attr,omitempty"`
	Results []ReportResult `xml:"result,omitempty"`
}
type ReportResult struct {
	ID               string          `xml:"id,attr,omitempty"`
	Name             string          `xml:"name,omitempty"`
	Owner            *ReportOwner    `xml:"owner,omitempty"`
	Comment          string          `xml:"comment,omitempty"`
	CreationTime     string          `xml:"creation_time,omitempty"`
	ModificationTime string          `xml:"modification_time,omitempty"`
	UserTags         *ReportUserTags `xml:"user_tags,omitempty"`
	Report           *ReportRef      `xml:"report,omitempty"`
	Task             *ReportTaskRef  `xml:"task,omitempty"`
	Host             *ReportHostRef  `xml:"host,omitempty"`
	Port             string          `xml:"port,omitempty"`
	Path             string          `xml:"path,omitempty"`
	NVT              *ReportNVT      `xml:"nvt,omitempty"`
	ScanNVTVersion   string          `xml:"scan_nvt_version,omitempty"`
	Threat           string          `xml:"threat,omitempty"`
	Severity         string          `xml:"severity,omitempty"`
	QOD              string          `xml:"qod,omitempty"`
	OriginalThreat   string          `xml:"original_threat,omitempty"`
	OriginalSeverity string          `xml:"original_severity,omitempty"`
	Compliance       string          `xml:"compliance,omitempty"`
	Description      string          `xml:"description,omitempty"`
	Delta            string          `xml:"delta,omitempty"`
	Detection        string          `xml:"detection,omitempty"`
	Notes            string          `xml:"notes,omitempty"`
	Overrides        string          `xml:"overrides,omitempty"`
	Tickets          string          `xml:"tickets,omitempty"`
}
type ReportRef struct {
	ID string `xml:"id,attr,omitempty"`
}
type ReportTaskRef struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}
type ReportHostRef struct {
	AssetID  string `xml:"asset_id,attr,omitempty"`
	Hostname string `xml:"hostname,omitempty"`
	IP       string `xml:",chardata"`
}
type ReportNVT struct {
	OID        string `xml:"oid,attr,omitempty"`
	Name       string `xml:"name,omitempty"`
	Type       string `xml:"type,omitempty"`
	Family     string `xml:"family,omitempty"`
	CVSSBase   string `xml:"cvss_base,omitempty"`
	Severities string `xml:"severities,omitempty"`
	CPE        string `xml:"cpe,omitempty"`
	Tags       string `xml:"tags,omitempty"`
	EPSS       string `xml:"epss,omitempty"`
	Refs       string `xml:"refs,omitempty"`
}

type ReportHosts struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
	Count int `xml:"count,omitempty"`
}

type ReportCount struct {
	Count int `xml:"count,omitempty"`
}

type ReportHost struct {
	IP              string             `xml:"ip,omitempty"`
	Asset           string             `xml:"asset,omitempty"`
	Start           string             `xml:"start,omitempty"`
	End             string             `xml:"end,omitempty"`
	PortCount       int                `xml:"port_count,omitempty"`
	ResultCount     int                `xml:"result_count,omitempty"`
	ComplianceCount int                `xml:"compliance_count,omitempty"`
	HostCompliance  string             `xml:"host_compliance,omitempty"`
	Details         []ReportHostDetail `xml:"detail,omitempty"`
}
type ReportHostDetail struct {
	Name   string `xml:"name,omitempty"`
	Value  string `xml:"value,omitempty"`
	Source string `xml:"source,omitempty"`
	Extra  string `xml:"extra,omitempty"`
}

type ReportTLSCertificates struct {
	Certificates []ReportTLSCertificate `xml:"tls_certificate,omitempty"`
}
type ReportTLSCertificate struct {
	Name              string `xml:"name,omitempty"`
	Certificate       string `xml:"certificate,omitempty"`
	Format            string `xml:"format,attr,omitempty"`
	SHA256Fingerprint string `xml:"sha256_fingerprint,omitempty"`
	MD5Fingerprint    string `xml:"md5_fingerprint,omitempty"`
	Valid             string `xml:"valid,omitempty"`
	ActivationTime    string `xml:"activation_time,omitempty"`
	ExpirationTime    string `xml:"expiration_time,omitempty"`
	SubjectDN         string `xml:"subject_dn,omitempty"`
	IssuerDN          string `xml:"issuer_dn,omitempty"`
	Serial            string `xml:"serial,omitempty"`
	Host              string `xml:"host,omitempty"`
	Ports             []int  `xml:"ports>port,omitempty"`
}

type ReportErrors struct {
	Count  int           `xml:"count,omitempty"`
	Errors []ReportError `xml:"error,omitempty"`
}
type ReportError struct {
	Host           string `xml:"host,omitempty"`
	Port           string `xml:"port,omitempty"`
	Description    string `xml:"description,omitempty"`
	NVT            string `xml:"nvt,omitempty"`
	ScanNVTVersion string `xml:"scan_nvt_version,omitempty"`
	Severity       string `xml:"severity,omitempty"`
}
