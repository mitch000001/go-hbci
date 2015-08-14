package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/swift"
)

func (a *AccountTransactionResponseSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], a)
	if err != nil {
		return err
	}
	a.Segment = seg
	if len(elements) >= 2 {
		a.BookedTransactions = &element.BinaryDataElement{}
		err = a.BookedTransactions.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
		messageExtractor := swift.NewMessageExtractor(a.BookedTransactions.Val())
		messages, err := messageExtractor.Extract()
		if err != nil {
			return err
		}
		for _, message := range messages {
			tr := &swift.MT940{}
			err = tr.Unmarshal(message)
			if err != nil {
				return err
			}
			a.bookedTransactions = append(a.bookedTransactions, tr)
		}
	}
	if len(elements) >= 3 {
		a.UnbookedTransactions = &element.BinaryDataElement{}
		err = a.UnbookedTransactions.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	return nil
}
