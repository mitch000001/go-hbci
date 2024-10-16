// Code generated by *generator.VersionedSegmentUnmarshalerGenerator; DO NOT EDIT.

package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

var (
	_	BankSegment	= &CommonBankParameterV2{}
	_	BankSegment	= &CommonBankParameterV3{}
)

func init() {
	v2 := CommonBankParameterV2{}
	KnownSegments.mustAddToIndex(VersionedSegment{v2.ID(), v2.Version()}, func() Segment { return &CommonBankParameterV2{} })
	v3 := CommonBankParameterV3{}
	KnownSegments.mustAddToIndex(VersionedSegment{v3.ID(), v3.Version()}, func() Segment { return &CommonBankParameterV3{} })
}

func (c *CommonBankParameterSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment commonBankParameterSegment
	switch header.Version.Val() {
	case 2:
		segment = &CommonBankParameterV2{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 3:
		segment = &CommonBankParameterV3{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown segment version: %d", header.Version.Val())
	}
	c.commonBankParameterSegment = segment
	return nil
}

func (c *CommonBankParameterV2) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], c)
	if err != nil {
		return err
	}
	c.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		c.BPDVersion = &element.NumberDataElement{}
		err = c.BPDVersion.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("error unmarshaling BPDVersion: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		c.BankID = &element.BankIdentificationDataElement{}
		err = c.BankID.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("error unmarshaling BankID: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		c.BankName = &element.AlphaNumericDataElement{}
		err = c.BankName.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("error unmarshaling BankName: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		c.BusinessTransactionCount = &element.NumberDataElement{}
		err = c.BusinessTransactionCount.UnmarshalHBCI(elements[4])
		if err != nil {
			return fmt.Errorf("error unmarshaling BusinessTransactionCount: %w", err)
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		c.SupportedLanguages = &element.SupportedLanguagesDataElement{}
		err = c.SupportedLanguages.UnmarshalHBCI(elements[5])
		if err != nil {
			return fmt.Errorf("error unmarshaling SupportedLanguages: %w", err)
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		c.SupportedHBCIVersions = &element.SupportedHBCIVersionsDataElement{}
		err = c.SupportedHBCIVersions.UnmarshalHBCI(elements[6])
		if err != nil {
			return fmt.Errorf("error unmarshaling SupportedHBCIVersions: %w", err)
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		c.MaxMessageSize = &element.NumberDataElement{}
		if len(elements)+1 > 7 {
			err = c.MaxMessageSize.UnmarshalHBCI(bytes.Join(elements[7:], []byte("+")))
		} else {
			err = c.MaxMessageSize.UnmarshalHBCI(elements[7])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling MaxMessageSize: %w", err)
		}
	}
	return nil
}

func (c *CommonBankParameterV3) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], c)
	if err != nil {
		return err
	}
	c.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		c.BPDVersion = &element.NumberDataElement{}
		err = c.BPDVersion.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("error unmarshaling BPDVersion: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		c.BankID = &element.BankIdentificationDataElement{}
		err = c.BankID.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("error unmarshaling BankID: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		c.BankName = &element.AlphaNumericDataElement{}
		err = c.BankName.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("error unmarshaling BankName: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		c.BusinessTransactionCount = &element.NumberDataElement{}
		err = c.BusinessTransactionCount.UnmarshalHBCI(elements[4])
		if err != nil {
			return fmt.Errorf("error unmarshaling BusinessTransactionCount: %w", err)
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		c.SupportedLanguages = &element.SupportedLanguagesDataElement{}
		err = c.SupportedLanguages.UnmarshalHBCI(elements[5])
		if err != nil {
			return fmt.Errorf("error unmarshaling SupportedLanguages: %w", err)
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		c.SupportedHBCIVersions = &element.SupportedHBCIVersionsDataElement{}
		err = c.SupportedHBCIVersions.UnmarshalHBCI(elements[6])
		if err != nil {
			return fmt.Errorf("error unmarshaling SupportedHBCIVersions: %w", err)
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		c.MaxMessageSize = &element.NumberDataElement{}
		err = c.MaxMessageSize.UnmarshalHBCI(elements[7])
		if err != nil {
			return fmt.Errorf("error unmarshaling MaxMessageSize: %w", err)
		}
	}
	if len(elements) > 8 && len(elements[8]) > 0 {
		c.MinTimeoutValue = &element.NumberDataElement{}
		err = c.MinTimeoutValue.UnmarshalHBCI(elements[8])
		if err != nil {
			return fmt.Errorf("error unmarshaling MinTimeoutValue: %w", err)
		}
	}
	if len(elements) > 9 && len(elements[9]) > 0 {
		c.MaxTimeoutValue = &element.NumberDataElement{}
		if len(elements)+1 > 9 {
			err = c.MaxTimeoutValue.UnmarshalHBCI(bytes.Join(elements[9:], []byte("+")))
		} else {
			err = c.MaxTimeoutValue.UnmarshalHBCI(elements[9])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling MaxTimeoutValue: %w", err)
		}
	}
	return nil
}
