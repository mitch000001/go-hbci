package message

import (
	"crypto/rand"
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

type EncryptionProvider interface {
	SetClientSystemID(clientSystemId string)
	Encrypt(message []byte) ([]byte, error)
	WriteEncryptionHeader(message *EncryptedMessage)
	EncryptWithInitialKeyName(message []byte) (*EncryptedMessage, error)
}

const encryptionInitializationVector = "\x00\x00\x00\x00\x00\x00\x00\x00"

func GenerateMessageKey() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewEncryptedPinTanMessage(clientSystemId string, keyName domain.KeyName, encryptedMessage []byte) *EncryptedMessage {
	e := &EncryptedMessage{
		EncryptionHeader: segment.NewPinTanEncryptionHeaderSegment(clientSystemId, keyName),
		EncryptedData:    segment.NewEncryptedDataSegment(encryptedMessage),
	}
	e.Message = NewBasicMessage(e)
	return e
}

func NewEncryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment) *EncryptedMessage {
	e := &EncryptedMessage{}
	e.Message = NewBasicMessageWithHeaderAndEnd(header, end, e)
	return e
}

type EncryptedMessage struct {
	Message
	EncryptionHeader *segment.EncryptionHeaderSegment
	EncryptedData    *segment.EncryptedDataSegment
}

func (e *EncryptedMessage) HBCISegments() []segment.Segment {
	return []segment.Segment{
		e.EncryptionHeader,
		e.EncryptedData,
	}
}

func (e *EncryptedMessage) SetNumbers() {
	panic(fmt.Errorf("SetNumbers: Operation not allowed on encrypted messages"))
}

func NewPinTanEncryptionProvider(key *domain.PinKey, clientSystemId string) *PinTanEncryptionProvider {
	return &PinTanEncryptionProvider{
		key:            key,
		clientSystemId: clientSystemId,
	}
}

type PinTanEncryptionProvider struct {
	key            *domain.PinKey
	clientSystemId string
}

func (p *PinTanEncryptionProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanEncryptionProvider) Encrypt(message []byte) ([]byte, error) {
	return p.key.Encrypt(message)
}

func (p *PinTanEncryptionProvider) EncryptWithInitialKeyName(message []byte) (*EncryptedMessage, error) {
	keyName := p.key.KeyName()
	keyName.SetInitial()
	encryptedBytes, _ := p.key.Encrypt(message)
	encryptedMessage := NewEncryptedPinTanMessage(p.clientSystemId, keyName, encryptedBytes)
	return encryptedMessage, nil
}

func (p *PinTanEncryptionProvider) WriteEncryptionHeader(message *EncryptedMessage) {
	message.EncryptionHeader = segment.NewPinTanEncryptionHeaderSegment(p.clientSystemId, p.key.KeyName())
}
