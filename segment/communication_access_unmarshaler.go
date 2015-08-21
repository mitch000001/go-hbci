package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (c *CommunicationAccessResponseSegment) UnmarshalHBCI(value []byte) error {
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
		c.BankID = &element.BankIdentificationDataElement{}
		err = c.BankID.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		c.StandardLanguage = &element.NumberDataElement{}
		err = c.StandardLanguage.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		c.CommunicationParams = &element.CommunicationParameterDataElement{}
		if len(elements)+1 > 3 {
			err = c.CommunicationParams.UnmarshalHBCI(bytes.Join(elements[3:], []byte("+")))
		} else {
			err = c.CommunicationParams.UnmarshalHBCI(elements[3])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
