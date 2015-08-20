package segment

import "github.com/mitch000001/go-hbci/element"

func NewSynchronisationSegment(modus int) *SynchronisationRequestSegment {
	s := &SynchronisationRequestSegment{
		SyncModus: element.NewNumber(modus, 1),
	}
	s.Segment = NewBasicSegment(5, s)
	return s
}

type SynchronisationRequestSegment struct {
	Segment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.NumberDataElement
}

func (s *SynchronisationRequestSegment) Version() int         { return 2 }
func (s *SynchronisationRequestSegment) ID() string           { return "HKSYN" }
func (s *SynchronisationRequestSegment) referencedId() string { return "" }
func (s *SynchronisationRequestSegment) sender() string       { return senderUser }

func (s *SynchronisationRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.SyncModus,
	}
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment SynchronisationResponseSegment

type SynchronisationResponseSegment struct {
	Segment
	ClientSystemID *element.IdentificationDataElement
	MessageNumber  *element.NumberDataElement
	SignatureID    *element.NumberDataElement
}

func (s *SynchronisationResponseSegment) Version() int         { return 3 }
func (s *SynchronisationResponseSegment) ID() string           { return "HISYN" }
func (s *SynchronisationResponseSegment) referencedId() string { return "HKSYN" }
func (s *SynchronisationResponseSegment) sender() string       { return senderBank }

func (s *SynchronisationResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.ClientSystemID,
		s.MessageNumber,
		s.SignatureID,
	}
}
