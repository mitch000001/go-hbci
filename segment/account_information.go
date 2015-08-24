package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

type AccountInformationRequest interface {
	ClientSegment
	SetContinuationMark(continuationMark string)
}

func NewAccountInformationRequestSegmentV1(account domain.AccountConnection, allAccounts bool) AccountInformationRequest {
	a := &AccountInformationRequestSegmentV1{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegmentV1 struct {
	ClientSegment
	AccountConnection *element.AccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegmentV1) Version() int         { return 1 }
func (a *AccountInformationRequestSegmentV1) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegmentV1) referencedId() string { return "" }
func (a *AccountInformationRequestSegmentV1) sender() string       { return senderUser }

func (a *AccountInformationRequestSegmentV1) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

func (a *AccountInformationRequestSegmentV1) SetContinuationMark(continuationMark string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountInformationRequestSegmentV2(account domain.AccountConnection, allAccounts bool) AccountInformationRequest {
	a := &AccountInformationRequestSegmentV2{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegmentV2 struct {
	ClientSegment
	AccountConnection *element.AccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegmentV2) Version() int         { return 2 }
func (a *AccountInformationRequestSegmentV2) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegmentV2) referencedId() string { return "" }
func (a *AccountInformationRequestSegmentV2) sender() string       { return senderUser }

func (a *AccountInformationRequestSegmentV2) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

func (a *AccountInformationRequestSegmentV2) SetContinuationMark(continuationMark string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountInformationRequestSegmentV3(account domain.InternationalAccountConnection, allAccounts bool) AccountInformationRequest {
	a := &AccountInformationRequestSegmentV3{
		AccountConnection: element.NewInternationalAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegmentV3 struct {
	ClientSegment
	AccountConnection *element.InternationalAccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegmentV3) Version() int         { return 3 }
func (a *AccountInformationRequestSegmentV3) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegmentV3) referencedId() string { return "" }
func (a *AccountInformationRequestSegmentV3) sender() string       { return senderUser }

func (a *AccountInformationRequestSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

func (a *AccountInformationRequestSegmentV3) SetContinuationMark(continuationMark string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountInformationRequestSegmentV4(account domain.InternationalAccountConnection, allAccounts bool) AccountInformationRequest {
	a := &AccountInformationRequestSegmentV4{
		AccountConnection: element.NewInternationalAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegmentV4 struct {
	ClientSegment
	AccountConnection *element.InternationalAccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegmentV4) Version() int         { return 4 }
func (a *AccountInformationRequestSegmentV4) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegmentV4) referencedId() string { return "" }
func (a *AccountInformationRequestSegmentV4) sender() string       { return senderUser }

func (a *AccountInformationRequestSegmentV4) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

func (a *AccountInformationRequestSegmentV4) SetContinuationMark(continuationMark string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountInformationRequestSegmentV5(account domain.InternationalAccountConnection, allAccounts bool) AccountInformationRequest {
	a := &AccountInformationRequestSegmentV5{
		AccountConnection: element.NewInternationalAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegmentV5 struct {
	ClientSegment
	AccountConnection *element.InternationalAccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegmentV5) Version() int         { return 5 }
func (a *AccountInformationRequestSegmentV5) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegmentV5) referencedId() string { return "" }
func (a *AccountInformationRequestSegmentV5) sender() string       { return senderUser }

func (a *AccountInformationRequestSegmentV5) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

func (a *AccountInformationRequestSegmentV5) SetContinuationMark(continuationMark string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountInformationRequestSegmentV6(account domain.InternationalAccountConnection, allAccounts bool) AccountInformationRequest {
	a := &AccountInformationRequestSegmentV6{
		AccountConnection: element.NewInternationalAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountInformationRequestSegmentV6 struct {
	ClientSegment
	AccountConnection *element.InternationalAccountConnectionDataElement
	AllAccounts       *element.BooleanDataElement
	MaxEntries        *element.NumberDataElement
	Aufsetzpunkt      *element.AlphaNumericDataElement
}

func (a *AccountInformationRequestSegmentV6) Version() int         { return 6 }
func (a *AccountInformationRequestSegmentV6) ID() string           { return "HKKIF" }
func (a *AccountInformationRequestSegmentV6) referencedId() string { return "" }
func (a *AccountInformationRequestSegmentV6) sender() string       { return senderUser }

func (a *AccountInformationRequestSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.Aufsetzpunkt,
	}
}

func (a *AccountInformationRequestSegmentV6) SetContinuationMark(continuationMark string) {
	a.Aufsetzpunkt = element.NewAlphaNumeric(continuationMark, 35)
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountInformationResponseSegment

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
	DisposalEligiblePersons          *element.DisposalEligiblePersonsDataElement
}

func (a *AccountInformationResponseSegment) Version() int         { return 1 }
func (a *AccountInformationResponseSegment) ID() string           { return "HIKIF" }
func (a *AccountInformationResponseSegment) referencedId() string { return "HKKIF" }
func (a *AccountInformationResponseSegment) sender() string       { return senderBank }

func (a *AccountInformationResponseSegment) elements() []element.DataElement {
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
		a.DisposalEligiblePersons,
	}
}
