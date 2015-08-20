package segment

import (
	"bytes"
	"reflect"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

type SegmentSequence map[string]Segment

func (s SegmentSequence) Header() *element.SegmentHeader { return nil }

func (s SegmentSequence) SetNumber(numberFn func() int) {
	for _, segment := range s {
		if !reflect.ValueOf(segment).IsNil() {
			segment.SetNumber(numberFn)
		}
	}
}

func (s SegmentSequence) DataElements() []element.DataElement {
	var elements []element.DataElement
	for _, segment := range s {
		if !reflect.ValueOf(segment).IsNil() {
			elements = append(elements, segment.DataElements()...)
		}
	}
	return elements
}

func (s SegmentSequence) String() string {
	var buf bytes.Buffer
	for _, segment := range s {
		if !reflect.ValueOf(segment).IsNil() {
			buf.WriteString(segment.String())
		}
	}
	return buf.String()
}

func NewReferencingMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int, referencingMessage domain.ReferencingMessage) *MessageHeaderSegment {
	m := NewMessageHeaderSegment(size, hbciVersion, dialogId, number)
	m.Ref = element.NewReferencingMessage(referencingMessage.DialogID, referencingMessage.MessageNumber)
	return m
}

func NewMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	m := &MessageHeaderSegment{
		Size:        element.NewDigit(size, 12),
		HBCIVersion: element.NewNumber(hbciVersion, 3),
		DialogID:    element.NewIdentification(dialogId),
		Number:      element.NewNumber(number, 4),
	}
	m.Segment = NewBasicSegment(1, m)
	return m
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment MessageHeaderSegment

type MessageHeaderSegment struct {
	Segment
	Size        *element.DigitDataElement
	HBCIVersion *element.NumberDataElement
	DialogID    *element.IdentificationDataElement
	Number      *element.NumberDataElement
	Ref         *element.ReferencingMessageDataElement
}

func (m *MessageHeaderSegment) Version() int         { return 3 }
func (m *MessageHeaderSegment) ID() string           { return "HNHBK" }
func (m *MessageHeaderSegment) referencedId() string { return "" }
func (m *MessageHeaderSegment) sender() string       { return senderBoth }

func (m *MessageHeaderSegment) ReferencingMessage() domain.ReferencingMessage {
	var reference domain.ReferencingMessage
	if m.Ref != nil {
		reference = m.Ref.Val()
	}
	return reference
}

func (m *MessageHeaderSegment) SetMessageNumber(messageNumber int) {
	m.Number = element.NewNumber(messageNumber, 4)
}

func (m *MessageHeaderSegment) elements() []element.DataElement {
	return []element.DataElement{
		m.Size,
		m.HBCIVersion,
		m.DialogID,
		m.Number,
		m.Ref,
	}
}

func (m *MessageHeaderSegment) SetSize(size int) {
	*m.Size = *element.NewDigit(size, 12)
}

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	end := &MessageEndSegment{
		Number: element.NewNumber(messageNumber, 4),
	}
	end.Segment = NewBasicSegment(segmentNumber, end)
	return end
}
