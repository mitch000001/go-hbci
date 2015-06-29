package message

import "github.com/mitch000001/go-hbci/segment"

type InitialPublicKeyRenewalMessage struct {
	*BasicMessage
	Identification             *segment.IdentificationSegment
	PublicSigningKeyRequest    *segment.PublicKeyRequestSegment
	PublicEncryptionKeyRequest *segment.PublicKeyRequestSegment
}
