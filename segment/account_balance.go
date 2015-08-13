package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewAccountBalanceRequestSegment(account domain.AccountConnection, allAccounts bool) *AccountBalanceRequestSegment {
	a := &AccountBalanceRequestSegment{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.Segment = NewBasicSegment(1, a)
	return a
}

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

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountBalanceResponseSegment

type AccountBalanceResponseSegment struct {
	Segment
	AccountConnection  *element.AccountConnectionDataElement
	AccountProductName *element.AlphaNumericDataElement
	AccountCurrency    *element.CurrencyDataElement
	BookedBalance      *element.BalanceDataElement
	EarmarkedBalance   *element.BalanceDataElement
	CreditLimit        *element.AmountDataElement
	AvailableAmount    *element.AmountDataElement
	UsedAmount         *element.AmountDataElement
	BookingDate        *element.DateDataElement
	BookingTime        *element.TimeDataElement
	DueDate            *element.DateDataElement
}

func (a *AccountBalanceResponseSegment) version() int         { return 5 }
func (a *AccountBalanceResponseSegment) id() string           { return "HISAL" }
func (a *AccountBalanceResponseSegment) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegment) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegment) AccountBalance() domain.AccountBalance {
	balance := domain.AccountBalance{
		Account:       a.AccountConnection.Val(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance
}

func (a *AccountBalanceResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AccountProductName,
		a.AccountCurrency,
		a.BookedBalance,
		a.EarmarkedBalance,
		a.CreditLimit,
		a.AvailableAmount,
		a.UsedAmount,
		a.BookingDate,
		a.BookingTime,
		a.DueDate,
	}
}
