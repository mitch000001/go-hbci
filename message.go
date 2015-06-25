package hbci

import (
	"bytes"
	"fmt"
	"reflect"
)

var bankSegments = map[string]Segment{
	"HIRMG": &MessageAcknowledgement{},
}

type Message interface {
	MarshalHBCI() ([]byte, error)
}

type HBCIMessage interface {
	HBCISegments() []Segment
}

type SignedHBCIMessage interface {
	HBCIMessage
	SignatureBeginSegment() *SignatureHeaderSegment
	SignatureEndSegment() *SignatureEndSegment
}

func newBasicMessage(message HBCIMessage) *basicMessage {
	b := &basicMessage{
		HBCIMessage: message,
	}
	return b
}

type basicMessage struct {
	Header *MessageHeaderSegment
	End    *MessageEndSegment
	HBCIMessage
	marshaledContent []byte
}

func (b *basicMessage) SetNumbers() {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	n := 0
	num := func() int {
		n += 1
		return n
	}
	b.Header.SetNumber(num())
	switch msg := b.HBCIMessage.(type) {
	case SignedHBCIMessage:
		msg.SignatureBeginSegment().SetNumber(num())
		for _, segment := range msg.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				segment.SetNumber(num())
			}
		}
		msg.SignatureEndSegment().SetNumber(num())
	default:
		for _, segment := range b.HBCIMessage.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				segment.SetNumber(num())
			}
		}
	}
	b.End.SetNumber(num())
}

func (b *basicMessage) SetSize() {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	var buffer bytes.Buffer
	buffer.WriteString(b.Header.String())
	switch msg := b.HBCIMessage.(type) {
	case SignedHBCIMessage:
		buffer.WriteString(msg.SignatureBeginSegment().String())
		for _, segment := range msg.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				buffer.WriteString(segment.String())
			}
		}
		buffer.WriteString(msg.SignatureEndSegment().String())
	default:
		for _, segment := range b.HBCIMessage.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				buffer.WriteString(segment.String())
			}
		}
	}
	buffer.WriteString(b.End.String())
	b.Header.SetSize(buffer.Len())
}

func (b *basicMessage) MarshalHBCI() ([]byte, error) {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	if len(b.marshaledContent) == 0 {
		var buffer bytes.Buffer
		buffer.WriteString(b.Header.String())
		switch msg := b.HBCIMessage.(type) {
		case SignedHBCIMessage:
			buffer.WriteString(msg.SignatureBeginSegment().String())
			for _, segment := range msg.HBCISegments() {
				if !reflect.ValueOf(segment).IsNil() {
					buffer.WriteString(segment.String())
				}
			}
			buffer.WriteString(msg.SignatureEndSegment().String())
		default:
			for _, segment := range b.HBCIMessage.HBCISegments() {
				if !reflect.ValueOf(segment).IsNil() {
					buffer.WriteString(segment.String())
				}
			}
		}
		buffer.WriteString(b.End.String())
		b.marshaledContent = buffer.Bytes()
	}
	return b.marshaledContent, nil
}

func (b *basicMessage) Encrypt(provider EncryptionProvider) (*EncryptedMessage, error) {
	var buffer bytes.Buffer
	switch msg := b.HBCIMessage.(type) {
	case SignedHBCIMessage:
		buffer.WriteString(msg.SignatureBeginSegment().String())
		for _, segment := range msg.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				buffer.WriteString(segment.String())
			}
		}
		buffer.WriteString(msg.SignatureEndSegment().String())
	default:
		for _, segment := range b.HBCIMessage.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				buffer.WriteString(segment.String())
			}
		}
	}
	encryptedMessage, err := provider.Encrypt(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	encryptedMessage.Header = b.Header
	encryptedMessage.End = b.End
	return encryptedMessage, nil
}

func newBasicSignedMessage(message HBCIMessage) *basicSignedMessage {
	b := &basicSignedMessage{
		HBCIMessage: message,
	}
	b.basicMessage = newBasicMessage(b)
	return b
}

type basicSignedMessage struct {
	*basicMessage
	SignatureBegin *SignatureHeaderSegment
	SignatureEnd   *SignatureEndSegment
	HBCIMessage
}

func (b *basicSignedMessage) SetSignatureBeginSegment(sigBegin *SignatureHeaderSegment) {
	b.SignatureBegin = sigBegin
}

