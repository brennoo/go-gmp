package commands

import "encoding/xml"

// DescribeAuth represents a describe_auth command request.
type DescribeAuth struct {
	XMLName xml.Name `xml:"describe_auth"`
}

// DescribeAuthResponse represents a describe_auth command response.
type DescribeAuthResponse struct {
	XMLName    xml.Name            `xml:"describe_auth_response"`
	Status     string              `xml:"status,attr"`
	StatusText string              `xml:"status_text,attr"`
	Groups     []DescribeAuthGroup `xml:"group"`
}

type DescribeAuthGroup struct {
	Name     string                `xml:"name,attr"`
	Settings []DescribeAuthSetting `xml:"auth_conf_setting"`
}

type DescribeAuthSetting struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}
