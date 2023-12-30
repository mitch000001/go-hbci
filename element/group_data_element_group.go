package element

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

// NewAmount returns a new AmountDataElement
func NewAmount(value float64, currency string) *AmountDataElement {
	a := &AmountDataElement{
		Amount:   NewValue(value),
		Currency: NewCurrency(currency),
	}
	a.DataElement = NewGroupDataElementGroup(amountGDEG, 2, a)
	return a
}

// An AmountDataElement represents a value with a currency
type AmountDataElement struct {
	DataElement
	Amount   *ValueDataElement
	Currency *CurrencyDataElement
}

// Elements returns the child elements of the group
func (a *AmountDataElement) Elements() []DataElement {
	return []DataElement{
		a.Amount,
		a.Currency,
	}
}

// Val returns the value of a as a domain.Amount
func (a *AmountDataElement) Val() domain.Amount {
	return domain.Amount{
		Amount:   a.Amount.Val(),
		Currency: a.Currency.Val(),
	}
}

// UnmarshalHBCI unmarshals value into a
func (a *AmountDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("malformed marshaled value")
	}
	a.Amount = &ValueDataElement{}
	err = a.Amount.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	a.Currency = &CurrencyDataElement{}
	err = a.Currency.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	a.DataElement = NewGroupDataElementGroup(amountGDEG, 2, a)
	return nil
}

// NewBankIdentification returns a new BankIndentificationDataElement
func NewBankIdentification(bankID domain.BankID) *BankIdentificationDataElement {
	b := &BankIdentificationDataElement{
		CountryCode: NewCountryCode(bankID.CountryCode),
		BankID:      NewAlphaNumeric(bankID.ID, 30),
	}
	b.DataElement = NewGroupDataElementGroup(bankIdentificationGDEG, 2, b)
	return b
}

// A BankIdentificationDataElement represents the identification for a bank institute
type BankIdentificationDataElement struct {
	DataElement
	CountryCode *CountryCodeDataElement
	BankID      *AlphaNumericDataElement
}

// Val returns the value of b as domain.BankID
func (b *BankIdentificationDataElement) Val() domain.BankID {
	return domain.BankID{
		CountryCode: b.CountryCode.Val(),
		ID:          b.BankID.Val(),
	}
}

// Elements returns all child elements of b
func (b *BankIdentificationDataElement) Elements() []DataElement {
	return []DataElement{
		b.CountryCode,
		b.BankID,
	}
}

// UnmarshalHBCI unmarshals value into b
func (b *BankIdentificationDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 2 {
		return fmt.Errorf("malformed marshaled value: less than 2 elements: %s", elements)
	}
	countryCode := &CountryCodeDataElement{}
	err = countryCode.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	b.CountryCode = countryCode
	b.BankID = NewAlphaNumeric(charset.ToUTF8(elements[1]), 30)
	return nil
}

// NewAccountConnection returns a new AccountConnectionDataElement
func NewAccountConnection(conn domain.AccountConnection) *AccountConnectionDataElement {
	a := &AccountConnectionDataElement{
		AccountID:                 NewIdentification(conn.AccountID),
		SubAccountCharacteristics: NewIdentification(conn.SubAccountCharacteristics),
		CountryCode:               NewCountryCode(conn.CountryCode),
		BankID:                    NewAlphaNumeric(conn.BankID, 30),
	}
	a.DataElement = NewGroupDataElementGroup(accountConnectionGDEG, 4, a)
	return a
}

// AccountConnectionDataElement represents a bank account
type AccountConnectionDataElement struct {
	DataElement
	AccountID                 *IdentificationDataElement
	SubAccountCharacteristics *IdentificationDataElement
	CountryCode               *CountryCodeDataElement
	BankID                    *AlphaNumericDataElement
}

// Elements returns all child elements of a
func (a *AccountConnectionDataElement) Elements() []DataElement {
	return []DataElement{
		a.AccountID,
		a.SubAccountCharacteristics,
		a.CountryCode,
		a.BankID,
	}
}

// UnmarshalHBCI unmarshals value into a
func (a *AccountConnectionDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 4 {
		return fmt.Errorf("malformed AccountConnection")
	}
	countryCode, err := strconv.Atoi(charset.ToUTF8(elements[2]))
	if err != nil {
		return fmt.Errorf("%T: Malformed CountryCode: %q", a, elements[2])
	}
	accountConnection := domain.AccountConnection{
		AccountID:                 charset.ToUTF8(elements[0]),
		SubAccountCharacteristics: charset.ToUTF8(elements[1]),
		CountryCode:               countryCode,
		BankID:                    charset.ToUTF8(elements[3]),
	}
	*a = *NewAccountConnection(accountConnection)
	return nil
}

