package test_files

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (v *VersionedTestSegmentCustomInterface) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment versionedTestSegmentCustomInterface
	switch header.Version.Val() {
	case 1:
		segment = &VersionedTestSegmentCustomInterfaceV1{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	v.versionedTestSegmentCustomInterface = segment
	return nil
}

func (v *VersionedTestSegmentCustomInterfaceV1) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], v)
	if err != nil {
		return err
	}
	v.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		v.Abc = &element.AlphaNumericDataElement{}
		err = v.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		v.Def = &element.NumberDataElement{}
		if len(elements)+1 > 2 {
			err = v.Def.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = v.Def.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
