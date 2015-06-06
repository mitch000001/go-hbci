package hbci

import "fmt"
import "github.com/mitch000001/go-hbci/token"

func NewTokenLexer(name string, input token.Tokens) *TokenLexer {
	t := &TokenLexer{
		name:   name,
		input:  input,
		state:  lexStart,
		tokens: make(chan token.Token, 2), // Two token sufficient.
	}
	return t
}

type TokenLexer struct {
	name   string // the name of the input; used only for error reports.
	input  []token.Token
	state  tokenLexerStateFn
	pos    int // current position in the input.
	start  int // start position of this token.
	tokens chan token.Token
}

// Next returns the next item from the input.
func (l *TokenLexer) Next() token.Token {
	for {
		select {
		case item, ok := <-l.tokens:
			if ok {
				return item
			} else {
				panic("No items left")
			}
		default:
			l.state = l.state(l)
			if l.state == nil {
				close(l.tokens)
			}
		}
	}
	panic("not reached")
}

// HasNext returns true if there are tokens left, false if EOF has reached
func (l *TokenLexer) HasNext() bool {
	return l.state != nil
}

// emit passes a token back to the client.
func (l *TokenLexer) emit(t token.TokenType) {
	l.tokens <- token.NewGroupToken(t, l.input[l.start:l.pos]...)
	l.start = l.pos
}

// emitToken passes the provided token dorectly back to the client
// without wrapping into a GroupToken.
func (l *TokenLexer) emitToken(t token.Token) {
	l.tokens <- t
	l.start = l.pos
}

// next returns the next Token in the input.
func (l *TokenLexer) next() token.Token {
	if l.pos >= len(l.input) {
		return token.NewToken(token.EOF, "", l.pos)
	}
	t := l.input[l.pos]
	l.pos += 1
	return t
}

// ignore skips over the pending input before this point.
func (l *TokenLexer) ignore() {
	l.start = l.pos
}

// backup steps back one Token.
// Can be called only once per call of next.
func (l *TokenLexer) backup() {
	l.pos -= 1
}

// reset steps back until the last emited Token.
func (l *TokenLexer) reset() {
	l.pos = l.start
}

// peek returns but does not consume
// the next Token in the input.
func (l *TokenLexer) peek() token.Token {
	t := l.next()
	l.backup()
	return t
}

// accept consumes the next Token
// if it's from the valid set.
func (l *TokenLexer) accept(valid ...token.TokenType) bool {
	if includes(l.next(), valid...) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of Tokens from the valid set.
func (l *TokenLexer) acceptRun(valid ...token.TokenType) {
	for includes(l.next(), valid...) {
	}
	l.backup()
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *TokenLexer) errorf(format string, args ...interface{}) tokenLexerStateFn {
	l.tokens <- token.NewGroupToken(token.ERROR, token.NewToken(token.ERROR, fmt.Sprintf(format, args...), l.start))
	return nil
}

func includes(t token.Token, tokens ...token.TokenType) bool {
	for _, typ := range tokens {
		if typ == t.Type() {
			return true
		}
	}
	return false
}

type tokenLexerStateFn func(*TokenLexer) tokenLexerStateFn

// The first state function which is called
func lexStart(l *TokenLexer) tokenLexerStateFn {
	// Perform one run to see if there are escape sequences within the token set
	includesEscapeSequence := false
	for _, t := range l.input {
		if includes(t, token.ESCAPE_SEQUENCE) {
			includesEscapeSequence = true
			break
		}
	}
	if includesEscapeSequence {
		return lexEscapeSequenceToken
	} else {
		if t := l.peek(); t.Type() == token.DATA_ELEMENT_GROUP {
			return lexSegmentHeader
		} else {
			return lexTokens
		}
	}
}

func lexTokens(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == token.GROUP_DATA_ELEMENT_SEPARATOR:
			l.reset()
			return lexGroupDataElement
		case t.Type() == token.DATA_ELEMENT_SEPARATOR:
			l.reset()
			return lexDataElement
		case t.Type() == token.GROUP_DATA_ELEMENT:
			l.reset()
			return lexDataElementGroup
		case t.Type() == token.DATA_ELEMENT:
			l.emitToken(t)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		case t.Type() == token.EOF:
			l.emit(token.EOF)
			return nil
		}
	}
	return l.errorf("Syntax error")
}

