package message

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

var bankSegments = map[string]segment.Segment{
	"HIRMG": &segment.MessageAcknowledgement{},
}

type Message interface {
	MessageHeader() *segment.MessageHeaderSegment
	MessageEnd() *segment.MessageEndSegment
	FindSegment(segmentID string) []byte
	FindSegments(segmentID string) [][]byte
}

type ClientMessage interface {
	Message
	SetNumbers()
	SetSize()
	MarshalHBCI() ([]byte, error)
	Encrypt(provider CryptoProvider) (*EncryptedMessage, error)
}

type BankMessage interface {
	Message
	Acknowledgements() []domain.Acknowledgement
}

type HBCIMessage interface {
	HBCISegments() []segment.Segment
}

type SignedHBCIMessage interface {
	HBCIMessage
	SignatureBeginSegment() *segment.SignatureHeaderSegment
	SignatureEndSegment() *segment.SignatureEndSegment
}

func NewBasicMessageWithHeaderAndEnd(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, message HBCIMessage) *BasicMessage {
	b := &BasicMessage{
		Header:      header,
		End:         end,
		HBCIMessage: message,
	}
	return b
}

func NewBasicMessage(message HBCIMessage) *BasicMessage {
	b := &BasicMessage{
		HBCIMessage: message,
	}
	return b
}

type BasicMessage struct {
	Header *segment.MessageHeaderSegment
	End    *segment.MessageEndSegment
	HBCIMessage
	marshaledContent []byte
}

func (b *BasicMessage) SetNumbers() {
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

func (b *BasicMessage) SetSize() {
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

func (b *BasicMessage) MarshalHBCI() ([]byte, error) {
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

func (b *BasicMessage) Encrypt(provider CryptoProvider) (*EncryptedMessage, error) {
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
	encryptionMessage := NewEncryptedMessage(b.Header, b.End)
	provider.WriteEncryptionHeader(encryptionMessage)
	encryptionMessage.EncryptedData = segment.NewEncryptedDataSegment(encryptedMessage)
	return encryptionMessage, nil
}

func (b *BasicMessage) MessageHeader() *segment.MessageHeaderSegment {
	return b.Header
}

func (b *BasicMessage) MessageEnd() *segment.MessageEndSegment {
	return b.End
}

func (b *BasicMessage) FindSegment(segmentID string) []byte {
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			return []byte(segment.String())
		}
	}
	return nil
}

func (b *BasicMessage) FindSegments(segmentID string) [][]byte {
	var segments [][]byte
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			segments = append(segments, []byte(segment.String()))
		}
	}
	return segments
}

func newBasicSignedMessage(message HBCIMessage) *basicSignedMessage {
	b := &basicSignedMessage{
		HBCIMessage: message,
	}
	b.BasicMessage = NewBasicMessage(b)
	return b
}

type basicSignedMessage struct {
	*BasicMessage
	SignatureBegin *segment.SignatureHeaderSegment
	SignatureEnd   *segment.SignatureEndSegment
	HBCIMessage
}

func (b *basicSignedMessage) SetSignatureBeginSegment(sigBegin *segment.SignatureHeaderSegment) {
	b.SignatureBegin = sigBegin
}

func (b *basicSignedMessage) SetSignatureEndSegment(sigEnd *segment.SignatureEndSegment) {
	b.SignatureEnd = sigEnd
}

func (b *basicSignedMessage) SignatureBeginSegment() *segment.SignatureHeaderSegment {
	return b.SignatureBegin
}

func (b *basicSignedMessage) SignatureEndSegment() *segment.SignatureEndSegment {
	return b.SignatureEnd
}

func (b *basicSignedMessage) Sign(provider SignatureProvider) error {
	if b.BasicMessage == nil {
		panic(fmt.Errorf("BasicMessage must be set"))
	}
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	return provider.SignMessage(b)
}

type clientMessage interface {
	jobs() segment.SegmentSequence
}

type bankMessage interface {
	dataSegments() segment.SegmentSequence
}

type basicBankMessage struct {
	*basicSignedMessage
	bankMessage
	MessageAcknowledgements *segment.MessageAcknowledgement
	SegmentAcknowledgements *segment.SegmentAcknowledgement
}

func NewBasicClientMessage(clientMessage clientMessage) *BasicClientMessage {
	b := &BasicClientMessage{
		clientMessage: clientMessage,
	}
	b.basicSignedMessage = newBasicSignedMessage(b)
	return b
}

type BasicClientMessage struct {
	*basicSignedMessage
	clientMessage
}

func (b *BasicClientMessage) HBCISegments() []segment.Segment {
	return b.clientMessage.jobs()
}
