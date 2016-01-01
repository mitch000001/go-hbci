package message

import (
	"crypto/rand"

	"github.com/mitch000001/go-hbci/segment"
)

const encryptionInitializationVector = "\x00\x00\x00\x00\x00\x00\x00\x00"

func GenerateMessageKey() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewEncryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, hbciVersion segment.HBCIVersion) *EncryptedMessage {
	e := &EncryptedMessage{
		hbciVersion: hbciVersion,
	}
	e.ClientMessage = NewBasicMessageWithHeaderAndEnd(header, end, e)
	return e
}

type EncryptedMessage struct {
	ClientMessage
	EncryptionHeader segment.EncryptionHeader
	EncryptedData    *segment.EncryptedDataSegment
	hbciVersion      segment.HBCIVersion
}

func (e *EncryptedMessage) HBCIVersion() segment.HBCIVersion {
	return e.hbciVersion
}

func (e *EncryptedMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		e.EncryptionHeader,
		e.EncryptedData,
	}
}

func (e *EncryptedMessage) Decrypt(provider CryptoProvider) (*DecryptedMessage, error) {
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
