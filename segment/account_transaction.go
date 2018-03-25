package segment

import (
	"fmt"
	"sort"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

var accountTransactionRequests = map[int](func(account domain.AccountConnection, allAccounts bool) *AccountTransactionRequestSegment){
	6: NewAccountTransactionRequestSegmentV6,
	5: NewAccountTransactionRequestSegmentV5,
}

func AccountTransactionRequestBuilder(versions []int) (func(account domain.AccountConnection, allAccounts bool) *AccountTransactionRequestSegment, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		builder, ok := accountTransactionRequests[version]
		if ok {
			return builder, nil
		}
	}
	return nil, fmt.Errorf("unsupported versions %v", versions)
}

func SepaAccountTransactionRequestBuilder(versions []int) (func(account domain.InternationalAccountConnection, allAccounts bool) *AccountTransactionRequestSegment, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		if version > 6 {
			continue
		}
		switch version {
		case 7:
			return NewAccountTransactionRequestSegmentV7, nil
		}
	}
	return nil, fmt.Errorf("unsupported versions %v", versions)
}

type AccountTransactionRequestSegment struct {
	accountTransactionRequestSegment
}

type accountTransactionRequestSegment interface {
	ClientSegment
	SetContinuationReference(string)
	SetTransactionRange(domain.Timeframe)
}

func NewAccountTransactionRequestSegmentV5(account domain.AccountConnection, allAccounts bool) *AccountTransactionRequestSegment {
	v5 := &AccountTransactionRequestV5{
		Account:     element.NewAccountConnection(account),
		AllAccounts: element.NewBoolean(allAccounts),
	}
	v5.ClientSegment = NewBasicSegment(1, v5)

	segment := &AccountTransactionRequestSegment{
		accountTransactionRequestSegment: v5,
	}
	return segment
}

type AccountTransactionRequestV5 struct {
	ClientSegment
	Account               *element.AccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	From                  *element.DateDataElement
	To                    *element.DateDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountTransactionRequestV5) SetContinuationReference(aufsetzpoint string) {
	a.ContinuationReference = element.NewAlphaNumeric(aufsetzpoint, len(aufsetzpoint))
}

func (a *AccountTransactionRequestV5) SetTransactionRange(timeframe domain.Timeframe) {
	from := timeframe.StartDate
	to := timeframe.EndDate
	if to.IsZero() {
		to = domain.NewShortDate(time.Now())
	}
	if from.IsZero() { // use sane defaults
		from = domain.NewShortDate(time.Now().AddDate(0, -1, 0))
	}
	a.From = element.NewDate(from.Time)
	a.To = element.NewDate(to.Time)
}

func (a *AccountTransactionRequestV5) Version() int         { return 5 }
func (a *AccountTransactionRequestV5) ID() string           { return "HKKAZ" }
func (a *AccountTransactionRequestV5) referencedId() string { return "" }
func (a *AccountTransactionRequestV5) sender() string       { return senderUser }

