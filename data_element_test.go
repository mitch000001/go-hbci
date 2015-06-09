package hbci

import (
	"reflect"
	"testing"
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
		{1, NumberDE, 3, true, 1, NumberDE, 1, "1"},
		{1234, NumberDE, 3, false, 1234, NumberDE, 4, "1234"},
	}
	for _, test := range tests {
		d := NewDataElement(test.inType, test.inValue, test.inMaxLength)

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

		valid := d.Valid()

		if valid != test.valid {
			t.Logf("Input: %+#v\n", test)
			if test.valid {
				t.Logf("Expected DataElement to be valid, was not\n")
			} else {
				t.Logf("Expected DataElement to be invalid, was valid\n")
			}
			t.Logf("Expected DataElement to be valid, was not\n", expectedString, actualString)
			t.Fail()
		}
	}
}

func TestNewAlphaNumericDataElement(t *testing.T) {
	dataElement := NewAlphaNumericDataElement("abc", 5)

	expectedType := AlphaNumericDE

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
		d := NewDigitDataElement(test.in, test.inMaxLength)
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

		valid := d.Valid()

		if valid != test.valid {
			t.Logf("Input: %+#v\n", test)
			if test.valid {
				t.Logf("Expected DataElement to be valid, was not\n")
			} else {
				t.Logf("Expected DataElement to be invalid, was valid\n")
			}
			t.Logf("Expected DataElement to be valid, was not\n", expectedString, actualString)
			t.Fail()
		}
	}
}

func TestDigitDataElementValue(t *testing.T) {
	d := NewDigitDataElement(1, 2)

	var expected interface{} = 1

	actual := d.Value()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}

func TestDigitDataElementType(t *testing.T) {
	d := NewDigitDataElement(1, 2)

	expected := DigitDE

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
		d := NewNumberDataElement(test.in, test.inMaxLength)
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

		valid := d.Valid()

		if valid != test.valid {
			t.Logf("Input: %+#v\n", test)
			if test.valid {
				t.Logf("Expected DataElement to be valid, was not\n")
			} else {
				t.Logf("Expected DataElement to be invalid, was valid\n")
			}
			t.Logf("Expected DataElement to be valid, was not\n", expectedString, actualString)
			t.Fail()
		}
	}
}

func TestNumberDataElementValue(t *testing.T) {
	d := NewNumberDataElement(1, 2)

	var expected interface{} = 1

	actual := d.Value()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}

func TestNumberDataElementType(t *testing.T) {
	d := NewNumberDataElement(1, 2)

	expected := NumberDE

	actual := d.Type()

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Value() to return %v, got %v\n", expected, actual)
		t.Fail()
	}
}
