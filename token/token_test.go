package token

import (
	"reflect"
	"sort"
	"testing"
)

func testTokens(types ...TokenType) Tokens {
	var tokens Tokens
	for _, typ := range types {
		tokens = append(tokens, testToken(typ))
	}
	return tokens
}

type testToken int

func (t testToken) Type() TokenType {
	return TokenType(t)
}

func (t testToken) Value() string {
	return ""
}

func (t testToken) Pos() int {
	return 0
}

func (t testToken) IsSyntaxSymbol() bool {
	return TokenType(t) == GROUP_DATA_ELEMENT_SEPARATOR || TokenType(t) == DATA_ELEMENT_SEPARATOR || TokenType(t) == SEGMENT_END_MARKER
}

func (t testToken) Children() Tokens {
	return Tokens{}
}

func (t testToken) RawTokens() Tokens {
	return Tokens{t}
}

func (t testToken) String() string {
	return TokenType(t).String()
}

func TestTokenValue(t *testing.T) {
	token := NewToken(ALPHA_NUMERIC, "abc", 0)

	if !reflect.DeepEqual(token.Value(), "abc") {
		t.Logf("Expected Value to return %q, got %q\n", "abc", token.Value())
		t.Fail()
	}
}

func TestTokenType(t *testing.T) {
	token := NewToken(ALPHA_NUMERIC, "abc", 0)

	if !reflect.DeepEqual(token.Type(), ALPHA_NUMERIC) {
		t.Logf("Expected Type to return %q, got %q\n", ALPHA_NUMERIC, token.Type())
		t.Fail()
	}
}

func TestTokenPos(t *testing.T) {
	token := NewToken(ALPHA_NUMERIC, "abc", 0)

	if !reflect.DeepEqual(token.Pos(), 0) {
		t.Logf("Expected Type to return %d, got %d\n", 0, token.Pos())
		t.Fail()
	}
}

func TestGroupTokensChildren(t *testing.T) {
	type testData struct {
		children Tokens
		types    []TokenType
	}

	tests := []testData{
		{
			testTokens(ALPHA_NUMERIC),
			[]TokenType{ALPHA_NUMERIC},
		},
		{
			testTokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC),
			[]TokenType{ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC},
		},
	}

	for _, test := range tests {
		gt := NewGroupToken(GROUP_DATA_ELEMENT, test.children...)

		children := gt.Children()

		expectedChildrenTypes := test.types
		actualChildrenTypes := children.Types()

		if !reflect.DeepEqual(expectedChildrenTypes, actualChildrenTypes) {
			t.Logf("Expected Children to equal\n%s\n\tgot:\n%s\n", expectedChildrenTypes, actualChildrenTypes)
			t.Fail()
		}
	}
}

func TestGroupTokenRawTokens(t *testing.T) {
	type testData struct {
		childTokens Tokens
		output      Tokens
	}
	tests := []testData{
		{
			// First level
			testTokens(ALPHA_NUMERIC),
			testTokens(ALPHA_NUMERIC),
		},
		{
			// Second level
			Tokens{NewGroupToken(GROUP_DATA_ELEMENT, testTokens(ALPHA_NUMERIC)...)},
			testTokens(ALPHA_NUMERIC),
		},
		{
			// Third level
			Tokens{NewGroupToken(GROUP_DATA_ELEMENT, NewGroupToken(GROUP_DATA_ELEMENT, testTokens(ALPHA_NUMERIC)...))},
			testTokens(ALPHA_NUMERIC),
		},
		{
			// Mixed levels
			Tokens{NewToken(NUMERIC, "2", 0), NewGroupToken(GROUP_DATA_ELEMENT, testTokens(DIGIT)...), NewGroupToken(GROUP_DATA_ELEMENT, NewGroupToken(GROUP_DATA_ELEMENT, testTokens(ALPHA_NUMERIC)...))},
			testTokens(NUMERIC, DIGIT, ALPHA_NUMERIC),
		},
	}

	for idx, test := range tests {
		gt := NewGroupToken(GROUP_DATA_ELEMENT, test.childTokens...)

		expectedRawTokenTypes := test.output.Types()
		sort.Sort(TokenTypes(expectedRawTokenTypes))

		rawTokens := gt.RawTokens()
		actualRawTypes := rawTokens.Types()
		sort.Sort(TokenTypes(actualRawTypes))

		if !reflect.DeepEqual(expectedRawTokenTypes, actualRawTypes) {
			t.Logf("Children (%d):\n%s", idx, test.childTokens.Types())
			t.Logf("Expected RawTokens to equal\n%s\n\tgot:\n%s\n", expectedRawTokenTypes, actualRawTypes)
			t.Fail()
		}
	}
}
