package message

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

// Message represents a basic message
type Message interface {
	MessageHeader() *segment.MessageHeaderSegment
	MessageEnd() *segment.MessageEndSegment
	FindMarshaledSegment(segmentID string) []byte
	FindMarshaledSegments(segmentID string) [][]byte
	FindSegment(segmentID string) segment.Segment
	FindSegments(segmentID string) []segment.Segment
	SegmentPosition(segmentID string) int
}

// ClientMessage represents a message composed by the client
type ClientMessage interface {
	Message
	MarshalHBCI() ([]byte, error)
	Encrypt(provider CryptoProvider) (*EncryptedMessage, error)
	SetMessageNumber(messageNumber int)
}

// BankMessage represents a message composed by the bank
type BankMessage interface {
	Message
	Segments() []segment.Segment
	Acknowledgements() map[int]domain.Acknowledgement
	SupportedSegments() []segment.VersionedSegment
}

// HBCIMessage represents a basic set of message for introspecting HBCI messages
type HBCIMessage interface {
	HBCIVersion() segment.HBCIVersion
	HBCISegments() []segment.ClientSegment
}

// A SignedHBCIMessage represents a HBCI message that can be signed
type SignedHBCIMessage interface {
	HBCIMessage
	SetSegmentPositions()
	SetSignatureHeader(*segment.SignatureHeaderSegment)
	SetSignatureEnd(*segment.SignatureEndSegment)
}

// NewHBCIMessage creates a new hbci message for the given version and adds the
// segments as message body
func NewHBCIMessage(hbciVersion segment.HBCIVersion, segments ...segment.ClientSegment) HBCIMessage {
	return &hbciMessage{hbciSegments: segments, hbciVersion: hbciVersion}
}

type hbciMessage struct {
	hbciSegments []segment.ClientSegment
	hbciVersion  segment.HBCIVersion
}

func (h *hbciMessage) HBCIVersion() segment.HBCIVersion {
	return h.hbciVersion
}

func (h *hbciMessage) HBCISegments() []segment.ClientSegment {
	return h.hbciSegments
}

type hbciSegmentClientMessage []segment.ClientSegment

func (h hbciSegmentClientMessage) jobs() []segment.ClientSegment {
	return h
}

// NewBasicMessageWithHeaderAndEnd creates a new BasicMessage with the provided
// header and end and embodies message
func NewBasicMessageWithHeaderAndEnd(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, message HBCIMessage) *BasicMessage {
	b := &BasicMessage{
		Header:      header,
		End:         end,
		HBCIMessage: message,
		hbciVersion: message.HBCIVersion(),
	}
	return b
}

// NewBasicMessage creates a new BasicMessage from the HBCIMessage
func NewBasicMessage(message HBCIMessage) *BasicMessage {
	b := &BasicMessage{
		HBCIMessage: message,
		hbciVersion: message.HBCIVersion(),
	}
	return b
}

// BasicMessage represents a basic HBCI message with all necessary components
// such as MessageHeader, MessageEnd, Signature and message body
type BasicMessage struct {
	Header         *segment.MessageHeaderSegment
	End            *segment.MessageEndSegment
	SignatureBegin *segment.SignatureHeaderSegment
	SignatureEnd   *segment.SignatureEndSegment
	HBCIMessage
	hbciVersion      segment.HBCIVersion
	marshaledContent []byte
}

// SetSegmentPositions sets the message number on every segment within the message
func (b *BasicMessage) SetSegmentPositions() {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	n := 0
	num := func() int {
		n++
		return n
	}
	b.Header.SetPosition(num)
	if b.SignatureBegin != nil {
		b.SignatureBegin.SetPosition(num)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segment.SetPosition(num)
		}
	}
	if b.SignatureEnd != nil {
		b.SignatureEnd.SetPosition(num)
	}
	b.End.SetPosition(num)
}

