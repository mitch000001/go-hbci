package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (c *CommonUserParameterDataSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], c)
	if err != nil {
		return err
	}
	c.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		c.UserID = &element.IdentificationDataElement{}
		err = c.UserID.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		c.UPDVersion = &element.NumberDataElement{}
		err = c.UPDVersion.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		c.UPDUsage = &element.NumberDataElement{}
		err = c.UPDUsage.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	return nil
}
