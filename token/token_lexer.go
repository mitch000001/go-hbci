package token

import "fmt"

func NewTokenLexer(name string, input Tokens) *TokenLexer {
	t := &TokenLexer{
		name:   name,
		input:  input,
		state:  lexStart,
		tokens: make(chan Token, 2), // Two token sufficient.
	}
	return t
}

type TokenLexer struct {
	name   string // the name of the input; used only for error reports.
	input  []Token
	state  tokenLexerStateFn
	pos    int // current position in the input.
	start  int // start position of this token.
	tokens chan Token
}

// Next returns the next item from the input.
func (l *TokenLexer) Next() Token {
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
}

// HasNext returns true if there are tokens left, false if EOF has reached
func (l *TokenLexer) HasNext() bool {
	return l.state != nil
}

// emit passes a token back to the client.
func (l *TokenLexer) emit(t TokenType) {
	l.tokens <- NewGroupToken(t, l.input[l.start:l.pos]...)
	l.start = l.pos
}

// emitToken passes the provided token dorectly back to the client
// without wrapping into a GroupToken.
func (l *TokenLexer) emitToken(t Token) {
	l.tokens <- t
	l.start = l.pos
}

// next returns the next Token in the input.
func (l *TokenLexer) next() Token {
	if l.pos >= len(l.input) {
		return NewToken(EOF, "", l.pos)
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
func (l *TokenLexer) peek() Token {
	t := l.next()
	l.backup()
	return t
}

// accept consumes the next Token
// if it's from the valid set.
func (l *TokenLexer) accept(valid ...TokenType) bool {
	if includes(l.next(), valid...) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of Tokens from the valid set.
func (l *TokenLexer) acceptRun(valid ...TokenType) {
	for includes(l.next(), valid...) {
	}
	l.backup()
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *TokenLexer) errorf(format string, args ...interface{}) tokenLexerStateFn {
	l.tokens <- NewGroupToken(ERROR, NewToken(ERROR, fmt.Sprintf(format, args...), l.start))
	return nil
}

func includes(t Token, tokens ...TokenType) bool {
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
	if t := l.peek(); t.Type() == DATA_ELEMENT_GROUP {
		return lexSegmentHeader
	} else {
		return lexTokens
	}
}

func lexTokens(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == GROUP_DATA_ELEMENT_SEPARATOR:
			l.reset()
			return lexGroupDataElement
		case t.Type() == DATA_ELEMENT_SEPARATOR:
			l.reset()
			return lexDataElement
		case t.Type() == GROUP_DATA_ELEMENT:
			l.reset()
			return lexDataElementGroup
		case t.Type() == DATA_ELEMENT:
			l.emitToken(t)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		case t.Type() == EOF:
			l.emit(EOF)
			return nil
		}
	}
}

func lexGroupDataElement(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == GROUP_DATA_ELEMENT_SEPARATOR:
			l.backup()
			l.emit(GROUP_DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexGroupDataElement, l)
		case t.Type() == DATA_ELEMENT_SEPARATOR:
			l.backup()
			l.emit(GROUP_DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexDataElement, l)
		case t.Type() == SEGMENT_END_MARKER:
			l.backup()
			l.emit(GROUP_DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		}
	}
}

func lexDataElement(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == DATA_ELEMENT_SEPARATOR:
			l.backup()
			l.emit(DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexDataElement, l)
		case t.Type() == SEGMENT_END_MARKER:
			l.backup()
			l.emit(DATA_ELEMENT)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		}
	}
}

func lexDataElementGroup(l *TokenLexer) tokenLexerStateFn {
	l.acceptRun(GROUP_DATA_ELEMENT, GROUP_DATA_ELEMENT_SEPARATOR)
	l.emit(DATA_ELEMENT_GROUP)
	return lexSyntaxSymbolWithContext(lexTokens, l)
}

func lexSegmentHeader(l *TokenLexer) tokenLexerStateFn {
	// Token is a DATA_ELEMENT_GROUP
	t := l.next()
	rawTokens := t.RawTokens()
	var tokensWithoutSeparators Tokens
	for _, tok := range rawTokens {
		if !tok.IsSyntaxSymbol() {
			tokensWithoutSeparators = append(tokensWithoutSeparators, tok)
		}
	}
	iter := NewTokenIterator(tokensWithoutSeparators)
	if acceptTokenSequence(iter, ALPHA_NUMERIC, NUMERIC, NUMERIC) {
		acceptToken(iter, NUMERIC)
		if iter.HasNext() {
			return l.errorf("Malformed Segment Header")
		} else {
			l.emit(SEGMENT_HEADER)
			return lexSyntaxSymbolWithContext(lexTokens, l)
		}
	} else {
		return l.errorf("Malformed Segment Header")
	}
}

func acceptTokenSequence(tokens *TokenIterator, validSequence ...TokenType) bool {
	for _, typ := range validSequence {
		token := tokens.Next()
		if typ != token.Type() {
			return false
		}
	}
	return true
}

func acceptToken(tokens *TokenIterator, valid ...TokenType) bool {
	if includes(tokens.Next(), valid...) {
		return true
	}
	tokens.Backup()
	return false
}

func lexSyntaxSymbolWithContext(followingStateFn tokenLexerStateFn, l *TokenLexer) tokenLexerStateFn {
	return func(*TokenLexer) tokenLexerStateFn {
		switch t := l.next(); {
		case t.Type() == GROUP_DATA_ELEMENT_SEPARATOR:
			l.emitToken(t)
		case t.Type() == DATA_ELEMENT_SEPARATOR:
			l.emitToken(t)
		case t.Type() == SEGMENT_END_MARKER:
			l.emitToken(t)
		default:
			l.errorf("Syntax error")
		}
		return followingStateFn
	}
}
