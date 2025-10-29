package commands

import (
	"encoding/xml"
	"time"
)

// GetScanners represents a get_scanners command request.
type GetScanners struct {
	XMLName   xml.Name `xml:"get_scanners"`
	ScannerID string   `xml:"scanner_id,attr,omitempty"`
	Filter    string   `xml:"filter,attr,omitempty"`
	FiltID    string   `xml:"filt_id,attr,omitempty"`
	Trash     string   `xml:"trash,attr,omitempty"`
	Details   string   `xml:"details,attr,omitempty"`
}

// GetScannersResponse represents a get_scanners command response.
type GetScannersResponse struct {
	XMLName    xml.Name        `xml:"get_scanners_response"`
	Status     string          `xml:"status,attr"`
	StatusText string          `xml:"status_text,attr"`
	Scanner    []Scanner       `xml:"scanner"`
	Filters    ScannerFilters  `xml:"filters"`
	Sort       ScannerSort     `xml:"sort"`
	Scanners   ScannerScanners `xml:"scanners"`
	Count      ScannerCount    `xml:"scanner_count"`
	Truncate   string          `xml:"truncate,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetScannersResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Scanner represents a <scanner> element in the get_scanners response.
type Scanner struct {
	ID               string                 `xml:"id,attr"`
	Owner            ScannerOwner           `xml:"owner"`
	Name             string                 `xml:"name"`
	Comment          string                 `xml:"comment"`
	Copy             string                 `xml:"copy"`
	CreationTime     time.Time              `xml:"creation_time"`
	ModificationTime time.Time              `xml:"modification_time"`
	Writable         bool                   `xml:"writable"`
	InUse            bool                   `xml:"in_use"`
	Host             string                 `xml:"host"`
	Port             string                 `xml:"port"`
	Type             string                 `xml:"type"`
	CAPub            string                 `xml:"ca_pub"`
	KeyPub           string                 `xml:"key_pub"`
	Tasks            ScannerTasks           `xml:"tasks"`
	Truncated        string                 `xml:"truncated,omitempty"`
	Permissions      ScannerPermissions     `xml:"permissions"`
	UserTags         ScannerUserTags        `xml:"user_tags"`
	CAPubInfo        ScannerCertificateInfo `xml:"ca_pub_info"`
	CertificateInfo  ScannerCertificateInfo `xml:"certificate_info"`
	Credential       ScannerCredential      `xml:"credential"`
	RelayHost        string                 `xml:"relay_host"`
	RelayPort        string                 `xml:"relay_port"`
	Configs          ScannerConfigs         `xml:"configs"`
	Info             ScannerInfo            `xml:"info"`
}

// ScannerOwner represents the owner of a scanner.
type ScannerOwner struct {
	Name string `xml:"name"`
}

// ScannerPermissions represents permissions for a scanner.
type ScannerPermissions struct {
	Permission []ScannerPermission `xml:"permission"`
}

// ScannerPermission represents a single permission for a scanner.
type ScannerPermission struct {
	Name string `xml:"name"`
}

// ScannerUserTags represents user tags attached to a scanner.
type ScannerUserTags struct {
	Count int          `xml:"count"`
	Tag   []ScannerTag `xml:"tag"`
}

// ScannerTag represents a single user tag for a scanner.
type ScannerTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// ScannerCredential represents a credential for a scanner.
type ScannerCredential struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Login string `xml:"login"`
	Trash bool   `xml:"trash"`
}

// ScannerConfigs represents configs attached to a scanner.
type ScannerConfigs struct {
	Config []ScannerConfig `xml:"config"`
}

// ScannerConfig represents a single config for a scanner.
type ScannerConfig struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions"`
}

// ScannerTasks represents tasks attached to a scanner.
type ScannerTasks struct {
	Task []ScannerTask `xml:"task"`
}

// ScannerTask represents a single task for a scanner.
type ScannerTask struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions"`
}

// ScannerInfo represents info about a scanner.
type ScannerInfo struct {
	Scanner     ScannerInfoScanner  `xml:"scanner"`
	Daemon      ScannerInfoDaemon   `xml:"daemon"`
	Protocol    ScannerInfoProtocol `xml:"protocol"`
	Description string              `xml:"description"`
	Params      ScannerInfoParams   `xml:"params"`
}

// ScannerInfoScanner represents the <scanner> element in scanner info.
type ScannerInfoScanner struct {
	Name    string `xml:"name"`
	Version string `xml:"version"`
}

// ScannerInfoDaemon represents the <daemon> element in scanner info.
type ScannerInfoDaemon struct {
	Name    string `xml:"name"`
	Version string `xml:"version"`
}

// ScannerInfoProtocol represents the <protocol> element in scanner info.
type ScannerInfoProtocol struct {
	Name    string `xml:"name"`
	Version string `xml:"version"`
}

// ScannerInfoParams represents the <params> element in scanner info.
type ScannerInfoParams struct {
	Param []string `xml:"param"`
}

// ScannerFilters represents filters in a get_scanners response.
type ScannerFilters struct {
	ID       string                 `xml:"id,attr"`
	Term     string                 `xml:"term"`
	Name     string                 `xml:"name"`
	Keywords ScannerFiltersKeywords `xml:"keywords"`
}

// ScannerFiltersKeywords represents keywords in filters.
type ScannerFiltersKeywords struct {
	Keyword []ScannerFiltersKeyword `xml:"keyword"`
}

// ScannerFiltersKeyword represents a single keyword in filters.
type ScannerFiltersKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// ScannerSort represents sorting information in a get_scanners response.
type ScannerSort struct {
	Value string           `xml:",chardata"`
	Field ScannerSortField `xml:"field"`
}

// ScannerSortField represents a sort field in sorting information.
type ScannerSortField struct {
	Value string `xml:",chardata"`
	Order string `xml:"order"`
}

// ScannerScanners represents meta information (pagination) in a get_scanners response.
type ScannerScanners struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// ScannerCount represents the scanner count in a get_scanners response.
type ScannerCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}

// ScannerCertificateInfo represents certificate information for a scanner.
type ScannerCertificateInfo struct {
	//	XMLName        xml.Name `xml:"certificate_info"`
	TimeStatus        string    `xml:"time_status"`
	ActivationTime    time.Time `xml:"activation_time"`
	ExpirationTime    time.Time `xml:"expiration_time"`
	Issuer            string    `xml:"issuer"`
	MD5Fingerprint    string    `xml:"md5_fingerprint"`
	SHA256Fingerprint string    `xml:"sha256_fingerprint,omitempty"`
	Subject           string    `xml:"subject,omitempty"`
	Serial            string    `xml:"serial,omitempty"`
}
