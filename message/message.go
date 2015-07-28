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
	SegmentNumber(segmentID string) int
}

type ClientMessage interface {
	Message
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
	SetNumbers()
	SetSignatureHeader(*segment.SignatureHeaderSegment)
	SetSignatureEnd(*segment.SignatureEndSegment)
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
	Header         *segment.MessageHeaderSegment
	End            *segment.MessageEndSegment
	SignatureBegin *segment.SignatureHeaderSegment
	SignatureEnd   *segment.SignatureEndSegment
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
	b.Header.SetNumber(num)
	if b.SignatureBegin != nil {
		b.SignatureBegin.SetNumber(num)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segment.SetNumber(num)
		}
	}
	if b.SignatureEnd != nil {
		b.SignatureEnd.SetNumber(num)
	}
	b.End.SetNumber(num)
}

func (b *BasicMessage) SetSize() {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	var buffer bytes.Buffer
	buffer.WriteString(b.Header.String())
	if b.SignatureBegin != nil {
		buffer.WriteString(b.SignatureBegin.String())
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			buffer.WriteString(segment.String())
		}
	}
	if b.SignatureEnd != nil {
		buffer.WriteString(b.SignatureEnd.String())
	}
	buffer.WriteString(b.End.String())
	b.Header.SetSize(buffer.Len())
}

func (b *BasicMessage) MarshalHBCI() ([]byte, error) {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	b.SetSize()
	if len(b.marshaledContent) == 0 {
		var buffer bytes.Buffer
		buffer.WriteString(b.Header.String())
		if b.SignatureBegin != nil {
			buffer.WriteString(b.SignatureBegin.String())
		}
		for _, segment := range b.HBCIMessage.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				buffer.WriteString(segment.String())
			}
		}
		if b.SignatureEnd != nil {
			buffer.WriteString(b.SignatureEnd.String())
		}
		buffer.WriteString(b.End.String())
		b.marshaledContent = buffer.Bytes()
	}
	return b.marshaledContent, nil
}

func (b *BasicMessage) Sign(provider SignatureProvider) (*BasicSignedMessage, error) {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	signedMessage := NewBasicSignedMessage(b)
	err := provider.SignMessage(signedMessage)
	if err != nil {
		return nil, err
	}
	return signedMessage, nil
}

func (b *BasicMessage) Encrypt(provider CryptoProvider) (*EncryptedMessage, error) {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	var buffer bytes.Buffer
	if b.SignatureBegin != nil {
		buffer.WriteString(b.SignatureBegin.String())
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			buffer.WriteString(segment.String())
		}
	}
	if b.SignatureEnd != nil {
		buffer.WriteString(b.SignatureEnd.String())
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

func (b *BasicMessage) SegmentNumber(segmentID string) int {
	idx := -1
	//for i, segment := range b.HBCIMessage.HBCISegments() {
	//if
	//}
	return idx
}

func NewBasicSignedMessage(message *BasicMessage) *BasicSignedMessage {
	b := &BasicSignedMessage{
		message: message,
	}
	return b
}

type BasicSignedMessage struct {
	message *BasicMessage
}

func (b *BasicSignedMessage) SetNumbers() {
	if b.message.SignatureBegin == nil || b.message.SignatureEnd == nil {
		panic(fmt.Errorf("Cannot call set Numbers when signature is not set"))
	}
	b.message.SetNumbers()
}

func (b *BasicSignedMessage) SetSignatureHeader(sigBegin *segment.SignatureHeaderSegment) {
	b.message.SignatureBegin = sigBegin
}

func (b *BasicSignedMessage) SetSignatureEnd(sigEnd *segment.SignatureEndSegment) {
	b.message.SignatureEnd = sigEnd
}

func (b *BasicSignedMessage) HBCISegments() []segment.Segment {
	return b.message.HBCISegments()
}

func (b *BasicSignedMessage) MarshalHBCI() ([]byte, error) {
	return b.message.MarshalHBCI()
}

func (b *BasicSignedMessage) Encrypt(provider CryptoProvider) (*EncryptedMessage, error) {
	return b.message.Encrypt(provider)
}

type bankMessage interface {
	dataSegments() []segment.Segment
}

type basicBankMessage struct {
	*BasicMessage
	bankMessage
	MessageAcknowledgements *segment.MessageAcknowledgement
	SegmentAcknowledgements *segment.SegmentAcknowledgement
}

type clientMessage interface {
	jobs() []segment.Segment
}

func NewBasicClientMessage(clientMessage clientMessage) *BasicClientMessage {
	b := &BasicClientMessage{
		clientMessage: clientMessage,
	}
	b.BasicMessage = NewBasicMessage(b)
	return b
}

type BasicClientMessage struct {
	*BasicMessage
	clientMessage
}

func (b *BasicClientMessage) HBCISegments() []segment.Segment {
	return b.clientMessage.jobs()
}
