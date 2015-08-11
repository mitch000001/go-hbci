package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

func ExtractElements(dataElementGroup []byte) ([][]byte, error) {
	extractor := NewElementExtractor(dataElementGroup)
	return extractor.Extract()
}

func NewElementExtractor(dataElementGroup []byte) *ElementExtractor {
	// TODO: workaround to get the lexer work properly for us. Maybe we should adopt the lexer?
	if bytes.LastIndexFunc(dataElementGroup, func(r rune) bool {
		return r == '+' || r == '\''
	}) == -1 {
		dataElementGroup = append(dataElementGroup, '+')
	}
	return &ElementExtractor{
		rawDataElementGroup: dataElementGroup,
	}
}

type ElementExtractor struct {
	rawDataElementGroup []byte
	elements            [][]byte
}

func (e *ElementExtractor) Extract() ([][]byte, error) {
	var current string
	lexer := token.NewStringLexer("ElementExtractor", string(e.rawDataElementGroup))
	for lexer.HasNext() {
		t := lexer.Next()
		currentType := t.Type()
		if currentType == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", e, t.Pos(), t.Value(), e.rawDataElementGroup)
		}
		if currentType == token.SEGMENT_END_MARKER || currentType == token.DATA_ELEMENT_SEPARATOR || currentType == token.GROUP_DATA_ELEMENT_SEPARATOR {
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
