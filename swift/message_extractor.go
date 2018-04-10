package swift

import (
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

// NewMessageExtractor returns a message extractor feasable extracting
// S.W.I.F.T. messages from the given input
func NewMessageExtractor(swiftMessage []byte) *MessageExtractor {
	lexer := token.NewSwiftLexer("MessageExtractor", swiftMessage)
	return &MessageExtractor{
		lexer:           lexer,
		rawSwiftMessage: swiftMessage,
	}
}

// MessageExtractor represents an extractor for S.W.I.F.T. messages
type MessageExtractor struct {
	lexer             *token.SwiftLexer
	rawSwiftMessage   []byte
	extractedMessages [][]byte
}

// Extract extracts raw S.W.I.F.T. messages from the given input
func (m *MessageExtractor) Extract() ([][]byte, error) {
	var current []byte
	for m.lexer.HasNext() {
		tok := m.lexer.Next()
		if tok.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", m, tok.Pos(), tok.Value(), m.rawSwiftMessage)
		}
		current = append(current, tok.Value()...)
		if tok.Type() == token.SWIFT_MESSAGE_SEPARATOR {
			m.extractedMessages = append(m.extractedMessages, current)
			current = []byte{}
		}
	}
	result := make([][]byte, len(m.extractedMessages))
	copy(result, m.extractedMessages)
	return result, nil
}
