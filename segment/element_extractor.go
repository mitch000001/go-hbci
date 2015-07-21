package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

func ExtractElements(segment []byte) ([][]byte, error) {
	extractor := NewElementExtractor(segment)
	return extractor.Extract()
}

func NewElementExtractor(segment []byte) *ElementExtractor {
	return &ElementExtractor{
		rawSegment: segment,
	}
}

type ElementExtractor struct {
	rawSegment []byte
	elements   [][]byte
}

func (e *ElementExtractor) Extract() ([][]byte, error) {
	var current string
	lexer := token.NewStringLexer("ElementExtractor", string(e.rawSegment))
	for lexer.HasNext() {
		t := lexer.Next()
		if t.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", e, t.Pos(), t.Value(), e.rawSegment)
		}
		if t.Type() == token.SEGMENT_END_MARKER || t.Type() == token.DATA_ELEMENT_SEPARATOR {
			e.elements = append(e.elements, []byte(current))
			current = ""
		} else {
			current += t.Value()
		}
	}
	result := make([][]byte, len(e.elements))
	copy(result, e.elements)
	return result, nil
}
