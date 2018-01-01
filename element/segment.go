package element

import "fmt"

// NewReferencingSegmentHeader returns a new SegmentHeader with a reference to
// another segment
func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	header := NewSegmentHeader(id, number, version)
	header.Ref = NewNumber(reference, 3)
	return header
}

// NewSegmentHeader returns a new SegmentHeader for the id, number and version
func NewSegmentHeader(id string, number, version int) *SegmentHeader {
	header := &SegmentHeader{
		ID:      NewAlphaNumeric(id, 6),
		Number:  NewNumber(number, 3),
		Version: NewNumber(version, 3),
	}
	header.DataElement = NewDataElementGroup(segmentHeaderDEG, 4, header)
	return header
}

// A SegmentHeader represents the metadata of a given segment such as ID or version
type SegmentHeader struct {
	DataElement
	ID      *AlphaNumericDataElement
	Number  *NumberDataElement
	Version *NumberDataElement
	Ref     *NumberDataElement
}

// Type returns the DataElementType of s
func (s *SegmentHeader) Type() DataElementType {
	return segmentHeaderDEG
}

// SetNumber sets the number of the segment
func (s *SegmentHeader) SetNumber(number int) {
	*s.Number = *NewNumber(number, 3)
}

// SetReference sets the reference to another segment
func (s *SegmentHeader) SetReference(ref int) {
	*s.Ref = *NewNumber(ref, 3)
}

// ReferencingSegment returns the reference number of the referenced segment
func (s *SegmentHeader) ReferencingSegment() int {
	if s.Ref != nil {
		return s.Ref.Val()
	}
	return -1
}

// IsValid returns true if the DataElement and all its grouped elements
// are valid, false otherwise
func (s *SegmentHeader) IsValid() bool {
	if s.ID == nil || s.Number == nil || s.Version == nil {
		return false
	}
	return s.DataElement.IsValid()
}

// GroupDataElements returns the grouped DataElements
func (s *SegmentHeader) GroupDataElements() []DataElement {
	return []DataElement{
		s.ID,
		s.Number,
		s.Version,
		s.Ref,
	}
}

// UnmarshalHBCI unmarshals value into s
func (s *SegmentHeader) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 3 {
		return fmt.Errorf("Malformed segment header")
	}
	s.DataElement = NewDataElementGroup(segmentHeaderDEG, 4, s)
	if len(elements) > 0 {
		s.ID = &AlphaNumericDataElement{}
		err = s.ID.UnmarshalHBCI(elements[0])
		if err != nil {
			return err
		}
	}
	if len(elements) > 1 {
		s.Number = &NumberDataElement{}
		err = s.Number.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("Malformed segment header number: %v", err)
		}
	}
	if len(elements) > 2 {
		s.Version = &NumberDataElement{}
		err = s.Version.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("Malformed segment header version: %v", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		s.Ref = &NumberDataElement{}
		err = s.Ref.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("Malformed segment header reference: %v", err)
		}
	}
	return nil
}
