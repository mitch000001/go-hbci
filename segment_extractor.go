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
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", s, t.Pos(), t.Value(), s.rawMessage)
		}
		current += t.Value()
		if t.Type() == token.SEGMENT_END_MARKER {
			s.segments = append(s.segments, []byte(current))
			current = ""
		}
	}
	result := make([][]byte, len(s.segments))
	copy(result, s.segments)
	return result, nil
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

func (s *SegmentExtractor) FindSegments(id string) [][]byte {
	segmentMap := make(map[string][][]byte)
	for _, segment := range s.segments {
		segmentId := bytes.SplitN(segment, []byte(":"), 2)[0]
		mappedSegments, ok := segmentMap[string(segmentId)]
		if !ok {
			mappedSegments = make([][]byte, 0)
		}
		mappedSegments = append(mappedSegments, segment)
		segmentMap[string(segmentId)] = mappedSegments
	}
	return segmentMap[id]
}

type sortedSegments [][]byte

// Len is the number of elements in the collection.
func (s sortedSegments) Len() int { return len(s) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (s sortedSegments) Less(i, j int) bool { return string(s[i]) < string(s[j]) }

// Swap swaps the elements with indexes i and j.
func (s sortedSegments) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
