package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewAccountInformationRequestSegment(account domain.AccountConnection, allAccounts bool) *AccountInformationRequestSegment {
	a := &AccountInformationRequestSegment{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.Segment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegment struct {
	Segment
	AccountConnection *element.AccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegment) Version() int         { return 1 }
func (a *AccountInformationRequestSegment) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegment) referencedId() string { return "" }
func (a *AccountInformationRequestSegment) sender() string       { return senderUser }

func (a *AccountInformationRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

type AccountInformationResponseSegment struct {
	Segment
	AccountConnection                *element.AccountConnectionDataElement
	AccountKind                      *element.NumberDataElement
	Name1                            *element.AlphaNumericDataElement
	Name2                            *element.AlphaNumericDataElement
	AccountProductID                 *element.AlphaNumericDataElement
	AccountCurrency                  *element.CurrencyDataElement
	OpeningDate                      *element.DateDataElement
	DebitInterest                    *element.ValueDataElement
	CreditInterest                   *element.ValueDataElement
	OverDebitInterest                *element.ValueDataElement
	CreditLimit                      *element.AmountDataElement
	ReferenceAccount                 *element.AccountConnectionDataElement
	AccountStatementShippingType     *element.NumberDataElement
	AccountStatementShippingRotation *element.NumberDataElement
	AdditionalInformation            *element.TextDataElement
	DisposalEligiblePersons          []*element.DisposalEligiblePersonDataElement
}

func (a *AccountInformationResponseSegment) version() int         { return 1 }
func (a *AccountInformationResponseSegment) id() string           { return "HIKIF" }
func (a *AccountInformationResponseSegment) referencedId() string { return "HKKIF" }
func (a *AccountInformationResponseSegment) sender() string       { return senderBank }

func (a *AccountInformationResponseSegment) elements() []element.DataElement {
	dataElements := make([]element.DataElement, len(a.DisposalEligiblePersons))
	for _, de := range a.DisposalEligiblePersons {
		dataElements = append(dataElements, de)
	}
	arrayDataElement := element.NewArrayElementGroup(element.DisposalEligiblePersonDEG, 0, 9, dataElements)
	return []element.DataElement{
		a.AccountConnection,
		a.AccountKind,
		a.Name1,
		a.Name2,
		a.AccountProductID,
		a.AccountCurrency,
		a.OpeningDate,
		a.DebitInterest,
		a.CreditInterest,
		a.OverDebitInterest,
		a.CreditLimit,
		a.ReferenceAccount,
		a.AccountStatementShippingType,
		a.AccountStatementShippingRotation,
		a.AdditionalInformation,
		arrayDataElement,
	}
}
