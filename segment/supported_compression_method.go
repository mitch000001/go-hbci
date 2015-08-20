package segment

import "github.com/mitch000001/go-hbci/element"

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment CompressionMethodSegment

type CompressionMethodSegment struct {
	Segment
	SupportedCompressionMethods *element.SupportedCompressionMethodsDataElement
}

func (c *CompressionMethodSegment) version() int         { return 1 }
func (c *CompressionMethodSegment) id() string           { return "HIKPV" }
func (c *CompressionMethodSegment) referencedId() string { return "HKVVB" }
func (c *CompressionMethodSegment) sender() string       { return senderBank }

func (c *CompressionMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.SupportedCompressionMethods,
	}
}
