package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (m *MessageHeaderSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], m)
	if err != nil {
		return err
	}
	m.Segment = seg
	if len(elements) > 1 {
		m.Size = &element.DigitDataElement{}
		err = m.Size.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 {
		m.HBCIVersion = &element.NumberDataElement{}
		err = m.HBCIVersion.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 {
		m.DialogID = &element.IdentificationDataElement{}
		err = m.DialogID.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 {
		m.Number = &element.NumberDataElement{}
		err = m.Number.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 {
		m.Ref = &element.ReferencingMessageDataElement{}
		err = m.Ref.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	return nil
}
