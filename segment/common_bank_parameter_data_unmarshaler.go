package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (c *CommonBankParameterSegment) UnmarshalHBCI(value []byte) error {
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
	if len(elements) > 1 {
		c.BPDVersion = &element.NumberDataElement{}
		err = c.BPDVersion.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 {
		c.BankID = &element.BankIdentificationDataElement{}
		err = c.BankID.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 {
		c.BankName = &element.AlphaNumericDataElement{}
		err = c.BankName.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 {
		c.BusinessTransactionCount = &element.NumberDataElement{}
		err = c.BusinessTransactionCount.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 {
		c.SupportedLanguages = &element.SupportedLanguagesDataElement{}
		err = c.SupportedLanguages.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) > 6 {
		c.SupportedHBCIVersions = &element.SupportedHBCIVersionsDataElement{}
		err = c.SupportedHBCIVersions.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	if len(elements) > 7 {
		c.MaxMessageSize = &element.NumberDataElement{}
		err = c.MaxMessageSize.UnmarshalHBCI(elements[7])
		if err != nil {
			return err
		}
	}
	return nil
}
