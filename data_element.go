package hbci

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type DataElement interface {
	Value() interface{}
	Type() DataElementType
	Valid() bool
	Length() int
	String() string
}

type DataElementGroup interface {
	DataElement
	GroupDataElements() []DataElement
}

type DataElementType int

const (
	AlphaNumericDE DataElementType = iota << 1
	TextDE
	NumberDE
	DigitDE
	FloatDE
	DTAUSCharsetDE
	BinaryDE
	// Derived types
	BooleanDE
	DateDE
	VirtualDateDE
	TimeDE
	IdentificationDE
	CountryCodeDE
	CurrencyDE
	ValueDE
	// Multiple used element
	AmountGDEG
	BankIdentificationGDEG
	AccountConnectionGDEG
	BalanceGDEG
	AddressGDEG
	// DataElementGroups
	SegmentHeaderDEG
)

var typeName = map[DataElementType]string{
	AlphaNumericDE: "an",
	TextDE:         "txt",
	NumberDE:       "num",
	DigitDE:        "dig",
	FloatDE:        "float",
	DTAUSCharsetDE: "dta",
	BinaryDE:       "bin",
	// Derived types
	BooleanDE:        "jn",
	DateDE:           "dat",
	VirtualDateDE:    "vdat",
	TimeDE:           "tim",
	IdentificationDE: "id",
	CountryCodeDE:    "ctr",
	CurrencyDE:       "cur",
	ValueDE:          "wrt",
	// Multiple used element
	AmountGDEG:             "btg",
	BankIdentificationGDEG: "kik",
	AccountConnectionGDEG:  "ktv",
	BalanceGDEG:            "sdo",
	AddressGDEG:            "addr",
	// DataElementGroups
	SegmentHeaderDEG: "Segmentkopf",
}

func (d DataElementType) String() string {
	s := typeName[d]
	if s == "" {
		return fmt.Sprintf("DataElementType%d", int(d))
	}
	return s
}

func NewGroupDataElementGroup(typ DataElementType, elementCount int, elements ...DataElement) *GroupDataElementGroup {
	return &GroupDataElementGroup{elements: elements, elementCount: elementCount, typ: typ}
}

// GroupDataElementGroup defines a group of GroupDataElements.
// It implements the DataElement and the DataElementGroup interface
type GroupDataElementGroup struct {
	elements     []DataElement
	typ          DataElementType
	elementCount int
}

// Value returns the values of all GroupDataElements as []interface{}
func (g *GroupDataElementGroup) Value() interface{} {
	values := make([]interface{}, len(g.elements))
	for i, elem := range g.elements {
		values[i] = elem.Value()
	}
	return values
}
func (g *GroupDataElementGroup) Type() DataElementType { return g.typ }
func (g *GroupDataElementGroup) Valid() bool {
	if g.elementCount != len(g.elements) {
		return false
	}
	for _, elem := range g.elements {
		if !elem.Valid() {
			return false
		}
	}
	return true
}

func (g *GroupDataElementGroup) Length() int {
	length := 0
	for _, elem := range g.elements {
		length += elem.Length()
	}
	return length
}
func (g *GroupDataElementGroup) String() string {
	elementStrings := make([]string, len(g.elements))
	for i, e := range g.elements {
		elementStrings[i] = e.String()
	}
	return strings.Join(elementStrings, ":")
}

func (g *GroupDataElementGroup) GroupDataElements() []DataElement {
	return g.elements
}

func NewDataElement(typ DataElementType, value interface{}, maxLength int) DataElement {
	return &dataElement{value, typ, maxLength}
}

type dataElement struct {
	val       interface{}
	typ       DataElementType
	maxLength int
}

func (d *dataElement) Value() interface{}    { return d.val }
func (d *dataElement) Type() DataElementType { return d.typ }
func (d *dataElement) Valid() bool           { return d.Length() <= d.maxLength }
func (d *dataElement) Length() int           { return len(d.String()) }
func (d *dataElement) String() string        { return fmt.Sprintf("%v", d.val) }

func NewAlphaNumericDataElement(val string, maxLength int) *AlphaNumericDataElement {
	return &AlphaNumericDataElement{&dataElement{val, AlphaNumericDE, maxLength}}
}

type AlphaNumericDataElement struct {
	*dataElement
}

func (a *AlphaNumericDataElement) Val() string { return a.val.(string) }

func NewDigitDataElement(val, maxLength int) *DigitDataElement {
	return &DigitDataElement{&dataElement{val, DigitDE, maxLength}}
}

