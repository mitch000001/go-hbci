package hbci

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

func NewSegmentExtractor(messageBytes []byte) *SegmentExtractor {
	return &SegmentExtractor{rawMessage: messageBytes}
}

type SegmentExtractor struct {
	rawMessage []byte
	segments   [][]byte
}

func (s *SegmentExtractor) Extract() ([][]byte, error) {
	var current string
	lexer := NewStringLexer("SegmentExtractor", string(s.rawMessage))
	for lexer.HasNext() {
		t := lexer.Next()
		if t.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError: %q", s, t.Value())
		}
		if t.Type() == token.SEGMENT_END_MARKER {
			s.segments = append(s.segments, []byte(current))
			current = ""
		} else {
			current += t.Value()
		}
	}
	return s.segments, nil
}

func (s *SegmentExtractor) FindSegment(id string) []byte {
	byteId := []byte(id)
	for _, segment := range s.segments {
		if bytes.HasPrefix(segment, byteId) {
			return segment
		}
	}
	return nil
}

type sortedByteArrays [][]byte

// Len is the number of elements in the collection.
func (s sortedByteArrays) Len() int { return len(s) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (s sortedByteArrays) Less(i, j int) bool { return string(s[i]) < string(s[j]) }

// Swap swaps the elements with indexes i and j.
func (s sortedByteArrays) Swap(i, j int) { s[j], s[i] = s[i], s[j] }