// Val returns the value of a as domain.AccountConnection
func (a *AccountConnectionDataElement) Val() domain.AccountConnection {
	return domain.AccountConnection{
		AccountID:                 a.AccountID.Val(),
		SubAccountCharacteristics: a.SubAccountCharacteristics.Val(),
		CountryCode:               a.CountryCode.Val(),
		BankID:                    a.BankID.Val(),
	}
}

// NewInternationalAccountConnection returns a new InternationalAccountConnectionDataElement
func NewInternationalAccountConnection(conn domain.InternationalAccountConnection) *InternationalAccountConnectionDataElement {
	i := &InternationalAccountConnectionDataElement{
		IBAN:                      NewAlphaNumeric(conn.IBAN, 34),
		BIC:                       NewAlphaNumeric(conn.BIC, 11),
		AccountID:                 NewIdentification(conn.AccountID),
		SubAccountCharacteristics: NewIdentification(conn.SubAccountCharacteristics),
		BankID:                    NewBankIdentification(conn.BankID),
	}
	i.DataElement = NewGroupDataElementGroup(internationalAccountConnectionGDEG, 5, i)
	return i
}

// An InternationalAccountConnectionDataElement represents an international bank account
type InternationalAccountConnectionDataElement struct {
	DataElement
	IBAN                      *AlphaNumericDataElement
	BIC                       *AlphaNumericDataElement
	AccountID                 *IdentificationDataElement
	SubAccountCharacteristics *IdentificationDataElement
	BankID                    *BankIdentificationDataElement
}

// Elements returns the child elements of i
func (i *InternationalAccountConnectionDataElement) Elements() []DataElement {
	return []DataElement{
		i.IBAN,
		i.BIC,
		i.AccountID,
		i.SubAccountCharacteristics,
		i.BankID,
	}
}

