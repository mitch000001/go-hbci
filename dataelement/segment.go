package dataelement

func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	header := NewSegmentHeader(id, number, version)
	header.Ref = NewNumberDataElement(reference, 3)
	return header
}

func NewSegmentHeader(id string, number, version int) *SegmentHeader {
	header := &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
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
	*s.Number = *NewNumberDataElement(number, 3)
}

func (s *SegmentHeader) SetReference(ref int) {
	*s.Ref = *NewNumberDataElement(ref, 3)
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
