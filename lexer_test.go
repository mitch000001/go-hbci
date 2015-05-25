package hbci

import "testing"

type testData struct {
	text  string
	typ   itemType
	value string
}

func TestLexAlphaNumeric(t *testing.T) {
	tests := []testData{
		{"ab'", itemAlphaNumeric, "ab"},
		{"ab123'", itemAlphaNumeric, "ab123"},
		{"ab!)'", itemAlphaNumeric, "ab!)"},
		{"ab??'", itemAlphaNumeric, "ab??"},
		{"ab?+'", itemAlphaNumeric, "ab?+"},
		{"ab?:'", itemAlphaNumeric, "ab?:"},
		{"ab?''", itemAlphaNumeric, "ab?'"},
		{"ab?'", itemError, "Unexpected end of input"},
		{"ab?a", itemError, "Unexpected '?' at pos 3"},
	}
	for _, test := range tests {
		l := NewLexer("", test.text)
		item := l.NextItem()
		if item.typ != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %d, got %d\n", test.typ, item.typ)
			t.Fail()
		}
		if item.val != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.val)
			t.Fail()
		}
	}
}

func TestLexSyntaxSymbol(t *testing.T) {
	tests := []testData{
		{"'", itemSegmentEnd, "'"},
		{"+", itemDataElementSeparator, "+"},
		{":", itemGroupDataElementSeparator, ":"},
	}
	for _, test := range tests {
		l := NewLexer("", test.text)
		item := l.NextItem()
		if item.typ != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %d, got %d\n", test.typ, item.typ)
			t.Fail()
		}
		if item.val != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.val)
			t.Fail()
		}
	}
}

func TestLexBinaryData(t *testing.T) {
	tests := []testData{
		{"@2@ab'", itemBinaryData, "@2@ab"},
		{"@@ab'", itemError, "Binary length can't be empty"},
		{"@2@a'", itemError, "Expected syntax symbol after binary data"},
		{"@2@abc'", itemError, "Expected syntax symbol after binary data"},
		{"@2x@ab'", itemError, "Binary length must contain of digits only"},
	}
	for _, test := range tests {
		l := NewLexer("", test.text)
		item := l.NextItem()
		if item.typ != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %d, got %d\n", test.typ, item.typ)
			t.Fail()
		}
		if item.val != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.val)
			t.Fail()
		}
	}
}

func TestLexDigit(t *testing.T) {
	tests := []testData{
		{"123'", itemNumeric, "123"},
		{"0123'", itemDigit, "0123"},
		{"0,123'", itemFloat, "0,123"},
		{"1,23'", itemFloat, "1,23"},
		{"1,''", itemFloat, "1,"},
		{"0'", itemNumeric, "0"},
		{"0,'", itemFloat, "0,"},
		{"01,23'", itemError, "Malformed float"},
		{"0,12a'", itemError, "Malformed float"},
		{"1,23a'", itemError, "Malformed float"},
		{"012a'", itemError, "Malformed digit"},
		{"12a'", itemError, "Malformed numeric"},
	}
	for _, test := range tests {
		l := NewLexer("", test.text)
		item := l.NextItem()
		if item.typ != test.typ {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected type to equal %d, got %d\n", test.typ, item.typ)
			t.Fail()
		}
		if item.val != test.value {
			t.Logf("Input: %q\n", test.text)
			t.Logf("Expected val to equal %q, got %q\n", test.value, item.val)
			t.Fail()
		}
	}
}
