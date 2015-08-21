package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (m *MessageEndSegment) UnmarshalHBCI(value []byte) error {
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
	m.ClientSegment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		m.Number = &element.NumberDataElement{}
		if len(elements)+1 > 1 {
			err = m.Number.UnmarshalHBCI(bytes.Join(elements[1:], []byte("+")))
		} else {
			err = m.Number.UnmarshalHBCI(elements[1])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
