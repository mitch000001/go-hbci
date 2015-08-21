package test_files

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (m *MultipleVersionedTestSegmentCustomInterfaces) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment BankSegment
	switch header.Version.Val() {
	case 1:
		segment = &MultipleVersionedTestSegmentCustomInterfacesV1{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 2:
		segment = &MultipleVersionedTestSegmentCustomInterfacesV2{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	m.BankSegment = segment
	return nil
}

func (m *MultipleVersionedTestSegmentCustomInterfacesV1) UnmarshalHBCI(value []byte) error {
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
	m.versionInterface1 = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		m.Abc = &element.AlphaNumericDataElement{}
		err = m.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		m.Def = &element.NumberDataElement{}
		if len(elements)+1 > 2 {
			err = m.Def.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = m.Def.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MultipleVersionedTestSegmentCustomInterfacesV2) UnmarshalHBCI(value []byte) error {
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
	m.versionInterface2 = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		m.Abc = &element.AlphaNumericDataElement{}
		err = m.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		m.Def = &element.NumberDataElement{}
		if len(elements)+1 > 2 {
			err = m.Def.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = m.Def.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
