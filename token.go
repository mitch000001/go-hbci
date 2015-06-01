package hbci

import "fmt"

type Token interface {
	Type() TokenType
	Value() string
	Pos() int
	IsSyntaxSymbol() bool
	Children() Tokens
	RawTokens() Tokens
}

type Tokens []Token

func (t Tokens) Types() []TokenType {
	var types []TokenType
	for _, token := range t {
		types = append(types, token.Type())
	}
	return types
}

func (t Tokens) RawTokens() Tokens {
	var tokens Tokens
	for _, token := range t {
		tokens = append(tokens, token.RawTokens()...)
	}
	return tokens
}

func NewTokenIterator(tokens Tokens) *TokenIterator {
	return &TokenIterator{tokens: tokens, pos: 0}
}

type TokenIterator struct {
	tokens Tokens
	pos    int
}

func (t *TokenIterator) HasNext() bool {
	return t.pos < len(t.tokens)
}

func (t *TokenIterator) Next() Token {
	if t.pos >= len(t.tokens) {
		return NewToken(EOF, "", t.pos)
	}
	token := t.tokens[t.pos]
	t.pos += 1
	return token
}

func (t *TokenIterator) Backup() {
	t.pos -= 1
}

func NewGroupToken(typ TokenType, tokens ...Token) Token {
	groupToken := groupToken{elementToken: elementToken{typ: typ}, tokens: tokens}
	val := ""
	for _, token := range tokens {
		val += token.Value()
	}
	groupToken.val = val
	return groupToken
}

type groupToken struct {
	elementToken
	tokens Tokens
}

func (g groupToken) Children() Tokens {
	return g.tokens
}

func (g groupToken) RawTokens() Tokens {
	var tokens Tokens
	for _, token := range g.Children() {
		if len(token.Children()) > 0 {
			tokens = append(tokens, token.RawTokens()...)
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func NewToken(typ TokenType, val string, pos int) Token {
	return elementToken{typ, val, pos}
}

// elementToken represents a token returned from the scanner.
type elementToken struct {
	typ TokenType // Type, such as FLOAT
	val string    // Value, such as "23.2".
	pos int       // position of token in input
}

func (e elementToken) Type() TokenType {
	return e.typ
}

func (e elementToken) Value() string {
	return e.val
}
func (e elementToken) Pos() int {
	return e.pos
}

func (e elementToken) IsSyntaxSymbol() bool {
	switch e.typ {
	case GROUP_DATA_ELEMENT_SEPARATOR:
		return true
	case DATA_ELEMENT_SEPARATOR:
		return true
	case SEGMENT_END_MARKER:
		return true
	default:
		return false
	}
}

func (e elementToken) Children() Tokens {
	return Tokens{}
}

func (e elementToken) RawTokens() Tokens {
	return Tokens{e}
}

func (t elementToken) String() string {
	switch t.typ {
	case EOF:
		return "EOF"
	case ERROR:
		return t.val
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

const (
	ERROR TokenType = iota // error occurred;
	// value is text of error
	DATA_ELEMENT                       // Datenelement (DE)
	DATA_ELEMENT_SEPARATOR             // Datenelement (DE)-Trennzeichen
	DATA_ELEMENT_GROUP                 // Datenelementgruppe (DEG)
	GROUP_DATA_ELEMENT                 // Gruppendatenelement (GD)
	GROUP_DATA_ELEMENT_SEPARATOR       // Gruppendatenelement (GD)-Trennzeichen
	SEGMENT                            // Segment
	SEGMENT_HEADER                     // Segmentende-Zeichen
	SEGMENT_END_MARKER                 // Segmentende-Zeichen
	ESCAPE_SEQUENCE                    // Freigabeabfolge
	ESCAPE_CHARACTER                   // Freigabezeichen
	ESCAPED_CHARACTER                  // Freigegebenes Zeichen
	BINARY_DATA_LENGTH                 // Binärdaten Länge
	BINARY_DATA                        // Binärdaten
	BINARY_DATA_MARKER                 // Binärdatenkennzeichen
	ALPHA_NUMERIC                      // an
	ALPHA_NUMERIC_WITH_ESCAPE_SEQUENCE // an with an escape sequence
	TEXT                               // txt
	TEXT_WITH_ESCAPE_SEQUENCE          // txt with an escape sequence
	DTAUS_CHARSET                      // dta
	NUMERIC                            // num: 0-9 without leading 0
	DIGIT                              // dig: 0-9 with optional leading 0
	FLOAT                              // float
	YES_NO                             // jn
	DATE                               // dat
	VIRTUAL_DATE                       // vdat
	TIME                               // tim
	IDENTIFICATION                     // id
	COUNTRY_CODE                       // ctr: ISO 3166-1 numeric
	CURRENCY                           // cur: ISO 4217
	VALUE                              // wrt
	EOF
)

type TokenTypes []TokenType

func (t TokenTypes) Len() int           { return len(t) }
func (t TokenTypes) Less(i, j int) bool { return t[i] < t[j] }
func (t TokenTypes) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// TokenType identifies the type of lex tokens.
type TokenType int

var tokenName = map[TokenType]string{
	ERROR: "error",
	// value is text of error
	DATA_ELEMENT:                 "dataElement",
	DATA_ELEMENT_SEPARATOR:       "dataElementSeparator",
	DATA_ELEMENT_GROUP:           "dataElementGroup",
	GROUP_DATA_ELEMENT:           "groupDataElement",
	GROUP_DATA_ELEMENT_SEPARATOR: "groupDataElementSeparator",
	SEGMENT:                            "segment",
	SEGMENT_HEADER:                     "segmentHeader",
	SEGMENT_END_MARKER:                 "segmentEndMarker",
	ESCAPE_SEQUENCE:                    "escapeSequence",
	ESCAPE_CHARACTER:                   "escapeCharacter",
	ESCAPED_CHARACTER:                  "escapedCharacter",
	BINARY_DATA_LENGTH:                 "binaryDataLength",
	BINARY_DATA:                        "binaryData",
	BINARY_DATA_MARKER:                 "binaryDataMarker",
	ALPHA_NUMERIC:                      "alphaNumeric",
	ALPHA_NUMERIC_WITH_ESCAPE_SEQUENCE: "alphaNumericWithEscapeSequence",
	TEXT: "text",
	TEXT_WITH_ESCAPE_SEQUENCE: "textWithEscapeSequence",
	DTAUS_CHARSET:             "dtausCharset",
	NUMERIC:                   "numeric",
	DIGIT:                     "digit",
	FLOAT:                     "float",
	YES_NO:                    "yesNo",
	DATE:                      "date",
	VIRTUAL_DATE:              "virtualDate",
	TIME:                      "time",
	IDENTIFICATION:            "identification",
	COUNTRY_CODE:              "countryCode",
	CURRENCY:                  "currency",
	VALUE:                     "value",
	EOF:                       "eof",
}

func (t TokenType) String() string {
	s := tokenName[t]
	if s == "" {
		return fmt.Sprintf("Token%d", int(t))
	}
	return s
}
