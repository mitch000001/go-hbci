// Code generated by *generator.VersionedSegmentUnmarshalerGenerator; DO NOT EDIT.

package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

var (
	_ BankSegment = &PinTanBankParameterV1{}
)

func init() {
	v1 := PinTanBankParameterV1{}
	KnownSegments.mustAddToIndex(VersionedSegment{v1.ID(), v1.Version()}, func() Segment { return &PinTanBankParameterV1{} })
}

func (p *PinTanBankParameterSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment PinTanBankParameter
	switch header.Version.Val() {
	case 1:
		segment = &PinTanBankParameterV1{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown segment version: %d", header.Version.Val())
	}
	p.PinTanBankParameter = segment
	return nil
}

func (p *PinTanBankParameterV1) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], p)
	if err != nil {
		return err
	}
	p.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		p.MaxJobs = &element.NumberDataElement{}
		err = p.MaxJobs.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("error unmarshaling MaxJobs: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		p.MinSignatures = &element.NumberDataElement{}
		err = p.MinSignatures.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("error unmarshaling MinSignatures: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		p.XXX_Unknown = &element.NumberDataElement{}
		err = p.XXX_Unknown.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("error unmarshaling XXX_Unknown: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		p.PinTanSpecificParams = &element.PinTanSpecificParamDataElement{}
		if len(elements)+1 > 4 {
			err = p.PinTanSpecificParams.UnmarshalHBCI(bytes.Join(elements[4:], []byte("+")))
		} else {
			err = p.PinTanSpecificParams.UnmarshalHBCI(elements[4])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling PinTanSpecificParams: %w", err)
		}
	}
	return nil
}
