package commands

import "encoding/xml"

// ModifyAuthCommand represents a modify_auth command request.
type ModifyAuthCommand struct {
	XMLName xml.Name        `xml:"modify_auth"`
	Group   ModifyAuthGroup `xml:"group"`
}

type ModifyAuthGroup struct {
	Name     string            `xml:"name,attr"`
	Settings []AuthConfSetting `xml:"auth_conf_setting"`
}

type AuthConfSetting struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

// ModifyAuthResponse represents a modify_auth command response.
type ModifyAuthResponse struct {
	XMLName    xml.Name `xml:"modify_auth_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
