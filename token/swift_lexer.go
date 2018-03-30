package token

import (
	"strings"
)

// NewSwiftLexer returns a SwiftLexer ready for parsing the given input string
func NewSwiftLexer(name, input string) *SwiftLexer {
	lexer := NewLexer(name, input)
	lexer.SetEntryPoint(lexSwiftEntryPoint)
	return &SwiftLexer{lexer}
}

// A SwiftLexer parses the given input and emits SWIFT tokens
type SwiftLexer struct {
	*Lexer
}

func lexSwiftEntryPoint(l *Lexer) LexerStateFn {
	if l.accept(string(carriageReturn)) && l.accept(string(lineFeed)) {
		l.emit(SWIFT_DATASET_START)
		return lexSwiftStart
	}
	return l.errorf("Malformed swift dataset")
}

func lexSwiftStart(l *Lexer) LexerStateFn {
	r := l.next()
	switch {
	case ('0' <= r && r <= '9'):
		return lexSwiftDigit
	case ('A' <= r && r <= 'Z'):
		return lexSwiftAlpha
	case r == carriageReturn:
		l.backup()
		return lexSwiftSyntaxSymbol
	case r == eof:
		// Correctly reached EOF.
		l.emit(EOF)
		return nil
	default:
		return lexSwiftAlphaNumeric
	}
}

func lexSwiftSyntaxSymbol(l *Lexer) LexerStateFn {
	if l.accept(string(carriageReturn)) && l.accept(string(lineFeed)) {
		p := l.peek()
		switch {
		case p == eof:
			return l.errorf("Unexpected end of input")
		case p == dash:
			l.next()
			l.emit(SWIFT_MESSAGE_SEPARATOR)
		case p != dash:
			l.emit(SWIFT_TAG_SEPARATOR)
		}
		return lexSwiftStart
	}
	return l.errorf("Malformed syntax symbol")
}

func lexSwiftDigit(l *Lexer) LexerStateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(",") {
		digits := "0123456789"
		l.acceptRun(digits)
		if p := l.peek(); p == carriageReturn {
			l.emit(SWIFT_DECIMAL)
			return lexSwiftStart
		}
		return lexSwiftAlphaNumeric
	}
	if isTagBoundary(l) {
		l.emit(SWIFT_NUMERIC)
		return lexSwiftStart
	}
	return lexSwiftCharacter
}

func lexSwiftAlpha(l *Lexer) LexerStateFn {
	switch r := l.next(); {
	case ('A' <= r && r <= 'Z'):
		return lexSwiftAlpha
	case ('0' <= r && r <= '9'):
		return lexSwiftCharacter
	case r == carriageReturn:
		l.backup()
		l.emit(SWIFT_ALPHA)
		return lexSwiftStart
	case r == eof:
		return l.errorf("Unexpected end of input")
	case isSwiftAlphaNumeric(r):
		return lexSwiftAlphaNumeric
	default:
		return lexSwiftAlphaNumeric
	}
}

func lexSwiftCharacter(l *Lexer) LexerStateFn {
	switch r := l.next(); {
	case ('A' <= r && r <= 'Z'):
		return lexSwiftCharacter
	case ('0' <= r && r <= '9'):
		return lexSwiftCharacter
	case r == carriageReturn:
		l.backup()
		l.emit(SWIFT_CHARACTER)
		return lexSwiftStart
	case r == eof:
		return l.errorf("Unexpected end of input")
	case isSwiftAlphaNumeric(r):
		return lexSwiftAlphaNumeric
	default:
		return lexSwiftAlphaNumeric
	}
}

func lexSwiftAlphaNumeric(l *Lexer) LexerStateFn {
	r := l.next()
	switch {
	case r == carriageReturn:
		l.backup()
		if isTagBoundary(l) { // are we really on a tag boundary
			l.emit(SWIFT_ALPHANUMERIC)
			return lexSwiftSyntaxSymbol
		}
		l.next()
		return lexSwiftAlphaNumeric
	case r == eof:
		return l.errorf("Unexpected end of input")
	case isSwiftAlphaNumeric(r):
		return lexSwiftAlphaNumeric
	default:
		return lexSwiftAlphaNumeric
	}
}

func isTagBoundary(s *Lexer) bool {
	oneOf := func(fn ...func() bool) bool {
		for _, f := range fn {
			if ok := f(); ok {
				return true
			}
		}
		return false
	}
	currentPos := s.pos
	isTagBoundary := s.accept(string(carriageReturn)) &&
		s.accept(string(lineFeed)) && oneOf(
		func() bool {
			return strings.HasPrefix(s.input[s.pos:], string(dash)+string(tagIdentifier))
		},
		func() bool {
			return strings.HasPrefix(s.input[s.pos:], string(dash)+string(carriageReturn))
		},
		func() bool {
			return strings.HasPrefix(s.input[s.pos:], string(tagIdentifier))
		},
		func() bool {
			return s.input[s.pos:] == string(dash)
		},
	)
	s.pos = currentPos
	return isTagBoundary
}

func isSwiftAlphaNumeric(r rune) bool {
	return r == dash || r == lineFeed || r == ' ' || ('\'' <= r && r <= ')') || ('+' <= r && r <= ':') || r == '?' || ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func peekAfterSequence(l *Lexer, valid string) rune {
	currentPos := l.pos
	l.acceptRun(valid)
	r := l.peek()
	l.pos = currentPos
	return r
}

const (
	carriageReturn           = '\r'
	lineFeed                 = '\n'
	dash                     = '-'
	tagIdentifier            = ':'
	tagSeparatorSequence     = "\r\n"
	messageSeparatorSequence = "\r\n-"
)

const (
	SWIFT_ALPHA        = EOF + iota + 1 // 'A' - 'Z'
	SWIFT_CHARACTER                     // 'A' - 'Z', '0' - '9'
	SWIFT_DECIMAL                       // '0' - '9', ','
	SWIFT_NUMERIC                       // '0' - '9'
	SWIFT_ALPHANUMERIC                  // all characters from charset
	SWIFT_DATASET_START
	SWIFT_TAG_SEPARATOR
	SWIFT_MESSAGE_SEPARATOR
)

var swiftTokenName = map[Type]string{
	SWIFT_ALPHA:             "a",
	SWIFT_CHARACTER:         "c",
	SWIFT_DECIMAL:           "d",
	SWIFT_NUMERIC:           "n",
	SWIFT_ALPHANUMERIC:      "an",
	SWIFT_DATASET_START:     "datasetStart",
	SWIFT_TAG_SEPARATOR:     "tagSeparator",
	SWIFT_MESSAGE_SEPARATOR: "messageSeparator",
}

func init() {
	for k, v := range swiftTokenName {
		tokenName[k] = v
	}
}
