package hbci

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type DataElement interface {
	Value() interface{}
	Type() DataElementType
	IsValid() bool
	Length() int
	String() string
}

type OptionalDataElement interface {
	DataElement
	Optional() bool
}

type DataElementGroup interface {
	DataElement
	GroupDataElements() []DataElement
}

type GroupDataElementGroup interface {
	Elements() []DataElement
}

type DataElementType int

const (
	// DataElements
	AlphaNumericDE DataElementType = iota + 1
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
	SecurityMethodVersionGDEG
	AcknowlegdementParamsGDEG
	// DataElementGroups
	SegmentHeaderDEG
	ReferenceMessageDEG
	AcknowledgementDEG
	SecurityIdentificationDEG
	SecurityDateDEG
	HashAlgorithmDEG
	SignatureAlgorithmDEG
	EncryptionAlgorithmDEG
	KeyNameDEG
	CertificateDEG
	PublicKeyDEG
	SupportedLanguagesDEG
	SupportedHBCIVersionDEG
	CommunicationParameterDEG
	SupportedSecurityMethodDEG
	PinTanDEG
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
	AmountGDEG:                "btg",
	BankIdentificationGDEG:    "kik",
	AccountConnectionGDEG:     "ktv",
	BalanceGDEG:               "sdo",
	AddressGDEG:               "addr",
	SecurityMethodVersionGDEG: "Unterstützte Sicherheitsverfahren",
	AcknowlegdementParamsGDEG: "Rückmeldungsparameter",
	// DataElementGroups
	SegmentHeaderDEG:          "Segmentkopf",
	ReferenceMessageDEG:       "Bezugsnachricht",
	AcknowledgementDEG:        "Rückmeldung",
	SecurityIdentificationDEG: "Sicherheitsidentifikation, Details",
	SecurityDateDEG:           "Sicherheitsdatum und -uhrzeit",
	HashAlgorithmDEG:          "Hashalgorithmus",
	SignatureAlgorithmDEG:     "Signaturalgorithmus",
	EncryptionAlgorithmDEG:    "Verschlüsselungsalgorithmus",
	KeyNameDEG:                "Schlüsselname",
	CertificateDEG:            "Zertifikat",
	PublicKeyDEG:              "Öffentlicher Schlüssel",
	SupportedLanguagesDEG:     "Unterstützte Sprachen",
	SupportedHBCIVersionDEG:   "Unterstützte HBCI-Versionen",
	CommunicationParameterDEG: "Kommunikationsparameter",
	PinTanDEG:                 "PIN-TAN",
}

func (d DataElementType) String() string {
	s := typeName[d]
	if s == "" {
		return fmt.Sprintf("DataElementType%d", int(d))
	}
	return s
}

func NewDataElement(typ DataElementType, value interface{}, maxLength int) DataElement {
	return &dataElement{value, typ, maxLength, false}
}

type dataElement struct {
	val       interface{}
	typ       DataElementType
	maxLength int
	optional  bool
}

func (d *dataElement) Value() interface{}    { return d.val }
func (d *dataElement) Type() DataElementType { return d.typ }
func (d *dataElement) IsValid() bool         { return d.Length() <= d.maxLength }
func (d *dataElement) Length() int           { return len(d.String()) }
func (d *dataElement) String() string        { return fmt.Sprintf("%v", d.val) }
func (d *dataElement) Optional() bool        { return d.optional }
func (d *dataElement) SetOptional()          { d.optional = true }

func NewDataElementGroup(typ DataElementType, elementCount int, group DataElementGroup) DataElement {
	return &elementGroup{elements: group.GroupDataElements(), elementCount: elementCount, typ: typ}
}

func NewGroupDataElementGroup(typ DataElementType, elementCount int, group GroupDataElementGroup) DataElement {
	return &elementGroup{elements: group.Elements(), elementCount: elementCount, typ: typ}
}

// groupDataElementGroup defines a group of GroupDataElements.
// It implements the DataElement and the DataElementGroup interface
type elementGroup struct {
	elements     []DataElement
	typ          DataElementType
	elementCount int
	optional     bool
}

