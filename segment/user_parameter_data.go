package segment

import "github.com/mitch000001/go-hbci/element"

type CommonUserParameterDataSegment struct {
	Segment
	UserId     *element.IdentificationDataElement
	UPDVersion *element.NumberDataElement
	// Status |￼Beschreibung
	// -----------------------------------------------------------------
	// 0	  | Die nicht aufgeführten Geschäftsvorfälle sind gesperrt
	//		  | (die aufgeführten Geschäftsvorfälle sind zugelassen).
	// 1 ￼ ￼  | Bei den nicht aufgeführten Geschäftsvorfällen ist anhand
	//        | der UPD keine Aussage darüber möglich, ob diese erlaubt
	//        | oder gesperrt sind. Diese Prüfung kann nur online vom
	//        | Kreditinstitutssystem vorgenommen werden.
	UPDUsage *element.NumberDataElement
}

func (c *CommonUserParameterDataSegment) version() int         { return 2 }
func (c *CommonUserParameterDataSegment) id() string           { return "HIUPA" }
func (c *CommonUserParameterDataSegment) referencedId() string { return "HKVVB" }

func (c *CommonUserParameterDataSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.UserId,
		c.UPDVersion,
		c.UPDUsage,
	}
}

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

func (a *AccountInformationSegment) version() int         { return 4 }
func (a *AccountInformationSegment) id() string           { return "HIUPD" }
func (a *AccountInformationSegment) referencedId() string { return "HKVVB" }

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
