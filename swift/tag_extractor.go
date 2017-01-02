package swift

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/token"
)

func ExtractTagElements(tag []byte) ([][]byte, error) {
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

func ExtractTagID(tag []byte) ([]byte, error) {
	elements, err := ExtractTagElements(tag)
	if err != nil {
		return nil, err
	} else {
		return elements[0], nil
	}
}

func NewTagExtractor(swiftMessage []byte) *TagExtractor {
	lexer := token.NewSwiftLexer("TagExtractor", charset.ToUTF8(swiftMessage))
	return &TagExtractor{
		lexer:           lexer,
		rawSwiftMessage: swiftMessage,
	}
}

type TagExtractor struct {
	lexer           *token.SwiftLexer
	rawSwiftMessage []byte
	extractedTags   [][]byte
}

func (t *TagExtractor) Extract() ([][]byte, error) {
	var current string
	for t.lexer.HasNext() {
		tok := t.lexer.Next()
		if tok.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", t, tok.Pos(), tok.Value(), t.rawSwiftMessage)
		}
		if tok.Type() == token.SWIFT_TAG_SEPARATOR || tok.Type() == token.SWIFT_MESSAGE_SEPARATOR {
			t.extractedTags = append(t.extractedTags, []byte(current))
			current = ""
		} else {
			if tok.Type() != token.SWIFT_DATASET_START {
				current += tok.Value()
			}
		}
	}
	result := make([][]byte, len(t.extractedTags))
	copy(result, t.extractedTags)
	return result, nil
}
