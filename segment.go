package hbci

import "strings"

type Segment struct{}

func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	return &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
		Ref:     NewNumberDataElement(reference, 3),
	}
}

func NewSegmentHeader(id string, number, version int) *SegmentHeader {
	return &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
	}
}

// TODO: How to retrieve the concrete or casted values?
type SegmentHeader struct {
	ID      DataElement
	Number  DataElement
	Version DataElement
	Ref     DataElement
}

func (s *SegmentHeader) Type() DataElementType {
	return SegmentHeaderType
}

func (s *SegmentHeader) Valid() bool {
	for _, d := range s.GroupDataElements() {
		if !d.Valid() {
			return false
		}
	}
	if s.ID == nil || s.Number == nil || s.Version == nil {
		return false
	}
	return true
}

func (s *SegmentHeader) Value() interface{} {
	return s
}

func (s *SegmentHeader) Length() int {
	length := 0
	for _, d := range s.GroupDataElements() {
		length += d.Length()
	}
	return length
}

func (s *SegmentHeader) String() string {
	elementStrings := make([]string, len(s.GroupDataElements()))
	for i, d := range s.GroupDataElements() {
		elementStrings[i] = d.String()
	}
	return strings.Join(elementStrings, ":")
}

func (s *SegmentHeader) GroupDataElements() []DataElement {
	var groupDataElements []DataElement
	groupDataElements = append(groupDataElements, s.ID)
	groupDataElements = append(groupDataElements, s.Number)
	groupDataElements = append(groupDataElements, s.Version)
	if s.Ref != nil {
		groupDataElements = append(groupDataElements, s.Ref)
	}
	return groupDataElements
}
