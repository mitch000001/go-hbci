package hbci

import "strings"

func NewSegment(header *SegmentHeader, dataElements ...DataElement) *Segment {
	return &Segment{Header: header, elements: dataElements}
}

type Segment struct {
	Header   *SegmentHeader
	elements []DataElement
}

func (s *Segment) String() string {
	elementStrings := make([]string, len(s.elements)+1)
	elementStrings[0] = s.Header.String()
	for i, de := range s.elements {
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
	return &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
	}
}

type SegmentHeader struct {
	ID      *AlphaNumericDataElement
	Number  *NumberDataElement
	Version *NumberDataElement
	Ref     *NumberDataElement
}

func (s *SegmentHeader) Type() DataElementType {
	return SegmentHeaderDEG
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
	returnStr := strings.Join(elementStrings, ":")
	if s.Ref == nil {
		returnStr += ":"
	}
	return returnStr
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
