package hbci

import "strings"

type Segment interface {
	DataElements() []DataElement
}

func NewSegment(header *SegmentHeader, seg Segment) *segment {
	return &segment{Header: header, segment: seg}
}

type segment struct {
	segment Segment
	Header  *SegmentHeader
}

func (s *segment) String() string {
	elementStrings := make([]string, len(s.segment.DataElements())+1)
	elementStrings[0] = s.Header.String()
	for i, de := range s.segment.DataElements() {
		elementStrings[i+1] = de.String()
	}
	return strings.Join(elementStrings, "+") + "'"
}

func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	return &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
		Ref:     NewNumberDataElement(reference, 3),
	}
}

func NewSegmentHeader(id string, number, version int) *SegmentHeader {
	header := &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
	}
	header.elementGroup = NewDataElementGroup(SegmentHeaderDEG, 4, header)
	return header
}

type SegmentHeader struct {
	*elementGroup
	ID      *AlphaNumericDataElement
	Number  *NumberDataElement
	Version *NumberDataElement
	Ref     *NumberDataElement
}

func (s *SegmentHeader) Valid() bool {
	if s.ID == nil || s.Number == nil || s.Version == nil {
		return false
	} else {
		return s.elementGroup.Valid()
	}
}

func (s *SegmentHeader) Value() interface{} {
	return s
}

func (s *SegmentHeader) GroupDataElements() []DataElement {
	return []DataElement{
		s.ID,
		s.Number,
		s.Version,
		s.Ref,
	}
}
