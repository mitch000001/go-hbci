package element

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/charset"
)

type DataElement interface {
	Value() interface{}
	Type() DataElementType
	IsValid() bool
	Length() int
	String() string
	MarshalHBCI() ([]byte, error)
	UnmarshalHBCI([]byte) error
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

func New(typ DataElementType, value interface{}, maxLength int) DataElement {
	return &basicDataElement{value, typ, maxLength, false}
}

type basicDataElement struct {
	val       interface{}
	typ       DataElementType
	maxLength int
	optional  bool
}

func (d *basicDataElement) Value() interface{}    { return d.val }
func (d *basicDataElement) Type() DataElementType { return d.typ }
func (d *basicDataElement) IsValid() bool         { return d.Length() <= d.maxLength }
func (d *basicDataElement) Length() int           { return len(d.String()) }
func (d *basicDataElement) String() string        { return fmt.Sprintf("%v", d.val) }
func (d *basicDataElement) Optional() bool        { return d.optional }
func (d *basicDataElement) SetOptional()          { d.optional = true }
func (d *basicDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(d.String()), nil
}
func (d *basicDataElement) UnmarshalHBCI(value []byte) error {
	return fmt.Errorf("Not implemented")
}

func NewDataElementGroup(typ DataElementType, elementCount int, group DataElementGroup) DataElement {
	return &elementGroup{elements: group.GroupDataElements, elementCount: elementCount, typ: typ}
}

func NewGroupDataElementGroup(typ DataElementType, elementCount int, group GroupDataElementGroup) DataElement {
	return &elementGroup{elements: group.Elements, elementCount: elementCount, typ: typ}
}

// groupDataElementGroup defines a group of GroupDataElements.
// It implements the DataElement and the DataElementGroup interface
type elementGroup struct {
	elements     func() []DataElement
	typ          DataElementType
	elementCount int
	optional     bool
}

// Value returns the values of all GroupDataElements as []interface{}
func (g *elementGroup) Value() interface{} {
	values := make([]interface{}, len(g.elements()))
	for i, elem := range g.elements() {
		if !reflect.ValueOf(elem).IsNil() {
			values[i] = elem.Value()
		}
	}
	return values
}

func (g *elementGroup) Type() DataElementType { return g.typ }

func (g *elementGroup) IsValid() bool {
	if g.elementCount != len(g.elements()) {
		return false
	}
	for _, elem := range g.elements() {
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
	for _, elem := range g.elements() {
		if !reflect.ValueOf(elem).IsNil() {
			length += elem.Length()
		}
	}
	return length
}

func (g *elementGroup) String() string {
	elementStrings := make([]string, len(g.elements()))
	for i, e := range g.elements() {
		if !reflect.ValueOf(e).IsNil() {
			elementStrings[i] = e.String()
		}
	}
	return strings.Join(elementStrings, ":")
}

func (g *elementGroup) MarshalHBCI() ([]byte, error) {
	elementBytes := make([][]byte, len(g.elements()))
	for i, e := range g.elements() {
		if !reflect.ValueOf(e).IsNil() {
			marshaled, err := e.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			elementBytes[i] = marshaled
		}
	}
	return bytes.Join(elementBytes, []byte(":")), nil
}

func (g *elementGroup) UnmarshalHBCI(value []byte) error {
	marshaledElements := bytes.Split(value, []byte(":"))
	var errors []string
	elements := g.elements()
	for i, elem := range marshaledElements {
		if len(elem) == 0 {
			continue
		}
		elemToMarshal := reflect.New(reflect.TypeOf(elements[i]).Elem())
		err := elemToMarshal.Interface().(DataElement).UnmarshalHBCI(elem)
		if err == nil {
			elements[i] = elemToMarshal.Interface().(DataElement)
		} else {
			errors = append(errors, fmt.Sprintf("%T:%q", err, err))
		}
	}
	if len(errors) != 0 {
		return fmt.Errorf("Errors while unmarshaling elements:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

func (g *elementGroup) Optional() bool {
	return g.optional
}

func (g *elementGroup) SetOptional() {
	g.optional = true
}

func NewArrayElementGroup(typ DataElementType, min, max int, elems []DataElement) *arrayElementGroup {
	e := &arrayElementGroup{
		array:     elems,
		minLength: min,
		maxLength: max,
	}
	e.DataElement = NewDataElementGroup(typ, len(elems), e)
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

func NewAlphaNumeric(val string, maxLength int) *AlphaNumericDataElement {
	return &AlphaNumericDataElement{&basicDataElement{val, AlphaNumericDE, maxLength, false}}
}

type AlphaNumericDataElement struct {
	*basicDataElement
}

func (a *AlphaNumericDataElement) Val() string { return a.val.(string) }

func (a *AlphaNumericDataElement) IsValid() bool {
	if strings.ContainsAny(a.Val(), "\n & \r") {
		return false
	} else {
		return a.basicDataElement.IsValid()
	}
}

func (a *AlphaNumericDataElement) String() string {
	return escape(a.basicDataElement.String())
}

func (a *AlphaNumericDataElement) MarshalHBCI() ([]byte, error) {
	val := charset.ToISO8859_1(escape(a.basicDataElement.String()))
	return val, nil
}

func (a *AlphaNumericDataElement) UnmarshalHBCI(value []byte) error {
	decoded := charset.ToUTF8(value)
	unescaped := unescape(decoded)
	*a = AlphaNumericDataElement{&basicDataElement{unescaped, AlphaNumericDE, len(unescaped), false}}
	return nil
}

func NewText(val string, maxLength int) *TextDataElement {
	return &TextDataElement{&basicDataElement{val, TextDE, maxLength, false}}
}

type TextDataElement struct {
	*basicDataElement
}

func (a *TextDataElement) Val() string { return a.val.(string) }
func (a *TextDataElement) String() string {
	return escape(a.basicDataElement.String())
}

func (a *TextDataElement) MarshalHBCI() ([]byte, error) {
	val := charset.ToISO8859_1(escape(a.basicDataElement.String()))
	return val, nil
}

func (a *TextDataElement) UnmarshalHBCI(value []byte) error {
	decoded := charset.ToUTF8(value)
	unescaped := unescape(decoded)
	*a = TextDataElement{&basicDataElement{unescaped, TextDE, len(unescaped), false}}
	return nil
}

func NewDigit(val, maxLength int) *DigitDataElement {
	return &DigitDataElement{&basicDataElement{val, DigitDE, maxLength, false}}
}

type DigitDataElement struct {
	*basicDataElement
}

func (d *DigitDataElement) Val() int { return d.val.(int) }

func (d *DigitDataElement) String() string {
	fmtString := fmt.Sprintf("%%0%dd", d.maxLength)
	return fmt.Sprintf(fmtString, d.Val())
}

func (d *DigitDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(d.String()), nil
}

func (d *DigitDataElement) UnmarshalHBCI(value []byte) error {
	val, err := strconv.Atoi(charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*d = DigitDataElement{&basicDataElement{val, DigitDE, len(value), false}}
	return nil
}

func NewNumber(val, maxLength int) *NumberDataElement {
	return &NumberDataElement{&basicDataElement{val, NumberDE, maxLength, false}}
}

type NumberDataElement struct {
	*basicDataElement
}

func (n *NumberDataElement) Val() int { return n.val.(int) }

func (n *NumberDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(n.String()), nil
}

func (n *NumberDataElement) UnmarshalHBCI(value []byte) error {
	val, err := strconv.Atoi(charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*n = NumberDataElement{&basicDataElement{val, NumberDE, len(value), false}}
	return nil
}

func NewFloat(val float64, maxLength int) *FloatDataElement {
	return &FloatDataElement{&basicDataElement{val, FloatDE, maxLength, false}}
}

type FloatDataElement struct {
	*basicDataElement
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

func (f *FloatDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(f.String()), nil
}

func (f *FloatDataElement) UnmarshalHBCI(value []byte) error {
	str := strings.Replace(charset.ToUTF8(value), ",", ".", 1)
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*f = FloatDataElement{&basicDataElement{val, FloatDE, len(value), false}}
	return nil
}

func NewDtausCharset(data []byte, maxLength int) *DtausCharsetDataElement {
	b := NewBinary(data, maxLength)
	b.typ = DTAUSCharsetDE
	return &DtausCharsetDataElement{b}
}

type DtausCharsetDataElement struct {
	*BinaryDataElement
}

func (d *DtausCharsetDataElement) UnmarshalHBCI(value []byte) error {
	var bin BinaryDataElement
	err := bin.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*d = DtausCharsetDataElement{&bin}
	return nil
}

func NewBinary(data []byte, maxLength int) *BinaryDataElement {
	return &BinaryDataElement{&basicDataElement{data, BinaryDE, maxLength, false}}
}

type BinaryDataElement struct {
	*basicDataElement
}

func (b *BinaryDataElement) Val() []byte {
	return b.Value().([]byte)
}

func (b *BinaryDataElement) String() string {
	return fmt.Sprintf("@%d@%s", len(b.Val()), b.Val())
}

func (b *BinaryDataElement) MarshalHBCI() ([]byte, error) {
	var buf []byte
	binaryFormat := fmt.Sprintf("@%d@", len(b.Val()))
	buf = append(buf, charset.ToISO8859_1(binaryFormat)...)
	buf = append(buf, b.Val()...)
	return buf, nil
}

func (b *BinaryDataElement) UnmarshalHBCI(value []byte) error {
	buf := bytes.NewBuffer(value)
	r, _, err := buf.ReadRune()
	if err != nil {
		return err
	}
	if r != '@' {
		return fmt.Errorf("BinaryDataElement: Malformed input")
	}
	binSizeStr, err := buf.ReadString('@')
	if err != nil {
		return err
	}
	binSize, err := strconv.Atoi(binSizeStr[:len(binSizeStr)-1])
	if err != nil {
		return fmt.Errorf("Error while parsing binary size: %T:%v", err, err)
	}
	binData := make([]byte, binSize)
	_, err = buf.Read(binData)
	if err != nil {
		return fmt.Errorf("Error while reading binary data: %T:%v", err, err)
	}
	*b = BinaryDataElement{&basicDataElement{binData, BinaryDE, binSize, false}}
	return nil
}

func NewBoolean(val bool) *BooleanDataElement {
	return &BooleanDataElement{&basicDataElement{val, BooleanDE, 1, false}}
}

type BooleanDataElement struct {
	*basicDataElement
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

func (b *BooleanDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(b.String()), nil
}

func (b *BooleanDataElement) UnmarshalHBCI(value []byte) error {
	val := charset.ToUTF8(value)
	if val == "J" {
		*b = BooleanDataElement{&basicDataElement{true, BooleanDE, 1, false}}
	} else if val == "N" {
		*b = BooleanDataElement{&basicDataElement{false, BooleanDE, 1, false}}
	} else {
		return fmt.Errorf("BooleanDataElement: Malformed input")
	}
	return nil
}

func NewCode(val string, maxLength int, validSet []string) *CodeDataElement {
	sort.Strings(validSet)
	return &CodeDataElement{
		AlphaNumericDataElement: NewAlphaNumeric(val, maxLength),
		validSet:                validSet,
	}
}

type CodeDataElement struct {
	*AlphaNumericDataElement
	validSet []string
}

func (c *CodeDataElement) Type() DataElementType {
	return CodeDE
}

func (c *CodeDataElement) IsValid() bool {
	if sort.SearchStrings(c.validSet, c.Val()) >= len(c.validSet) {
		return false
	} else {
		return c.AlphaNumericDataElement.IsValid()
	}
}

func (a *CodeDataElement) MarshalHBCI() ([]byte, error) {
	return a.AlphaNumericDataElement.MarshalHBCI()
}

func (a *CodeDataElement) UnmarshalHBCI(value []byte) error {
	a.AlphaNumericDataElement = &AlphaNumericDataElement{}
	return a.AlphaNumericDataElement.UnmarshalHBCI(value)
}

func NewDate(date time.Time) *DateDataElement {
	return &DateDataElement{&basicDataElement{date, DateDE, 8, false}}
}

type DateDataElement struct {
	*basicDataElement
}

func (d *DateDataElement) Val() time.Time {
	return d.Value().(time.Time)
}

func (d *DateDataElement) String() string {
	return d.Val().Format("20060102")
}

func (d *DateDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(d.String()), nil
}

func (d *DateDataElement) UnmarshalHBCI(value []byte) error {
	t, err := time.Parse("20060102", charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*d = DateDataElement{&basicDataElement{t, DateDE, 8, false}}
	return nil
}

func (d *DateDataElement) IsValid() bool {
	return !d.Val().IsZero()
}

func NewVirtualDate(date int) *VirtualDateDataElement {
	n := NewNumber(date, 8)
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

func NewTime(date time.Time) *TimeDataElement {
	return &TimeDataElement{&basicDataElement{date, DateDE, 6, false}}
}

type TimeDataElement struct {
	*basicDataElement
}

func (t *TimeDataElement) Val() time.Time {
	return t.Value().(time.Time)
}

func (t *TimeDataElement) String() string {
	return t.Val().Format("150405")
}

func (t *TimeDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(t.String()), nil
}

func (t *TimeDataElement) UnmarshalHBCI(value []byte) error {
	parsedTime, err := time.Parse("150405", charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*t = TimeDataElement{&basicDataElement{parsedTime, TimeDE, 6, false}}
	return nil
}

func (t *TimeDataElement) IsValid() bool {
	return !t.Val().IsZero()
}

func NewIdentification(id string) *IdentificationDataElement {
	an := NewAlphaNumeric(id, 30)
	an.typ = IdentificationDE
	return &IdentificationDataElement{an}
}

type IdentificationDataElement struct {
	*AlphaNumericDataElement
}

func (i *IdentificationDataElement) Type() DataElementType {
	return IdentificationDE
}

func (i *IdentificationDataElement) UnmarshalHBCI(value []byte) error {
	var alpha AlphaNumericDataElement
	err := alpha.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*i = IdentificationDataElement{&alpha}
	return nil
}

func NewCountryCode(code int) *CountryCodeDataElement {
	d := NewDigit(code, 3)
	d.typ = CountryCodeDE
	return &CountryCodeDataElement{d}
}

type CountryCodeDataElement struct {
	*DigitDataElement
}

func (c *CountryCodeDataElement) Type() DataElementType {
	return CountryCodeDE
}

func (c *CountryCodeDataElement) UnmarshalHBCI(value []byte) error {
	var d DigitDataElement
	err := d.UnmarshalHBCI(value)
	if err != nil {
		return nil
	}
	*c = CountryCodeDataElement{&d}
	return nil
}

func NewCurrency(cur string) *CurrencyDataElement {
	an := NewAlphaNumeric(cur, 3)
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

func (c *CurrencyDataElement) UnmarshalHBCI(value []byte) error {
	var a AlphaNumericDataElement
	err := a.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*c = CurrencyDataElement{&a}
	return nil
}

func NewValue(val float64) *ValueDataElement {
	f := NewFloat(val, 15)
	f.typ = ValueDE
	return &ValueDataElement{f}
}

type ValueDataElement struct {
	*FloatDataElement
}

func (v *ValueDataElement) Type() DataElementType {
	return ValueDE
}

func (v *ValueDataElement) UnmarshalHBCI(value []byte) error {
	var f FloatDataElement
	err := f.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*v = ValueDataElement{&f}
	return nil
}