// SetSize writes the size of the marshaled message into the message header
func (b *BasicMessage) SetSize() error {
	if b.HBCIMessage == nil {
		return fmt.Errorf("HBCIMessage must be set")
	}
	var buffer bytes.Buffer
	headerBytes, err := b.Header.MarshalHBCI()
	if err != nil {
		return err
	}
	buffer.Write(headerBytes)
	if b.SignatureBegin != nil {
		sigBytes, err := b.SignatureBegin.MarshalHBCI()
		if err != nil {
			return err
		}
		buffer.Write(sigBytes)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segBytes, err := segment.MarshalHBCI()
			if err != nil {
				return err
			}
			buffer.Write(segBytes)
		}
	}
	if b.SignatureEnd != nil {
		sigEndBytes, err := b.SignatureEnd.MarshalHBCI()
		if err != nil {
			return err
		}
		buffer.Write(sigEndBytes)
	}
	endBytes, err := b.End.MarshalHBCI()
	if err != nil {
		return err
	}
	buffer.Write(endBytes)
	b.Header.SetSize(buffer.Len())
	return nil
}

// SetMessageNumber sets the message number in the MessageHeader
func (b *BasicMessage) SetMessageNumber(messageNumber int) {
	b.Header.SetMessageNumber(messageNumber)
}

// MarshalHBCI marshals b into HBCI wire format
func (b *BasicMessage) MarshalHBCI() ([]byte, error) {
	if b.HBCIMessage == nil {
		return nil, fmt.Errorf("HBCIMessage must be set")
	}
	err := b.SetSize()
	if err != nil {
		return nil, err
	}
	if len(b.marshaledContent) == 0 {
		var buffer bytes.Buffer
		headerBytes, err := b.Header.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		buffer.Write(headerBytes)
		if b.SignatureBegin != nil {
			sigBytes, err := b.SignatureBegin.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			buffer.Write(sigBytes)
		}
		for _, segment := range b.HBCIMessage.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				segBytes, err := segment.MarshalHBCI()
				if err != nil {
					return nil, err
				}
				buffer.Write(segBytes)
			}
		}
		if b.SignatureEnd != nil {
			sigEndBytes, err := b.SignatureEnd.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			buffer.Write(sigEndBytes)
		}
		endBytes, err := b.End.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		buffer.Write(endBytes)
		b.marshaledContent = buffer.Bytes()
	}
	return b.marshaledContent, nil
}

// Sign signs b using the SignatureProvider
func (b *BasicMessage) Sign(provider SignatureProvider) (*BasicSignedMessage, error) {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	// TODO: fix only PinTan segments!!!
	b.SignatureBegin = b.hbciVersion.SignatureHeader()
	provider.WriteSignatureHeader(b.SignatureBegin)
	b.SignatureEnd = b.hbciVersion.SignatureEnd()
	b.SetSegmentPositions()
	var buffer bytes.Buffer
	buffer.WriteString(b.SignatureBegin.String())
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			buffer.WriteString(segment.String())
		}
	}
	sig, err := provider.Sign(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	provider.WriteSignature(b.SignatureEnd, sig)
	signedMessage := NewBasicSignedMessage(b)
	return signedMessage, nil
}

// Encrypt encrypts the message using the CryptoProvider
func (b *BasicMessage) Encrypt(provider CryptoProvider) (*EncryptedMessage, error) {
	if b.HBCIMessage == nil {
		return nil, fmt.Errorf("HBCIMessage must be set")
	}
	var messageBytes []byte
	if b.SignatureBegin != nil {
		sigBytes, err := b.SignatureBegin.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		messageBytes = append(messageBytes, sigBytes...)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segBytes, err := segment.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			messageBytes = append(messageBytes, segBytes...)
		}
	}
	if b.SignatureEnd != nil {
		sigEndBytes, err := b.SignatureEnd.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		messageBytes = append(messageBytes, sigEndBytes...)
	}
	encryptedMessage, err := provider.Encrypt(messageBytes)
	if err != nil {
		return nil, err
	}
	encryptionMessage := NewEncryptedMessage(b.Header, b.End, b.hbciVersion)
	encryptionMessage.EncryptionHeader = b.hbciVersion.PinTanEncryptionHeader("", domain.KeyName{})
	provider.WriteEncryptionHeader(encryptionMessage.EncryptionHeader)
	encryptionMessage.EncryptedData = segment.NewEncryptedDataSegment(encryptedMessage)
	return encryptionMessage, nil
}

