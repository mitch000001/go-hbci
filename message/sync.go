package message

import (
	"github.com/mitch000001/go-hbci/segment"
)

func NewSynchronisationMessage() *SynchronisationMessage {
	s := new(SynchronisationMessage)
	s.BasicClientMessage = NewBasicClientMessage(s)
	return s
}

type SynchronisationMessage struct {
	*BasicClientMessage
	Identification             *segment.IdentificationSegment
	ProcessingPreparation      *segment.ProcessingPreparationSegment
	PublicSigningKeyRequest    *segment.PublicKeyRequestSegment
	PublicEncryptionKeyRequest *segment.PublicKeyRequestSegment
	Sync                       *segment.SynchronisationRequestSegment
}

func (s *SynchronisationMessage) jobs() []segment.Segment {
	return []segment.Segment{
		s.Identification,
		s.ProcessingPreparation,
		s.PublicSigningKeyRequest,
		s.PublicEncryptionKeyRequest,
		s.Sync,
	}
}
