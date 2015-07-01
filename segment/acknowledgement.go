package segment

import "github.com/mitch000001/go-hbci/element"

type MessageAcknowledgement struct {
	Segment
	Acknowledgements []*element.AcknowledgementDataElement
}

func (m *MessageAcknowledgement) elements() []element.DataElement {
	m.Acknowledgements = make([]*element.AcknowledgementDataElement, 99)
	dataElements := make([]element.DataElement, 99)
	for i, _ := range dataElements {
		dataElements[i] = m.Acknowledgements[i]
	}
	return dataElements
}

type SegmentAcknowledgement struct {
	Segment
	Acknowledgements []*element.AcknowledgementDataElement
}

func (s *SegmentAcknowledgement) elements() []element.DataElement {
	s.Acknowledgements = make([]*element.AcknowledgementDataElement, 99)
	dataElements := make([]element.DataElement, 99)
	for i, _ := range dataElements {
		dataElements[i] = s.Acknowledgements[i]
	}
	return dataElements
}
