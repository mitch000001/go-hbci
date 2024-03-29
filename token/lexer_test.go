package token

import (
	"bytes"
	"reflect"
	"testing"
)

type testData struct {
	text  string
	typ   Type
	value string
}

func TestLexer(t *testing.T) {
	testInput := "ab??cd\ref+12345+@2@ab'"
	l := NewLexer("", []byte(testInput))
	var items []Token
	for l.HasNext() {
		item := l.Next()
		items = append(items, item)
	}
	var itemTypes []Type
	for _, item := range items {
		itemTypes = append(itemTypes, item.Type())
	}
	expectedItemTypes := []Type{
		TEXT,
		DATA_ELEMENT_SEPARATOR,
		NUMERIC,
		DATA_ELEMENT_SEPARATOR,
		BINARY_DATA,
		SEGMENT_END_MARKER,
		EOF,
	}
	if !reflect.DeepEqual(expectedItemTypes, itemTypes) {
		t.Logf("Expected types to equal \n\t'%s' \ngot: \n\t'%s'\n", expectedItemTypes, itemTypes)
		t.Fail()
	}
}

func TestLexText(t *testing.T) {
	tests := []testData{
		{"ab\rcd'", TEXT, "ab\rcd"},
		{"ab\ncd'", TEXT, "ab\ncd"},
		{"ab\r\ncd'", TEXT, "ab\r\ncd"},
		{"ab\n\rcd'", TEXT, "ab\n\rcd"},
		{"ab\n\rcd", ERROR, "Unexpected end of input"},
	}
	for _, test := range tests {
		l := NewLexer("", []byte(test.text))
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if !bytes.Equal(item.Value(), []byte(test.value)) {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexAlphaNumeric(t *testing.T) {
	tests := []testData{
		{"ab'", ALPHA_NUMERIC, "ab"},
		{"ab123'", ALPHA_NUMERIC, "ab123"},
		{"ab!)'", ALPHA_NUMERIC, "ab!)"},
		{"ab!)'", ALPHA_NUMERIC, "ab!)"},
		{"ab!):", ALPHA_NUMERIC, "ab!)"},
		{"ab!)+", ALPHA_NUMERIC, "ab!)"},
		{"ab?''", ALPHA_NUMERIC, "ab?'"},
		{"ab?@'", ALPHA_NUMERIC, "ab?@"},
		{"ab?e", ERROR, "Unexpected escape character"},
		{"ab", ERROR, "Unexpected end of input"},
		{"ab??", ERROR, "Unexpected end of input"},
	}
	for _, test := range tests {
		l := NewLexer("", []byte(test.text))
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if !bytes.Equal(item.Value(), []byte(test.value)) {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexSyntaxSymbol(t *testing.T) {
	tests := []testData{
		{"'", SEGMENT_END_MARKER, "'"},
		{"+", DATA_ELEMENT_SEPARATOR, "+"},
		{":", GROUP_DATA_ELEMENT_SEPARATOR, ":"},
	}
	for _, test := range tests {
		l := NewLexer("", []byte(test.text))
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if !bytes.Equal(item.Value(), []byte(test.value)) {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexBinaryData(t *testing.T) {
	tests := []testData{
		{"@2@ab'", BINARY_DATA, "@2@ab"},
		{"@@ab'", ERROR, "Binary length can't be empty"},
		{"@2@a1@", ERROR, "Expected syntax symbol after binary data"},
		{"@2@abc'", ERROR, "Expected syntax symbol after binary data"},
		{"@2x@ab'", ERROR, "Binary length must contain of digits only"},
		{"@2@ab", ERROR, "Unexpected end of input"},
	}
	for _, test := range tests {
		l := NewLexer("", []byte(test.text))
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if !bytes.Equal(item.Value(), []byte(test.value)) {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexDigit(t *testing.T) {
	tests := []testData{
		{"123'", NUMERIC, "123"},
		{"0123'", DIGIT, "0123"},
		{"0,123'", FLOAT, "0,123"},
		{"1,23'", FLOAT, "1,23"},
		{"1,''", FLOAT, "1,"},
		{"0'", NUMERIC, "0"},
		{"0,'", FLOAT, "0,"},
		{"01,23'", ERROR, "Malformed float"},
		{"0,12a'", ALPHA_NUMERIC, "0,12a"},
		{"1,23a'", ALPHA_NUMERIC, "1,23a"},
		{"012a'", ALPHA_NUMERIC, "012a"},
		{"12a'", ALPHA_NUMERIC, "12a"},
		{"12?+'", ALPHA_NUMERIC, "12?+"},
		{"12", ERROR, "Unexpected end of input"},
	}
	for _, test := range tests {
		l := NewLexer("", []byte(test.text))
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if !bytes.Equal(item.Value(), []byte(test.value)) {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}
