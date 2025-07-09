package commands

import "encoding/xml"

// ModifyTargetCommand represents a modify_target command request.
type ModifyTargetCommand struct {
	XMLName              xml.Name                          `xml:"modify_target"`
	TargetID             string                            `xml:"target_id,attr"`
	Name                 string                            `xml:"name,omitempty"`
	Comment              string                            `xml:"comment,omitempty"`
	Hosts                string                            `xml:"hosts,omitempty"`
	HostsOrdering        string                            `xml:"hosts_ordering,omitempty"`
	ExcludeHosts         string                            `xml:"exclude_hosts,omitempty"`
	SSHCredential        *ModifyTargetSSHCredential        `xml:"ssh_credential,omitempty"`
	SSHElevateCredential *ModifyTargetSSHElevateCredential `xml:"ssh_elevate_credential,omitempty"`
	Krb5Credential       *ModifyTargetKrb5Credential       `xml:"krb5_credential,omitempty"`
	SMBCredential        *ModifyTargetSMBCredential        `xml:"smb_credential,omitempty"`
	ESXICredential       *ModifyTargetESXICredential       `xml:"esxi_credential,omitempty"`
	SNMPCredential       *ModifyTargetSNMPCredential       `xml:"snmp_credential,omitempty"`
	SSHLSCCredential     *ModifyTargetSSHLSCCredential     `xml:"ssh_lsc_credential,omitempty"`
	SMBLSCCredential     *ModifyTargetSMBLSCCredential     `xml:"smb_lsc_credential,omitempty"`
	ESXILSCCredential    *ModifyTargetESXILSCCredential    `xml:"esxi_lsc_credential,omitempty"`
	PortList             *ModifyTargetPortList             `xml:"port_list,omitempty"`
	AliveTests           string                            `xml:"alive_tests,omitempty"`
	ReverseLookupOnly    bool                              `xml:"reverse_lookup_only,omitempty"`
	ReverseLookupUnify   bool                              `xml:"reverse_lookup_unify,omitempty"`
	AllowSimultaneousIPs bool                              `xml:"allow_simultaneous_ips,omitempty"`
}

type ModifyTargetSSHCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetSSHElevateCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetKrb5Credential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetSMBCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetESXICredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetSNMPCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetSSHLSCCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetSMBLSCCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTargetESXILSCCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

// ModifyTargetPortList represents the port_list element for modify_target.
type ModifyTargetPortList struct {
	ID string `xml:"id,attr,omitempty"`
}

// ModifyTargetResponse represents a modify_target command response.
type ModifyTargetResponse struct {
	XMLName    xml.Name `xml:"modify_target_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
