package token

import (
	"bytes"
)

const unexpectedEOF = "unexpected end of input"

func IsUnexpectedEndOfInput(t Token) bool {
	if t.Type() != ERROR {
		return false
	}
	return string(t.Value()) == unexpectedEOF
}

// NewSwiftLexer returns a SwiftLexer ready for parsing the given input string
func NewSwiftLexer(name string, input []byte) *SwiftLexer {
	lexer := NewLexer(name, input)
	lexer.SetEntryPoint(lexSwiftEntryPoint)
	return &SwiftLexer{lexer}
}

// A SwiftLexer parses the given input and emits SWIFT tokens
type SwiftLexer struct {
	*Lexer
}

func lexSwiftEntryPoint(l *Lexer) LexerStateFn {
	if l.accept(carriageReturn) && l.accept(lineFeed) {
		l.emit(SWIFT_DATASET_START)
		return lexTagID
	}
	return l.errorf("Malformed swift dataset")
}

func lexTagID(l *Lexer) LexerStateFn {
	l.accept(tagIdentifier)
	digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	if !l.accept(digits...) {
		return l.errorf("malformed swift tag identifier: must start with a digit")
	}
	l.acceptRun(digits)
	if p := l.peek(); 'A' <= p && p <= 'Z' {
		l.next()
	}
	if !l.accept(tagIdentifier) {
		return l.errorf("malformed swift tag identifier: must be enclodes by ':'")
	}
	l.emit(SWIFT_TAG)
	return lexSwiftStart
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
	if l.accept(carriageReturn) && l.accept(lineFeed) {
		p := l.peek()
		switch {
		case p == eof:
			return l.errorf(unexpectedEOF)
		case p == dash:
			l.next()
			l.emit(SWIFT_MESSAGE_SEPARATOR)
		case p != dash:
			l.emit(SWIFT_TAG_SEPARATOR)
			return lexTagID
		}
		return lexSwiftStart
	}
	return l.errorf("Malformed syntax symbol")
}

func lexSwiftDigit(l *Lexer) LexerStateFn {
	digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	l.acceptRun(digits)
	if l.accept(',') {
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
		return l.errorf(unexpectedEOF)
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
		return l.errorf(unexpectedEOF)
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
		return l.errorf(unexpectedEOF)
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
	isTagBoundary := s.accept(carriageReturn) &&
		s.accept(lineFeed) && oneOf(
		func() bool {
			return bytes.HasPrefix(s.input[s.pos:], []byte{dash, tagIdentifier})
		},
		func() bool {
			return bytes.HasPrefix(s.input[s.pos:], []byte{dash, carriageReturn, lineFeed})
		},
		func() bool {
			if len(s.input[s.pos:]) < 3 {
				return false
			}
			tagIDStart := s.input[s.pos+1]
			return bytes.HasPrefix(s.input[s.pos:], []byte{tagIdentifier}) &&
				'0' <= tagIDStart && tagIDStart <= '9'
		},
		func() bool {
			return bytes.Equal(s.input[s.pos:], []byte{dash})
		},
	)
	s.pos = currentPos
	return isTagBoundary
}

func isSwiftAlphaNumeric(r byte) bool {
	return r == dash || r == lineFeed || r == ' ' || ('\'' <= r && r <= ')') || ('+' <= r && r <= ':') || r == '?' || ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
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
	SWIFT_TAG
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
	SWIFT_TAG:               "tag",
	SWIFT_MESSAGE_SEPARATOR: "messageSeparator",
}

func init() {
	for k, v := range swiftTokenName {
		tokenName[k] = v
	}
}
