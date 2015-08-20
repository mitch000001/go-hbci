package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (c *CompressionMethodSegment) UnmarshalHBCI(value []byte) error {
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
		c.SupportedCompressionMethods = &element.SupportedCompressionMethodsDataElement{}
		err = c.SupportedCompressionMethods.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	return nil
}
