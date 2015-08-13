package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (a *AccountBalanceResponseSegment) UnmarshalHBCI(value []byte) error {
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
		a.AccountConnection = &element.AccountConnectionDataElement{}
		err = a.AccountConnection.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 3 {
		a.AccountProductName = &element.AlphaNumericDataElement{}
		err = a.AccountProductName.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 4 {
		a.AccountCurrency = &element.CurrencyDataElement{}
		err = a.AccountCurrency.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 5 {
		a.BookedBalance = &element.BalanceDataElement{}
		err = a.BookedBalance.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 6 {
		a.EarmarkedBalance = &element.BalanceDataElement{}
		err = a.EarmarkedBalance.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 7 {
		a.CreditLimit = &element.AmountDataElement{}
		err = a.CreditLimit.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 8 {
		a.AvailableAmount = &element.AmountDataElement{}
		err = a.AvailableAmount.UnmarshalHBCI(elements[7])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 9 {
		a.UsedAmount = &element.AmountDataElement{}
		err = a.UsedAmount.UnmarshalHBCI(elements[8])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 10 {
		a.BookingDate = &element.DateDataElement{}
		err = a.BookingDate.UnmarshalHBCI(elements[9])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 11 {
		a.BookingTime = &element.TimeDataElement{}
		err = a.BookingTime.UnmarshalHBCI(elements[10])
		if err != nil {
			return err
		}
	}
	if len(elements) >= 12 {
		a.DueDate = &element.DateDataElement{}
		err = a.DueDate.UnmarshalHBCI(elements[11])
		if err != nil {
			return err
		}
	}
	return nil
}
