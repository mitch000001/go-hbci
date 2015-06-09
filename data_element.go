package hbci

import (
	"fmt"
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
	AlphaNumeric DataElementType = iota << 1
	Text
	Number
	Digit
	Float
	DTAUSCharset
	Binary
	// Derived types
	Boolean
	Date
	VirtualDate
	Time
	Identification
	CountryCode
	Currency
	Value
	// Multiple used element
	Amount
	BankIdentification
	AccountConnection
	Balance
	Address
	// DataElementGroups
	SegmentHeaderType
)

var typeName = map[DataElementType]string{
	AlphaNumeric: "an",
	Text:         "txt",
	Number:       "num",
	Digit:        "dig",
	Float:        "float",
	DTAUSCharset: "dta",
	Binary:       "bin",
	// Derived types
	Boolean:        "jn",
	Date:           "dat",
	VirtualDate:    "vdat",
	Time:           "tim",
	Identification: "id",
	CountryCode:    "ctr",
	Currency:       "cur",
	Value:          "wrt",
	// Multiple used element
	Amount:             "btg",
	BankIdentification: "kik",
	AccountConnection:  "ktv",
	Balance:            "sdo",
	Address:            "addr",
	// DataElementGroups
	SegmentHeaderType: "Segmentkopf",
}

func (d DataElementType) String() string {
	s := typeName[d]
	if s == "" {
		return fmt.Sprintf("DataElementType%d", int(d))
	}
	return s
}

func Optional(dataElement DataElement) *OptionalDataElement {
	return &OptionalDataElement{dataElement}
}

type OptionalDataElement struct {
	DataElement
}

func (o *OptionalDataElement) Value() interface{}    { return o.Value() }
func (o *OptionalDataElement) Type() DataElementType { return o.Type() }
func (o *OptionalDataElement) Valid() bool           { return o.Valid() }
func (o *OptionalDataElement) Length() int           { return o.Length() }
func (o *OptionalDataElement) String() string        { return o.String() }

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
	return &AlphaNumericDataElement{&dataElement{val, AlphaNumeric, maxLength}}
}

type AlphaNumericDataElement struct {
	*dataElement
}

func (a *AlphaNumericDataElement) Val() string { return a.val.(string) }

func NewDigitDataElement(val, maxLength int) *DigitDataElement {
	return &DigitDataElement{&dataElement{val, Digit, maxLength}}
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
	return &NumberDataElement{&dataElement{val, Number, maxLength}}
}

type NumberDataElement struct {
	*dataElement
}

func NewFloatDataElement(val float64, maxLength int) *FloatDataElement {
	return &FloatDataElement{&dataElement{val, Float, maxLength}}
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
	return &DtausCharsetDataElement{&dataElement{data, DTAUSCharset, maxLength}}
}

type BinaryDataElement struct {
	*dataElement
}

func NewBinaryDataElement(data []byte, maxLength int) *BinaryDataElement {
	return &BinaryDataElement{&dataElement{data, Binary, maxLength}}
}

func (b *BinaryDataElement) Val() []byte {
	return b.Value().([]byte)
}

func (b *BinaryDataElement) String() string {
	return fmt.Sprintf("@%d@%s", len(b.Val()), b.Val())
}

func NewBooleanDataElement(val bool) *BooleanDataElement {
	return &BooleanDataElement{&dataElement{val, Boolean, 1}}
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
	return &DateDataElement{&dataElement{date, Date, 8}}
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
	n.typ = VirtualDate
	return &VirtualDateDataElement{n}
}

type VirtualDateDataElement struct {
	*NumberDataElement
}

func (v *VirtualDateDataElement) Valid() bool {
	return v.Length() == 8
}

func NewTimeDataElement(date time.Time) *TimeDataElement {
	return &TimeDataElement{&dataElement{date, Date, 6}}
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
	an.typ = Identification
	return &IdentificationDataElement{an}
}

type IdentificationDataElement struct {
	*AlphaNumericDataElement
}

func NewCountryCodeDataElement(code int) *CountryCodeDataElement {
	d := NewDigitDataElement(code, 3)
	d.typ = CountryCode
	return &CountryCodeDataElement{d}
}

type CountryCodeDataElement struct {
	*DigitDataElement
}

func NewCurrencyDataElement(cur string) *CurrencyDataElement {
	an := NewAlphaNumericDataElement(cur, 3)
	an.typ = Currency
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
	f.typ = Value
	return &ValueDataElement{f}
}

type ValueDataElement struct {
	*FloatDataElement
}

func NewAmountDataElement(value float64, currency string) *AmountDataElement {
	g := NewGroupDataElementGroup(Amount, 2, NewValueDataElement(value), NewCurrencyDataElement(currency))
	return &AmountDataElement{g}
}

type AmountDataElement struct {
	*GroupDataElementGroup
}

func (a *AmountDataElement) Val() (value float64, currency string) {
	return a.elements[0].Value().(float64), a.elements[1].Value().(string)
}

func NewBankIndentificationDataElement(countryCode int) *BankIdentificationDataElement {
	g := NewGroupDataElementGroup(BankIdentification, 2, NewCountryCodeDataElement(countryCode))
	return &BankIdentificationDataElement{g}
}

func NewBankIndentificationDataElementWithBankId(countryCode int, bankId string) *BankIdentificationDataElement {
	g := NewGroupDataElementGroup(BankIdentification, 2, NewCountryCodeDataElement(countryCode), NewAlphaNumericDataElement(bankId, 30))
	return &BankIdentificationDataElement{g}
}

type BankIdentificationDataElement struct {
	*GroupDataElementGroup
}

func (b *BankIdentificationDataElement) Valid() bool {
	if len(b.elements) == 1 {
		elem := b.elements[0]
		return elem.Type() == CountryCode && elem.Valid()
	} else {
		return b.GroupDataElementGroup.Valid()
	}
}

type AccountConnectionDataElement struct {
	*GroupDataElementGroup
}