// UnmarshalHBCI unmarshals value into i
func (i *InternationalAccountConnectionDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed AccountConnection")
	}
	i.DataElement = NewGroupDataElementGroup(internationalAccountConnectionGDEG, 5, i)
	if len(elements) > 0 && len(elements[0]) > 0 {
		i.IBAN = &AlphaNumericDataElement{}
		err = i.IBAN.UnmarshalHBCI(elements[0])
		if err != nil {
			return fmt.Errorf("malformed IBAN: %w", err)
		}
	}
	if len(elements) > 1 && len(elements[1]) > 0 {
		i.BIC = &AlphaNumericDataElement{}
		err = i.BIC.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("malformed BIC: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		i.AccountID = &IdentificationDataElement{}
		err = i.AccountID.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("malformed AccountID: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		i.SubAccountCharacteristics = &IdentificationDataElement{}
		err = i.SubAccountCharacteristics.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("malformed SubAccountCharacteristics: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		i.BankID = &BankIdentificationDataElement{}
		err = i.BankID.UnmarshalHBCI(bytes.Join(elements[4:], []byte(":")))
		if err != nil {
			return fmt.Errorf("malformed BankID: %w", err)
		}
	}
	return nil
}

// Val returns the value of i as domain.InternationalAccountConnection
func (i *InternationalAccountConnectionDataElement) Val() domain.InternationalAccountConnection {
	conn := domain.InternationalAccountConnection{}
	if i.IBAN != nil {
		conn.IBAN = i.IBAN.Val()
	}
	if i.BIC != nil {
		conn.BIC = i.BIC.Val()
	}
	if i.AccountID != nil {
		conn.AccountID = i.AccountID.Val()
	}
	if i.SubAccountCharacteristics != nil {
		conn.SubAccountCharacteristics = i.SubAccountCharacteristics.Val()
	}
	if i.BankID != nil {
		conn.BankID = i.BankID.Val()
	}
	return conn
}

// NewBalance returns a new BalanceDataElement
func NewBalance(amount domain.Amount, date time.Time, withTime bool) *BalanceDataElement {
	var debitCredit string
	if amount.Amount < 0 {
		debitCredit = "D"
	} else {
		debitCredit = "C"
	}
	b := &BalanceDataElement{
		DebitCreditIndicator: NewAlphaNumeric(debitCredit, 1),
		Amount:               NewValue(math.Abs(amount.Amount)),
		Currency:             NewCurrency(amount.Currency),
		TransmissionDate:     NewDate(date),
	}
	if withTime {
		b.TransmissionTime = NewTime(date)
	}
	b.DataElement = NewGroupDataElementGroup(balanceGDEG, 5, b)
	return b
}

// A BalanceDataElement represents an account balance to a given date
type BalanceDataElement struct {
	DataElement
	DebitCreditIndicator *AlphaNumericDataElement
	Amount               *ValueDataElement
	Currency             *CurrencyDataElement
	TransmissionDate     *DateDataElement
	TransmissionTime     *TimeDataElement
}

// Elements returns all child elements of b
func (b *BalanceDataElement) Elements() []DataElement {
	return []DataElement{
		b.DebitCreditIndicator,
		b.Amount,
		b.Currency,
		b.TransmissionDate,
		b.TransmissionTime,
	}
}

// Balance returns the balance as domain.Balance
func (b *BalanceDataElement) Balance() domain.Balance {
	sign := b.DebitCreditIndicator.Val()
	val := b.Amount.Val()
	if sign == "D" {
		val = -val
	}
	currency := b.Currency.Val()
	amount := domain.Amount{
		Amount:   val,
		Currency: currency,
	}
	balance := domain.Balance{
		Amount:           amount,
		TransmissionDate: b.TransmissionDate.Val(),
	}
	if transmissionTime := b.TransmissionTime; transmissionTime != nil {
		val := transmissionTime.Val()
		balance.TransmissionTime = &val
	}
	return balance
}

// UnmarshalHBCI unmarshals value into b
func (b *BalanceDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 4 {
		return fmt.Errorf("%T: Malformed marshaled value", b)
	}
	b.DebitCreditIndicator = &AlphaNumericDataElement{}
	err = b.DebitCreditIndicator.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	b.Amount = &ValueDataElement{}
	err = b.Amount.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	b.Currency = &CurrencyDataElement{}
	err = b.Currency.UnmarshalHBCI(elements[2])
	if err != nil {
		return err
	}
	b.TransmissionDate = &DateDataElement{}
	err = b.TransmissionDate.UnmarshalHBCI(elements[3])
	if err != nil {
		return err
	}
	if len(elements) == 5 {
		b.TransmissionTime = &TimeDataElement{}
		err = b.TransmissionTime.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	b.DataElement = NewGroupDataElementGroup(balanceGDEG, 5, b)
	return nil
}

// Date returns the transmission date of the balance
func (b *BalanceDataElement) Date() time.Time {
	return b.TransmissionDate.Val()
}

// NewAddress creates a new AddressDataElement from address
func NewAddress(address domain.Address) *AddressDataElement {
	a := &AddressDataElement{
		Name1:       NewAlphaNumeric(address.Name1, 35),
		Name2:       NewAlphaNumeric(address.Name2, 35),
		Street:      NewAlphaNumeric(address.Street, 35),
		PLZ:         NewAlphaNumeric(address.PLZ, 10),
		City:        NewAlphaNumeric(address.City, 35),
		CountryCode: NewCountryCode(address.CountryCode),
		Phone:       NewAlphaNumeric(address.Phone, 35),
		Fax:         NewAlphaNumeric(address.Fax, 35),
		Email:       NewAlphaNumeric(address.Email, 35),
	}
	a.DataElement = NewGroupDataElementGroup(addressGDEG, 9, a)
	return a
}

// AddressDataElement represents an address
type AddressDataElement struct {
	DataElement
	Name1       *AlphaNumericDataElement
	Name2       *AlphaNumericDataElement
	Street      *AlphaNumericDataElement
	PLZ         *AlphaNumericDataElement
	City        *AlphaNumericDataElement
	CountryCode *CountryCodeDataElement
	Phone       *AlphaNumericDataElement
	Fax         *AlphaNumericDataElement
	Email       *AlphaNumericDataElement
}

// Elements returns all child elements of a
func (a *AddressDataElement) Elements() []DataElement {
	return []DataElement{
		a.Name1,
		a.Name2,
		a.Street,
		a.PLZ,
		a.City,
		a.CountryCode,
		a.Phone,
		a.Fax,
		a.Email,
	}
}

// Address returns the address as a domain.Address
func (a *AddressDataElement) Address() domain.Address {
	return domain.Address{
		Name1:       a.Name1.Val(),
		Name2:       a.Name2.Val(),
		Street:      a.Street.Val(),
		PLZ:         a.PLZ.Val(),
		City:        a.City.Val(),
		CountryCode: a.CountryCode.Val(),
		Phone:       a.Phone.Val(),
		Fax:         a.Fax.Val(),
		Email:       a.Email.Val(),
	}
}
