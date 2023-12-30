package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
)

// NewAccountLimit creates a new Account limit
func NewAccountLimit(kind string, amount float64, currency string, days int) *AccountLimitDataElement {
	a := &AccountLimitDataElement{
		Kind:   NewAlphaNumeric(kind, 1),
		Amount: NewAmount(amount, currency),
		Days:   NewNumber(days, 3),
	}
	a.DataElement = NewDataElementGroup(accountLimitDEG, 3, a)
	return a
}

// AccountLimitDataElement defines an Account Limit
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

// Val returns the limit as domain.AccountLimit
func (a *AccountLimitDataElement) Val() domain.AccountLimit {
	return domain.AccountLimit{
		Kind:   a.Kind.Val(),
		Amount: a.Amount.Val(),
		Days:   a.Days.Val(),
	}
}

// GroupDataElements returns the grouped DataElements
func (a *AccountLimitDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.Kind,
		a.Amount,
		a.Days,
	}
}

// NewAllowedBusinessTransactions creates a new
// AllowedBusinessTransactionDataElement for the given transactions
func NewAllowedBusinessTransactions(transactions ...domain.BusinessTransaction) *AllowedBusinessTransactionsDataElement {
	var transactionDEs []DataElement
	for _, tr := range transactions {
		transactionDEs = append(transactionDEs, NewAllowedBusinessTransaction(tr))
	}
	a := &AllowedBusinessTransactionsDataElement{
		arrayElementGroup: newArrayElementGroup(allowedBusinessTransactionDEG, 0, 999, transactionDEs),
	}
	return a
}

// AllowedBusinessTransactionsDataElement represents allowed business
// transactions for a bank institute
type AllowedBusinessTransactionsDataElement struct {
	*arrayElementGroup
}

// UnmarshalHBCI unmarshals value into a
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
		arrayElementGroup: newArrayElementGroup(allowedBusinessTransactionDEG, 0, 999, transactions),
	}
	return nil
}

// AllowedBusinessTransactions returns all allowed business transactions
func (a *AllowedBusinessTransactionsDataElement) AllowedBusinessTransactions() []domain.BusinessTransaction {
	businessTransactions := make([]domain.BusinessTransaction, len(a.array))
	for i, de := range a.array {
		businessTransactions[i] = de.(*AllowedBusinessTransactionDataElement).Val()
	}
	return businessTransactions
}

// NewAllowedBusinessTransaction creates a new allowed business transaction
func NewAllowedBusinessTransaction(businessTransaction domain.BusinessTransaction) *AllowedBusinessTransactionDataElement {
	a := &AllowedBusinessTransactionDataElement{
		BusinessTransactionID: NewAlphaNumeric(businessTransaction.ID, 6),
		NeededSignatures:      NewNumber(businessTransaction.NeededSignatures, 2),
	}
	if businessTransaction.Limit != nil {
		a.Kind = NewAlphaNumeric(businessTransaction.Limit.Kind, 1)
		a.Amount = NewAmount(businessTransaction.Limit.Amount.Amount, businessTransaction.Limit.Amount.Currency)
		a.Days = NewNumber(businessTransaction.Limit.Days, 3)
	}
	a.DataElement = NewDataElementGroup(allowedBusinessTransactionDEG, 5, a)
	return a
}

// AllowedBusinessTransactionDataElement represents an allowed business
// transaction for a bank institute
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

// UnmarshalHBCI unmarshals value into a
func (a *AllowedBusinessTransactionDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 2 {
		return fmt.Errorf("malformed marshaled value")
	}
	a.DataElement = NewDataElementGroup(allowedBusinessTransactionDEG, 5, a)
	a.BusinessTransactionID = &AlphaNumericDataElement{}
	err = a.BusinessTransactionID.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	a.NeededSignatures = &NumberDataElement{}
	err = a.NeededSignatures.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		a.Kind = &AlphaNumericDataElement{}
		err = a.Kind.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		a.Amount = &AmountDataElement{}
		err = a.Amount.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		a.Days = &NumberDataElement{}
		err = a.Days.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	return nil
}

// Val returns the allowed business transaction
func (a *AllowedBusinessTransactionDataElement) Val() domain.BusinessTransaction {
	tr := domain.BusinessTransaction{
		ID:               a.BusinessTransactionID.Val(),
		NeededSignatures: a.NeededSignatures.Val(),
	}
	if a.Kind != nil {
		tr.Limit = &domain.AccountLimit{}
		tr.Limit.Kind = a.Kind.Val()
	}
	if a.Amount != nil && tr.Limit != nil {
		tr.Limit.Amount = a.Amount.Val()
	}
	if a.Days != nil && tr.Limit != nil {
		tr.Limit.Days = a.Days.Val()
	}
	return tr
}

// GroupDataElements returns the grouped DataElements
func (a *AllowedBusinessTransactionDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.BusinessTransactionID,
		a.NeededSignatures,
		a.Kind,
		a.Amount,
		a.Days,
	}
}
