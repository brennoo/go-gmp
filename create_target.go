package gmp

import "encoding/xml"

// CreateTargetCommand represents a create_target command request.
type CreateTargetCommand struct {
	XMLName              xml.Name                          `xml:"create_target"`
	Name                 string                            `xml:"name,omitempty"`
	Comment              string                            `xml:"comment,omitempty"`
	Copy                 string                            `xml:"copy,omitempty"`
	AssetHosts           *CreateTargetAssetHosts           `xml:"asset_hosts,omitempty"`
	Hosts                string                            `xml:"hosts,omitempty"`
	ExcludeHosts         string                            `xml:"exclude_hosts,omitempty"`
	SSHCredential        *CreateTargetSSHCredential        `xml:"ssh_credential,omitempty"`
	SMBCredential        *CreateTargetSMBCredential        `xml:"smb_credential,omitempty"`
	ESXICredential       *CreateTargetESXICredential       `xml:"esxi_credential,omitempty"`
	SNMPCredential       *CreateTargetSNMPCredential       `xml:"snmp_credential,omitempty"`
	SSHLSCCredential     *CreateTargetSSHLSCCredential     `xml:"ssh_lsc_credential,omitempty"`
	SMBLSCCredential     *CreateTargetSMBLSCCredential     `xml:"smb_lsc_credential,omitempty"`
	ESXILSCCredential    *CreateTargetESXILSCCredential    `xml:"esxi_lsc_credential,omitempty"`
	AliveTests           string                            `xml:"alive_tests,omitempty"`
	ReverseLookupOnly    bool                              `xml:"reverse_lookup_only,omitempty"`
	ReverseLookupUnify   bool                              `xml:"reverse_lookup_unify,omitempty"`
	PortRange            string                            `xml:"port_range,omitempty"`
	PortList             *CreateTargetPortList             `xml:"port_list,omitempty"`
	Krb5Credential       *CreateTargetKrb5Credential       `xml:"krb5_credential,omitempty"`
	SSHElevateCredential *CreateTargetSSHElevateCredential `xml:"ssh_elevate_credential,omitempty"`
	AllowSimultaneousIPs bool                              `xml:"allow_simultaneous_ips,omitempty"`
}

// CreateTargetAssetHosts represents the asset_hosts element for create_target.
type CreateTargetAssetHosts struct {
	Filter string `xml:"filter,attr,omitempty"`
}

// CreateTargetSSHCredential represents the ssh_credential element for create_target.
type CreateTargetSSHCredential struct {
	ID   string `xml:"id,attr,omitempty"`
	Port string `xml:"port,omitempty"`
}

// CreateTargetSMBCredential represents the smb_credential element for create_target.
type CreateTargetSMBCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetESXICredential represents the esxi_credential element for create_target.
type CreateTargetESXICredential struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetSNMPCredential represents the snmp_credential element for create_target.
type CreateTargetSNMPCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetSSHLSCCredential represents the ssh_lsc_credential element for create_target.
type CreateTargetSSHLSCCredential struct {
	ID   string `xml:"id,attr,omitempty"`
	Port string `xml:"port,omitempty"`
}

// CreateTargetSMBLSCCredential represents the smb_lsc_credential element for create_target.
type CreateTargetSMBLSCCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetESXILSCCredential represents the esxi_lsc_credential element for create_target.
type CreateTargetESXILSCCredential struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetPortList represents the port_list element for create_target.
type CreateTargetPortList struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetResponse represents a create_target command response.
type CreateTargetResponse struct {
	XMLName    xml.Name `xml:"create_target_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// CreateTargetKrb5Credential represents the krb5_credential element for create_target.
type CreateTargetKrb5Credential struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTargetSSHElevateCredential represents the ssh_elevate_credential element for create_target.
type CreateTargetSSHElevateCredential struct {
	ID string `xml:"id,attr,omitempty"`
}
