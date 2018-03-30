package segment

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

func (s *SegmentExtractor) Segments() [][]byte {
	result := make([][]byte, len(s.segments))
	copy(result, s.segments)
	return result
}
