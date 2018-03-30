package token

import "fmt"

// Token represents a HBCI token.
type Token interface {
	Type() Type
	Value() string
	Pos() int
	IsSyntaxSymbol() bool
}

// A Iterator iterates over a slice of Tokens
type Iterator struct {
	tokens []Token
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

func (e elementToken) String() string {
	switch e.typ {
	case EOF:
		return "EOF"
	case ERROR:
		return e.val
	}
	if len(e.val) > 10 {
		return fmt.Sprintf("%s(%.10q...)", e.typ.String(), e.val)
	}
	return fmt.Sprintf("%s(%q)", e.typ.String(), e.val)
}

const (
	ILLEGAL Type = iota // An illegal/unknown character
	ERROR               // error occurred;

	DATA_ELEMENT_SEPARATOR       // Datenelement (DE)-Trennzeichen
	GROUP_DATA_ELEMENT_SEPARATOR // Gruppendatenelement (GD)-Trennzeichen
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
	DATA_ELEMENT_SEPARATOR:       "dataElementSeparator",
	GROUP_DATA_ELEMENT_SEPARATOR: "groupDataElementSeparator",
	SEGMENT_END_MARKER:           "segmentEndMarker",
	BINARY_DATA_LENGTH:           "binaryDataLength",
	BINARY_DATA:                  "binaryData",
	BINARY_DATA_MARKER:           "binaryDataMarker",
	ALPHA_NUMERIC:                "alphaNumeric",
	TEXT:                         "text",
	DTAUS_CHARSET:                "dtausCharset",
	NUMERIC:                      "numeric",
	DIGIT:                        "digit",
	FLOAT:                        "float",
	YES_NO:                       "yesNo",
	DATE:                         "date",
	VIRTUAL_DATE:                 "virtualDate",
	TIME:                         "time",
	IDENTIFICATION:               "identification",
	COUNTRY_CODE:                 "countryCode",
	CURRENCY:                     "currency",
	VALUE:                        "value",
	EOF:                          "eof",
}

func (t Type) String() string {
	s := tokenName[t]
	if s == "" {
		return fmt.Sprintf("Token%d", int(t))
	}
	return s
}