type DigitDataElement struct {
	*dataElement
}

func (d *DigitDataElement) Val() int { return d.val.(int) }

func (d *DigitDataElement) String() string {
	fmtString := fmt.Sprintf("%%0%dd", d.maxLength)
	return fmt.Sprintf(fmtString, d.Val())
}

func NewNumberDataElement(val, maxLength int) *NumberDataElement {
	return &NumberDataElement{&dataElement{val, NumberDE, maxLength}}
}

type NumberDataElement struct {
	*dataElement
}

func NewFloatDataElement(val float64, maxLength int) *FloatDataElement {
	return &FloatDataElement{&dataElement{val, FloatDE, maxLength}}
}

type FloatDataElement struct {
	*dataElement
}

func (f *FloatDataElement) Val() float64 { return f.val.(float64) }
func (f *FloatDataElement) String() string {
	str := strconv.FormatFloat(f.Val(), 'f', -1, 64)
	str = strings.Replace(str, ".", ",", 1)
	if !strings.Contains(str, ",") {
		str = str + ","
	}
	return str
}

type DtausCharsetDataElement struct {
	*dataElement
}

func NewDtausCharsetDataElement(data []byte, maxLength int) *DtausCharsetDataElement {
	return &DtausCharsetDataElement{&dataElement{data, DTAUSCharsetDE, maxLength}}
}

type BinaryDataElement struct {
	*dataElement
}

func NewBinaryDataElement(data []byte, maxLength int) *BinaryDataElement {
	return &BinaryDataElement{&dataElement{data, BinaryDE, maxLength}}
}

func (b *BinaryDataElement) Val() []byte {
	return b.Value().([]byte)
}

func (b *BinaryDataElement) String() string {
	return fmt.Sprintf("@%d@%s", len(b.Val()), b.Val())
}

func NewBooleanDataElement(val bool) *BooleanDataElement {
	return &BooleanDataElement{&dataElement{val, BooleanDE, 1}}
}

type BooleanDataElement struct {
	*dataElement
}

func (b *BooleanDataElement) Val() bool {
	return b.Value().(bool)
}

func (b *BooleanDataElement) String() string {
	if b.Val() {
		return "J"
	} else {
		return "N"
	}
}

func NewDateDataElement(date time.Time) *DateDataElement {
	return &DateDataElement{&dataElement{date, DateDE, 8}}
}

type DateDataElement struct {
	*dataElement
}

func (d *DateDataElement) Val() time.Time {
	return d.Value().(time.Time)
}

func (d *DateDataElement) String() string {
	return d.Val().Format("20060102")
}

func (d *DateDataElement) Valid() bool {
	return !d.Val().IsZero()
}

func NewVirtualDateDataElement(date int) *VirtualDateDataElement {
	n := NewNumberDataElement(date, 8)
	n.typ = VirtualDateDE
	return &VirtualDateDataElement{n}
}

type VirtualDateDataElement struct {
	*NumberDataElement
}

func (v *VirtualDateDataElement) Valid() bool {
	return v.Length() == 8
}

func NewTimeDataElement(date time.Time) *TimeDataElement {
	return &TimeDataElement{&dataElement{date, DateDE, 6}}
}

type TimeDataElement struct {
	*dataElement
}

func (t *TimeDataElement) Val() time.Time {
	return t.Value().(time.Time)
}

func (t *TimeDataElement) String() string {
	return t.Val().Format("150405")
}

func (t *TimeDataElement) Valid() bool {
	return !t.Val().IsZero()
}

func NewIdentificationDataElement(id string) *IdentificationDataElement {
	an := NewAlphaNumericDataElement(id, 30)
	an.typ = IdentificationDE
	return &IdentificationDataElement{an}
}

type IdentificationDataElement struct {
	*AlphaNumericDataElement
}

func NewCountryCodeDataElement(code int) *CountryCodeDataElement {
	d := NewDigitDataElement(code, 3)
	d.typ = CountryCodeDE
	return &CountryCodeDataElement{d}
}

type CountryCodeDataElement struct {
	*DigitDataElement
}

func NewCurrencyDataElement(cur string) *CurrencyDataElement {
	an := NewAlphaNumericDataElement(cur, 3)
	an.typ = CurrencyDE
	return &CurrencyDataElement{an}
}

type CurrencyDataElement struct {
	*AlphaNumericDataElement
}

func (c *CurrencyDataElement) Valid() bool {
	return c.Length() == 3
}

