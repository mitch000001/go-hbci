package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

// ExtractElements extracts DataElements from a DataElementGroup
func ExtractElements(dataElementGroup []byte) ([][]byte, error) {
	extractor := NewElementExtractor(dataElementGroup)
	return extractor.Extract()
}

// NewElementExtractor creates a new GroupExtractor ready to use
func NewElementExtractor(dataElementGroup []byte) *GroupExtractor {
	// TODO: workaround to get the lexer work properly for us. Maybe we should adopt the lexer?
	if bytes.LastIndexFunc(dataElementGroup, func(r rune) bool {
		return r == '+' || r == '\''
	}) == -1 {
		dataElementGroup = append(dataElementGroup, '+')
	}
	return &GroupExtractor{
		rawDataElementGroup: dataElementGroup,
	}
}

// An GroupExtractor extracts DataElements from DataElementGroups
type GroupExtractor struct {
	rawDataElementGroup []byte
	elements            [][]byte
}

// Extract extracts DataElements from the underlying DataElementGroup
func (e *GroupExtractor) Extract() ([][]byte, error) {
	var current string
	lexer := token.NewLexer("ElementExtractor", string(e.rawDataElementGroup))
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
