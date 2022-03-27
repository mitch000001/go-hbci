package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const AccountInformationID string = "HIUPD"

type AccountInformation interface {
	BankSegment
	Account() domain.AccountInformation
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountInformationSegment -segment_interface accountInformationSegment -segment_versions="AccountInformationV4:4:Segment,AccountInformationV5:5:Segment,AccountInformationV6:6:Segment"

type AccountInformationSegment struct {
	accountInformationSegment
}

type accountInformationSegment interface {
	BankSegment
	Account() domain.AccountInformation
}

type AccountInformationV4 struct {
	Segment
	AccountConnection           *element.AccountConnectionDataElement
	UserID                      *element.IdentificationDataElement
	AccountCurrency             *element.CurrencyDataElement
	Name1                       *element.AlphaNumericDataElement
	Name2                       *element.AlphaNumericDataElement
	AccountProductID            *element.AlphaNumericDataElement
	AccountLimit                *element.AccountLimitDataElement
	AllowedBusinessTransactions *element.AllowedBusinessTransactionsDataElement
}

func (a *AccountInformationV4) Version() int         { return 4 }
func (a *AccountInformationV4) ID() string           { return AccountInformationID }
func (a *AccountInformationV4) referencedId() string { return "HKVVB" }
func (a *AccountInformationV4) sender() string       { return senderBank }

func (a *AccountInformationV4) Account() domain.AccountInformation {
	accountConnection := a.AccountConnection.Val()
	info := domain.AccountInformation{
		AccountConnection: accountConnection,
		UserID:            a.UserID.Val(),
		Currency:          a.AccountCurrency.Val(),
		Name1:             a.Name1.Val(),
	}
	if a.Name2 != nil {
		info.Name2 = a.Name2.Val()
	}
	if a.AccountProductID != nil {
		info.ProductID = a.AccountProductID.Val()
	}
	if a.AccountLimit != nil {
		limit := a.AccountLimit.Val()
		info.Limit = &limit
	}
	if a.AllowedBusinessTransactions != nil {
		info.AllowedBusinessTransactions = a.AllowedBusinessTransactions.AllowedBusinessTransactions()
	}
	return info
}

func (a *AccountInformationV4) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.UserID,
		a.AccountCurrency,
		a.Name1,
		a.Name2,
		a.AccountProductID,
		a.AccountLimit,
		a.AllowedBusinessTransactions,
	}
}

type AccountInformationV5 struct {
	Segment
	AccountConnection           *element.AccountConnectionDataElement
	UserID                      *element.IdentificationDataElement
	AccountType                 *element.NumberDataElement
	AccountCurrency             *element.CurrencyDataElement
	Name1                       *element.AlphaNumericDataElement
	Name2                       *element.AlphaNumericDataElement
	AccountProductID            *element.AlphaNumericDataElement
	AccountLimit                *element.AccountLimitDataElement
	AllowedBusinessTransactions *element.AllowedBusinessTransactionsDataElement
}

func (a *AccountInformationV5) Version() int         { return 5 }
func (a *AccountInformationV5) ID() string           { return AccountInformationID }
func (a *AccountInformationV5) referencedId() string { return "HKVVB" }
func (a *AccountInformationV5) sender() string       { return senderBank }

func (a *AccountInformationV5) Account() domain.AccountInformation {
	info := domain.AccountInformation{
		UserID: a.UserID.Val(),
		Name1:  a.Name1.Val(),
	}
	if a.AccountConnection != nil {
		info.AccountConnection = a.AccountConnection.Val()
	}
	if a.AccountCurrency != nil {
		info.Currency = a.AccountCurrency.Val()
	}
	if a.Name2 != nil {
		info.Name2 = a.Name2.Val()
	}
	if a.AccountProductID != nil {
		info.ProductID = a.AccountProductID.Val()
	}
	if a.AccountLimit != nil {
		limit := a.AccountLimit.Val()
		info.Limit = &limit
	}
	if a.AllowedBusinessTransactions != nil {
		info.AllowedBusinessTransactions = a.AllowedBusinessTransactions.AllowedBusinessTransactions()
	}
	return info
}

func (a *AccountInformationV5) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.UserID,
		a.AccountType,
		a.AccountCurrency,
		a.Name1,
		a.Name2,
		a.AccountProductID,
		a.AccountLimit,
		a.AllowedBusinessTransactions,
	}
}

type AccountInformationV6 struct {
	Segment
	AccountConnection           *element.AccountConnectionDataElement
	IBAN                        *element.AlphaNumericDataElement
	UserID                      *element.IdentificationDataElement
	AccountType                 *element.NumberDataElement
	AccountCurrency             *element.CurrencyDataElement
	Name1                       *element.AlphaNumericDataElement
	Name2                       *element.AlphaNumericDataElement
	AccountProductID            *element.AlphaNumericDataElement
	AccountLimit                *element.AccountLimitDataElement
	AllowedBusinessTransactions *element.AllowedBusinessTransactionsDataElement
	AccountExtensions           *element.AlphaNumericDataElement
}

func (a *AccountInformationV6) Version() int         { return 6 }
func (a *AccountInformationV6) ID() string           { return AccountInformationID }
func (a *AccountInformationV6) referencedId() string { return "HKVVB" }
func (a *AccountInformationV6) sender() string       { return senderBank }

func (a *AccountInformationV6) Account() domain.AccountInformation {
	info := domain.AccountInformation{
		UserID: a.UserID.Val(),
		Name1:  a.Name1.Val(),
	}
	if a.AccountConnection != nil {
		info.AccountConnection = a.AccountConnection.Val()
	}
	if a.AccountCurrency != nil {
		info.Currency = a.AccountCurrency.Val()
	}
	if a.Name2 != nil {
		info.Name2 = a.Name2.Val()
	}
	if a.AccountProductID != nil {
		info.ProductID = a.AccountProductID.Val()
	}
	if a.AccountLimit != nil {
		limit := a.AccountLimit.Val()
		info.Limit = &limit
	}
	if a.AllowedBusinessTransactions != nil {
		info.AllowedBusinessTransactions = a.AllowedBusinessTransactions.AllowedBusinessTransactions()
	}
	return info
}

func (a *AccountInformationV6) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.IBAN,
		a.UserID,
		a.AccountType,
		a.AccountCurrency,
		a.Name1,
		a.Name2,
		a.AccountProductID,
		a.AccountLimit,
		a.AllowedBusinessTransactions,
		a.AccountExtensions,
	}
}
