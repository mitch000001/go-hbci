package segment

import "github.com/mitch000001/go-hbci/element"

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	end := &MessageEndSegment{
		Number: element.NewNumber(messageNumber, 4),
	}
	end.ClientSegment = NewBasicSegment(segmentNumber, end)
	return end
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment MessageEndSegment

type MessageEndSegment struct {
	ClientSegment
	Number *element.NumberDataElement
}

func (m *MessageEndSegment) Version() int         { return 1 }
func (m *MessageEndSegment) ID() string           { return "HNHBS" }
func (m *MessageEndSegment) referencedId() string { return "" }
func (m *MessageEndSegment) sender() string       { return senderBoth }

func (m *MessageEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		m.Number,
	}
}
