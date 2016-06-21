package token

import (
	"reflect"
	"testing"
)

type tokenLexerTestData struct {
	example string
	input   Tokens
	output  Tokens
}

func TestTokenLexer(t *testing.T) {
	tests := []tokenLexerTestData{
		{
			"ab:23'",
			testTokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			testTokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab+23'",
			testTokens(ALPHA_NUMERIC, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			testTokens(DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:23:05+45'",
			testTokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, DIGIT, DATA_ELEMENT_SEPARATOR, SEGMENT_END_MARKER, EOF),
			testTokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab::01+5'",
			testTokens(ALPHA_NUMERIC, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT_SEPARATOR, DIGIT, DATA_ELEMENT_SEPARATOR, NUMERIC, SEGMENT_END_MARKER, EOF),
			testTokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:cd+ef'",
			testTokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			testTokens(DATA_ELEMENT_GROUP, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:cd+ef+gh'",
			testTokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			testTokens(DATA_ELEMENT_GROUP, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
		{
			"ab:cd:ef'",
			testTokens(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR, GROUP_DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
			testTokens(DATA_ELEMENT_GROUP, SEGMENT_END_MARKER, EOF),
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
			testTokens(SEGMENT_HEADER, SEGMENT_END_MARKER, EOF),
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
			testTokens(SEGMENT_HEADER, DATA_ELEMENT_SEPARATOR, DATA_ELEMENT, SEGMENT_END_MARKER, EOF),
		},
	}
	for idx, test := range tests {
		var emittedTokens Tokens
		lexer := NewTokenLexer("Test TokenLexer", test.input)
		for lexer.HasNext() {
			tok := lexer.Next()
			emittedTokens = append(emittedTokens, tok)
			if tok.Type() == ERROR {
				t.Logf("Lexer returned error: %#s\n", tok)
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
