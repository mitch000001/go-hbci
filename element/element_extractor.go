package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

// ExtractElements extracts DataElements from a DataElementGroup
func ExtractElements(dataElementGroup []byte) ([][]byte, error) {
	extractor := newElementExtractor(dataElementGroup)
	return extractor.Extract()
}

// newElementExtractor creates a new GroupExtractor ready to use
func newElementExtractor(dataElementGroup []byte) *groupExtractor {
	// TODO: workaround to get the lexer work properly for us. Maybe we should adopt the lexer?
	if bytes.LastIndexFunc(dataElementGroup, func(r rune) bool {
		return r == '+' || r == '\''
	}) == -1 {
		dataElementGroup = append(dataElementGroup, '+')
	}
	return &groupExtractor{
		rawDataElementGroup: dataElementGroup,
	}
}

// An groupExtractor extracts DataElements from DataElementGroups
type groupExtractor struct {
	rawDataElementGroup []byte
	elements            [][]byte
}

// Extract extracts DataElements from the underlying DataElementGroup
func (e *groupExtractor) Extract() ([][]byte, error) {
	var current []byte
	lexer := token.NewLexer("ElementExtractor", e.rawDataElementGroup)
	for lexer.HasNext() {
		t := lexer.Next()
		currentType := t.Type()
		if currentType == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", e, t.Pos(), t.Value(), e.rawDataElementGroup)
		}
		if currentType == token.SEGMENT_END_MARKER || currentType == token.DATA_ELEMENT_SEPARATOR || currentType == token.GROUP_DATA_ELEMENT_SEPARATOR {
			e.elements = append(e.elements, current)
			current = []byte{}
		} else {
			current = append(current, t.Value()...)
		}
	}
	result := make([][]byte, len(e.elements))
	copy(result, e.elements)
	return result, nil
}
