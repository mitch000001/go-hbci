package segment

import "github.com/mitch000001/go-hbci/element"

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment SecurityMethodSegment

type SecurityMethodSegment struct {
	Segment
	MixAllowed       *element.BooleanDataElement
	SupportedMethods *element.SupportedSecurityMethodDataElement
}

func (s *SecurityMethodSegment) Version() int         { return 2 }
func (s *SecurityMethodSegment) ID() string           { return "HISHV" }
func (s *SecurityMethodSegment) referencedId() string { return "HKVVB" }
func (s *SecurityMethodSegment) sender() string       { return senderBank }

func (s *SecurityMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.MixAllowed,
		s.SupportedMethods,
	}
}