// Value returns the values of all GroupDataElements as []interface{}
func (g *elementGroup) Value() interface{} {
	values := make([]interface{}, len(g.elements))
	for i, elem := range g.elements {
		if !reflect.ValueOf(elem).IsNil() {
			values[i] = elem.Value()
		}
	}
	return values
}

func (g *elementGroup) Type() DataElementType { return g.typ }

func (g *elementGroup) IsValid() bool {
	if g.elementCount != len(g.elements) {
		return false
	}
	for _, elem := range g.elements {
		if !reflect.ValueOf(elem).IsNil() {
			if !elem.IsValid() {
				return false
			}
		}
	}
	return true
}

func (g *elementGroup) Length() int {
	length := 0
	for _, elem := range g.elements {
		if !reflect.ValueOf(elem).IsNil() {
			length += elem.Length()
		}
	}
	return length
}

func (g *elementGroup) String() string {
	elementStrings := make([]string, g.elementCount)
	for i, e := range g.elements {
		if !reflect.ValueOf(e).IsNil() {
			elementStrings[i] = e.String()
		}
	}
	return strings.Join(elementStrings, ":")
}

func (g *elementGroup) Optional() bool {
	return g.optional
}

func (g *elementGroup) SetOptional() {
	g.optional = true
}

func NewArrayElementGroup(typ DataElementType, min, max int, elem ...DataElement) *arrayElementGroup {
	e := &arrayElementGroup{
		array: elem,
	}
	e.DataElement = NewDataElementGroup(typ, max, e)
	return e
}

type arrayElementGroup struct {
	DataElement
	minLength int
	maxLength int
	array     []DataElement
}

func (a *arrayElementGroup) IsValid() bool {
	if len(a.array) < a.minLength || len(a.array) > a.maxLength {
		return false
	} else {
		return a.DataElement.IsValid()
	}
}

func (a *arrayElementGroup) GroupDataElements() []DataElement {
	return a.array
}

func escape(in string) string {
	replacer := strings.NewReplacer("?", "??", "@", "?@", "'", "?'", ":", "?:", "+", "?+")
	return replacer.Replace(in)
}

func unescape(in string) string {
	replacer := strings.NewReplacer("??", "?", "?@", "@", "?'", "'", "?:", ":", "?+", "+")
	return replacer.Replace(in)
}

func NewAlphaNumericDataElement(val string, maxLength int) *AlphaNumericDataElement {
	return &AlphaNumericDataElement{&dataElement{val, AlphaNumericDE, maxLength, false}}
}

type AlphaNumericDataElement struct {
	*dataElement
}

func (a *AlphaNumericDataElement) Val() string { return a.val.(string) }

func (a *AlphaNumericDataElement) IsValid() bool {
	if strings.ContainsAny(a.Val(), "\n & \r") {
		return false
	} else {
		return a.dataElement.IsValid()
	}
}

func (a *AlphaNumericDataElement) String() string {
	return escape(a.dataElement.String())
}

func NewTextDataElement(val string, maxLength int) *TextDataElement {
	return &TextDataElement{&dataElement{val, TextDE, maxLength, false}}
}

type TextDataElement struct {
	*dataElement
}

func (a *TextDataElement) Val() string { return a.val.(string) }
func (a *TextDataElement) String() string {
	return escape(a.dataElement.String())
}

func NewDigitDataElement(val, maxLength int) *DigitDataElement {
	return &DigitDataElement{&dataElement{val, DigitDE, maxLength, false}}
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
	return &NumberDataElement{&dataElement{val, NumberDE, maxLength, false}}
}

type NumberDataElement struct {
	*dataElement
}

func (n *NumberDataElement) Val() int { return n.val.(int) }

