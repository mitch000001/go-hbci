package swift

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/token"
)

type rawTag struct {
	ID    string
	Value []byte
}

func extractRawTag(tag []byte) (*rawTag, error) {
	elements, err := extractTagElements(tag)
	if err != nil {
		return nil, err
	}
	if len(elements) != 2 {
		return nil, fmt.Errorf("Malformed marshaled tag")
	}
	return &rawTag{
		ID:    charset.ToUTF8(elements[0]),
		Value: elements[1],
	}, nil
}

func extractTagElements(tag []byte) ([][]byte, error) {
	buf := bytes.NewBuffer(tag)
	b, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != ':' {
		return nil, fmt.Errorf("Malformed tag id")
	}
	tagBytes := make([]byte, 1)
	tagBytes[0] = b
	tagID, err := buf.ReadBytes(':')
	if err != nil {
		return nil, err
	}
	tagBytes = append(tagBytes, tagID...)
	tagValue := buf.Bytes()
	return [][]byte{
		tagBytes,
		tagValue,
	}, nil
}

func extractTagID(tag []byte) ([]byte, error) {
	elements, err := extractTagElements(tag)
	if err != nil {
		return nil, err
	}
	return elements[0], nil
}

func newTagExtractor(swiftMessage []byte) *tagExtractor {
	lexer := token.NewSwiftLexer("TagExtractor", swiftMessage)
	return &tagExtractor{
		lexer:           lexer,
		rawSwiftMessage: swiftMessage,
	}
}

type tagExtractor struct {
	lexer           *token.SwiftLexer
	rawSwiftMessage []byte
	extractedTags   [][]byte
}

func (t *tagExtractor) Extract() ([][]byte, error) {
	var current []byte
	for t.lexer.HasNext() {
		tok := t.lexer.Next()
		if tok.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", t, tok.Pos(), tok.Value(), t.rawSwiftMessage)
		}
		if tok.Type() == token.SWIFT_TAG_SEPARATOR || tok.Type() == token.SWIFT_MESSAGE_SEPARATOR {
			t.extractedTags = append(t.extractedTags, current)
			current = []byte{}
		} else {
			if tok.Type() != token.SWIFT_DATASET_START {
				current = append(current, tok.Value()...)
			}
		}
	}
	result := make([][]byte, len(t.extractedTags))
	copy(result, t.extractedTags)
	return result, nil
}
