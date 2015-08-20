package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (s *SynchronisationResponseSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], s)
	if err != nil {
		return err
	}
	s.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		s.ClientSystemID = &element.IdentificationDataElement{}
		err = s.ClientSystemID.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		s.MessageNumber = &element.NumberDataElement{}
		err = s.MessageNumber.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		s.SignatureID = &element.NumberDataElement{}
		err = s.SignatureID.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	return nil
}
