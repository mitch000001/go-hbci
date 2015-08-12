package element

func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	header := NewSegmentHeader(id, number, version)
	header.Ref = NewNumber(reference, 3)
	return header
}

func NewSegmentHeader(id string, number, version int) *SegmentHeader {
	header := &SegmentHeader{
		ID:      NewAlphaNumeric(id, 6),
		Number:  NewNumber(number, 3),
		Version: NewNumber(version, 3),
	}
	header.DataElement = NewDataElementGroup(SegmentHeaderDEG, 4, header)
	return header
}

type SegmentHeader struct {
	DataElement
	ID      *AlphaNumericDataElement
	Number  *NumberDataElement
	Version *NumberDataElement
	Ref     *NumberDataElement
}

func (s *SegmentHeader) SetNumber(number int) {
	*s.Number = *NewNumber(number, 3)
}

func (s *SegmentHeader) SetReference(ref int) {
	*s.Ref = *NewNumber(ref, 3)
}

func (s *SegmentHeader) ReferencingSegment() int {
	if s.Ref != nil {
		return s.Ref.Val()
	} else {
		return -1
	}
}

func (s *SegmentHeader) IsValid() bool {
	if s.ID == nil || s.Number == nil || s.Version == nil {
		return false
	} else {
		return s.DataElement.IsValid()
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
