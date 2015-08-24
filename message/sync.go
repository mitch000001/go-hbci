package message

import (
	"github.com/mitch000001/go-hbci/segment"
)

func NewSynchronisationMessage(hbciVersion segment.HBCIVersion) *SynchronisationMessage {
	s := &SynchronisationMessage{
		hbciVersion: hbciVersion,
	}
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
	hbciVersion                segment.HBCIVersion
}

func (s *SynchronisationMessage) HBCIVersion() segment.HBCIVersion {
	return s.hbciVersion
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
