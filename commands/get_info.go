package commands

import (
	"encoding/xml"
	"time"
)

// GetInfo represents a get_info command request.
type GetInfo struct {
	XMLName xml.Name `xml:"get_info"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr,omitempty"`
	InfoID  string   `xml:"info_id,attr,omitempty"`
	Details string   `xml:"details,attr,omitempty"`
}

// GetInfoResponse represents a get_info command response.
type GetInfoResponse struct {
	XMLName    xml.Name    `xml:"get_info_response"`
	Status     string      `xml:"status,attr"`
	StatusText string      `xml:"status_text,attr"`
	Infos      []Info      `xml:"info"`
	Filters    InfoFilters `xml:"filters"`
	Sort       InfoSort    `xml:"sort"`
	InfoCount  InfoCount   `xml:"info_count"`
	Details    bool        `xml:"details"`
}

// GetStatus returns the status and status text from the response.
func (r *GetInfoResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Info represents an <info> element in the get_info response.
type Info struct {
	ID               string           `xml:"id,attr"`
	Name             string           `xml:"name"`
	Comment          string           `xml:"comment"`
	CreationTime     time.Time        `xml:"creation_time"`
	ModificationTime time.Time        `xml:"modification_time"`
	Writable         bool             `xml:"writable"`
	InUse            bool             `xml:"in_use"`
	UpdateTime       time.Time        `xml:"update_time"`
	Owner            *InfoOwner       `xml:"owner,omitempty"`
	Permissions      *InfoPermissions `xml:"permissions,omitempty"`
	UserTags         *InfoUserTags    `xml:"user_tags,omitempty"`
	CPE              *InfoCPE         `xml:"cpe,omitempty"`
	CVE              *InfoCVE         `xml:"cve,omitempty"`
	CERTBundAdv      *InfoCERTBundAdv `xml:"cert_bund_adv,omitempty"`
	NVT              string           `xml:"nvt,omitempty"`
}

// InfoOwner represents the owner element in an info.
type InfoOwner struct {
	Name string `xml:"name"`
}

// InfoPermissions represents the permissions element in an info.
type InfoPermissions struct {
	Permission []InfoPermission `xml:"permission"`
}

// InfoPermission represents a permission element in info permissions.
type InfoPermission struct {
	Name string `xml:"name"`
}

// InfoUserTags represents the user tags element in an info.
type InfoUserTags struct {
	Count int           `xml:"count"`
	Tags  []InfoUserTag `xml:"tag,omitempty"`
}

// InfoUserTag represents a tag element in info user tags.
type InfoUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// InfoCPE represents the cpe element in an info.
type InfoCPE struct {
	Title        string             `xml:"title"`
	Score        string             `xml:"score"`
	CVERefs      int                `xml:"cve_refs"`
	Status       string             `xml:"status"`
	Deprecated   bool               `xml:"deprecated"`
	DeprecatedBy []InfoDeprecatedBy `xml:"deprecated_by,omitempty"`
	CVEs         *InfoCVEs          `xml:"cves,omitempty"`
	References   *InfoReferences    `xml:"references,omitempty"`
	RawData      string             `xml:"raw_data,omitempty"`
}

// InfoDeprecatedBy represents a deprecated_by element in CPE.
type InfoDeprecatedBy struct {
	CPEID string `xml:"cpe_id,attr"`
}

// InfoCVEs represents the cves element in CPE.
type InfoCVEs struct {
	CVE []string `xml:"cve"`
}

// InfoReferences represents the references element in CPE.
type InfoReferences struct {
	Reference []InfoReference `xml:"reference"`
}

// InfoReference represents a reference element in CPE references.
type InfoReference struct {
	HREF  string `xml:"href,attr"`
	Value string `xml:",chardata"`
}

// InfoCVE represents the cve element in an info.
type InfoCVE struct {
	Score              string                  `xml:"score"`
	CVSSVector         string                  `xml:"cvss_vector"`
	Description        string                  `xml:"description"`
	Products           string                  `xml:"products"`
	EPSS               *InfoEPSS               `xml:"epss,omitempty"`
	NVTs               *InfoNVTs               `xml:"nvts,omitempty"`
	CERT               *InfoCERT               `xml:"cert,omitempty"`
	ConfigurationNodes *InfoConfigurationNodes `xml:"configuration_nodes,omitempty"`
	References         *InfoReferences         `xml:"references,omitempty"`
	RawData            string                  `xml:"raw_data,omitempty"`
}

// InfoEPSS represents the epss element in CVE.
type InfoEPSS struct {
	Value string `xml:",chardata"`
}

// InfoNVTs represents the nvts element in CVE.
type InfoNVTs struct {
	NVT []string `xml:"nvt"`
}

// InfoCERT represents the cert element in CVE.
type InfoCERT struct {
	Value string `xml:",chardata"`
}

// InfoConfigurationNodes represents the configuration_nodes element in CVE.
type InfoConfigurationNodes struct {
	Value string `xml:",chardata"`
}

// InfoCERTBundAdv represents the cert_bund_adv element in an info.
type InfoCERTBundAdv struct {
	Title    string  `xml:"title"`
	Summary  string  `xml:"summary"`
	Severity float64 `xml:"severity"`
	CVERefs  int     `xml:"cve_refs"`
	RawData  string  `xml:"raw_data,omitempty"`
}

// InfoFilters represents the filters element in the response.
type InfoFilters struct {
	ID       string       `xml:"id,attr"`
	Term     string       `xml:"term"`
	Name     string       `xml:"name"`
	Keywords InfoKeywords `xml:"keywords"`
}

// InfoKeywords represents the keywords element in filters.
type InfoKeywords struct {
	Keyword []InfoKeyword `xml:"keyword"`
}

// InfoKeyword represents a keyword element in filters keywords.
type InfoKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// InfoSort represents the sort element in the response.
type InfoSort struct {
	Value string        `xml:",chardata"`
	Field InfoSortField `xml:"field"`
}

// InfoSortField represents the field element in sort.
type InfoSortField struct {
	Order string `xml:"order"`
}

// InfoCount represents the info_count element in the response.
type InfoCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
