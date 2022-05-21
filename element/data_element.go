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

// DataElement represent the general interface of a DataElement used by HBCI
type DataElement interface {
	// Value returns the underlying value
	Value() interface{}
	// IsValid returns true if the DataElement and all its grouped elements
	// are valid, false otherwise
	IsValid() bool
	Length() int
	String() string
	MarshalHBCI() ([]byte, error)
	UnmarshalHBCI([]byte) error
}

// OptionalDataElement represents a DataElement that can be marked as optional or required.
type OptionalDataElement interface {
	DataElement
	Optional() bool
}

// DataElementGroup represents a DataElement which consists of subelements that
// represent a logical group.
type DataElementGroup interface {
	DataElement
	// GroupDataElements returns the grouped DataElements
	GroupDataElements() []DataElement
}

// GroupDataElementGroup represents a DataElement which is composed by subelements
type GroupDataElementGroup interface {
	Elements() []DataElement
}

// New returns a new DataElement with the provided type, value and maxLengt for validation
func New(typ DataElementType, value interface{}, maxLength int) DataElement {
	return &basicDataElement{value, typ, maxLength, false}
}

type basicDataElement struct {
	val       interface{}
	typ       DataElementType
	maxLength int
	optional  bool
}

// Value returns the underlying value
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
func (d *basicDataElement) MarshalYAML() (interface{}, error) {
	return d.val, nil
}

// NewDataElementGroup returns a new DataElement that embodies group
func NewDataElementGroup(typ DataElementType, elementCount int, group DataElementGroup) DataElement {
	return &elementGroup{elements: group.GroupDataElements, elementCount: elementCount, typ: typ}
}

// NewGroupDataElementGroup returns a new DataElement which embodies group
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
	marshaled := bytes.Join(elementBytes, []byte(":"))
	marshaled = bytes.TrimRight(marshaled, ":")
	return marshaled, nil
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

// newArrayElementGroup returns a new DataElement which has multiple occurrences in its parents
func newArrayElementGroup(typ DataElementType, min, max int, elems []DataElement) *arrayElementGroup {
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
	}
	return a.DataElement.IsValid()
}

// GroupDataElements returns the grouped DataElements
func (a *arrayElementGroup) GroupDataElements() []DataElement {
	return a.array
}

func (a *arrayElementGroup) MarshalYAML() (interface{}, error) {
	return a.array, nil
}

func escape(in string) string {
	replacer := strings.NewReplacer("?", "??", "@", "?@", "'", "?'", ":", "?:", "+", "?+")
	return replacer.Replace(in)
}

func unescape(in string) string {
	replacer := strings.NewReplacer("??", "?", "?@", "@", "?'", "'", "?:", ":", "?+", "+")
	return replacer.Replace(in)
}

// NewAlphaNumeric returns a new AlphaNumeric DataElement
func NewAlphaNumeric(val string, maxLength int) *AlphaNumericDataElement {
	return &AlphaNumericDataElement{&basicDataElement{val, alphaNumericDE, maxLength, false}}
}

// An AlphaNumericDataElement represents a DataElement which contains alphanumeric characters
type AlphaNumericDataElement struct {
	*basicDataElement
}

// Val returns the data as string
func (a *AlphaNumericDataElement) Val() string { return a.val.(string) }

// IsValid returns false if a contains '\n' and '\r', true otherwise
func (a *AlphaNumericDataElement) IsValid() bool {
	if strings.ContainsAny(a.Val(), "\n & \r") {
		return false
	}
	return a.basicDataElement.IsValid()
}

func (a *AlphaNumericDataElement) String() string {
	return escape(a.basicDataElement.String())
}

// MarshalHBCI marshals a into a byte representation ready to be sent over the wire.
func (a *AlphaNumericDataElement) MarshalHBCI() ([]byte, error) {
	val := charset.ToISO8859_1(escape(a.basicDataElement.String()))
	return val, nil
}

// UnmarshalHBCI unmarshals the given bytes into a
func (a *AlphaNumericDataElement) UnmarshalHBCI(value []byte) error {
	decoded := charset.ToUTF8(value)
	unescaped := unescape(decoded)
	*a = AlphaNumericDataElement{&basicDataElement{unescaped, alphaNumericDE, len(unescaped), false}}
	return nil
}

// NewText returns a new TextDataElement
func NewText(val string, maxLength int) *TextDataElement {
	return &TextDataElement{&basicDataElement{val, textDE, maxLength, false}}
}

// TextDataElement represents a DataElement that can hold alphanumeric
// characters, but also '\n' and '\r'
type TextDataElement struct {
	*basicDataElement
}

