package commands

import "encoding/xml"

// GMPAtoiBool is a boolean that marshals like GMPBool but unmarshals per
// gvmd's "boolean_atoi" wire type, which runs the text through C's atoi(3)
// before testing for zero: unlike "boolean", the empty string (and any text
// with no leading digits) is false, not true.
//
// Use *GMPAtoiBool (via AtoiBool) for request fields so an explicit false
// can be sent instead of being dropped by omitempty.
type GMPAtoiBool bool

func (b GMPAtoiBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := "0"
	if b {
		v = "1"
	}
	return e.EncodeElement(v, start)
}

func (b *GMPAtoiBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	*b = GMPAtoiBool(atoi(s) != 0)
	return nil
}

func (b GMPAtoiBool) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	v := "0"
	if b {
		v = "1"
	}
	return xml.Attr{Name: name, Value: v}, nil
}

func (b *GMPAtoiBool) UnmarshalXMLAttr(attr xml.Attr) error {
	*b = GMPAtoiBool(atoi(attr.Value) != 0)
	return nil
}

// AtoiBool returns a pointer to b as a GMPAtoiBool, for setting a
// *GMPAtoiBool request field to an explicit value (including false) rather
// than leaving it unset.
func AtoiBool(b bool) *GMPAtoiBool {
	v := GMPAtoiBool(b)
	return &v
}

// atoi mimics C's atoi(3): skip leading whitespace, an optional sign, then
// digits; the first non-digit ends the scan. No valid digits yields 0.
func atoi(s string) int {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\v' || s[i] == '\f' || s[i] == '\r') {
		i++
	}
	neg := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		neg = s[i] == '-'
		i++
	}
	n := 0
	for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
	}
	if neg {
		return -n
	}
	return n
}
