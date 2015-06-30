package segment

import "github.com/mitch000001/go-hbci/dataelement"

type SegmentSequence []Segment

var validHBCIVersions = []int{201, 210, 220}

func NewReferencingMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int, referencedMessage *dataelement.ReferenceMessage) *MessageHeaderSegment {
	m := NewMessageHeaderSegment(size, hbciVersion, dialogId, number)
	m.Ref = referencedMessage
	return m
}

func NewMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	m := &MessageHeaderSegment{
		Size:        dataelement.NewDigit(size, 12),
		HBCIVersion: dataelement.NewNumber(hbciVersion, 3),
		DialogID:    dataelement.NewIdentification(dialogId),
		Number:      dataelement.NewNumber(number, 4),
	}
	m.Segment = NewBasicSegment("HNHBK", 1, 3, m)
	return m
}

type MessageHeaderSegment struct {
	Segment
	Size        *dataelement.DigitDataElement
	HBCIVersion *dataelement.NumberDataElement
	DialogID    *dataelement.IdentificationDataElement
	Number      *dataelement.NumberDataElement
	Ref         *dataelement.ReferenceMessage
}

func (m *MessageHeaderSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		m.Size,
		m.HBCIVersion,
		m.DialogID,
		m.Number,
		m.Ref,
	}
}

func (m *MessageHeaderSegment) SetSize(size int) {
	m.Size = dataelement.NewDigit(size, 12)
}

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	end := &MessageEndSegment{
		Number: dataelement.NewNumber(messageNumber, 4),
	}
	end.Segment = NewBasicSegment("HNHBS", segmentNumber, 1, end)
	return end
}

type MessageEndSegment struct {
	Segment
	Number *dataelement.NumberDataElement
}

func (m *MessageEndSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		m.Number,
	}
}
