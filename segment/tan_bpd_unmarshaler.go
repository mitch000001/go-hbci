// Code generated by *generator.VersionedSegmentUnmarshalerGenerator; DO NOT EDIT.

package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (t *TanBankParameterSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment TanBankParameter
	switch header.Version.Val() {
	case 6:
		segment = &TanBankParameterV6{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown segment version: %d", header.Version.Val())
	}
	t.TanBankParameter = segment
	return nil
}

func (t *TanBankParameterV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], t)
	if err != nil {
		return err
	}
	t.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		t.MaxJobs = &element.NumberDataElement{}
		err = t.MaxJobs.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("error unmarshaling MaxJobs: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		t.MinSignatures = &element.NumberDataElement{}
		err = t.MinSignatures.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("error unmarshaling MinSignatures: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		t.SecurityClass = &element.CodeDataElement{}
		err = t.SecurityClass.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("error unmarshaling SecurityClass: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		t.Tan2StepSubmissionParameter = &element.Tan2StepSubmissionParameterV6{}
		if len(elements)+1 > 4 {
			err = t.Tan2StepSubmissionParameter.UnmarshalHBCI(bytes.Join(elements[4:], []byte("+")))
		} else {
			err = t.Tan2StepSubmissionParameter.UnmarshalHBCI(elements[4])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling Tan2StepSubmissionParameter: %w", err)
		}
	}
	return nil
}
