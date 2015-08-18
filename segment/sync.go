package segment

import "github.com/mitch000001/go-hbci/element"

func NewSynchronisationSegmentV2(modus int) *SynchronisationRequestSegment {
	s := &SynchronisationRequestV2{
		SyncModus: element.NewNumber(modus, 1),
	}
	s.ClientSegment = NewBasicSegment(5, s)

	segment := &SynchronisationRequestSegment{
		Segment: s,
	}
	return segment
}

type SynchronisationRequestSegment struct {
	Segment
}

type SynchronisationRequestV2 struct {
	ClientSegment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.NumberDataElement
}

func (s *SynchronisationRequestV2) Version() int         { return 2 }
func (s *SynchronisationRequestV2) ID() string           { return "HKSYN" }
func (s *SynchronisationRequestV2) referencedId() string { return "" }
func (s *SynchronisationRequestV2) sender() string       { return senderUser }

func (s *SynchronisationRequestV2) elements() []element.DataElement {
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
