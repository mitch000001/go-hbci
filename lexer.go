package hbci

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred;
	// value is text of error
	itemDataElement               // Datenelement (DE)
	itemDataElementSeparator      // Datenelement (DE)-Trennzeichen
	itemGroupDataElement          // Gruppendatenelement (GD)
	itemGroupDataElementSeparator // Gruppendatenelement (GD)-Trennzeichen
	itemSegment                   // Segment
	itemSegmentHeader             // Segmentende-Zeichen
	itemSegmentEnd                // Segmentende-Zeichen
	itemEscapeSequence            // Freigabezeichen
	itemBinaryLength              // Binärdaten Länge
	itemBinaryData                // Binärdaten
	itemBinaryMarker              // Binärdatenkennzeichen
	itemAlphaNumeric              // an
	itemText                      // txt
	itemDTAUSCharset              // dta
	itemNumeric                   // num: 0-9 without leading 0
	itemDigit                     // dig: 0-9 with optional leading 0
	itemFloat                     // float
	itemYesNo                     // jn
	itemDate                      // dat
	itemVirtualDate               // vdat
	itemTime                      // tim
	itemIdentification            // id
	itemCountryCode               // ctr: ISO 3166-1 numeric
	itemCurrency                  // cur: ISO 4217
	itemValue                     // wrt
	itemEOF
)

const eof = -1

// item represents a token returned from the scanner.
type item struct {
	typ itemType // Type, such as itemNumber.
	val string   // Value, such as "23.2".
	pos int      // position of item in input
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		state: lexText,
		items: make(chan item, 2), // Two items sufficient.
	}
	return l
}

type lexer struct {
	name  string    // the name of the input; used only for error reports.
	input string    // the string being scanned.
	state stateFn   // the next lexing function to enter
	pos   int       // current position in the input.
	start int       // start position of this item.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos], l.start}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...), l.start}
	return nil
}

// state functions

const (
	dataElementSeparator      = '+'
	groupDataElementSeparator = ':'
	segmentEnd                = '\''
	escapeCharacter           = '?'
	binaryIdentifier          = '@'
)

func lexText(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == escapeCharacter:
			if p := l.peek(); isSyntaxSymbol(p) {
				l.next()
			} else {
				return l.errorf("Unexpected '?' at pos %d", l.pos)
			}
		case r == dataElementSeparator:
			l.emit(itemDataElementSeparator)
			return lexText
		case r == segmentEnd:
			l.emit(itemSegmentEnd)
			return lexText
		case r == groupDataElementSeparator:
			l.emit(itemGroupDataElementSeparator)
			return lexText
		case r == binaryIdentifier:
			l.backup()
			return lexBinaryData
		case r == eof:
			return l.errorf("Unexpected end of input")
		case ('0' <= r && r <= '9'):
			l.backup()
			return lexDigit

		}
	}
	return nil
}

func lexAlphaNumeric(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == escapeCharacter:
			if p := l.peek(); isSyntaxSymbol(p) {
				l.next()
			} else {
				return l.errorf("Unexpected '?' at pos %d", l.pos)
			}
		case r == dataElementSeparator:
			l.backup()
			return lexSyntaxSymbol
		case r == segmentEnd:
			l.backup()
			return lexSyntaxSymbol
		case r == groupDataElementSeparator:
			l.backup()
			return lexSyntaxSymbol
		case r == binaryIdentifier:
			l.backup()
			return lexBinaryData
		case r == eof:
			return l.errorf("Unexpected end of input")
		}
	}
}

func lexSyntaxSymbol(l *lexer) stateFn {
	switch r := l.next(); {
	case r == dataElementSeparator:
		l.emit(itemDataElementSeparator)
		return lexText
	case r == segmentEnd:
		l.emit(itemSegmentEnd)
		return lexText
	case r == groupDataElementSeparator:
		l.emit(itemGroupDataElementSeparator)
		return lexText
	default:
		return l.errorf("Unexpected syntax symbol: %c\n", r)
	}
}

func lexBinaryData(l *lexer) stateFn {
	l.accept("@")
	digits := "0123456789"
	binaryLengthStart := l.pos
	l.acceptRun(digits)
	binaryLengthEnd := l.pos
	if binaryLengthEnd == binaryLengthStart {
		return l.errorf("Binary length can't be empty")
	}
	length, err := strconv.Atoi(l.input[binaryLengthStart:binaryLengthEnd])
	if err != nil {
		return l.errorf("Binary length must contain of digits only")
	}
	l.accept("@")
	l.pos += length
	if p := l.peek(); isSyntaxSymbol(p) {
		l.emit(itemBinaryData)
		return lexSyntaxSymbol
	} else {
		return l.errorf("Expected syntax symbol after binary data")
	}
}

func lexDigit(l *lexer) stateFn {
	// Is it a number?
	leadingZero := l.accept("0")
	if r := l.peek(); leadingZero && (r < '0' || '9' < r) {
		l.emit(itemNumeric)
		return lexText
	}
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(",") {
		l.acceptRun(digits)
		l.emit(itemFloat)
		return lexText
	}
	l.emit(itemDigit)
	return lexText
}

func isSyntaxSymbol(r rune) bool {
	return r == '+' || r == '\'' || r == ':' || r == '@' || r == '?'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
