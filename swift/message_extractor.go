package swift

import (
	"fmt"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/token"
)

func NewMessageExtractor(swiftMessage []byte) *MessageExtractor {
	lexer := token.NewSwiftLexer("MessageExtractor", charset.ToUTF8(swiftMessage))
	return &MessageExtractor{
		lexer:           lexer,
		rawSwiftMessage: swiftMessage,
	}
}

type MessageExtractor struct {
	lexer             *token.SwiftLexer
	rawSwiftMessage   []byte
	extractedMessages [][]byte
}

func (m *MessageExtractor) Extract() ([][]byte, error) {
	var current string
	for m.lexer.HasNext() {
		tok := m.lexer.Next()
		if tok.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", m, tok.Pos(), tok.Value(), m.rawSwiftMessage)
		}
		current += tok.Value()
		if tok.Type() == token.SWIFT_MESSAGE_SEPARATOR {
			m.extractedMessages = append(m.extractedMessages, []byte(current))
			current = ""
		}
	}
	result := make([][]byte, len(m.extractedMessages))
	copy(result, m.extractedMessages)
	return result, nil
}
