package gmp

type CertificateInfo struct {
	//	XMLName        xml.Name `xml:"certificate_info"`
	TimeStatus        string `xml:"time_status"`
	ActivationTime    string `xml:"activation_time"`
	ExpirationTime    string `xml:"expiration_time"`
	Issuer            string `xml:"issuer"`
	MD5Fingerprint    string `xml:"md5_fingerprint"`
	SHA256Fingerprint string `xml:"sha256_fingerprint,omitempty"`
	Subject           string `xml:"subject,omitempty"`
	Serial            string `xml:"serial,omitempty"`
}
