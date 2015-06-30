package hbci

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/token"
)

type testData struct {
	text  string
	typ   token.TokenType
	value string
}

func TestStringLexer(t *testing.T) {
	testInput := "ab??cd\ref+12345+@2@ab'"
	l := NewStringLexer("", testInput)
	var items []token.Token
	for l.HasNext() {
		item := l.Next()
		items = append(items, item)
	}
	var itemTypes []token.TokenType
	for _, item := range items {
		itemTypes = append(itemTypes, item.Type())
	}
	expectedItemTypes := []token.TokenType{
		token.TEXT,
		token.DATA_ELEMENT_SEPARATOR,
		token.NUMERIC,
		token.DATA_ELEMENT_SEPARATOR,
		token.BINARY_DATA,
		token.SEGMENT_END_MARKER,
		token.EOF,
	}
	if !reflect.DeepEqual(expectedItemTypes, itemTypes) {
		t.Logf("Expected types to equal \n\t'%s' \ngot: \n\t'%s'\n", expectedItemTypes, itemTypes)
		t.Fail()
	}
}

func TestLexText(t *testing.T) {
	tests := []testData{
		{"ab\rcd'", token.TEXT, "ab\rcd"},
		{"ab\ncd'", token.TEXT, "ab\ncd"},
		{"ab\r\ncd'", token.TEXT, "ab\r\ncd"},
		{"ab\n\rcd'", token.TEXT, "ab\n\rcd"},
	}
	for _, test := range tests {
		l := NewStringLexer("", test.text)
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if item.Value() != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexAlphaNumeric(t *testing.T) {
	tests := []testData{
		{"ab'", token.ALPHA_NUMERIC, "ab"},
		{"ab123'", token.ALPHA_NUMERIC, "ab123"},
		{"ab!)'", token.ALPHA_NUMERIC, "ab!)"},
		{"ab!)'", token.ALPHA_NUMERIC, "ab!)"},
		{"ab!):", token.ALPHA_NUMERIC, "ab!)"},
		{"ab!)+", token.ALPHA_NUMERIC, "ab!)"},
		{"ab?''", token.ALPHA_NUMERIC, "ab?'"},
		{"ab?e", token.ERROR, "Unexpected escape character"},
		{"ab", token.ERROR, "Unexpected end of input"},
	}
	for _, test := range tests {
		l := NewStringLexer("", test.text)
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if item.Value() != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexSyntaxSymbol(t *testing.T) {
	tests := []testData{
		{"'", token.SEGMENT_END_MARKER, "'"},
		{"+", token.DATA_ELEMENT_SEPARATOR, "+"},
		{":", token.GROUP_DATA_ELEMENT_SEPARATOR, ":"},
	}
	for _, test := range tests {
		l := NewStringLexer("", test.text)
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if item.Value() != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexBinaryData(t *testing.T) {
	tests := []testData{
		{"@2@ab'", token.BINARY_DATA, "@2@ab"},
		{"@@ab'", token.ERROR, "Binary length can't be empty"},
		{"@2@a'", token.ERROR, "Expected syntax symbol after binary data"},
		{"@2@abc'", token.ERROR, "Expected syntax symbol after binary data"},
		{"@2x@ab'", token.ERROR, "Binary length must contain of digits only"},
	}
	for _, test := range tests {
		l := NewStringLexer("", test.text)
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if item.Value() != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}

func TestLexDigit(t *testing.T) {
	tests := []testData{
		{"123'", token.NUMERIC, "123"},
		{"0123'", token.DIGIT, "0123"},
		{"0,123'", token.FLOAT, "0,123"},
		{"1,23'", token.FLOAT, "1,23"},
		{"1,''", token.FLOAT, "1,"},
		{"0'", token.NUMERIC, "0"},
		{"0,'", token.FLOAT, "0,"},
		{"01,23'", token.ERROR, "Malformed float"},
		{"0,12a'", token.ERROR, "Malformed float"},
		{"1,23a'", token.ERROR, "Malformed float"},
		{"012a'", token.ERROR, "Malformed digit"},
		{"12a'", token.ERROR, "Malformed numeric"},
	}
	for _, test := range tests {
		l := NewStringLexer("", test.text)
		item := l.Next()
		if item.Type() != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %s, got %s\n", test.typ, item.Type())
			t.Fail()
		}
		if item.Value() != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.Value())
			t.Fail()
		}
	}
}
