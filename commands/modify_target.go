package commands

import "encoding/xml"

// ModifyTarget represents a modify_target command request.
type ModifyTarget struct {
	XMLName              xml.Name          `xml:"modify_target"`
	TargetID             string            `xml:"target_id,attr"`
	Name                 string            `xml:"name,omitempty"`
	Comment              string            `xml:"comment,omitempty"`
	Hosts                string            `xml:"hosts,omitempty"`
	HostsOrdering        string            `xml:"hosts_ordering,omitempty"`
	ExcludeHosts         string            `xml:"exclude_hosts,omitempty"`
	SSHCredential        *TargetCredential `xml:"ssh_credential,omitempty"`
	SSHElevateCredential *TargetCredential `xml:"ssh_elevate_credential,omitempty"`
	Krb5Credential       *TargetCredential `xml:"krb5_credential,omitempty"`
	SMBCredential        *TargetCredential `xml:"smb_credential,omitempty"`
	ESXICredential       *TargetCredential `xml:"esxi_credential,omitempty"`
	SNMPCredential       *TargetCredential `xml:"snmp_credential,omitempty"`
	SSHLSCCredential     *TargetCredential `xml:"ssh_lsc_credential,omitempty"`
	SMBLSCCredential     *TargetCredential `xml:"smb_lsc_credential,omitempty"`
	ESXILSCCredential    *TargetCredential `xml:"esxi_lsc_credential,omitempty"`
	PortList             *TargetPortList   `xml:"port_list,omitempty"`
	AliveTests           string            `xml:"alive_tests,omitempty"`
	ReverseLookupOnly    string            `xml:"reverse_lookup_only,omitempty"`
	ReverseLookupUnify   string            `xml:"reverse_lookup_unify,omitempty"`
	AllowSimultaneousIPs string            `xml:"allow_simultaneous_ips,omitempty"`
}

// ModifyTargetResponse represents a modify_target command response.
type ModifyTargetResponse struct {
	XMLName    xml.Name `xml:"modify_target_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
