package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (b *BankAnnouncementSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], b)
	if err != nil {
		return err
	}
	b.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		b.Subject = &element.AlphaNumericDataElement{}
		err = b.Subject.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		b.Body = &element.TextDataElement{}
		if len(elements)+1 > 2 {
			err = b.Body.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = b.Body.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
