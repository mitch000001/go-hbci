package hbci

import (
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

func NewSegmentExtractor(messageBytes []byte) *SegmentExtractor {
	return &SegmentExtractor{rawMessage: messageBytes}
}

type SegmentExtractor struct {
	rawMessage []byte
}

func (s *SegmentExtractor) Extract() ([][]byte, error) {
	var segments [][]byte
	var current string
	lexer := NewStringLexer("SegmentExtractor", string(s.rawMessage))
	for lexer.HasNext() {
		t := lexer.Next()
		if t.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError: %q", s, t.Value())
		}
		if t.Type() == token.SEGMENT_END_MARKER {
			segments = append(segments, []byte(current))
			current = ""
		} else {
			current += t.Value()
		}
	}
	return segments, nil
}
