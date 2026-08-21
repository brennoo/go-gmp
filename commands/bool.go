package commands

import "encoding/xml"

// GMPBool is a boolean that marshals to gvmd's lexical "boolean" wire
// representation: "1" for true, "0" for false. Unmarshal treats "0" as
// false and anything else — including the empty string — as true, per
// gvmd's schema note that empty text is true.
//
// Use *GMPBool (via Bool) for request fields so an explicit false can be
// sent instead of being dropped by omitempty.
type GMPBool bool

func (b GMPBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := "0"
	if b {
		v = "1"
	}
	return e.EncodeElement(v, start)
}

func (b *GMPBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	*b = GMPBool(s != "0")
	return nil
}

func (b GMPBool) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	v := "0"
	if b {
		v = "1"
	}
	return xml.Attr{Name: name, Value: v}, nil
}

func (b *GMPBool) UnmarshalXMLAttr(attr xml.Attr) error {
	*b = GMPBool(attr.Value != "0")
	return nil
}

// Bool returns a pointer to b as a GMPBool, for setting a *GMPBool request
// field to an explicit value (including false) rather than leaving it unset.
func Bool(b bool) *GMPBool {
	v := GMPBool(b)
	return &v
}
