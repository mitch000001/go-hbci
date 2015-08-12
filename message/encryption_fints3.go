package message

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

func NewFINTS3EncryptedPinTanMessage(clientSystemId string, keyName domain.KeyName, encryptedMessage []byte) *EncryptedFINTS3Message {
	e := &EncryptedFINTS3Message{
		EncryptionHeader: segment.NewFINTS3PinTanEncryptionHeaderSegment(clientSystemId, keyName),
		EncryptedMessage: &EncryptedMessage{
			EncryptedData: segment.NewEncryptedDataSegment(encryptedMessage),
		},
	}
	e.ClientMessage = NewBasicMessage(e)
	return e
}
