package segment

import "github.com/mitch000001/go-hbci/element"

type MessageAcknowledgement struct {
	Segment
	Acknowledgements *element.AcknowledgementsDataElement
}

func (m *MessageAcknowledgement) version() int         { return 2 }
func (m *MessageAcknowledgement) id() string           { return "HIRMG" }
func (m *MessageAcknowledgement) referencedId() string { return "" }
func (m *MessageAcknowledgement) sender() string       { return senderBank }

func (m *MessageAcknowledgement) elements() []element.DataElement {
	return []element.DataElement{
		m.Acknowledgements,
	}
}

type SegmentAcknowledgement struct {
	Segment
	Acknowledgements *element.AcknowledgementsDataElement
}

func (m *SegmentAcknowledgement) version() int         { return 2 }
func (m *SegmentAcknowledgement) id() string           { return "HIRMS" }
func (m *SegmentAcknowledgement) referencedId() string { return "" }
func (m *SegmentAcknowledgement) sender() string       { return senderBank }

func (s *SegmentAcknowledgement) elements() []element.DataElement {
	return []element.DataElement{
		s.Acknowledgements,
	}
}
