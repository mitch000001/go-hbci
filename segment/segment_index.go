package segment

import "fmt"

type SegmentIndex map[string]func() Segment

func (s SegmentIndex) UnmarshalerForSegment(segmentId string) (Unmarshaler, error) {
	segmentFn, ok := s[segmentId]
	if ok {
		unmarshaler, ok := segmentFn().(Unmarshaler)
		if ok {
			return unmarshaler, nil
		} else {
			return nil, fmt.Errorf("Segment does not implement the Unmarshaler interface")
		}
	} else {
		return nil, fmt.Errorf("Segment not in index: %q", segmentId)
	}
}

func (s SegmentIndex) IsIndexed(segmentId string) bool {
	_, ok := s[segmentId]
	return ok
}

func (s SegmentIndex) IsUnmarshaler(segmentId string) bool {
	segmentFn, ok := s[segmentId]
	if ok {
		_, ok := segmentFn().(Unmarshaler)
		return ok
	} else {
		return false
	}
}

var knownSegments = SegmentIndex{
	"HNHBK": func() Segment { return &MessageHeaderSegment{} },
	"HNHBS": func() Segment { return &MessageEndSegment{} },
	"HNVSK": func() Segment { return &EncryptionHeaderSegment{} },
	"HNVSD": func() Segment { return &EncryptedDataSegment{} },
	"HIRMG": func() Segment { return &MessageAcknowledgement{} },
	"HIRMS": func() Segment { return &SegmentAcknowledgement{} },
}
