package element

import "github.com/mitch000001/go-hbci/domain"

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

func (a *AccountLimitDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.Kind,
		a.Amount,
		a.Days,
	}
}

func AllowedBusinessTransactions(transactions ...domain.BusinessTransaction) *AllowedBusinessTransactionsDataElement {
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

func (a *AllowedBusinessTransactionsDataElement) AllowedBusinessTransactions() []domain.BusinessTransaction {
	businessTransactions := make([]domain.BusinessTransaction, len(a.array))
	for i, de := range a.array {
		businessTransactions[i] = de.Value().(domain.BusinessTransaction)
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

func (a *AllowedBusinessTransactionDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.BusinessTransactionID,
		a.NeededSignatures,
		a.Kind,
		a.Amount,
		a.Days,
	}
}