func (a *AccountTransactionRequestV5) elements() []element.DataElement {
	return []element.DataElement{
		a.Account,
		a.AllAccounts,
		a.From,
		a.To,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func NewAccountTransactionRequestSegmentV6(account domain.AccountConnection, allAccounts bool) *AccountTransactionRequestSegment {
	v6 := &AccountTransactionRequestV6{
		Account:     element.NewAccountConnection(account),
		AllAccounts: element.NewBoolean(allAccounts),
	}
	v6.ClientSegment = NewBasicSegment(1, v6)

	segment := &AccountTransactionRequestSegment{
		accountTransactionRequestSegment: v6,
	}
	return segment
}

type AccountTransactionRequestV6 struct {
	ClientSegment
	Account               *element.AccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	From                  *element.DateDataElement
	To                    *element.DateDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountTransactionRequestV6) SetContinuationReference(aufsetzpoint string) {
	a.ContinuationReference = element.NewAlphaNumeric(aufsetzpoint, len(aufsetzpoint))
}

func (a *AccountTransactionRequestV6) SetTransactionRange(timeframe domain.Timeframe) {
	from := timeframe.StartDate
	to := timeframe.EndDate
	if to.IsZero() {
		to = domain.NewShortDate(time.Now())
	}
	if from.IsZero() { // use sane defaults
		from = domain.NewShortDate(time.Now().AddDate(0, -1, 0))
	}
	a.From = element.NewDate(from.Time)
	a.To = element.NewDate(to.Time)
}

func (a *AccountTransactionRequestV6) Version() int         { return 6 }
func (a *AccountTransactionRequestV6) ID() string           { return "HKKAZ" }
func (a *AccountTransactionRequestV6) referencedId() string { return "" }
func (a *AccountTransactionRequestV6) sender() string       { return senderUser }

func (a *AccountTransactionRequestV6) elements() []element.DataElement {
	return []element.DataElement{
		a.Account,
		a.AllAccounts,
		a.From,
		a.To,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func NewAccountTransactionRequestSegmentV7(account domain.InternationalAccountConnection, allAccounts bool) *AccountTransactionRequestSegment {
	v7 := &AccountTransactionRequestV7{
		InternationalAccount: element.NewInternationalAccountConnection(account),
		AllAccounts:          element.NewBoolean(allAccounts),
	}
	v7.ClientSegment = NewBasicSegment(1, v7)

	segment := &AccountTransactionRequestSegment{
		accountTransactionRequestSegment: v7,
	}
	return segment
}

type AccountTransactionRequestV7 struct {
	ClientSegment
	InternationalAccount  *element.InternationalAccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	From                  *element.DateDataElement
	To                    *element.DateDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountTransactionRequestV7) SetContinuationReference(aufsetzpoint string) {
	a.ContinuationReference = element.NewAlphaNumeric(aufsetzpoint, len(aufsetzpoint))
}

func (a *AccountTransactionRequestV7) SetTransactionRange(timeframe domain.Timeframe) {
	from := timeframe.StartDate
	to := timeframe.EndDate
	if to.IsZero() {
		to = domain.NewShortDate(time.Now())
	}
	if from.IsZero() { // use sane defaults
		from = domain.NewShortDate(time.Now().AddDate(0, -1, 0))
	}
	a.From = element.NewDate(from.Time)
	a.To = element.NewDate(to.Time)
}

func (a *AccountTransactionRequestV7) Version() int         { return 7 }
func (a *AccountTransactionRequestV7) ID() string           { return "HKKAZ" }
func (a *AccountTransactionRequestV7) referencedId() string { return "" }
func (a *AccountTransactionRequestV7) sender() string       { return senderUser }

func (a *AccountTransactionRequestV7) elements() []element.DataElement {
	return []element.DataElement{
		a.InternationalAccount,
		a.AllAccounts,
		a.From,
		a.To,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

type AccountTransactionResponse interface {
	BankSegment
	Transactions() []domain.AccountTransaction
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountTransactionResponseSegment -segment_interface AccountTransactionResponse -segment_versions="AccountTransactionResponseSegmentV5:5:Segment,AccountTransactionResponseSegmentV6:6:Segment,AccountTransactionResponseSegmentV7:7:Segment"

type AccountTransactionResponseSegment struct {
	AccountTransactionResponse
}

type AccountTransactionResponseSegmentV5 struct {
	Segment
	BookedTransactions   *element.SwiftMT940DataElement
	UnbookedTransactions *element.BinaryDataElement
}

func (a *AccountTransactionResponseSegmentV5) Transactions() []domain.AccountTransaction {
	return a.BookedTransactions.Val()
}

func (a *AccountTransactionResponseSegmentV5) Version() int         { return 5 }
func (a *AccountTransactionResponseSegmentV5) ID() string           { return "HIKAZ" }
func (a *AccountTransactionResponseSegmentV5) referencedId() string { return "HKKAZ" }
func (a *AccountTransactionResponseSegmentV5) sender() string       { return senderBank }

func (a *AccountTransactionResponseSegmentV5) elements() []element.DataElement {
	return []element.DataElement{
		a.BookedTransactions,
		a.UnbookedTransactions,
	}
}

type AccountTransactionResponseSegmentV6 struct {
	Segment
	BookedTransactions   *element.SwiftMT940DataElement
	UnbookedTransactions *element.BinaryDataElement
}

func (a *AccountTransactionResponseSegmentV6) Transactions() []domain.AccountTransaction {
	return a.BookedTransactions.Val()
}

func (a *AccountTransactionResponseSegmentV6) Version() int         { return 6 }
func (a *AccountTransactionResponseSegmentV6) ID() string           { return "HIKAZ" }
func (a *AccountTransactionResponseSegmentV6) referencedId() string { return "HKKAZ" }
func (a *AccountTransactionResponseSegmentV6) sender() string       { return senderBank }

func (a *AccountTransactionResponseSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		a.BookedTransactions,
		a.UnbookedTransactions,
	}
}

type AccountTransactionResponseSegmentV7 struct {
	Segment
	BookedTransactions   *element.SwiftMT940DataElement
	UnbookedTransactions *element.BinaryDataElement
}

func (a *AccountTransactionResponseSegmentV7) Transactions() []domain.AccountTransaction {
	return a.BookedTransactions.Val()
}

func (a *AccountTransactionResponseSegmentV7) Version() int         { return 7 }
func (a *AccountTransactionResponseSegmentV7) ID() string           { return "HIKAZ" }
func (a *AccountTransactionResponseSegmentV7) referencedId() string { return "HKKAZ" }
func (a *AccountTransactionResponseSegmentV7) sender() string       { return senderBank }

func (a *AccountTransactionResponseSegmentV7) elements() []element.DataElement {
	return []element.DataElement{
		a.BookedTransactions,
		a.UnbookedTransactions,
	}
}
