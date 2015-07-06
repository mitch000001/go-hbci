package segment

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/element"
)

type MessageAcknowledgement struct {
	Segment
	Acknowledgements []*element.AcknowledgementDataElement
}

func (m *MessageAcknowledgement) init() {
	m.Acknowledgements = make([]*element.AcknowledgementDataElement, 99)
}
func (m *MessageAcknowledgement) version() int         { return 2 }
func (m *MessageAcknowledgement) id() string           { return "HIRMG" }
func (m *MessageAcknowledgement) referencedId() string { return "" }
func (m *MessageAcknowledgement) sender() string       { return senderBank }

func (m *MessageAcknowledgement) UnmarshalHBCI(value []byte) error {
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
	m.Acknowledgements = acknowledgements
	return nil
}

func (m *MessageAcknowledgement) elements() []element.DataElement {
	dataElements := make([]element.DataElement, len(m.Acknowledgements))
	for i, ack := range m.Acknowledgements {
		dataElements[i] = ack
	}
	return dataElements
}

type SegmentAcknowledgement struct {
	Segment
	Acknowledgements []*element.AcknowledgementDataElement
}

func (s *SegmentAcknowledgement) init() {
	s.Acknowledgements = make([]*element.AcknowledgementDataElement, 99)
}
func (s *SegmentAcknowledgement) version() int         { return 2 }
func (s *SegmentAcknowledgement) id() string           { return "HIRMS" }
func (s *SegmentAcknowledgement) referencedId() string { return "" }
func (s *SegmentAcknowledgement) sender() string       { return senderBank }

func (s *SegmentAcknowledgement) UnmarshalHBCI(value []byte) error {
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
	s.Acknowledgements = acknowledgements
	return nil
}

func (s *SegmentAcknowledgement) elements() []element.DataElement {
	dataElements := make([]element.DataElement, len(s.Acknowledgements))
	for i, ack := range s.Acknowledgements {
		dataElements[i] = ack
	}
	return dataElements
}
