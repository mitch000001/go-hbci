package token

func NewSwiftLexer(name, input string) *SwiftLexer {
	lexer := NewStringLexer(name, input)
	lexer.SetEntryPoint(lexSwiftEntryPoint)
	return &SwiftLexer{lexer}
}

type SwiftLexer struct {
	*StringLexer
}

func lexSwiftEntryPoint(l *StringLexer) StringLexerStateFn {
	if l.accept(string(carriageReturn)) && l.accept(string(lineFeed)) {
		l.emit(SWIFT_DATASET_START)
		return lexSwiftStart
	} else {
		return l.errorf("Malformed swift dataset")
	}
}

func lexSwiftStart(l *StringLexer) StringLexerStateFn {
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

func lexSwiftSyntaxSymbol(l *StringLexer) StringLexerStateFn {
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
	} else {
		return l.errorf("Malformed syntax symbol")
	}
}

func lexSwiftDigit(l *StringLexer) StringLexerStateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(",") {
		digits := "0123456789"
		l.acceptRun(digits)
		if p := l.peek(); p == carriageReturn {
			l.emit(SWIFT_DECIMAL)
			return lexSwiftStart
		} else {
			return lexSwiftAlphaNumeric
		}
	} else {
		if isTagBoundary(l) {
			l.emit(SWIFT_NUMERIC)
			return lexSwiftStart
		} else {
			return lexSwiftCharacter
		}
	}
}

func lexSwiftAlpha(l *StringLexer) StringLexerStateFn {
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

func lexSwiftCharacter(l *StringLexer) StringLexerStateFn {
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

func lexSwiftAlphaNumeric(l *StringLexer) StringLexerStateFn {
	r := l.next()
	switch {
	case r == carriageReturn:
		l.backup()
		if isTagBoundary(l) { // are we really on a tag boundary
			l.emit(SWIFT_ALPHANUMERIC)
			return lexSwiftSyntaxSymbol
		} else {
			l.next()
			return lexSwiftAlphaNumeric
		}
	case r == eof:
		return l.errorf("Unexpected end of input")
	case isSwiftAlphaNumeric(r):
		return lexSwiftAlphaNumeric
	default:
		return lexSwiftAlphaNumeric
	}
}

func isTagBoundary(s *StringLexer) bool {
	currentPos := s.pos
	isTagBoundary := s.accept(string(carriageReturn)) && s.accept(string(lineFeed)) && s.accept(string(dash)+string(tagIdentifier))
	s.pos = currentPos
	return isTagBoundary
}

func isSwiftAlphaNumeric(r rune) bool {
	return r == dash || r == lineFeed || r == ' ' || ('\'' <= r && r <= ')') || ('+' <= r && r <= ':') || r == '?' || ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func peekAfterSequence(l *StringLexer, valid string) rune {
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

var swiftTokenName = map[TokenType]string{
	SWIFT_ALPHA:     "a",
	SWIFT_CHARACTER: "c",
	SWIFT_DECIMAL:   "d",
	SWIFT_NUMERIC:   "n",
}

func init() {
	for k, v := range swiftTokenName {
		tokenName[k] = v
	}
}
