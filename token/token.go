package token

import "fmt"

// Token represents a HBCI token.
type Token interface {
	Type() Type
	Value() string
	Pos() int
	IsSyntaxSymbol() bool
	Children() Tokens
	RawTokens() Tokens
}

// Tokens represent a collection of Token.
// It defines convenient methods on top of the collection.
type Tokens []Token

// Types returns a slice of TokenType in the order of the Tokens withing the
// Token slice.
func (t Tokens) Types() []Type {
	var types []Type
	for _, token := range t {
		types = append(types, token.Type())
	}
	return types
}

// RawTokens returns all raw Tokens returned b RawTokens.
// All RawTokens are appended to a big Tokens slice.
func (t Tokens) RawTokens() Tokens {
	var tokens Tokens
	for _, token := range t {
		tokens = append(tokens, token.RawTokens()...)
	}
	return tokens
}

// NewIterator returns a fully populated TokenIterator
func NewIterator(tokens Tokens) *Iterator {
	return &Iterator{tokens: tokens, pos: 0}
}

// A Iterator iterates over a slice of Tokens
type Iterator struct {
	tokens Tokens
	pos    int
}

// HasNext returns true if there are tokens less to emit, false otherwise.
func (t *Iterator) HasNext() bool {
	return t.pos < len(t.tokens)
}

// Next returns the next Token within the iterator. If there are no more tokens
// it will return an EOF Token signalling that the iterator has reached the
// last element.
func (t *Iterator) Next() Token {
	if t.pos >= len(t.tokens) {
		return New(EOF, "", t.pos)
	}
	token := t.tokens[t.pos]
	t.pos++
	return token
}

// Backup moves the iterator one position back.
func (t *Iterator) Backup() {
	t.pos--
}

// NewGroupToken returns a Token composed of a group of sub tokens and with the
// given type typ.
// The Value method of such a Token returns the values of all sub tokens
// appended in the order provided by tokens.
func NewGroupToken(typ Type, tokens ...Token) Token {
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

// New creates a Token with the given type, value and position
func New(typ Type, val string, pos int) Token {
	return elementToken{typ, val, pos}
}

// elementToken represents a token returned from the scanner.
type elementToken struct {
	typ Type   // Type, such as FLOAT
	val string // Value, such as "23.2".
	pos int    // position of token in input
}

func (e elementToken) Type() Type {
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

func (e elementToken) String() string {
	switch e.typ {
	case EOF:
		return "EOF"
	case ERROR:
		return e.val
	}
	if len(e.val) > 10 {
		return fmt.Sprintf("%.10q...", e.val)
	}
	return fmt.Sprintf("%q", e.val)
}

const (
	ERROR Type = iota // error occurred;

	DATA_ELEMENT                 // Datenelement (DE)
	DATA_ELEMENT_SEPARATOR       // Datenelement (DE)-Trennzeichen
	DATA_ELEMENT_GROUP           // Datenelementgruppe (DEG)
	GROUP_DATA_ELEMENT           // Gruppendatenelement (GD)
	GROUP_DATA_ELEMENT_SEPARATOR // Gruppendatenelement (GD)-Trennzeichen
	SEGMENT                      // Segment
	SEGMENT_HEADER               // Segmentende-Zeichen
	SEGMENT_END_MARKER           // Segmentende-Zeichen
	BINARY_DATA_LENGTH           // Binärdaten Länge
	BINARY_DATA                  // Binärdaten
	BINARY_DATA_MARKER           // Binärdatenkennzeichen
	ALPHA_NUMERIC                // an
	TEXT                         // txt
	DTAUS_CHARSET                // dta
	NUMERIC                      // num: 0-9 without leading 0
	DIGIT                        // dig: 0-9 with optional leading 0
	FLOAT                        // float
	YES_NO                       // jn
	DATE                         // dat
	VIRTUAL_DATE                 // vdat
	TIME                         // tim
	IDENTIFICATION               // id
	COUNTRY_CODE                 // ctr: ISO 3166-1 numeric
	CURRENCY                     // cur: ISO 4217
	VALUE                        // wrt
	EOF
)

// Types implements the sort.Sort interface
type Types []Type

func (t Types) Len() int           { return len(t) }
func (t Types) Less(i, j int) bool { return t[i] < t[j] }
func (t Types) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// Type identifies the type of lex tokens.
type Type int

var tokenName = map[Type]string{
	ERROR: "error",
	// value is text of error
	DATA_ELEMENT:                 "dataElement",
	DATA_ELEMENT_SEPARATOR:       "dataElementSeparator",
	DATA_ELEMENT_GROUP:           "dataElementGroup",
	GROUP_DATA_ELEMENT:           "groupDataElement",
	GROUP_DATA_ELEMENT_SEPARATOR: "groupDataElementSeparator",
	SEGMENT:            "segment",
	SEGMENT_HEADER:     "segmentHeader",
	SEGMENT_END_MARKER: "segmentEndMarker",
	BINARY_DATA_LENGTH: "binaryDataLength",
	BINARY_DATA:        "binaryData",
	BINARY_DATA_MARKER: "binaryDataMarker",
	ALPHA_NUMERIC:      "alphaNumeric",
	TEXT:               "text",
	DTAUS_CHARSET:      "dtausCharset",
	NUMERIC:            "numeric",
	DIGIT:              "digit",
	FLOAT:              "float",
	YES_NO:             "yesNo",
	DATE:               "date",
	VIRTUAL_DATE:       "virtualDate",
	TIME:               "time",
	IDENTIFICATION:     "identification",
	COUNTRY_CODE:       "countryCode",
	CURRENCY:           "currency",
	VALUE:              "value",
	EOF:                "eof",
}

func (t Type) String() string {
	s := tokenName[t]
	if s == "" {
		return fmt.Sprintf("Token%d", int(t))
	}
	return s
}
