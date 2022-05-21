// Code generated by *generator.VersionedSegmentUnmarshalerGenerator; DO NOT EDIT.

package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (t *TanResponseSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment TanResponse
	switch header.Version.Val() {
	case 6:
		segment = &TanResponseSegmentV6{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	t.TanResponse = segment
	return nil
}

func (t *TanResponseSegmentV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], t)
	if err != nil {
		return err
	}
	t.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		t.TANProcess = &element.AlphaNumericDataElement{}
		err = t.TANProcess.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		t.JobHash = &element.BinaryDataElement{}
		err = t.JobHash.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		t.JobReference = &element.AlphaNumericDataElement{}
		if len(elements)+1 > 3 {
			err = t.JobReference.UnmarshalHBCI(bytes.Join(elements[3:], []byte("+")))
		} else {
			err = t.JobReference.UnmarshalHBCI(elements[3])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