func NewFloatDataElement(val float64, maxLength int) *FloatDataElement {
	return &FloatDataElement{&dataElement{val, FloatDE, maxLength, false}}
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

func NewDtausCharsetDataElement(data []byte, maxLength int) *DtausCharsetDataElement {
	b := NewBinaryDataElement(data, maxLength)
	b.typ = DTAUSCharsetDE
	return &DtausCharsetDataElement{b}
}

type DtausCharsetDataElement struct {
	*BinaryDataElement
}

func NewDtausCharsetDataElement(data []byte, maxLength int) *DtausCharsetDataElement {
	return &DtausCharsetDataElement{&dataElement{data, DTAUSCharsetDE, maxLength, false}}
}

type BinaryDataElement struct {
	*dataElement
}

func NewBinaryDataElement(data []byte, maxLength int) *BinaryDataElement {
	return &BinaryDataElement{&dataElement{data, BinaryDE, maxLength, false}}
}

func (b *BinaryDataElement) Val() []byte {
	return b.Value().([]byte)
}

func (b *BinaryDataElement) String() string {
	return fmt.Sprintf("@%d@%s", len(b.Val()), b.Val())
}

func NewBooleanDataElement(val bool) *BooleanDataElement {
	return &BooleanDataElement{&dataElement{val, BooleanDE, 1, false}}
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
	return &DateDataElement{&dataElement{date, DateDE, 8, false}}
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

func (d *DateDataElement) IsValid() bool {
	return !d.Val().IsZero()
}

func NewVirtualDateDataElement(date int) *VirtualDateDataElement {
	n := NewNumberDataElement(date, 8)
	n.typ = VirtualDateDE
	return &VirtualDateDataElement{n}
}

// TODO: modelling a virtual date?!
type VirtualDateDataElement struct {
	*NumberDataElement
}

func (v *VirtualDateDataElement) Valid() bool {
	return v.Length() == 8
}

func NewTimeDataElement(date time.Time) *TimeDataElement {
	return &TimeDataElement{&dataElement{date, DateDE, 6, false}}
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

func (t *TimeDataElement) IsValid() bool {
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

func (i *IdentificationDataElement) Type() DataElementType {
	return IdentificationDE
}

func NewCountryCodeDataElement(code int) *CountryCodeDataElement {
	d := NewDigitDataElement(code, 3)
	d.typ = CountryCodeDE
	return &CountryCodeDataElement{d}
}

type CountryCodeDataElement struct {
	*DigitDataElement
}

func (c *CountryCodeDataElement) Type() DataElementType {
	return CountryCodeDE
}

func NewCurrencyDataElement(cur string) *CurrencyDataElement {
	an := NewAlphaNumericDataElement(cur, 3)
	an.typ = CurrencyDE
	return &CurrencyDataElement{an}
}

type CurrencyDataElement struct {
	*AlphaNumericDataElement
}

func (c *CurrencyDataElement) IsValid() bool {
	return c.Length() == 3
}

func (c *CurrencyDataElement) Type() DataElementType {
	return CurrencyDE
}

func NewValueDataElement(val float64) *ValueDataElement {
	f := NewFloatDataElement(val, 15)
	f.typ = ValueDE
	return &ValueDataElement{f}
}

type ValueDataElement struct {
	*FloatDataElement
}

func (v *ValueDataElement) Type() DataElementType {
	return ValueDE
}

// GroupDataElementGroups

func NewAmountDataElement(value float64, currency string) *AmountDataElement {
	a := &AmountDataElement{
		Amount:   NewValueDataElement(value),
		Currency: NewCurrencyDataElement(currency),
	}
	a.DataElement = NewGroupDataElementGroup(AmountGDEG, 2, a)
	return a
}

type AmountDataElement struct {
	DataElement
	Amount   *ValueDataElement
	Currency *CurrencyDataElement
}

func (a *AmountDataElement) Elements() []DataElement {
	return []DataElement{
		a.Amount,
		a.Currency,
	}
}

func (a *AmountDataElement) Val() (value float64, currency string) {
	return a.Amount.Val(), a.Currency.Val()
}

type BankId struct {
	CountryCode int
	ID          string
}

func NewBankIndentificationDataElement(bankId BankId) *BankIdentificationDataElement {
	b := &BankIdentificationDataElement{
		CountryCode: NewCountryCodeDataElement(bankId.CountryCode),
		BankID:      NewAlphaNumericDataElement(bankId.ID, 30),
	}
	b.DataElement = NewGroupDataElementGroup(BankIdentificationGDEG, 2, b)
	return b
}

type BankIdentificationDataElement struct {
	DataElement
	CountryCode *CountryCodeDataElement
	BankID      *AlphaNumericDataElement
}

func (b *BankIdentificationDataElement) Elements() []DataElement {
	return []DataElement{
		b.CountryCode,
		b.BankID,
	}
}

func NewAccountConnectionDataElement(accountId string, subAccountCharacteristic string, countryCode int, bankId string) *AccountConnectionDataElement {
	a := &AccountConnectionDataElement{
		AccountId:                 NewIdentificationDataElement(accountId),
		SubAccountCharacteristics: NewIdentificationDataElement(subAccountCharacteristic),
		CountryCode:               NewCountryCodeDataElement(countryCode),
		BankId:                    NewAlphaNumericDataElement(bankId, 30),
	}
	a.DataElement = NewGroupDataElementGroup(AccountConnectionGDEG, 4, a)
	return a
}

type AccountConnectionDataElement struct {
	DataElement
	AccountId                 *IdentificationDataElement
	SubAccountCharacteristics *IdentificationDataElement
	CountryCode               *CountryCodeDataElement
	BankId                    *AlphaNumericDataElement
}

func (a *AccountConnectionDataElement) Elements() []DataElement {
	return []DataElement{
		a.AccountId,
		a.SubAccountCharacteristics,
		a.CountryCode,
		a.BankId,
	}
}

type Balance struct {
	Amount   float64
	Currency string
}

func NewBalanceDataElement(balance Balance, date time.Time) *BalanceDataElement {
	var debitCredit string
	if balance.Amount < 0 {
		debitCredit = "D"
	} else {
		debitCredit = "C"
	}
	b := &BalanceDataElement{
		DebitCreditIndicator: NewAlphaNumericDataElement(debitCredit, 1),
		Amount:               NewValueDataElement(math.Abs(balance.Amount)),
		Currency:             NewCurrencyDataElement(balance.Currency),
		TransmissionDate:     NewDateDataElement(date),
		TransmissionTime:     NewTimeDataElement(date),
	}
	b.DataElement = NewGroupDataElementGroup(BalanceGDEG, 5, b)
	return b
}

type BalanceDataElement struct {
	DataElement
	DebitCreditIndicator *AlphaNumericDataElement
	Amount               *ValueDataElement
	Currency             *CurrencyDataElement
	TransmissionDate     *DateDataElement
	TransmissionTime     *TimeDataElement
}

func (b *BalanceDataElement) Elements() []DataElement {
	return []DataElement{
		b.DebitCreditIndicator,
		b.Amount,
		b.Currency,
		b.TransmissionDate,
		b.TransmissionTime,
	}
}

func (b *BalanceDataElement) Balance() Balance {
	sign := b.DebitCreditIndicator.Val()
	val := b.Amount.Val()
	if sign == "D" {
		val = -val
	}
	currency := b.Currency.Val()
	balance := Balance{
		Amount:   val,
		Currency: currency,
	}
	return balance
}

func (b *BalanceDataElement) Date() time.Time {
	return b.TransmissionDate.Val()
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
	a := &AddressDataElement{
		Name1:       NewAlphaNumericDataElement(address.Name1, 35),
		Name2:       NewAlphaNumericDataElement(address.Name2, 35),
		Street:      NewAlphaNumericDataElement(address.Street, 35),
		PLZ:         NewAlphaNumericDataElement(address.PLZ, 10),
		City:        NewAlphaNumericDataElement(address.City, 35),
		CountryCode: NewCountryCodeDataElement(address.CountryCode),
		Phone:       NewAlphaNumericDataElement(address.Phone, 35),
		Fax:         NewAlphaNumericDataElement(address.Fax, 35),
		Email:       NewAlphaNumericDataElement(address.Email, 35),
	}
	a.DataElement = NewGroupDataElementGroup(AddressGDEG, 9, a)
	return a
}

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

func (a *AddressDataElement) Address() Address {
	return Address{
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
