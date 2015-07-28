package message

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

type CryptoProvider interface {
	SetClientSystemID(clientSystemId string)
	Encrypt(message []byte) ([]byte, error)
	Decrypt(encryptedMessage []byte) ([]byte, error)
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
	e.ClientMessage = NewBasicMessage(e)
	return e
}

func NewEncryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment) *EncryptedMessage {
	e := &EncryptedMessage{}
	e.ClientMessage = NewBasicMessageWithHeaderAndEnd(header, end, e)
	return e
}

type EncryptedMessage struct {
	ClientMessage
	EncryptionHeader *segment.EncryptionHeaderSegment
	EncryptedData    *segment.EncryptedDataSegment
}

func (e *EncryptedMessage) HBCISegments() []segment.Segment {
	return []segment.Segment{
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

func NewDecryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, rawMessage []byte) (*DecryptedMessage, error) {
	segmentExtractor := segment.NewSegmentExtractor(rawMessage)
	_, err := segmentExtractor.Extract()
	if err != nil {
		return nil, fmt.Errorf("Malformed decrypted message bytes: %v", err)
	}
	messageAcknowledgementBytes := segmentExtractor.FindSegment("HIRMG")
	if messageAcknowledgementBytes == nil {
		return nil, fmt.Errorf("Malformed decrypted message: missing MessageAcknowledgement")
	}
	messageAcknowledgement := &segment.MessageAcknowledgement{}
	err = messageAcknowledgement.UnmarshalHBCI(messageAcknowledgementBytes)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling MessageAcknowledgement: %v", err)
	}
	acknowledgements := messageAcknowledgement.Acknowledgements()
	rawSegmentAcknowledgements := segmentExtractor.FindSegments("HIRMS")
	for _, segmentAcknowledgementBytes := range rawSegmentAcknowledgements {
		segmentAcknowledgement := &segment.SegmentAcknowledgement{}
		err = segmentAcknowledgement.UnmarshalHBCI(segmentAcknowledgementBytes)
		if err != nil {
			return nil, fmt.Errorf("Error while unmarshaling SegmentAcknowledgement: %v", err)
		}
		acknowledgements = append(acknowledgements, segmentAcknowledgement.Acknowledgements()...)
	}
	decryptedMessage := &DecryptedMessage{
		rawMessage:       rawMessage,
		acknowledgements: acknowledgements,
		segmentExtractor: segmentExtractor,
		unmarshaler:      NewUnmarshaler(rawMessage),
	}
	// TODO: set hbci message appropriate, if possible
	decryptedMessage.message = NewBasicMessageWithHeaderAndEnd(header, end, nil)
	return decryptedMessage, nil
}

type DecryptedMessage struct {
	rawMessage       []byte
	message          Message
	acknowledgements []domain.Acknowledgement
	segmentExtractor *segment.SegmentExtractor
	unmarshaler      *Unmarshaler
}

func (d *DecryptedMessage) MarshalHBCI() ([]byte, error) {
	return d.rawMessage, nil
}
func (d *DecryptedMessage) MessageHeader() *segment.MessageHeaderSegment {
	return d.message.MessageHeader()
}
func (d *DecryptedMessage) MessageEnd() *segment.MessageEndSegment {
	return d.message.MessageEnd()
}

func (d *DecryptedMessage) FindSegment(segmentID string) []byte {
	return d.segmentExtractor.FindSegment(segmentID)
}

func (d *DecryptedMessage) FindSegments(segmentID string) [][]byte {
	return d.segmentExtractor.FindSegments(segmentID)
}

func (d *DecryptedMessage) SegmentNumber(segmentID string) int {
	segmentBytes := d.segmentExtractor.FindSegment(segmentID)
	if segmentBytes == nil {
		return -1
	}
	elements, err := segment.ExtractElements(segmentBytes)
	if err != nil {
		return -1
	}
	if len(elements) < 1 {
		return -1
	}
	header := bytes.Split(elements[0], []byte(":"))
	numStr := header[1]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return -1
	}
	return num
}

func (d *DecryptedMessage) Acknowledgements() []domain.Acknowledgement {
	return d.acknowledgements
}

func NewPinTanCryptoProvider(key *domain.PinKey, clientSystemId string) *PinTanCryptoProvider {
	return &PinTanCryptoProvider{
		key:            key,
		clientSystemId: clientSystemId,
	}
}

type PinTanCryptoProvider struct {
	key            *domain.PinKey
	clientSystemId string
}

func (p *PinTanCryptoProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanCryptoProvider) Encrypt(message []byte) ([]byte, error) {
	return p.key.Encrypt(message)
}

func (p *PinTanCryptoProvider) Decrypt(encryptedMessage []byte) ([]byte, error) {
	return p.key.Decrypt(encryptedMessage)
}

func (p *PinTanCryptoProvider) EncryptWithInitialKeyName(message []byte) (*EncryptedMessage, error) {
	keyName := p.key.KeyName()
	keyName.SetInitial()
	encryptedBytes, _ := p.key.Encrypt(message)
	encryptedMessage := NewEncryptedPinTanMessage(p.clientSystemId, keyName, encryptedBytes)
	return encryptedMessage, nil
}

func (p *PinTanCryptoProvider) WriteEncryptionHeader(message *EncryptedMessage) {
	message.EncryptionHeader = segment.NewPinTanEncryptionHeaderSegment(p.clientSystemId, p.key.KeyName())
}
