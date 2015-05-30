package hbci

import (
	"reflect"
	"testing"
)

type tokenLexerTestData struct {
	example string
	input   Tokens
	output  Tokens
}

func tokens(types ...TokenType) Tokens {
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

func TestTokenLexer(t *testing.T) {
	tests := []tokenLexerTestData{
		{
			"ab??2'",
			tokens(ALPHA_NUMERIC, ESCAPE_SEQUENCE, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(ALPHA_NUMERIC_WITH_ESCAPE_SEQUENCE, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:23'",
			tokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab+23'",
			tokens(ALPHA_NUMERIC, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:23:05+45'",
			tokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, DIGIT, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab::01+5'",
			tokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT_SEPARATOR, DIGIT, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:cd+ef'",
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			tokens(DATA_ELEMENT_GROUP, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:cd+ef+gh'",
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			tokens(DATA_ELEMENT_GROUP, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:cd:ef'",
			tokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			tokens(DATA_ELEMENT_GROUP, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:1:2:3'",
			Tokens{
				NewGroupToken(DATA_ELEMENT_GROUP,
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(ALPHA_NUMERIC)),
					testToken(GROUP_DATA_ELEMENT_SEPARATOR),
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(NUMERIC)),
					testToken(GROUP_DATA_ELEMENT_SEPARATOR),
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(NUMERIC)),
					testToken(GROUP_DATA_ELEMENT_SEPARATOR),
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(NUMERIC)),
				),
				testToken(SEGMENT_END_MARKER),
				testToken(EOF),
			},
			tokens(SEGMENT_HEADER, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:1:2:3+4'",
			Tokens{
				NewGroupToken(DATA_ELEMENT_GROUP,
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(ALPHA_NUMERIC)),
					testToken(GROUP_DATA_ELEMENT_SEPARATOR),
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(NUMERIC)),
					testToken(GROUP_DATA_ELEMENT_SEPARATOR),
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(NUMERIC)),
					testToken(GROUP_DATA_ELEMENT_SEPARATOR),
					NewGroupToken(GROUP_DATA_ELEMENT, testToken(NUMERIC)),
				),
				testToken(DATA_ELEMENT_SEPARATOR),
				testToken(DATA_ELEMENT),
				testToken(SEGMENT_END_MARKER),
				testToken(EOF),
			},
			tokens(SEGMENT_HEADER, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
	}
	for idx, test := range tests {
		var emittedTokens Tokens
		lexer := NewTokenLexer("Test TokenLexer", test.input)
		for lexer.HasNext() {
			token := lexer.Next()
			emittedTokens = append(emittedTokens, token)
			if token.Type() == ERROR {
				t.Logf("Lexer returned error: %#s\n", token)
			}
		}
		emittedTokenTypes := emittedTokens.Types()
		if !reflect.DeepEqual(emittedTokenTypes, test.output.Types()) {
			t.Logf("Input (%d): %s\n", idx, test.input)
			t.Logf("Expected output to equal \n%s\n\tgot:\n%s\n", test.output.Types(), emittedTokenTypes)
			t.Fail()
		}
	}
}
