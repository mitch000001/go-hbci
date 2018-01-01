package element

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/charset"
)

type testDataElementData struct {
	inValue     interface{}
	inType      DataElementType
	inMaxLength int
	valid       bool
	outValue    interface{}
	outType     DataElementType
	outLength   int
	outString   string
}

func TestNewDataElement(t *testing.T) {
	tests := []testDataElementData{
		{1, numberDE, 3, true, 1, numberDE, 1, "1"},
		{1234, numberDE, 3, false, 1234, numberDE, 4, "1234"},
	}
	for _, test := range tests {
		d := New(test.inType, test.inValue, test.inMaxLength)

		expectedOut := test.outValue

		actualOut := d.Value()

		if !reflect.DeepEqual(expectedOut, actualOut) {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected Value() to return %v, got %v\n", expectedOut, actualOut)
			t.Fail()
		}

		expectedLength := test.outLength

		actualLength := d.Length()
		if actualLength != expectedLength {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected Length() to return %d, got %d\n", expectedLength, actualLength)
			t.Fail()
		}

		expectedString := test.outString

		actualString := d.String()

		if actualString != expectedString {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected String() to return %q, got %q\n", expectedString, actualString)
			t.Fail()
		}

		valid := d.IsValid()

		if valid != test.valid {
			t.Logf("Input: %+#v\n", test)
			if test.valid {
				t.Logf("Expected DataElement to be valid, was not\n")
			} else {
				t.Logf("Expected DataElement to be invalid, was valid\n")
			}
			t.Logf("Expected DataElement to be valid, was not\n")
			t.Fail()
		}
	}
}

func TestNewAlphaNumericDataElement(t *testing.T) {
	dataElement := NewAlphaNumeric("abc", 5)

	expectedType := alphaNumericDE

	actualType := dataElement.Type()

	if expectedType != actualType {
		t.Logf("Expected Type to equal %v, got %v\n", expectedType, actualType)
		t.Fail()
	}

	expectedLength := len("abc")

	actualLength := dataElement.Length()

	if expectedLength != actualLength {
		t.Logf("Expected Length() to return %d, got %d\n", expectedLength, actualLength)
		t.Fail()
	}

	expectedString := "abc"

	actualString := dataElement.String()

	if actualString != expectedString {
		t.Logf("Expected String() to return %q, got %q\n", expectedString, actualString)
		t.Fail()
	}
}

func TestAlphaNumericDataElementMarshalHBCI(t *testing.T) {
	type testData struct {
		unmarshaled string
		marshaled   []byte
		err         error
	}
	tests := []testData{
		{
			"abc",
			charset.ToISO8859_1("abc"),
			nil,
		},
		{
			"ab:c",
			charset.ToISO8859_1("ab?:c"),
			nil,
		},
		{
			"abüc",
			charset.ToISO8859_1("abüc"),
			nil,
		},
	}
	for _, test := range tests {
		element := NewAlphaNumeric(test.unmarshaled, len(test.unmarshaled))

		marshaled, err := element.MarshalHBCI()

		if !reflect.DeepEqual(test.err, err) {
			t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", test.err, err)
			t.Fail()
		}

		if !bytes.Equal(test.marshaled, marshaled) {
			t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", test.marshaled, marshaled)
			t.Fail()
		}
	}
}

func TestAlphaNumericDataElementUnmarshalHBCI(t *testing.T) {
	type testData struct {
		marshaled   []byte
		unmarshaled string
		err         error
	}
	tests := []testData{
		{
			charset.ToISO8859_1("abc"),
			NewAlphaNumeric("abc", 3).Val(),
			nil,
		},
		{
			charset.ToISO8859_1("ab?:c"),
			NewAlphaNumeric("ab:c", 3).Val(),
			nil,
		},
		{
			charset.ToISO8859_1("abüc"),
			NewAlphaNumeric("abüc", 3).Val(),
			nil,
		},
	}
	for _, test := range tests {
		element := &AlphaNumericDataElement{}

		err := element.UnmarshalHBCI(test.marshaled)

		if !reflect.DeepEqual(test.err, err) {
			t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", test.err, err)
			t.Fail()
		}

		actual := element.Val()

		if test.unmarshaled != actual {
			t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", test.unmarshaled, actual)
			t.Fail()
		}
	}
}

type testDigitDataElementData struct {
	in          int
	inMaxLength int
	valid       bool
	outLength   int
	outString   string
}

func TestNewDigitDataElement(t *testing.T) {
	tests := []testDigitDataElementData{
		{1, 4, true, 1, "0001"},
		{10, 4, true, 2, "0010"},
		{1000, 4, true, 4, "1000"},
		{10000, 4, false, 5, "10000"},
	}

	for _, test := range tests {
		d := NewDigit(test.in, test.inMaxLength)
		expectedLength := test.outLength

		actualLength := d.Length()

		if actualLength != expectedLength {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected Length() to return %d, got %d\n", expectedLength, actualLength)
			t.Fail()
		}

		expectedString := test.outString

		actualString := d.String()

		if actualString != expectedString {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected String() to return %q, got %q\n", expectedString, actualString)
			t.Fail()
		}

		valid := d.IsValid()

		if valid != test.valid {
			t.Logf("Input: %+#v\n", test)
			if test.valid {
				t.Logf("Expected DataElement to be valid, was not\n")
			} else {
				t.Logf("Expected DataElement to be invalid, was valid\n")
			}
			t.Logf("Expected DataElement to be valid, was not\n")
			t.Fail()
		}
	}
}

