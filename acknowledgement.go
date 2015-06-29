package hbci

import "github.com/mitch000001/go-hbci/dataelement"

type MessageAcknowledgement struct {
	Segment
	Acknowledgements []*dataelement.AcknowledgementDataElement
}

func (m *MessageAcknowledgement) elements() []dataelement.DataElement {
	m.Acknowledgements = make([]*dataelement.AcknowledgementDataElement, 99)
	dataElements := make([]dataelement.DataElement, 99)
	for i, _ := range dataElements {
		dataElements[i] = m.Acknowledgements[i]
	}
	return dataElements
}

type SegmentAcknowledgement struct {
	Segment
	Acknowledgements []*dataelement.AcknowledgementDataElement
}

func (s *SegmentAcknowledgement) elements() []dataelement.DataElement {
	s.Acknowledgements = make([]*dataelement.AcknowledgementDataElement, 99)
	dataElements := make([]dataelement.DataElement, 99)
	for i, _ := range dataElements {
		dataElements[i] = s.Acknowledgements[i]
	}
	return dataElements
}
