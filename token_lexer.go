package hbci

import "fmt"

func NewTokenLexer(name string, input []Token) *TokenLexer {
	t := &TokenLexer{
		name:   name,
		input:  input,
		state:  lex,
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
	panic("not reached")
}

// HasNext returns true if there are tokens left, false if EOF has reached
func (l *TokenLexer) HasNext() bool {
	return l.state != nil
}

// emit passes an item back to the client.
func (l *TokenLexer) emit(t TokenType) {
	l.tokens <- NewGroupToken(t, l.input[l.start:l.pos]...)
	l.start = l.pos
}

// next returns the next Token in the input.
func (l *TokenLexer) next() Token {
	if l.pos >= len(l.input) {
		return NewElementToken(EOF, "", l.pos)
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
	if includes(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of Tokens from the valid set.
func (l *TokenLexer) acceptRun(valid ...TokenType) {
	for includes(valid, l.next()) {
	}
	l.backup()
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *TokenLexer) errorf(format string, args ...interface{}) tokenLexerStateFn {
	l.tokens <- NewGroupToken(ERROR, ElementToken{ERROR, fmt.Sprintf(format, args...), l.start})
	return nil
}

func includes(tokens []TokenType, t Token) bool {
	for _, typ := range tokens {
		if typ == t.Type() {
			return true
		}
	}
	return false
}

type tokenLexerStateFn func(*TokenLexer) tokenLexerStateFn

func lex(l *TokenLexer) tokenLexerStateFn {
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
		case t.Type() == EOF:
			l.emit(EOF)
			return nil
		}
	}
	return l.errorf("Syntax error")
}

func lexGroupDataElement(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == GROUP_DATA_ELEMENT_SEPARATOR:
			l.emit(GROUP_DATA_ELEMENT)
			return lexGroupDataElement
		case t.Type() == DATA_ELEMENT_SEPARATOR:
			l.emit(GROUP_DATA_ELEMENT)
			return lexDataElement
		case t.Type() == SEGMENT_END_MARKER:
			l.emit(GROUP_DATA_ELEMENT)
			return lex
		}
	}
	return l.errorf("Syntax error")
}

func lexDataElement(l *TokenLexer) tokenLexerStateFn {
	for {
		switch t := l.next(); {
		case t.Type() == DATA_ELEMENT_SEPARATOR:
			l.emit(DATA_ELEMENT)
			return lexDataElement
		case t.Type() == SEGMENT_END_MARKER:
			l.emit(DATA_ELEMENT)
			return lex
		}
	}
	return l.errorf("Syntax error")
}

func lexDataElementGroup(l *TokenLexer) tokenLexerStateFn {
	l.acceptRun(GROUP_DATA_ELEMENT)
	l.emit(DATA_ELEMENT_GROUP)
	return lex
}
