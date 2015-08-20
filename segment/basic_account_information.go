package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountInformationSegment

type AccountInformationSegment struct {
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

func (a *AccountInformationSegment) Version() int         { return 4 }
func (a *AccountInformationSegment) ID() string           { return "HIUPD" }
func (a *AccountInformationSegment) referencedId() string { return "HKVVB" }
func (a *AccountInformationSegment) sender() string       { return senderBank }

func (a *AccountInformationSegment) Account() domain.AccountInformation {
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

func (a *AccountInformationSegment) elements() []element.DataElement {
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
