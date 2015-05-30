package hbci

import (
	"reflect"
	"testing"
)

type tokenLexerTestData struct {
	input  Tokens
	output Tokens
}

func tokens(types ...TokenType) []Token {
	var tokens []Token
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

func (t testToken) String() string {
	return TokenType(t).String()
}

func TestTokenLexer(t *testing.T) {
	tests := []tokenLexerTestData{
		{
			tokens(ALPHA_NUMERIC, ESCAPE_SEQUENCE, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(ALPHA_NUMERIC_WITH_ESCAPE_SEQUENCE, SEGMENT_END_MARKER, EOF),
		},
		{
			tokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			tokens(ALPHA_NUMERIC, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			tokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, DIGIT, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			tokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT_SEPARATOR, DIGIT, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			tokens(DATA_ELEMENT_GROUP, EOF),
		},
	}
	for idx, test := range tests {
		var emittedTokens Tokens
		lexer := NewTokenLexer("Test TokenLexer", test.input)
		for lexer.HasNext() {
			token := lexer.Next()
			emittedTokens = append(emittedTokens, token)
		}
		if !reflect.DeepEqual(emittedTokens.Types(), test.output.Types()) {
			t.Logf("Input (%d): %s\n", idx, test.input)
			t.Logf("Expected output to equal \n%s\n\tgot:\n%s\n", test.output.Types(), emittedTokens.Types())
			t.Fail()
		}
	}
}
