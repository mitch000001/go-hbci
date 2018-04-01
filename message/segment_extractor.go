package message

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/token"
)

// NewSegmentExtractor returns a new SegmentExtractor which extracts segments
// from messageBytes
func NewSegmentExtractor(messageBytes []byte) *SegmentExtractor {
	return &SegmentExtractor{rawMessage: messageBytes}
}

// A SegmentExtractor extracts segments from a raw message and caches its
// results
type SegmentExtractor struct {
	rawMessage []byte
	segments   [][]byte
}

// Extract extract segment byte slices from the message. It will return a copy
// of the result so it is safe to manipulate it
func (s *SegmentExtractor) Extract() ([][]byte, error) {
	var current []byte
	lexer := token.NewLexer("SegmentExtractor", s.rawMessage)
	for lexer.HasNext() {
		t := lexer.Next()
		if t.Type() == token.ERROR {
			return nil, fmt.Errorf("%T: SyntaxError at position %d: %q\n(%q)", s, t.Pos(), t.Value(), s.rawMessage)
		}
		current = append(current, t.Value()...)
		if t.Type() == token.SEGMENT_END_MARKER {
			s.segments = append(s.segments, current)
			current = []byte{}
		}
	}
	result := make([][]byte, len(s.segments))
	copy(result, s.segments)
	return result, nil
}

// FindSegment searches for a given ID and returns the first appearing segment
// bytes if present
func (s *SegmentExtractor) FindSegment(id string) []byte {
	byteID := []byte(id)
	for _, segment := range s.segments {
		if bytes.HasPrefix(segment, byteID) {
			return segment
		}
	}
	return nil
}

// FindSegments finds all segments for a given ID and returns them if present
func (s *SegmentExtractor) FindSegments(id string) [][]byte {
	segmentMap := make(map[string][][]byte)
	for _, segment := range s.segments {
		segmentID := bytes.SplitN(segment, []byte(":"), 2)[0]
		mappedSegments, ok := segmentMap[string(segmentID)]
		if !ok {
			mappedSegments = make([][]byte, 0)
		}
		mappedSegments = append(mappedSegments, segment)
		segmentMap[string(segmentID)] = mappedSegments
	}
	return segmentMap[id]
}

// Segments returns all found segments in order of appearance. It is safe to
// manipulate the result as it is a copy.
func (s *SegmentExtractor) Segments() [][]byte {
	result := make([][]byte, len(s.segments))
	copy(result, s.segments)
	return result
}