// Val returns the value of a as a string
func (a *TextDataElement) Val() string { return a.val.(string) }
func (a *TextDataElement) String() string {
	return escape(a.basicDataElement.String())
}

// MarshalHBCI marshals a into the HBCI wire format
func (a *TextDataElement) MarshalHBCI() ([]byte, error) {
	val := charset.ToISO8859_1(escape(a.basicDataElement.String()))
	return val, nil
}

// UnmarshalHBCI unmarshals value into a
func (a *TextDataElement) UnmarshalHBCI(value []byte) error {
	decoded := charset.ToUTF8(value)
	unescaped := unescape(decoded)
	*a = TextDataElement{&basicDataElement{unescaped, textDE, len(unescaped), false}}
	return nil
}

// NewDigit returns a new DigitDataElement
func NewDigit(val, maxLength int) *DigitDataElement {
	return &DigitDataElement{&basicDataElement{val, digitDE, maxLength, false}}
}

// DigitDataElement represents numbers from 0 to 9 with leading zeros
type DigitDataElement struct {
	*basicDataElement
}

// Val returns the digit as int
func (d *DigitDataElement) Val() int { return d.val.(int) }

func (d *DigitDataElement) String() string {
	fmtString := fmt.Sprintf("%%0%dd", d.maxLength)
	return fmt.Sprintf(fmtString, d.Val())
}

// MarshalHBCI marshals d into the HBCI wire format
func (d *DigitDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(d.String()), nil
}