// MessageHeader returns the MessageHeader
func (b *BasicMessage) MessageHeader() *segment.MessageHeaderSegment {
	return b.Header
}

// MessageEnd returns the MessageEnd
func (b *BasicMessage) MessageEnd() *segment.MessageEndSegment {
	return b.End
}

// FindSegment returns the segment found for the segmentID, or nil if not
// found. The first matching segment will be returned.
func (b *BasicMessage) FindSegment(segmentID string) segment.Segment {
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			return segment
		}
	}
	return nil
}

// FindSegments returns all segments found for the segmentID, or nil if nothing found
func (b *BasicMessage) FindSegments(segmentID string) []segment.Segment {
	var segments []segment.Segment
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			segments = append(segments, segment)
		}
	}
	return segments
}

// FindMarshaledSegment returns the first segment found for segmentID as []byte, or
// nil if nothing found
func (b *BasicMessage) FindMarshaledSegment(segmentID string) []byte {
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			return []byte(segment.String())
		}
	}
	return nil
}

// FindMarshaledSegments returns all segments found for segmentID as []byte, or
// nil if nothing found
func (b *BasicMessage) FindMarshaledSegments(segmentID string) [][]byte {
	var segments [][]byte
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			segments = append(segments, []byte(segment.String()))
		}
	}
	return segments
}

// SegmentPosition returns the segment position for the given segmentID
func (b *BasicMessage) SegmentPosition(segmentID string) int {
	idx := -1
	// TODO: implement
	//for i, segment := range b.HBCIMessage.HBCISegments() {
	//if
	//}
	return idx
}

// NewBasicSignedMessage creates a new signed message from message
func NewBasicSignedMessage(message *BasicMessage) *BasicSignedMessage {
	b := &BasicSignedMessage{
		message: message,
	}
	return b
}

// BasicSignedMessage represents a basic message which can be signed
type BasicSignedMessage struct {
	message *BasicMessage
}

// SetSegmentPositions sets the message numbers in all segments
func (b *BasicSignedMessage) SetSegmentPositions() {
	if b.message.SignatureBegin == nil || b.message.SignatureEnd == nil {
		panic(fmt.Errorf("Cannot set segment positions when signature is not set"))
	}
	b.message.SetSegmentPositions()
}

// SetSignatureHeader sets the SignatureHeader
func (b *BasicSignedMessage) SetSignatureHeader(sigBegin *segment.SignatureHeaderSegment) {
	b.message.SignatureBegin = sigBegin
}

// SetSignatureEnd sets the SignatureEnd
func (b *BasicSignedMessage) SetSignatureEnd(sigEnd *segment.SignatureEndSegment) {
	b.message.SignatureEnd = sigEnd
}

// HBCIVersion returns the HBCI version of the message
func (b *BasicSignedMessage) HBCIVersion() segment.HBCIVersion {
	return b.message.HBCIVersion()
}

// HBCISegments returns all segments of the message
func (b *BasicSignedMessage) HBCISegments() []segment.ClientSegment {
	return b.message.HBCISegments()
}

// MarshalHBCI marshals the message into HBCI wire format
func (b *BasicSignedMessage) MarshalHBCI() ([]byte, error) {
	return b.message.MarshalHBCI()
}

// Encrypt encrypts the message using the CryptoProvider
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
	jobs() []segment.ClientSegment
}
