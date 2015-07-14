package element

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/domain"
)

func NewAccountLimit(kind string, amount float64, currency string, days int) *AccountLimitDataElement {
	a := &AccountLimitDataElement{
		Kind:   NewAlphaNumeric(kind, 1),
		Amount: NewAmount(amount, currency),
		Days:   NewNumber(days, 3),
	}
	a.DataElement = NewDataElementGroup(AccountLimitDEG, 3, a)
	return a
}

type AccountLimitDataElement struct {
	DataElement
	// Code | Beschreibung
	// --------------------------
	// E	| Einzelauftragslimit
	// T	| Tageslimit
	// W	| Wochenlimit
	// M	| Monatslimit
	// Z ￼	| Zeitlimit
	Kind   *AlphaNumericDataElement
	Amount *AmountDataElement
	Days   *NumberDataElement
}

func (a *AccountLimitDataElement) Val() domain.AccountLimit {
	return domain.AccountLimit{
		Kind:   a.Kind.Val(),
		Amount: a.Amount.Val(),
		Days:   a.Days.Val(),
	}
}

func (a *AccountLimitDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.Kind,
		a.Amount,
		a.Days,
	}
}

func NewAllowedBusinessTransactions(transactions ...domain.BusinessTransaction) *AllowedBusinessTransactionsDataElement {
	var transactionDEs []DataElement
	for _, tr := range transactions {
		transactionDEs = append(transactionDEs, NewAllowedBusinessTransaction(tr))
	}
	a := &AllowedBusinessTransactionsDataElement{
		arrayElementGroup: NewArrayElementGroup(AllowedBusinessTransactionDEG, 0, 999, transactionDEs...),
	}
	return a
}

type AllowedBusinessTransactionsDataElement struct {
	*arrayElementGroup
}

func (a *AllowedBusinessTransactionsDataElement) UnmarshalHBCI(value []byte) error {
	elements := bytes.Split(value, []byte("+"))
	transactions := make([]DataElement, len(elements))
	for i, elem := range elements {
		tr := &AllowedBusinessTransactionDataElement{}
		err := tr.UnmarshalHBCI(elem)
		if err != nil {
			return err
		}
		transactions[i] = tr
	}
	*a = AllowedBusinessTransactionsDataElement{
		arrayElementGroup: NewArrayElementGroup(AllowedBusinessTransactionDEG, 0, 999, transactions...),
	}
	return nil
}

func (a *AllowedBusinessTransactionsDataElement) AllowedBusinessTransactions() []domain.BusinessTransaction {
	businessTransactions := make([]domain.BusinessTransaction, len(a.array))
	for i, de := range a.array {
		businessTransactions[i] = de.(*AllowedBusinessTransactionDataElement).Val()
	}
	return businessTransactions
}

func NewAllowedBusinessTransaction(businessTransaction domain.BusinessTransaction) *AllowedBusinessTransactionDataElement {
	a := &AllowedBusinessTransactionDataElement{
		BusinessTransactionID: NewAlphaNumeric(businessTransaction.ID, 6),
		NeededSignatures:      NewNumber(businessTransaction.NeededSignatures, 2),
		Kind:                  NewAlphaNumeric(businessTransaction.LimitKind, 1),
		Amount:                NewAmount(businessTransaction.LimitAmount.Amount, businessTransaction.LimitAmount.Currency),
		Days:                  NewNumber(businessTransaction.LimitDays, 3),
	}
	a.DataElement = NewDataElementGroup(AllowedBusinessTransactionDEG, 5, a)
	return a
}

type AllowedBusinessTransactionDataElement struct {
	DataElement
	BusinessTransactionID *AlphaNumericDataElement
	NeededSignatures      *NumberDataElement
	// Code | Beschreibung
	// --------------------------
	// E	| Einzelauftragslimit
	// T	| Tageslimit
	// W	| Wochenlimit
	// M	| Monatslimit
	// Z ￼	| Zeitlimit
	Kind   *AlphaNumericDataElement
	Amount *AmountDataElement
	Days   *NumberDataElement
}

func (a *AllowedBusinessTransactionDataElement) UnmarshalHBCI(value []byte) error {
	elements := bytes.Split(value, []byte(":"))
	businessTransaction := domain.BusinessTransaction{}
	businessTransaction.ID = string(elements[0])
	neededSignatures, err := strconv.Atoi(string(elements[1]))
	if err != nil {
		return fmt.Errorf("%T: Error while unmarshaling NeededSignatures: %T:%v", a, err, err)
	}
	businessTransaction.NeededSignatures = neededSignatures
	if len(elements) >= 3 {
		businessTransaction.LimitKind = string(elements[2])
	}
	if len(elements) >= 5 {
		amountVal, err := strconv.ParseFloat(string(elements[3]), 64)
		if err != nil {
			return fmt.Errorf("%T: Error while unmarshaling Amount: %T:%v", a, err, err)
		}
		currency := string(elements[4])
		businessTransaction.LimitAmount = domain.Amount{amountVal, currency}
	}
	if len(elements) == 6 {
		days, err := strconv.Atoi(string(elements[5]))
		if err != nil {
			return fmt.Errorf("%T: Error while unmarshaling LimitDays: %T:%v", a, err, err)
		}
		businessTransaction.LimitDays = days
	}
	*a = *NewAllowedBusinessTransaction(businessTransaction)
	return nil
}

func (a *AllowedBusinessTransactionDataElement) Val() domain.BusinessTransaction {
	tr := domain.BusinessTransaction{
		ID:               a.BusinessTransactionID.Val(),
		NeededSignatures: a.NeededSignatures.Val(),
	}
	if a.Kind != nil {
		tr.LimitKind = a.Kind.Val()
	}
	if a.Amount != nil {
		tr.LimitAmount = a.Amount.Val()
	}
	if a.Days != nil {
		tr.LimitDays = a.Days.Val()
	}
	return tr
}

func (a *AllowedBusinessTransactionDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.BusinessTransactionID,
		a.NeededSignatures,
		a.Kind,
		a.Amount,
		a.Days,
	}
}
