package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/swift"
)

type AccountTransactionRequestSegment struct {
	accountTransactionRequestSegment
}

type accountTransactionRequestSegment interface {
	ClientSegment
	SetAufsetzpunkt(string)
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
	Account      *element.AccountConnectionDataElement
	AllAccounts  *element.BooleanDataElement
	From         *element.DateDataElement
	To           *element.DateDataElement
	MaxEntries   *element.NumberDataElement
	Aufsetzpunkt *element.AlphaNumericDataElement
}

func (a *AccountTransactionRequestV5) SetAufsetzpunkt(aufsetzpoint string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(aufsetzpoint, len(aufsetzpoint))
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
		a.Aufsetzpunkt,
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
	Account      *element.AccountConnectionDataElement
	AllAccounts  *element.BooleanDataElement
	From         *element.DateDataElement
	To           *element.DateDataElement
	MaxEntries   *element.NumberDataElement
	Aufsetzpunkt *element.AlphaNumericDataElement
}

func (a *AccountTransactionRequestV6) SetAufsetzpunkt(aufsetzpoint string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(aufsetzpoint, len(aufsetzpoint))
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
		a.Aufsetzpunkt,
	}
}

type AccountTransactionResponseSegment struct {
	Segment
	BookedTransactions   *element.BinaryDataElement
	UnbookedTransactions *element.BinaryDataElement
	bookedTransactions   []*swift.MT940
}

func (a *AccountTransactionResponseSegment) Transactions() []domain.AccountTransaction {
	var transactions []domain.AccountTransaction
	for _, bookedTr := range a.bookedTransactions {
		transactions = append(transactions, bookedTr.AccountTransactions()...)
	}
	return transactions
}

func (a *AccountTransactionResponseSegment) Version() int         { return 5 }
func (a *AccountTransactionResponseSegment) ID() string           { return "HIKAZ" }
func (a *AccountTransactionResponseSegment) referencedId() string { return "HKKAZ" }
func (a *AccountTransactionResponseSegment) sender() string       { return senderBank }

func (a *AccountTransactionResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.BookedTransactions,
		a.UnbookedTransactions,
	}
}
