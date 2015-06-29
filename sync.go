package hbci

import "github.com/mitch000001/go-hbci/dataelement"

func NewSynchronisationMessage() *SynchronisationMessage {
	s := new(SynchronisationMessage)
	s.basicClientMessage = newBasicClientMessage(s)
	return s
}

type SynchronisationMessage struct {
	*basicClientMessage
	Identification             *IdentificationSegment
	ProcessingPreparation      *ProcessingPreparationSegment
	PublicSigningKeyRequest    *PublicKeyRequestSegment
	PublicEncryptionKeyRequest *PublicKeyRequestSegment
	Sync                       *SynchronisationSegment
}

func (s *SynchronisationMessage) Jobs() SegmentSequence {
	return SegmentSequence{
		s.Identification,
		s.ProcessingPreparation,
		s.PublicSigningKeyRequest,
		s.PublicEncryptionKeyRequest,
		s.Sync,
	}
}

func NewSynchronisationSegment(modus int) *SynchronisationSegment {
	s := &SynchronisationSegment{
		SyncModus: dataelement.NewNumberDataElement(modus, 1),
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
	SyncModus *dataelement.NumberDataElement
}

func (s *SynchronisationSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		s.SyncModus,
	}
}
