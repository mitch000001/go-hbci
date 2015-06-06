package hbci

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/token"
)

type testToken int

func (t testToken) Type() token.TokenType {
	return token.TokenType(t)
}

func (t testToken) Value() string {
	return ""
}

func (t testToken) Pos() int {
	return 0
}

func (t testToken) IsSyntaxSymbol() bool {
	return token.TokenType(t) == token.GROUP_DATA_ELEMENT_SEPARATOR || token.TokenType(t) == token.DATA_ELEMENT_SEPARATOR || token.TokenType(t) == token.SEGMENT_END_MARKER
}

func (t testToken) Children() token.Tokens {
	return token.Tokens{}
}

func (t testToken) RawTokens() token.Tokens {
	return token.Tokens{t}
}

func (t testToken) String() string {
	return token.TokenType(t).String()
}

func testTokens(types ...token.TokenType) token.Tokens {
	var tokens token.Tokens
	for _, typ := range types {
		tokens = append(tokens, testToken(typ))
	}
	return tokens
}

type tokenLexerTestData struct {
	example string
	input   token.Tokens
	output  token.Tokens
}

func TestTokenLexer(t *testing.T) {
	tests := []tokenLexerTestData{
		{
			"ab??2'",
			testTokens(token.ALPHA_NUMERIC, token.ESCAPE_SEQUENCE, token.NUMERIC, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.ALPHA_NUMERIC_WITH_ESCAPE_SEQUENCE, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:23'",
			testTokens(token.ALPHA_NUMERIC, token.GROUP_DATA_ELEMENT_SEPARATOR, token.NUMERIC, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab+23'",
			testTokens(token.ALPHA_NUMERIC, token.DATA_ELEMENT_SEPARATOR, token.NUMERIC, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:23:05+45'",
			testTokens(token.ALPHA_NUMERIC, token.GROUP_DATA_ELEMENT_SEPARATOR, token.NUMERIC, token.GROUP_DATA_ELEMENT_SEPARATOR, token.DIGIT, token.DATA_ELEMENT_SEPARATOR, token.NUMERIC, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab::01+5'",
			testTokens(token.ALPHA_NUMERIC, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT_SEPARATOR, token.DIGIT, token.DATA_ELEMENT_SEPARATOR, token.NUMERIC, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:cd+ef'",
			testTokens(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.DATA_ELEMENT_GROUP, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:cd+ef+gh'",
			testTokens(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.DATA_ELEMENT_GROUP, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:cd:ef'",
			testTokens(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR, token.GROUP_DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
			testTokens(token.DATA_ELEMENT_GROUP, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:1:2:3'",
			token.Tokens{
				token.NewGroupToken(token.DATA_ELEMENT_GROUP,
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.ALPHA_NUMERIC)),
					testToken(token.GROUP_DATA_ELEMENT_SEPARATOR),
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.NUMERIC)),
					testToken(token.GROUP_DATA_ELEMENT_SEPARATOR),
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.NUMERIC)),
					testToken(token.GROUP_DATA_ELEMENT_SEPARATOR),
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.NUMERIC)),
				),
				testToken(token.SEGMENT_END_MARKER),
				testToken(token.EOF),
			},
			testTokens(token.SEGMENT_HEADER, token.SEGMENT_END_MARKER, token.EOF),
		},
		{
			"ab:1:2:3+4'",
			token.Tokens{
				token.NewGroupToken(token.DATA_ELEMENT_GROUP,
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.ALPHA_NUMERIC)),
					testToken(token.GROUP_DATA_ELEMENT_SEPARATOR),
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.NUMERIC)),
					testToken(token.GROUP_DATA_ELEMENT_SEPARATOR),
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.NUMERIC)),
					testToken(token.GROUP_DATA_ELEMENT_SEPARATOR),
					token.NewGroupToken(token.GROUP_DATA_ELEMENT, testToken(token.NUMERIC)),
				),
				testToken(token.DATA_ELEMENT_SEPARATOR),
				testToken(token.DATA_ELEMENT),
				testToken(token.SEGMENT_END_MARKER),
				testToken(token.EOF),
			},
			testTokens(token.SEGMENT_HEADER, token.DATA_ELEMENT_SEPARATOR, token.DATA_ELEMENT, token.SEGMENT_END_MARKER, token.EOF),
		},
	}
	for idx, test := range tests {
		var emittedTokens token.Tokens
		lexer := NewTokenLexer("Test TokenLexer", test.input)
		for lexer.HasNext() {
			tok := lexer.Next()
			emittedTokens = append(emittedTokens, tok)
			if tok.Type() == token.ERROR {
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
