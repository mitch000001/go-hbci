package segment

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

type MessageAcknowledgement struct {
	Segment
	acknowledgements []*element.AcknowledgementDataElement
}

func (m *MessageAcknowledgement) version() int         { return 2 }
func (m *MessageAcknowledgement) id() string           { return "HIRMG" }
func (m *MessageAcknowledgement) referencedId() string { return "" }
func (m *MessageAcknowledgement) sender() string       { return senderBank }

func (m *MessageAcknowledgement) UnmarshalHBCI(value []byte) error {
	value = bytes.TrimSuffix(value, []byte("'"))
	elements := bytes.Split(value, []byte("+"))
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	header := elements[0]
	numStr := bytes.Split(header, []byte(":"))[1]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return fmt.Errorf("Malformed segment header")
	}
	m.Segment = NewBasicSegment(num, m)
	elements = elements[1:]
	acknowledgements := make([]*element.AcknowledgementDataElement, len(elements))
	for i, elem := range elements {
		ack := new(element.AcknowledgementDataElement)
		err := ack.UnmarshalHBCI(elem)
		if err != nil {
			return err
		}
		acknowledgements[i] = ack
	}
	m.acknowledgements = acknowledgements
	return nil
}

func (m *MessageAcknowledgement) Acknowledgements() []domain.Acknowledgement {
	acknowledgements := make([]domain.Acknowledgement, len(m.acknowledgements))
	for i, ackDe := range m.acknowledgements {
		acknowledgements[i] = ackDe.Val()
	}
	return acknowledgements
}

func (m *MessageAcknowledgement) elements() []element.DataElement {
	dataElements := make([]element.DataElement, len(m.acknowledgements))
	for i, ack := range m.acknowledgements {
		dataElements[i] = ack
	}
	return dataElements
}

func (m *MessageAcknowledgement) Value() []domain.Acknowledgement {
	acks := make([]domain.Acknowledgement, len(m.acknowledgements))
	for i, de := range m.acknowledgements {
		acks[i] = de.Val()
	}
	return acks
}

type SegmentAcknowledgement struct {
	Segment
	acknowledgements []*element.AcknowledgementDataElement
}

func (s *SegmentAcknowledgement) version() int         { return 2 }
func (s *SegmentAcknowledgement) id() string           { return "HIRMS" }
func (s *SegmentAcknowledgement) referencedId() string { return "" }
func (s *SegmentAcknowledgement) sender() string       { return senderBank }

func (s *SegmentAcknowledgement) UnmarshalHBCI(value []byte) error {
	value = bytes.TrimSuffix(value, []byte("'"))
	elements := bytes.Split(value, []byte("+"))
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	header := elements[0]
	numStr := bytes.Split(header, []byte(":"))[1]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return fmt.Errorf("Malformed segment header")
	}
	s.Segment = NewBasicSegment(num, s)
	elements = elements[1:]
	acknowledgements := make([]*element.AcknowledgementDataElement, len(elements))
	for i, elem := range elements {
		ack := new(element.AcknowledgementDataElement)
		err := ack.UnmarshalHBCI(elem)
		if err != nil {
			return err
		}
		acknowledgements[i] = ack
	}
	s.acknowledgements = acknowledgements
	return nil
}

func (s *SegmentAcknowledgement) Acknowledgements() []domain.Acknowledgement {
	acknowledgements := make([]domain.Acknowledgement, len(s.acknowledgements))
	for i, ackDe := range s.acknowledgements {
		acknowledgements[i] = ackDe.Val()
	}
	return acknowledgements
}

func (s *SegmentAcknowledgement) elements() []element.DataElement {
	dataElements := make([]element.DataElement, len(s.acknowledgements))
	for i, ack := range s.acknowledgements {
		dataElements[i] = ack
	}
	return dataElements
}
