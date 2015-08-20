package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (a *AccountInformationResponseSegment) UnmarshalHBCI(value []byte) error {
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
	if len(elements) > 1 && len(elements[1]) > 0 {
		a.AccountConnection = &element.AccountConnectionDataElement{}
		err = a.AccountConnection.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		a.AccountKind = &element.NumberDataElement{}
		err = a.AccountKind.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		a.Name1 = &element.AlphaNumericDataElement{}
		err = a.Name1.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		a.Name2 = &element.AlphaNumericDataElement{}
		err = a.Name2.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		a.AccountProductID = &element.AlphaNumericDataElement{}
		err = a.AccountProductID.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		a.AccountCurrency = &element.CurrencyDataElement{}
		err = a.AccountCurrency.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		a.OpeningDate = &element.DateDataElement{}
		err = a.OpeningDate.UnmarshalHBCI(elements[7])
		if err != nil {
			return err
		}
	}
	if len(elements) > 8 && len(elements[8]) > 0 {
		a.DebitInterest = &element.ValueDataElement{}
		err = a.DebitInterest.UnmarshalHBCI(elements[8])
		if err != nil {
			return err
		}
	}
	if len(elements) > 9 && len(elements[9]) > 0 {
		a.CreditInterest = &element.ValueDataElement{}
		err = a.CreditInterest.UnmarshalHBCI(elements[9])
		if err != nil {
			return err
		}
	}
	if len(elements) > 10 && len(elements[10]) > 0 {
		a.OverDebitInterest = &element.ValueDataElement{}
		err = a.OverDebitInterest.UnmarshalHBCI(elements[10])
		if err != nil {
			return err
		}
	}
	if len(elements) > 11 && len(elements[11]) > 0 {
		a.CreditLimit = &element.AmountDataElement{}
		err = a.CreditLimit.UnmarshalHBCI(elements[11])
		if err != nil {
			return err
		}
	}
	if len(elements) > 12 && len(elements[12]) > 0 {
		a.ReferenceAccount = &element.AccountConnectionDataElement{}
		err = a.ReferenceAccount.UnmarshalHBCI(elements[12])
		if err != nil {
			return err
		}
	}
	if len(elements) > 13 && len(elements[13]) > 0 {
		a.AccountStatementShippingType = &element.NumberDataElement{}
		err = a.AccountStatementShippingType.UnmarshalHBCI(elements[13])
		if err != nil {
			return err
		}
	}
	if len(elements) > 14 && len(elements[14]) > 0 {
		a.AccountStatementShippingRotation = &element.NumberDataElement{}
		err = a.AccountStatementShippingRotation.UnmarshalHBCI(elements[14])
		if err != nil {
			return err
		}
	}
	if len(elements) > 15 && len(elements[15]) > 0 {
		a.AdditionalInformation = &element.TextDataElement{}
		err = a.AdditionalInformation.UnmarshalHBCI(elements[15])
		if err != nil {
			return err
		}
	}
	if len(elements) > 16 && len(elements[16]) > 0 {
		a.DisposalEligiblePersons = &element.DisposalEligiblePersonsDataElement{}
		err = a.DisposalEligiblePersons.UnmarshalHBCI(elements[16])
		if err != nil {
			return err
		}
	}
	return nil
}
