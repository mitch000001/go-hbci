package segment

import "github.com/mitch000001/go-hbci/element"

func NewSynchronisationSegment(modus int) *SynchronisationSegment {
	s := &SynchronisationSegment{
		SyncModus: element.NewNumber(modus, 1),
	}
	s.Segment = NewBasicSegment("HKSYN", 5, 2, s)
	return s
}

type SynchronisationSegment struct {
	Segment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.NumberDataElement
}

func (s *SynchronisationSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.SyncModus,
	}
}
