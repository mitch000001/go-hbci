package message

import (
	"crypto/rand"
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

type CryptoProvider interface {
	SetClientSystemID(clientSystemId string)
	SetSecurityFunction(securityFn string)
	Encrypt(message []byte) ([]byte, error)
	Decrypt(encryptedMessage []byte) ([]byte, error)
	WriteEncryptionHeader(header segment.EncryptionHeader)
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

func NewDecryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, rawMessage []byte) (*DecryptedMessage, error) {
	unmarshaler := NewUnmarshaler(rawMessage)
	err := unmarshaler.Unmarshal()
	if err != nil {
		return nil, fmt.Errorf("Malformed decrypted message bytes: %v", err)
	}
	messageAcknowledgement := unmarshaler.SegmentById("HIRMG").(*segment.MessageAcknowledgement)
	messageAcknowledgement.SetReferencingMessage(header.ReferencingMessage())
	acknowledgements := messageAcknowledgement.Acknowledgements()
	for _, seg := range unmarshaler.SegmentsById("HIRMS") {
		if segmentAcknowledgement, ok := seg.(*segment.SegmentAcknowledgement); ok {
			acknowledgements = append(acknowledgements, segmentAcknowledgement.Acknowledgements()...)
		} else {
			panic(fmt.Errorf("Error while unmarshaling segments"))
		}
	}
	version, ok := segment.SupportedHBCIVersions[header.HBCIVersion.Val()]
	if !ok {
		return nil, fmt.Errorf("Unknown HBCI version: %d", header.HBCIVersion.Val())
	}
	decryptedMessage := &DecryptedMessage{
		rawMessage:       rawMessage,
		acknowledgements: acknowledgements,
		unmarshaler:      unmarshaler,
		hbciVersion:      version,
	}
	// TODO: set hbci message appropriate, if possible
	decryptedMessage.message = NewBasicMessageWithHeaderAndEnd(header, end, decryptedMessage)
	return decryptedMessage, nil
}

type DecryptedMessage struct {
	rawMessage       []byte
	message          Message
	acknowledgements []domain.Acknowledgement
	unmarshaler      *Unmarshaler
	hbciVersion      segment.HBCIVersion
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

func (d *DecryptedMessage) FindMarshaledSegment(segmentID string) []byte {
	return d.unmarshaler.MarshaledSegmentById(segmentID)
}

func (d *DecryptedMessage) FindMarshaledSegments(segmentID string) [][]byte {
	return d.unmarshaler.MarshaledSegmentsById(segmentID)
}

func (d *DecryptedMessage) MarshaledSegments() [][]byte {
	return d.unmarshaler.MarshaledSegments()
}

func (d *DecryptedMessage) FindSegment(segmentID string) segment.Segment {
	return d.unmarshaler.SegmentById(segmentID)
}

func (d *DecryptedMessage) FindSegments(segmentID string) []segment.Segment {
	return d.unmarshaler.SegmentsById(segmentID)
}

func (d *DecryptedMessage) SegmentNumber(segmentID string) int {
	seg := d.unmarshaler.MarshaledSegmentById(segmentID)
	if len(seg) == 0 {
		return -1
	}
	elements, err := segment.ExtractElements(seg)
	if err != nil {
		return -1
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return -1
	}
	return header.Number.Val()
}

func (d *DecryptedMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{}
}

func (d *DecryptedMessage) HBCIVersion() segment.HBCIVersion {
	return d.hbciVersion
}

func (d *DecryptedMessage) Acknowledgements() []domain.Acknowledgement {
	return d.acknowledgements
}

func NewPinTanCryptoProvider(key *domain.PinKey, clientSystemId string) *PinTanCryptoProvider {
	return &PinTanCryptoProvider{
		key:            key,
		clientSystemId: clientSystemId,
		securityFn:     "999",
	}
}

type PinTanCryptoProvider struct {
	key            *domain.PinKey
	clientSystemId string
	securityFn     string
}

func (p *PinTanCryptoProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanCryptoProvider) SetSecurityFunction(securityFn string) {
	p.securityFn = securityFn
}

func (p *PinTanCryptoProvider) Encrypt(message []byte) ([]byte, error) {
	if p.key.Pin() == "" {
		return nil, fmt.Errorf("Malformed PIN")
	}
	return p.key.Encrypt(message)
}

func (p *PinTanCryptoProvider) Decrypt(encryptedMessage []byte) ([]byte, error) {
	return p.key.Decrypt(encryptedMessage)
}

func (p *PinTanCryptoProvider) WriteEncryptionHeader(header segment.EncryptionHeader) {
	header.SetClientSystemID(p.clientSystemId)
	header.SetSecurityProfile(p.securityFn)
	header.SetEncryptionKeyName(p.key.KeyName())
	header.SetEncryptionAlgorithm(element.NewPinTanEncryptionAlgorithm())
}
