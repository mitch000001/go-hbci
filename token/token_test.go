package token

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTokenValue(t *testing.T) {
	token := New(ALPHA_NUMERIC, []byte("abc"), 0)

	if !bytes.Equal(token.Value(), []byte("abc")) {
		t.Logf("Expected Value to return %q, got %q\n", "abc", token.Value())
		t.Fail()
	}
}

func TestTokenType(t *testing.T) {
	token := New(ALPHA_NUMERIC, []byte("abc"), 0)

	if !reflect.DeepEqual(token.Type(), ALPHA_NUMERIC) {
		t.Logf("Expected Type to return %q, got %q\n", ALPHA_NUMERIC, token.Type())
		t.Fail()
	}
}

func TestTokenPos(t *testing.T) {
	token := New(ALPHA_NUMERIC, []byte("abc"), 0)

	if !reflect.DeepEqual(token.Pos(), 0) {
		t.Logf("Expected Type to return %d, got %d\n", 0, token.Pos())
		t.Fail()
	}
}
