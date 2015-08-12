package segment

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
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

type MessageHeaderSegment struct {
	Segment
	Size        *element.DigitDataElement
	HBCIVersion *element.NumberDataElement
	DialogID    *element.IdentificationDataElement
	Number      *element.NumberDataElement
	Ref         *element.ReferencingMessageDataElement
}

func (m *MessageHeaderSegment) version() int         { return 3 }
func (m *MessageHeaderSegment) id() string           { return "HNHBK" }
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

func (m *MessageHeaderSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	elementsLen := len(elements)
	if elementsLen == 0 || elementsLen < 5 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], m)
	if err != nil {
		return err
	}
	m.Segment = seg
	size, err := strconv.Atoi(charset.ToUtf8(elements[1]))
	if err != nil {
		return fmt.Errorf("Error while unmarshaling size: %v", err)
	}
	hbciVersion, err := strconv.Atoi(charset.ToUtf8(elements[2]))
	if err != nil {
		return fmt.Errorf("Error while unmarshaling hbci version: %v", err)
	}
	dialogId := charset.ToUtf8(elements[3])
	messageNum, err := strconv.Atoi(charset.ToUtf8(elements[4]))
	if err != nil {
		return fmt.Errorf("Error while unmarshaling message number: %v", err)
	}
	m.Size = element.NewDigit(size, 12)
	m.HBCIVersion = element.NewNumber(hbciVersion, 3)
	m.DialogID = element.NewIdentification(dialogId)
	m.Number = element.NewNumber(messageNum, 4)
	if elementsLen == 6 {
		referencedMessage := &element.ReferencingMessageDataElement{}
		err = referencedMessage.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
		m.Ref = referencedMessage
	}
	return nil
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

type MessageEndSegment struct {
	Segment
	Number *element.NumberDataElement
}

func (m *MessageEndSegment) version() int         { return 1 }
func (m *MessageEndSegment) id() string           { return "HNHBS" }
func (m *MessageEndSegment) referencedId() string { return "" }
func (m *MessageEndSegment) sender() string       { return senderBoth }

func (m *MessageEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		m.Number,
	}
}
