package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (s *SecurityMethodSegment) UnmarshalHBCI(value []byte) error {
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
		s.MixAllowed = &element.BooleanDataElement{}
		err = s.MixAllowed.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		s.SupportedMethods = &element.SupportedSecurityMethodDataElement{}
		err = s.SupportedMethods.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	return nil
}
