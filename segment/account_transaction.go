package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewAccountTransactionRequestSegment(account domain.AccountConnection, allAccounts bool) *AccountTransactionRequestSegment {
	a := &AccountTransactionRequestSegment{
		Account:     element.NewAccountConnection(account),
		AllAccounts: element.NewBoolean(allAccounts),
	}
	a.Segment = NewBasicSegment(1, a)
	return a
}

type AccountTransactionRequestSegment struct {
	Segment
	Account      *element.AccountConnectionDataElement
	AllAccounts  *element.BooleanDataElement
	From         *element.DateDataElement
	To           *element.DateDataElement
	MaxEntries   *element.NumberDataElement
	Aufsetzpunkt *element.AlphaNumericDataElement
}

func (a *AccountTransactionRequestSegment) SetTransactionRange(timeframe domain.Timeframe) {
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

func (a *AccountTransactionRequestSegment) version() int         { return 5 }
func (a *AccountTransactionRequestSegment) id() string           { return "HKKAZ" }
func (a *AccountTransactionRequestSegment) referencedId() string { return "" }
func (a *AccountTransactionRequestSegment) sender() string       { return senderUser }

func (a *AccountTransactionRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.Account,
		a.AllAccounts,
		a.From,
		a.To,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountTransactionResponseSegment

type AccountTransactionResponseSegment struct {
	Segment
	BookedTransactions   *element.BinaryDataElement
	UnbookedTransactions *element.BinaryDataElement
}

func (a *AccountTransactionResponseSegment) Transactions() []domain.AccountTransaction {
	return []domain.AccountTransaction{}
}

func (a *AccountTransactionResponseSegment) version() int         { return 5 }
func (a *AccountTransactionResponseSegment) id() string           { return "HIKAZ" }
func (a *AccountTransactionResponseSegment) referencedId() string { return "HKKAZ" }
func (a *AccountTransactionResponseSegment) sender() string       { return senderBank }

func (a *AccountTransactionResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.BookedTransactions,
		a.UnbookedTransactions,
	}
}
