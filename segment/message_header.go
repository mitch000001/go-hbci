package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

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
	m.ClientSegment = NewBasicSegment(1, m)
	return m
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment MessageHeaderSegment -segment_interface ClientSegment

type MessageHeaderSegment struct {
	ClientSegment
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
