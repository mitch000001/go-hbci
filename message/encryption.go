package message

import (
	"crypto/rand"

	"github.com/mitch000001/go-hbci/segment"
)

const encryptionInitializationVector = "\x00\x00\x00\x00\x00\x00\x00\x00"

// GenerateMessageKey generates a random key with 16 bytes
func GenerateMessageKey() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NewEncryptedMessage creates a new encrypted message
func NewEncryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, hbciVersion segment.HBCIVersion) *EncryptedMessage {
	e := &EncryptedMessage{
		hbciVersion: hbciVersion,
	}
	e.ClientMessage = NewBasicMessageWithHeaderAndEnd(header, end, e)
	return e
}

// EncryptedMessage represents an encrypted message
type EncryptedMessage struct {
	ClientMessage
	EncryptionHeader segment.EncryptionHeader
	EncryptedData    *segment.EncryptedDataSegment
	hbciVersion      segment.HBCIVersion
}

// HBCIVersion returns the HBCIVersion of this message
func (e *EncryptedMessage) HBCIVersion() segment.HBCIVersion {
	return e.hbciVersion
}

// HBCISegments returns all segments within the message
func (e *EncryptedMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		e.EncryptionHeader,
		e.EncryptedData,
	}
}

// Decrypt decrypts the message using the CryptoProvider
func (e *EncryptedMessage) Decrypt(provider CryptoProvider) (BankMessage, error) {
	decryptedMessageBytes, err := provider.Decrypt(e.EncryptedData.Data.Val())
	if err != nil {
		return nil, err
	}
	decryptedMessage, err := NewDecryptedMessage(e.MessageHeader(), e.MessageEnd(), decryptedMessageBytes)
	if err != nil {
		return nil, err
	}
	return decryptedMessage, nil
}
