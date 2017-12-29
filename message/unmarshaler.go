package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

func NewUnmarshaler(message []byte) *Unmarshaler {
	return &Unmarshaler{
		rawMessage:       message,
		segmentExtractor: segment.NewSegmentExtractor(message),
		segments:         make(map[string][]segment.Segment),
	}
}

type Unmarshaler struct {
	rawMessage       []byte
	segmentExtractor *segment.SegmentExtractor
	segments         map[string][]segment.Segment
}

func (u *Unmarshaler) CanUnmarshal(segmentId string, version int) bool {
	return segment.KnownSegments.IsUnmarshaler(
		segment.VersionedSegment{
			ID:      segmentId,
			Version: version,
		},
	)
}

func (u *Unmarshaler) Unmarshal() error {
	rawSegments, err := u.segmentExtractor.Extract()
	if err != nil {
		return err
	}
	for _, seg := range rawSegments {
		segmentId, err := extractVersionedSegmentIdentifier(seg)
		if err != nil {
			return err
		}
		if segment.KnownSegments.IsUnmarshaler(segmentId) {
			unmarshaler, err := segment.KnownSegments.UnmarshalerForSegment(segmentId)
			if err != nil {
				return err
			}
			err = unmarshaler.UnmarshalHBCI(seg)
			if err != nil {
				return err
			}
			segments, ok := u.segments[segmentId.ID]
			if !ok {
				segments = make([]segment.Segment, 0)
			}
			segments = append(segments, unmarshaler.(segment.Segment))
			u.segments[segmentId.ID] = segments
		}
	}
	return nil
}

func (u *Unmarshaler) UnmarshalSegment(segmentId string, version int) (segment.Segment, error) {
	unmarshaler, err := segment.KnownSegments.UnmarshalerForSegment(
		segment.VersionedSegment{ID: segmentId, Version: version},
	)
	if err != nil {
		return nil, err
	}
	segmentBytes, err := u.extractSegment(segmentId)
	if err != nil {
		return nil, err
	}
	err = unmarshaler.UnmarshalHBCI(segmentBytes)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling segment: %v", err)
	}
	return unmarshaler.(segment.Segment), nil
}

func (u *Unmarshaler) extractSegment(segmentId string) ([]byte, error) {
	_, err := u.segmentExtractor.Extract()
	if err != nil {
		return nil, err
	}
	segmentBytes := u.segmentExtractor.FindSegment(segmentId)
	if segmentBytes == nil {
		return nil, fmt.Errorf("Segment not found in message: %q", segmentId)
	}
	return segmentBytes, nil
}

func (u *Unmarshaler) SegmentsById(segmentId string) []segment.Segment {
	return u.segments[segmentId]
}

func (u *Unmarshaler) SegmentById(segmentId string) segment.Segment {
	segments, ok := u.segments[segmentId]
	if ok {
		return segments[0]
	}
	return nil
}

func (u *Unmarshaler) MarshaledSegmentsById(segmentId string) [][]byte {
	return u.segmentExtractor.FindSegments(segmentId)
}

func (u *Unmarshaler) MarshaledSegmentById(segmentId string) []byte {
	return u.segmentExtractor.FindSegment(segmentId)
}

func (u *Unmarshaler) MarshaledSegments() [][]byte {
	return u.segmentExtractor.Segments()
}

func extractVersionedSegmentIdentifier(segmentBytes []byte) (segment.VersionedSegment, error) {
	var id segment.VersionedSegment
	elements, err := segment.ExtractElements(segmentBytes)
	if err != nil {
		return id, err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return id, err
	}
	id = segment.VersionedSegment{ID: header.ID.Val(), Version: header.Version.Val()}
	return id, nil
}
