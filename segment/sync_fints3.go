package segment

import (
	"strconv"

	"github.com/mitch000001/go-hbci/element"
)

func NewSynchronisationSegmentV3(modus int) *SynchronisationRequestSegment {
	s := &SynchronisationRequestV3{
		SyncModus: element.NewCode(strconv.Itoa(modus), 1, []string{"0", "1", "2"}),
	}
	s.Segment = NewBasicSegment(5, s)

	segment := &SynchronisationRequestSegment{
		Segment: s,
	}
	return segment
}

type SynchronisationRequestV3 struct {
	Segment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.CodeDataElement
}

func (s *SynchronisationRequestV3) Version() int         { return 2 }
func (s *SynchronisationRequestV3) ID() string           { return "HKSYN" }
func (s *SynchronisationRequestV3) referencedId() string { return "" }
func (s *SynchronisationRequestV3) sender() string       { return senderUser }

func (s *SynchronisationRequestV3) elements() []element.DataElement {
	return []element.DataElement{
		s.SyncModus,
	}
}
