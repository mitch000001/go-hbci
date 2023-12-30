// Code generated by *generator.SegmentUnmarshalerGenerator; DO NOT EDIT.

package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (e *EncryptedDataSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], e)
	if err != nil {
		return err
	}
	e.ClientSegment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		e.Data = &element.BinaryDataElement{}
		if len(elements)+1 > 1 {
			err = e.Data.UnmarshalHBCI(bytes.Join(elements[1:], []byte("+")))
		} else {
			err = e.Data.UnmarshalHBCI(elements[1])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling Data: %w", err)
		}
	}
	return nil
}