func lexEscapeSequenceToken(l *TokenLexer) tokenLexerStateFn {
	t := l.next()
	switch t.Type() {
	case token.ALPHA_NUMERIC:
		if l.accept(token.ESCAPE_SEQUENCE) {
			l.accept(token.ALPHA_NUMERIC, token.TEXT, token.NUMERIC, token.DIGIT)
			l.emit(token.ALPHA_NUMERIC_WITH_ESCAPE_SEQUENCE)
		} else {
			l.emitToken(t)
		}
		return lexEscapeSequenceToken
	case token.TEXT:
		if l.accept(token.ESCAPE_SEQUENCE) {
			l.accept(token.ALPHA_NUMERIC, token.TEXT, token.NUMERIC, token.DIGIT)
			l.emit(token.TEXT_WITH_ESCAPE_SEQUENCE)
		} else {
			l.emitToken(t)
		}
		return lexEscapeSequenceToken
	case token.EOF:
		l.emit(token.EOF)
		return nil
	default:
		l.emitToken(t)
		return lexEscapeSequenceToken
	}
	return l.errorf("Syntax error")
}

func lexGroupDataElement(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == token.GROUP_DATA_ELEMENT_SEPARATOR:
			l.backup()
			l.emit(token.GROUP_DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexGroupDataElement, l)
		case t.Type() == token.DATA_ELEMENT_SEPARATOR:
			l.backup()
			l.emit(token.GROUP_DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexDataElement, l)
		case t.Type() == token.SEGMENT_END_MARKER:
			l.backup()
			l.emit(token.GROUP_DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		}
	}
	return l.errorf("Syntax error")
}

func lexDataElement(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == token.DATA_ELEMENT_SEPARATOR:
			l.backup()
			l.emit(token.DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexDataElement, l)
		case t.Type() == token.SEGMENT_END_MARKER:
			l.backup()
			l.emit(token.DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		}
	}
	return l.errorf("Syntax error")
}

func lexDataElementGroup(l *TokenLexer) tokenLexerStateFn {
	l.acceptRun(token.GROUP_DATA_ELEMENT, token.GROUP_DATA_ELEMENT_SEPARATOR)
	l.emit(token.DATA_ELEMENT_GROUP)
	return lexSyntaxSymbolWithContext(lexTokens, l)
}

func lexSegmentHeader(l *TokenLexer) tokenLexerStateFn {
	// Token is a DATA_ELEMENT_GROUP
	t := l.next()
	rawTokens := t.RawTokens()
	var tokensWithoutSeparators token.Tokens
	for _, tok := range rawTokens {
		if !tok.IsSyntaxSymbol() {
			tokensWithoutSeparators = append(tokensWithoutSeparators, tok)
		}
	}
	iter := token.NewTokenIterator(tokensWithoutSeparators)
	if acceptTokenSequence(iter, token.ALPHA_NUMERIC, token.NUMERIC, token.NUMERIC) {
		acceptToken(iter, token.NUMERIC)
		if iter.HasNext() {
			return l.errorf("Malformed Segment Header")
		} else {
			l.emit(token.SEGMENT_HEADER)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		}
	} else {
		return l.errorf("Malformed Segment Header")
	}
}

func acceptTokenSequence(tokens *token.TokenIterator, validSequence ...token.TokenType) bool {
	for _, typ := range validSequence {
		token := tokens.Next()
		if typ != token.Type() {
			return false
		}
	}
	return true
}

func acceptToken(tokens *token.TokenIterator, valid ...token.TokenType) bool {
	if includes(tokens.Next(), valid...) {
		return true
	}
	tokens.Backup()
	return false
}

func lexSyntaxSymbolWithContext(followingStateFn tokenLexerStateFn, l *TokenLexer) tokenLexerStateFn {
	return func(*TokenLexer) tokenLexerStateFn {
		switch t := l.next(); {
		case t.Type() == token.GROUP_DATA_ELEMENT_SEPARATOR:
			l.emitToken(t)
		case t.Type() == token.DATA_ELEMENT_SEPARATOR:
			l.emitToken(t)
		case t.Type() == token.SEGMENT_END_MARKER:
			l.emitToken(t)
		default:
			l.errorf("Syntax error")
		}
		return followingStateFn
	}
}
