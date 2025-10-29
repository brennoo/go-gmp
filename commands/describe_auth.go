package commands

import "encoding/xml"

// DescribeAuth represents a describe_auth command request.
type DescribeAuth struct {
	XMLName xml.Name `xml:"describe_auth"`
}

// DescribeAuthResponse represents a describe_auth command response.
type DescribeAuthResponse struct {
	XMLName    xml.Name    `xml:"describe_auth_response"`
	Status     string      `xml:"status,attr"`
	StatusText string      `xml:"status_text,attr"`
	Groups     []AuthGroup `xml:"group"`
}

// GetStatus returns the status and status text from the response.
func (r *DescribeAuthResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// AuthGroup represents a group of authentication configuration settings.
type AuthGroup struct {
	Name     string        `xml:"name,attr"`
	Settings []AuthSetting `xml:"auth_conf_setting"`
}

// AuthSetting represents a single authentication configuration setting.
type AuthSetting struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}
