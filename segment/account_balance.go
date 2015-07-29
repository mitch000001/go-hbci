package segment

import "github.com/mitch000001/go-hbci/element"

type AccountBalanceRequestSegment struct {
	Segment
	AccountConnection *element.AccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegment) version() int         { return 5 }
func (a *AccountBalanceRequestSegment) id() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegment) referencedId() string { return "" }
func (a *AccountBalanceRequestSegment) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

type AccountBalanceResponseSegment struct {
	Segment
	AccountConnection  *element.AccountConnectionDataElement
	AccountProductName *element.AlphaNumericDataElement
	AccountCurrency    *element.CurrencyDataElement
	BookedBalance      *element.BalanceDataElement
	EarMarkedBalance   *element.BalanceDataElement
	CreditLimit        *element.AmountDataElement
	AvailableAmount    *element.AmountDataElement
	UsedAmount         *element.AmountDataElement
	BookingDate        *element.DateDataElement
	BookingTime        *element.TimeDataElement
	BalancingDate      *element.DateDataElement
}

func (a *AccountBalanceResponseSegment) version() int         { return 5 }
func (a *AccountBalanceResponseSegment) id() string           { return "HISAL" }
func (a *AccountBalanceResponseSegment) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegment) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AccountProductName,
		a.AccountCurrency,
		a.BookedBalance,
		a.EarMarkedBalance,
		a.CreditLimit,
		a.AvailableAmount,
		a.UsedAmount,
		a.BookingDate,
		a.BookingTime,
		a.BalancingDate,
	}
}