// UnmarshalHBCI unmarshals value into d
func (d *DigitDataElement) UnmarshalHBCI(value []byte) error {
	val, err := strconv.Atoi(charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*d = DigitDataElement{&basicDataElement{val, digitDE, len(value), false}}
	return nil
}

// NewNumber returns a new NumberDataElement
func NewNumber(val, maxLength int) *NumberDataElement {
	return &NumberDataElement{&basicDataElement{val, numberDE, maxLength, false}}
}

// A NumberDataElement represents a number containing 0 - 9 without leading zeros
type NumberDataElement struct {
	*basicDataElement
}

// Val returns the value of n as int
func (n *NumberDataElement) Val() int { return n.val.(int) }

// MarshalHBCI marshals n into the HBCI wire format
func (n *NumberDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(n.String()), nil
}

// UnmarshalHBCI unmarshals value into n
func (n *NumberDataElement) UnmarshalHBCI(value []byte) error {
	val, err := strconv.Atoi(charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*n = NumberDataElement{&basicDataElement{val, numberDE, len(value), false}}
	return nil
}

// NewFloat returns a new FloatDataElement
func NewFloat(val float64, maxLength int) *FloatDataElement {
	return &FloatDataElement{&basicDataElement{val, floatDE, maxLength, false}}
}

// FloatDataElement represents a float in HBCI
type FloatDataElement struct {
	*basicDataElement
}

// Val returns the value of f as float64
func (f *FloatDataElement) Val() float64 { return f.val.(float64) }

func (f *FloatDataElement) String() string {
	str := strconv.FormatFloat(f.Val(), 'f', -1, 64)
	str = strings.Replace(str, ".", ",", 1)
	if !strings.Contains(str, ",") {
		str = str + ","
	}
	return str
}

// MarshalHBCI marshals f into HBCI wire format
func (f *FloatDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(f.String()), nil
}

// UnmarshalHBCI unmarshals value into f
func (f *FloatDataElement) UnmarshalHBCI(value []byte) error {
	str := strings.Replace(charset.ToUTF8(value), ",", ".", 1)
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*f = FloatDataElement{&basicDataElement{val, floatDE, len(value), false}}
	return nil
}

// NewDtausCharset returns a new DtausCharsetDataElement
func NewDtausCharset(data []byte, maxLength int) *DtausCharsetDataElement {
	b := NewBinary(data, maxLength)
	b.typ = dtausCharsetDE
	return &DtausCharsetDataElement{b}
}

// DtausCharsetDataElement represents binary data in DTAUS charset
type DtausCharsetDataElement struct {
	*BinaryDataElement
}

// UnmarshalHBCI unmarshals value into d
func (d *DtausCharsetDataElement) UnmarshalHBCI(value []byte) error {
	var bin BinaryDataElement
	err := bin.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*d = DtausCharsetDataElement{&bin}
	return nil
}

// NewBinary returns a new BinaryDataElement
func NewBinary(data []byte, maxLength int) *BinaryDataElement {
	return &BinaryDataElement{&basicDataElement{data, binaryDE, maxLength, false}}
}

// A BinaryDataElement embodies binary data into a DataElement
type BinaryDataElement struct {
	*basicDataElement
}

// Val returns b as []byte
func (b *BinaryDataElement) Val() []byte {
	return b.Value().([]byte)
}

func (b *BinaryDataElement) String() string {
	return fmt.Sprintf("@%d@%s", len(b.Val()), b.Val())
}

// MarshalHBCI marshals b into HBCI wire format
func (b *BinaryDataElement) MarshalHBCI() ([]byte, error) {
	var buf []byte
	binaryFormat := fmt.Sprintf("@%d@", len(b.Val()))
	buf = append(buf, charset.ToISO8859_1(binaryFormat)...)
	buf = append(buf, b.Val()...)
	return buf, nil
}

// UnmarshalHBCI unmarshals value into b
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
	*b = BinaryDataElement{&basicDataElement{binData, binaryDE, binSize, false}}
	return nil
}

// NewBoolean returns a new BooleanDataElement
func NewBoolean(val bool) *BooleanDataElement {
	return &BooleanDataElement{&basicDataElement{val, booleanDE, 1, false}}
}

// A BooleanDataElement represents a boolean in HBCI
type BooleanDataElement struct {
	*basicDataElement
}

// Val returns the value of b as bool
func (b *BooleanDataElement) Val() bool {
	return b.Value().(bool)
}

func (b *BooleanDataElement) String() string {
	if b.Val() {
		return "J"
	}
	return "N"
}

// MarshalHBCI marshals b into the HBCI wire format
func (b *BooleanDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(b.String()), nil
}

// UnmarshalHBCI unmarshals value into b
func (b *BooleanDataElement) UnmarshalHBCI(value []byte) error {
	val := charset.ToUTF8(value)
	if val == "J" {
		*b = BooleanDataElement{&basicDataElement{true, booleanDE, 1, false}}
	} else if val == "N" {
		*b = BooleanDataElement{&basicDataElement{false, booleanDE, 1, false}}
	} else {
		return fmt.Errorf("BooleanDataElement: Malformed input")
	}
	return nil
}

// NewCode returns a new CodeDataElement
func NewCode(val string, maxLength int, validSet []string) *CodeDataElement {
	sort.Strings(validSet)
	return &CodeDataElement{
		AlphaNumericDataElement: NewAlphaNumeric(val, maxLength),
		validSet:                validSet,
	}
}

// A CodeDataElement represents a value from a know bounded set of values
type CodeDataElement struct {
	*AlphaNumericDataElement
	validSet []string
}

// Type returns the DataElementType of c
func (c *CodeDataElement) Type() DataElementType {
	return codeDE
}

// IsValid returns true if the value is in the valid set and the underlying
// AlphaNumericDataElement is valid, false otherwise.
func (c *CodeDataElement) IsValid() bool {
	if sort.SearchStrings(c.validSet, c.Val()) >= len(c.validSet) {
		return false
	}
	return c.AlphaNumericDataElement.IsValid()
}

// UnmarshalHBCI unmarshals value into a. a will always be invalid, because the
// valid set is not known when unmarshaling
func (c *CodeDataElement) UnmarshalHBCI(value []byte) error {
	c.AlphaNumericDataElement = &AlphaNumericDataElement{}
	return c.AlphaNumericDataElement.UnmarshalHBCI(value)
}

// NewDate returns a new DateDataElement
func NewDate(date time.Time) *DateDataElement {
	return &DateDataElement{&basicDataElement{date, dateDE, 8, false}}
}

// A DateDataElement represents a date within HBCI. It does not include a time.
type DateDataElement struct {
	*basicDataElement
}

// Val returns the value of d as time.Time. The time component of it will always be zero.
func (d *DateDataElement) Val() time.Time {
	return d.Value().(time.Time)
}

func (d *DateDataElement) String() string {
	return d.Val().Format("20060102")
}

// MarshalHBCI marshals d into HBCI wire format
func (d *DateDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(d.String()), nil
}

// UnmarshalHBCI unmarshals value into d
func (d *DateDataElement) UnmarshalHBCI(value []byte) error {
	t, err := time.Parse("20060102", charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*d = DateDataElement{&basicDataElement{t, dateDE, 8, false}}
	return nil
}

// IsValid returns true if the underlying date is not Zero.
func (d *DateDataElement) IsValid() bool {
	return !d.Val().IsZero()
}

// NewVirtualDate returns a new VirtualDateDataElement
func NewVirtualDate(date int) *VirtualDateDataElement {
	n := NewNumber(date, 8)
	n.typ = virtualDateDE
	return &VirtualDateDataElement{n}
}

// VirtualDateDataElement represents a virtual date
// TODO: modelling a virtual date?!
type VirtualDateDataElement struct {
	*NumberDataElement
}

// Valid returns true if the length of v is 8, false otherwise
func (v *VirtualDateDataElement) Valid() bool {
	return v.Length() == 8
}

// NewTime returns a new TimeDataElement
func NewTime(date time.Time) *TimeDataElement {
	return &TimeDataElement{&basicDataElement{date, dateDE, 6, false}}
}

// A TimeDataElement represents a time of a date component. It always contains
// a time component, but the date components will always be 0, i.e. 0000-00-00.
type TimeDataElement struct {
	*basicDataElement
}

// Val returns the value of t as time.Time.
func (t *TimeDataElement) Val() time.Time {
	return t.Value().(time.Time)
}

func (t *TimeDataElement) String() string {
	return t.Val().Format("150405")
}

// MarshalHBCI marshals t into HBCI wire format
func (t *TimeDataElement) MarshalHBCI() ([]byte, error) {
	return charset.ToISO8859_1(t.String()), nil
}

// UnmarshalHBCI unmarshals value into t
func (t *TimeDataElement) UnmarshalHBCI(value []byte) error {
	parsedTime, err := time.Parse("150405", charset.ToUTF8(value))
	if err != nil {
		return err
	}
	*t = TimeDataElement{&basicDataElement{parsedTime, timeDE, 6, false}}
	return nil
}

// IsValid returns true if the underlying date is not Zero.
func (t *TimeDataElement) IsValid() bool {
	return !t.Val().IsZero()
}

// NewIdentification returns a new IdentificationDataElement
func NewIdentification(id string) *IdentificationDataElement {
	an := NewAlphaNumeric(id, 30)
	an.typ = identificationDE
	return &IdentificationDataElement{an}
}

// An IdentificationDataElement represents unique identifier. It is
// only valid if the underlying id is 30 characters or less.
type IdentificationDataElement struct {
	*AlphaNumericDataElement
}

// Type returns the DataElementType of i
func (i *IdentificationDataElement) Type() DataElementType {
	return identificationDE
}

// UnmarshalHBCI unmarshals value into i
func (i *IdentificationDataElement) UnmarshalHBCI(value []byte) error {
	var alpha AlphaNumericDataElement
	err := alpha.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*i = IdentificationDataElement{&alpha}
	return nil
}

// NewCountryCode returns a new CountryCodeDataElement
func NewCountryCode(code int) *CountryCodeDataElement {
	d := NewDigit(code, 3)
	d.typ = countryCodeDE
	return &CountryCodeDataElement{d}
}

// A CountryCodeDataElement represents a Country code as defined by ISO 3166-1
// (numeric code)
type CountryCodeDataElement struct {
	*DigitDataElement
}

// Type represents the DataElement Type of c
func (c *CountryCodeDataElement) Type() DataElementType {
	return countryCodeDE
}

// UnmarshalHBCI unmarshals value into c
func (c *CountryCodeDataElement) UnmarshalHBCI(value []byte) error {
	var d DigitDataElement
	err := d.UnmarshalHBCI(value)
	if err != nil {
		return nil
	}
	*c = CountryCodeDataElement{&d}
	return nil
}

// NewCurrency returns a new CurrencyDataElement
func NewCurrency(cur string) *CurrencyDataElement {
	an := NewAlphaNumeric(cur, 3)
	an.typ = currencyDE
	return &CurrencyDataElement{an}
}

// A CurrencyDataElement represents a currency code as defined by ISO 4217
// alpha-3 in upcase letters
type CurrencyDataElement struct {
	*AlphaNumericDataElement
}

// IsValid returns true if the currency format matches ISO 4217 alpha-3
func (c *CurrencyDataElement) IsValid() bool {
	return c.Length() == 3
}

// Type returns the DataElementType of c
func (c *CurrencyDataElement) Type() DataElementType {
	return currencyDE
}

// UnmarshalHBCI unmarshals value into c
func (c *CurrencyDataElement) UnmarshalHBCI(value []byte) error {
	var a AlphaNumericDataElement
	err := a.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*c = CurrencyDataElement{&a}
	return nil
}

// NewValue returns a new ValueDataElement
func NewValue(val float64) *ValueDataElement {
	f := NewFloat(val, 15)
	f.typ = valueDE
	return &ValueDataElement{f}
}

// A ValueDataElement represents a float which can have upto 15 characters
type ValueDataElement struct {
	*FloatDataElement
}

// Type returns the DataElementType of v
func (v *ValueDataElement) Type() DataElementType {
	return valueDE
}

// UnmarshalHBCI unmarshals value into v
func (v *ValueDataElement) UnmarshalHBCI(value []byte) error {
	var f FloatDataElement
	err := f.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	*v = ValueDataElement{&f}
	return nil
}
