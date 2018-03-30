package token

import (
	"reflect"
	"testing"
)

type testToken int

func (t testToken) Type() Type {
	return Type(t)
}

func (t testToken) Value() string {
	return ""
}

func (t testToken) Pos() int {
	return 0
}

func (t testToken) IsSyntaxSymbol() bool {
	return Type(t) == GROUP_DATA_ELEMENT_SEPARATOR || Type(t) == DATA_ELEMENT_SEPARATOR || Type(t) == SEGMENT_END_MARKER
}

func (t testToken) String() string {
	return Type(t).String()
}

func TestTokenValue(t *testing.T) {
	token := New(ALPHA_NUMERIC, "abc", 0)

	if !reflect.DeepEqual(token.Value(), "abc") {
		t.Logf("Expected Value to return %q, got %q\n", "abc", token.Value())
		t.Fail()
	}
}

func TestTokenType(t *testing.T) {
	token := New(ALPHA_NUMERIC, "abc", 0)

	if !reflect.DeepEqual(token.Type(), ALPHA_NUMERIC) {
		t.Logf("Expected Type to return %q, got %q\n", ALPHA_NUMERIC, token.Type())
		t.Fail()
	}
}

func TestTokenPos(t *testing.T) {
	token := New(ALPHA_NUMERIC, "abc", 0)

	if !reflect.DeepEqual(token.Pos(), 0) {
		t.Logf("Expected Type to return %d, got %d\n", 0, token.Pos())
		t.Fail()
	}
}
