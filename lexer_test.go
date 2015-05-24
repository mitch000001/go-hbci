package hbci

import "testing"

func TestLexBinaryData(t *testing.T) {
	type testData struct {
		text  string
		typ   itemType
		value string
	}
	tests := []testData{
		{"@2@ab'", itemBinaryData, "@2@ab"},
		{"@@ab'", itemError, "Binary length can't be empty"},
		{"@2@a'", itemError, "Expected syntax symbol after binary data"},
		{"@2@abc'", itemError, "Expected syntax symbol after binary data"},
	}
	for _, test := range tests {
		l := lex("", test.text)
		item := l.nextItem()
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