func (b *basicSignedMessage) SetSignatureEndSegment(sigEnd *SignatureEndSegment) {
	b.SignatureEnd = sigEnd
}

func (b *basicSignedMessage) SignatureBeginSegment() *SignatureHeaderSegment {
	return b.SignatureBegin
}

func (b *basicSignedMessage) SignatureEndSegment() *SignatureEndSegment {
	return b.SignatureEnd
}

func (b *basicSignedMessage) SignRDHMessage(signingKey *RSAKey) error {
	if b.basicMessage == nil {
		panic(fmt.Errorf("basicMessage must be set"))
	}
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	provider := NewRDHSignatureProvider(signingKey)
	return provider.SignMessage(b)
}

func (b *basicSignedMessage) Sign(provider SignatureProvider) error {
	return provider.SignMessage(b)
}

type ClientMessage interface {
	Jobs() SegmentSequence
}

type BankMessage interface {
	DataSegments() SegmentSequence
}

type basicBankMessage struct {
	*basicSignedMessage
	BankMessage
	MessageAcknowledgements *MessageAcknowledgement
	SegmentAcknowledgements *SegmentAcknowledgement
}

func newBasicClientMessage(clientMessage ClientMessage) *basicClientMessage {
	b := &basicClientMessage{
		ClientMessage: clientMessage,
	}
	b.basicSignedMessage = newBasicSignedMessage(b)
	return b
}

type basicClientMessage struct {
	*basicSignedMessage
	ClientMessage
}

func (b *basicClientMessage) HBCISegments() []Segment {
	return b.ClientMessage.Jobs()
}

type SegmentSequence []Segment

var validHBCIVersions = []int{201, 210, 220}

func NewReferencingMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int, referencedMessage *ReferenceMessage) *MessageHeaderSegment {
	m := NewMessageHeaderSegment(size, hbciVersion, dialogId, number)
	m.Ref = referencedMessage
	return m
}

func NewMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	m := &MessageHeaderSegment{
		Size:        NewDigitDataElement(size, 12),
		HBCIVersion: NewNumberDataElement(hbciVersion, 3),
		DialogID:    NewIdentificationDataElement(dialogId),
		Number:      NewNumberDataElement(number, 4),
	}
	m.basicSegment = NewBasicSegment("HNHBK", 1, 3, m)
	return m
}

type MessageHeaderSegment struct {
	*basicSegment
	Size        *DigitDataElement
	HBCIVersion *NumberDataElement
	DialogID    *IdentificationDataElement
	Number      *NumberDataElement
	Ref         *ReferenceMessage
}

func (m *MessageHeaderSegment) elements() []DataElement {
	return []DataElement{
		m.Size,
		m.HBCIVersion,
		m.DialogID,
		m.Number,
		m.Ref,
	}
}

func (m *MessageHeaderSegment) SetSize(size int) {
	m.Size = NewDigitDataElement(size, 12)
}

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	end := &MessageEndSegment{
		Number: NewNumberDataElement(messageNumber, 4),
	}
	end.basicSegment = NewBasicSegment("HNHBS", segmentNumber, 1, end)
	return end
}

type MessageEndSegment struct {
	*basicSegment
	Number *NumberDataElement
}

func (m *MessageEndSegment) elements() []DataElement {
	return []DataElement{
		m.Number,
	}
}

func NewReferenceMessage(dialogId string, messageNumber int) *ReferenceMessage {
	r := &ReferenceMessage{
		DialogID:      NewIdentificationDataElement(dialogId),
		MessageNumber: NewNumberDataElement(messageNumber, 4),
	}
	r.elementGroup = NewDataElementGroup(ReferenceMessageDEG, 2, r)
	return r
}

type ReferenceMessage struct {
	*elementGroup
	DialogID      *IdentificationDataElement
	MessageNumber *NumberDataElement
}

func (r *ReferenceMessage) IsValid() bool {
	if r.DialogID == nil || r.MessageNumber == nil {
		return false
	} else {
		return r.elementGroup.IsValid()
	}
}

func (r *ReferenceMessage) Value() interface{} {
	return r
}

func (r *ReferenceMessage) groupDataElements() []DataElement {
	return []DataElement{
		r.DialogID,
		r.MessageNumber,
	}
}
