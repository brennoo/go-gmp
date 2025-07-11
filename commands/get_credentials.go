package commands

import "encoding/xml"

// GetCredentials represents a get_credentials command request to GMP.
type GetCredentials struct {
	XMLName      xml.Name `xml:"get_credentials"`
	CredentialID string   `xml:"credential_id,attr,omitempty"`
	Filter       string   `xml:"filter,attr,omitempty"`
	FiltID       string   `xml:"filt_id,attr,omitempty"`
	Details      bool     `xml:"details,attr,omitempty"`
	Scanners     bool     `xml:"scanners,attr,omitempty"`
	Trash        bool     `xml:"trash,attr,omitempty"`
	Targets      bool     `xml:"targets,attr,omitempty"`
	Format       string   `xml:"format,attr,omitempty"`
}

// GetCredentialsResponse represents a get_credentials command response from GMP.
type GetCredentialsResponse struct {
	XMLName         xml.Name          `xml:"get_credentials_response"`
	Status          string            `xml:"status,attr"`
	StatusText      string            `xml:"status_text,attr"`
	Credentials     []Credential      `xml:"credential"`
	Filters         CredentialFilters `xml:"filters"`
	Sort            CredentialSort    `xml:"sort"`
	CredentialsInfo CredentialInfo    `xml:"credentials"`
	CredentialCount CredentialCount   `xml:"credential_count"`
	Truncated       string            `xml:"truncated,omitempty"`
}

// Credential represents a <credential> element in the get_credentials response.
type Credential struct {
	ID               string                     `xml:"id,attr"`
	Owner            CredentialOwner            `xml:"owner"`
	Name             string                     `xml:"name"`
	AllowInsecure    string                     `xml:"allow_insecure"`
	Login            string                     `xml:"login"`
	Comment          string                     `xml:"comment"`
	CreationTime     string                     `xml:"creation_time"`
	ModificationTime string                     `xml:"modification_time"`
	Writable         string                     `xml:"writable"`
	InUse            string                     `xml:"in_use"`
	Permissions      CredentialPermissions      `xml:"permissions"`
	UserTags         *CredentialUserTags        `xml:"user_tags,omitempty"`
	Type             string                     `xml:"type"`
	FullType         string                     `xml:"full_type"`
	Formats          *CredentialFormats         `xml:"formats,omitempty"`
	AuthAlgorithm    *CredentialAuthAlgorithm   `xml:"auth_algorithm,omitempty"`
	Privacy          *CredentialPrivacy         `xml:"privacy,omitempty"`
	CertificateInfo  *CredentialCertificateInfo `xml:"certificate_info,omitempty"`
	Scanners         *CredentialScanners        `xml:"scanners,omitempty"`
	Targets          *CredentialTargets         `xml:"targets,omitempty"`
	PublicKey        string                     `xml:"public_key,omitempty"`
	Package          *CredentialPackage         `xml:"package,omitempty"`
	Certificate      string                     `xml:"certificate,omitempty"`
	KDC              string                     `xml:"kdc,omitempty"`
	KDCS             *CredentialKDCS            `xml:"kdcs,omitempty"`
	Realm            string                     `xml:"realm,omitempty"`
}

// CredentialOwner represents the owner element in a credential.
type CredentialOwner struct {
	Name string `xml:"name"`
}

// CredentialPermissions represents the permissions element in a credential.
type CredentialPermissions struct {
	Permission []CredentialPermission `xml:"permission"`
}

// CredentialPermission represents a permission element in credential permissions.
type CredentialPermission struct {
	Name string `xml:"name"`
}

// CredentialUserTags represents the user tags element in a credential.
type CredentialUserTags struct {
	Count int                 `xml:"count"`
	Tags  []CredentialUserTag `xml:"tag,omitempty"`
}

// CredentialUserTag represents a tag element in credential user tags.
type CredentialUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// CredentialFormats represents the formats element in a credential.
type CredentialFormats struct {
	Formats []string `xml:"format"`
}

// CredentialAuthAlgorithm represents the auth_algorithm element in a credential.
type CredentialAuthAlgorithm struct {
	Value string `xml:",chardata"`
}

// CredentialPrivacy represents the privacy element in a credential.
type CredentialPrivacy struct {
	Value string `xml:",chardata"`
}

// CredentialCertificateInfo represents the certificate_info element in a credential.
type CredentialCertificateInfo struct {
	Value string `xml:",chardata"`
}

// CredentialScanners represents the scanners element in a credential.
type CredentialScanners struct {
	Scanner []CredentialScanner `xml:"scanner"`
}

// CredentialScanner represents a scanner element in credential scanners.
type CredentialScanner struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// CredentialTargets represents the targets element in a credential.
type CredentialTargets struct {
	Target []CredentialTarget `xml:"target"`
}

// CredentialTarget represents a target element in credential targets.
type CredentialTarget struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// CredentialPackage represents the package element in a credential.
type CredentialPackage struct {
	Format string `xml:"format,attr"`
	Value  string `xml:",chardata"`
}

// CredentialKDCS represents the kdcs element in a credential.
type CredentialKDCS struct {
	KDC []string `xml:"kdc"`
}

// CredentialFilters represents the filters element in the response.
type CredentialFilters struct {
	ID       string                    `xml:"id,attr"`
	Term     string                    `xml:"term"`
	Name     string                    `xml:"name"`
	Keywords CredentialFiltersKeywords `xml:"keywords"`
}

// CredentialFiltersKeywords represents the keywords element in filters.
type CredentialFiltersKeywords struct {
	Keyword []CredentialFilterKeyword `xml:"keyword"`
}

// CredentialFilterKeyword represents a keyword element in filters keywords.
type CredentialFilterKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// CredentialSort represents the sort element in the response.
type CredentialSort struct {
	Value string              `xml:",chardata"`
	Field CredentialSortField `xml:"field"`
}

// CredentialSortField represents the field element in sort.
type CredentialSortField struct {
	Order string `xml:"order"`
}

// CredentialInfo represents the credentials element in the response.
type CredentialInfo struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// CredentialCount represents the credential count element in the response.
type CredentialCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