func NewValueDataElement(val float64) *ValueDataElement {
	f := NewFloatDataElement(val, 15)
	f.typ = ValueDE
	return &ValueDataElement{f}
}

type ValueDataElement struct {
	*FloatDataElement
}

// GroupDataElementGroups

func NewAmountDataElement(value float64, currency string) *AmountDataElement {
	g := NewGroupDataElementGroup(AmountGDEG, 2, NewValueDataElement(value), NewCurrencyDataElement(currency))
	return &AmountDataElement{g}
}

type AmountDataElement struct {
	*GroupDataElementGroup
}

func (a *AmountDataElement) Val() (value float64, currency string) {
	return a.elements[0].Value().(float64), a.elements[1].Value().(string)
}

func NewBankIndentificationDataElementWithBankId(countryCode int, bankId string) *BankIdentificationDataElement {
	g := NewGroupDataElementGroup(BankIdentificationGDEG, 2, NewCountryCodeDataElement(countryCode), NewAlphaNumericDataElement(bankId, 30))
	return &BankIdentificationDataElement{g}
}

type BankIdentificationDataElement struct {
	*GroupDataElementGroup
}

func NewAccountConnectionDataElement(accountId string, subAccountCharacteristic string, countryCode int, bankId string) *AccountConnectionDataElement {
	g := NewGroupDataElementGroup(
		AccountConnectionGDEG,
		4,
		NewIdentificationDataElement(accountId),
		NewIdentificationDataElement(subAccountCharacteristic),
		NewCountryCodeDataElement(countryCode),
		NewAlphaNumericDataElement(bankId, 30),
	)
	return &AccountConnectionDataElement{g}
}

type AccountConnectionDataElement struct {
	*GroupDataElementGroup
}

type DebitCreditIndicator int

const (
	Debit  DebitCreditIndicator = iota // Soll
	Credit                             // Haben
)

type Balance struct {
	Value    float64
	Currency string
}

func NewBalanceDataElement(balance Balance, date time.Time) *BalanceDataElement {
	var debitCredit string
	if balance.Value < 0 {
		debitCredit = "D"
	} else {
		debitCredit = "C"
	}
	g := NewGroupDataElementGroup(
		BalanceGDEG,
		5,
		NewAlphaNumericDataElement(debitCredit, 1),
		NewValueDataElement(math.Abs(balance.Value)),
		NewCurrencyDataElement(balance.Currency),
		NewDateDataElement(date),
		NewTimeDataElement(date),
	)
	return &BalanceDataElement{g}
}

type BalanceDataElement struct {
	*GroupDataElementGroup
}

func (b *BalanceDataElement) Balance() Balance {
	sign := b.elements[0].Value().(string)
	val := b.elements[1].Value().(float64)
	if sign == "D" {
		val = -val
	}
	currency := b.elements[2].Value().(string)
	balance := Balance{
		Value:    val,
		Currency: currency,
	}
	return balance
}

func (b *BalanceDataElement) Date() time.Time {
	return b.elements[3].Value().(time.Time)
}

type Address struct {
	Name1       string
	Name2       string
	Street      string
	PLZ         string
	City        string
	CountryCode int
	Phone       string
	Fax         string
	Email       string
}

func NewAddressDataElement(address Address) *AddressDataElement {
	g := NewGroupDataElementGroup(
		AddressGDEG,
		9,
		NewAlphaNumericDataElement(address.Name1, 35),
		NewAlphaNumericDataElement(address.Name2, 35),
		NewAlphaNumericDataElement(address.Street, 35),
		NewAlphaNumericDataElement(address.PLZ, 10),
		NewAlphaNumericDataElement(address.City, 35),
		NewCountryCodeDataElement(address.CountryCode),
		NewAlphaNumericDataElement(address.Phone, 35),
		NewAlphaNumericDataElement(address.Fax, 35),
		NewAlphaNumericDataElement(address.Email, 35),
	)
	return &AddressDataElement{g}
}

type AddressDataElement struct {
	*GroupDataElementGroup
}

func (a *AddressDataElement) Address() Address {
	return Address{
		Name1:       a.elements[0].Value().(string),
		Name2:       a.elements[1].Value().(string),
		Street:      a.elements[2].Value().(string),
		PLZ:         a.elements[3].Value().(string),
		City:        a.elements[4].Value().(string),
		CountryCode: a.elements[5].Value().(int),
		Phone:       a.elements[6].Value().(string),
		Fax:         a.elements[7].Value().(string),
		Email:       a.elements[8].Value().(string),
	}
}
