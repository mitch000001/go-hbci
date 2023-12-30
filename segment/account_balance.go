package segment

import (
	"fmt"
	"sort"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

var accountBalanceRequests = map[int]func(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest{
	5: NewAccountBalanceRequestV5,
	6: NewAccountBalanceRequestV6,
}

// AccountBalanceRequestBuilder returns the highest matching versioned segment
func AccountBalanceRequestBuilder(versions []int) (func(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		builder, ok := accountBalanceRequests[version]
		if ok {
			return builder, nil
		}
	}
	return nil, fmt.Errorf("unsupported versions %v", versions)
}

type AccountBalanceRequest interface {
	ClientSegment
	SetContinuationMark(continuationMark string)
}

func NewAccountBalanceRequestV5(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest {
	a := &AccountBalanceRequestSegmentV5{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountBalanceRequestSegmentV5 struct {
	ClientSegment
	AccountConnection     *element.AccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegmentV5) Version() int         { return 5 }
func (a *AccountBalanceRequestSegmentV5) ID() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegmentV5) referencedId() string { return "" }
func (a *AccountBalanceRequestSegmentV5) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegmentV5) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func (a *AccountBalanceRequestSegmentV5) SetContinuationMark(continuationMark string) {
	a.ContinuationReference = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountBalanceRequestV6(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest {
	a := &AccountBalanceRequestSegmentV6{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountBalanceRequestSegmentV6 struct {
	ClientSegment
	AccountConnection     *element.AccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegmentV6) Version() int         { return 6 }
func (a *AccountBalanceRequestSegmentV6) ID() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegmentV6) referencedId() string { return "" }
func (a *AccountBalanceRequestSegmentV6) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func (a *AccountBalanceRequestSegmentV6) SetContinuationMark(continuationMark string) {
	a.ContinuationReference = element.NewAlphaNumeric(continuationMark, 35)
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountBalanceResponseSegmentV5

type AccountBalanceResponse interface {
	BankSegment
	AccountBalance() domain.AccountBalance
}

type AccountBalanceResponseSegmentV5 struct {
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

func (a *AccountBalanceResponseSegmentV5) Version() int         { return 5 }
func (a *AccountBalanceResponseSegmentV5) ID() string           { return "HISAL" }
func (a *AccountBalanceResponseSegmentV5) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegmentV5) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegmentV5) AccountBalance() domain.AccountBalance {
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
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance
}

func (a *AccountBalanceResponseSegmentV5) elements() []element.DataElement {
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
