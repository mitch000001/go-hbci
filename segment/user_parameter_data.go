package segment

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

type CommonUserParameterDataSegment struct {
	Segment
	UserID     *element.IdentificationDataElement
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
func (c *CommonUserParameterDataSegment) sender() string       { return senderBank }

func (c *CommonUserParameterDataSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.UserID,
		c.UPDVersion,
		c.UPDUsage,
	}
}

func (c *CommonUserParameterDataSegment) UserParameterData() domain.UserParameterData {
	return domain.UserParameterData{
		UserID:  c.UserID.Val(),
		Version: c.UPDVersion.Val(),
		Usage:   c.UPDUsage.Val(),
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
func (a *AccountInformationSegment) sender() string       { return senderBank }

func (a *AccountInformationSegment) Account() domain.AccountInformation {
	accountConnection := a.AccountConnection.Val()
	info := domain.AccountInformation{
		AccountConnection: &accountConnection,
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

func (a *AccountInformationSegment) UnmarshalHBCI(value []byte) error {
	value = bytes.TrimSuffix(value, []byte("'"))
	elements := bytes.Split(value, []byte("+"))
	header := elements[0]
	headerElems := bytes.Split(header, []byte(":"))
	num, err := strconv.Atoi(string(headerElems[1]))
	if err != nil {
		return fmt.Errorf("Malformed segment header")
	}
	if len(headerElems) == 4 {
		ref, err := strconv.Atoi(string(headerElems[3]))
		if err != nil {
			return fmt.Errorf("Malformed segment header reference: %v", err)
		}
		a.Segment = NewReferencingBasicSegment(num, ref, a)
	} else {
		a.Segment = NewBasicSegment(num, a)
	}
	elements = elements[1:]
	a.AccountConnection = &element.AccountConnectionDataElement{}
	err = a.AccountConnection.UnmarshalHBCI(elements[0])
	if err != nil {
		return fmt.Errorf("%T: Unmarshaling AccountConnection failed: %T:%v", a, err, err)
	}
	a.UserID = element.NewIdentification(string(elements[1]))
	a.AccountCurrency = element.NewCurrency(string(elements[2]))
	a.Name1 = element.NewAlphaNumeric(string(elements[3]), 27)
	if len(elements) > 4 {
		a.Name2 = element.NewAlphaNumeric(string(elements[4]), 27)
	}
	if len(elements) > 5 {
		a.AccountProductID = element.NewAlphaNumeric(string(elements[5]), 30)
	}
	if len(elements) > 6 {
		accountLimit := elements[6]
		if len(accountLimit) > 0 {
			a.AccountLimit = &element.AccountLimitDataElement{}
			err = a.AccountLimit.UnmarshalHBCI(accountLimit)
			if err != nil {
				return fmt.Errorf("%T: Unmarshaling AccountLimit failed: %T:%v", a, err, err)
			}
		}
	}
	if len(elements) > 7 {
		allowedBusinessTransactions := bytes.Join(elements[7:], []byte("+"))
		if len(allowedBusinessTransactions) > 0 {
			a.AllowedBusinessTransactions = &element.AllowedBusinessTransactionsDataElement{}
			err = a.AllowedBusinessTransactions.UnmarshalHBCI(allowedBusinessTransactions)
			if err != nil {
				return fmt.Errorf("%T: Unmarshaling AllowedBusinessTransactions failed: %T:%v", a, err, err)
			}
		}
	}
	return nil
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
