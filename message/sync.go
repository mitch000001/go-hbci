package message

import (
	"github.com/mitch000001/go-hbci/segment"
)

// NewSynchronisationMessage creates a new Message for Synchronizing client and server
func NewSynchronisationMessage(hbciVersion segment.HBCIVersion) *SynchronisationMessage {
	s := &SynchronisationMessage{
		hbciVersion: hbciVersion,
	}
	s.BasicMessage = NewBasicMessage(s)
	return s
}

// SynchronisationMessage serves the purpose of syncing the client and the server
type SynchronisationMessage struct {
	*BasicMessage
	Identification             *segment.IdentificationSegment
	ProcessingPreparation      *segment.ProcessingPreparationSegment
	TanRequest                 *segment.TanRequestSegment
	PublicSigningKeyRequest    *segment.PublicKeyRequestSegment
	PublicEncryptionKeyRequest *segment.PublicKeyRequestSegment
	PublicKeyRequest           *segment.PublicKeyRequestSegment
	Sync                       *segment.SynchronisationRequestSegment
	hbciVersion                segment.HBCIVersion
}

// HBCIVersion returns the HBCI version of the message
func (s *SynchronisationMessage) HBCIVersion() segment.HBCIVersion {
	return s.hbciVersion
}

// HBCISegments returns all segments of the message
func (s *SynchronisationMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		s.Identification,
		s.ProcessingPreparation,
		s.TanRequest,
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
		s.TanRequest,
		s.PublicSigningKeyRequest,
		s.PublicEncryptionKeyRequest,
		s.PublicKeyRequest,
		s.Sync,
	}
}
