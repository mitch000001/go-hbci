package message

import (
	"github.com/mitch000001/go-hbci/segment"
)

func NewSynchronisationMessage() *SynchronisationMessage {
	s := new(SynchronisationMessage)
	s.BasicMessage = NewBasicMessage(s)
	return s
}

type SynchronisationMessage struct {
	*BasicMessage
	Identification             *segment.IdentificationSegment
	ProcessingPreparation      *segment.ProcessingPreparationSegment
	PublicSigningKeyRequest    *segment.PublicKeyRequestSegment
	PublicEncryptionKeyRequest *segment.PublicKeyRequestSegment
	PublicKeyRequest           *segment.PublicKeyRequestSegment
	Sync                       *segment.SynchronisationRequestSegment
}

func (s *SynchronisationMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		s.Identification,
		s.ProcessingPreparation,
		s.PublicSigningKeyRequest,
		s.PublicEncryptionKeyRequest,
		s.PublicKeyRequest,
		s.Sync,
	}
}

func (s *SynchronisationMessage) jobs() []segment.ClientSegment {
	return []segment.ClientSegment{
		s.Identification,
		s.ProcessingPreparation,
		s.PublicSigningKeyRequest,
		s.PublicEncryptionKeyRequest,
		s.PublicKeyRequest,
		s.Sync,
	}
}
