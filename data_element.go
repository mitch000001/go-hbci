package hbci

import "fmt"

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
	NoType DataElementType = iota << 1
	AlphaNumeric
	Text
	Number
	Digit
	Float
	// Advanced types: DataElementGroups
	SegmentHeaderType
)

type dataElement struct {
	val       interface{}
	typ       DataElementType
	maxLength int
}

func NewDataElement(typ DataElementType, value interface{}, maxLength int) DataElement {
	return &dataElement{value, typ, maxLength}
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
	return fmt.Sprintf(fmtString, d.Value().(int))
}

func NewNumberDataElement(val, maxLength int) *NumberDataElement {
	return &NumberDataElement{&dataElement{val, Number, maxLength}}
}

type NumberDataElement struct {
	*dataElement
}