func TestDigitDataElementValue(t *testing.T) {
	d := NewDigit(1, 2)

	var expected interface{} = 1

	actual := d.Value()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}

func TestDigitDataElementType(t *testing.T) {
	d := NewDigit(1, 2)

	expected := digitDE

	actual := d.Type()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}

func TestNewNumberDataElement(t *testing.T) {
	tests := []testDigitDataElementData{
		{1, 4, true, 1, "1"},
		{10, 4, true, 2, "10"},
		{1000, 4, true, 4, "1000"},
		{10000, 4, false, 5, "10000"},
	}

	for _, test := range tests {
		d := NewNumber(test.in, test.inMaxLength)
		expectedLength := test.outLength

		actualLength := d.Length()

		if actualLength != expectedLength {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected Length() to return %d, got %d\n", expectedLength, actualLength)
			t.Fail()
		}

		expectedString := test.outString

		actualString := d.String()

		if actualString != expectedString {
			t.Logf("Input: %+#v\n", test)
			t.Logf("Expected String() to return %q, got %q\n", expectedString, actualString)
			t.Fail()
		}

		valid := d.IsValid()

		if valid != test.valid {
			t.Logf("Input: %+#v\n", test)
			if test.valid {
				t.Logf("Expected DataElement to be valid, was not\n")
			} else {
				t.Logf("Expected DataElement to be invalid, was valid\n")
			}
			t.Logf("Expected DataElement to be valid, was not\n")
			t.Fail()
		}
	}
}

func TestNumberDataElementValue(t *testing.T) {
	d := NewNumber(1, 2)

	var expected interface{} = 1

	actual := d.Value()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}

func TestNumberDataElementType(t *testing.T) {
	d := NewNumber(1, 2)

	expected := numberDE

	actual := d.Type()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}

func TestBinaryDataElementString(t *testing.T) {
	b := NewBinary([]byte("test123"), 7)

	expected := "@7@test123"

	actual := b.String()

	if expected != actual {
		t.Logf("Expected BinaryDataElement to serialize to %q, got %q\n", expected, actual)
		t.Fail()
	}
}

func TestBinaryDataElementUnmarshalHBCI(t *testing.T) {
	var b BinaryDataElement

	err := b.UnmarshalHBCI([]byte("@7@test123"))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	val := b.Val()
	expectedVal := []byte("test123")

	if !bytes.Equal(val, expectedVal) {
		t.Logf("Expected Val() to return %q, got %q\n", expectedVal, val)
		t.Fail()
	}
}

type testDataElement struct {
	alpha *AlphaNumericDataElement
	num   *NumberDataElement
}

func (t *testDataElement) groupDataElements() []DataElement {
	return []DataElement{t.alpha, t.num}
}

func (t *testDataElement) Elements() []DataElement {
	return []DataElement{t.alpha, t.num}
}

type testDataElementGroupData struct {
	alphaIn *AlphaNumericDataElement
	numIn   *NumberDataElement
	out     string
}

func TestGroupDataElementGroupString(t *testing.T) {
	tests := []testDataElementGroupData{
		{
			NewAlphaNumeric("abc", 3),
			NewNumber(123, 3),
			"abc:123",
		},
		{
			NewAlphaNumeric("abc", 3),
			nil,
			"abc:",
		},
		{
			nil,
			NewNumber(123, 3),
			":123",
		},
		{
			nil,
			nil,
			":",
		},
	}

	for _, test := range tests {
		testData := &testDataElement{
			alpha: test.alphaIn,
			num:   test.numIn,
		}

		group := NewGroupDataElementGroup(0, 2, testData)

		actualString := group.String()

		if test.out != actualString {
			t.Logf("Input: %#v\n", testData)
			t.Logf("Expected String() to return %q, got %q\n", test.out, actualString)
			t.Fail()
		}
	}
}

func TestGroupDataElementGroupUnmarshalHBCI(t *testing.T) {
	t.Skip("This test is broken due to necessary implementation changes.")
	type testDataElementGroupUnmarshalData struct {
		in       string
		alphaOut *AlphaNumericDataElement
		numOut   *NumberDataElement
	}

	tests := []testDataElementGroupUnmarshalData{
		{
			"abc:123",
			NewAlphaNumeric("abc", 3),
			NewNumber(123, 3),
		},
		{
			"abc:",
			NewAlphaNumeric("abc", 3),
			nil,
		},
		{
			":123",
			nil,
			NewNumber(123, 3),
		},
		{
			":",
			nil,
			nil,
		},
	}

	for _, test := range tests {
		tde := testDataElement{}
		group := new(elementGroup)
		group.elements = tde.Elements

		err := group.UnmarshalHBCI([]byte(test.in))

		if err != nil {
			t.Logf("Input: %q\n", test.in)
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		expectedArray := []DataElement{test.alphaOut, test.numOut}
		actualArray := group.elements()

		if !reflect.DeepEqual(expectedArray, actualArray) {
			t.Logf("Input: %q\n", test.in)
			t.Logf("Expected UnmarshalHBCI() to return \n%+#v\n\tgot \n%+#v\n", expectedArray, actualArray)
			t.Fail()
		}
	}
}
