package commands

import "encoding/xml"

// CreateTarget is a create_target command request.
type CreateTarget struct {
	XMLName              xml.Name          `xml:"create_target"`
	Name                 string            `xml:"name,omitempty"`
	Comment              string            `xml:"comment,omitempty"`
	Copy                 string            `xml:"copy,omitempty"`
	AssetHosts           *AssetHosts       `xml:"asset_hosts,omitempty"`
	Hosts                string            `xml:"hosts,omitempty"`
	ExcludeHosts         string            `xml:"exclude_hosts,omitempty"`
	SSHCredential        *TargetCredential `xml:"ssh_credential,omitempty"`
	SMBCredential        *TargetCredential `xml:"smb_credential,omitempty"`
	ESXICredential       *TargetCredential `xml:"esxi_credential,omitempty"`
	SNMPCredential       *TargetCredential `xml:"snmp_credential,omitempty"`
	SSHLSCCredential     *TargetCredential `xml:"ssh_lsc_credential,omitempty"`
	SMBLSCCredential     *TargetCredential `xml:"smb_lsc_credential,omitempty"`
	ESXILSCCredential    *TargetCredential `xml:"esxi_lsc_credential,omitempty"`
	AliveTests           string            `xml:"alive_tests,omitempty"`
	ReverseLookupOnly    bool              `xml:"reverse_lookup_only,omitempty"`
	ReverseLookupUnify   bool              `xml:"reverse_lookup_unify,omitempty"`
	PortRange            string            `xml:"port_range,omitempty"`
	PortList             *TargetPortList   `xml:"port_list,omitempty"`
	Krb5Credential       *TargetCredential `xml:"krb5_credential,omitempty"`
	SSHElevateCredential *TargetCredential `xml:"ssh_elevate_credential,omitempty"`
	AllowSimultaneousIPs bool              `xml:"allow_simultaneous_ips,omitempty"`
}

// CreateTargetResponse is a create_target command response.
type CreateTargetResponse struct {
	XMLName    xml.Name `xml:"create_target_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// AssetHosts is the <asset_hosts> element for create_target.
type AssetHosts struct {
	Filter string `xml:"filter,attr,omitempty"`
}
