package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

// NewUnmarshaler returns a new Unmarshaler for the message
func NewUnmarshaler(message []byte) *Unmarshaler {
	return &Unmarshaler{
		rawMessage:       message,
		segmentExtractor: NewSegmentExtractor(message),
		segments:         make(map[string][]segment.Segment),
	}
}

// Unmarshaler unmarshals a complete message
type Unmarshaler struct {
	rawMessage       []byte
	segmentExtractor *SegmentExtractor
	segments         map[string][]segment.Segment
}

// CanUnmarshal returns true if the segment with the ID and version can be
// unmarshaled, false otherwise
func (u *Unmarshaler) CanUnmarshal(segmentID string, version int) bool {
	return segment.KnownSegments.IsUnmarshaler(
		segment.VersionedSegment{
			ID:      segmentID,
			Version: version,
		},
	)
}

// Unmarshal unmarshals the raw message
func (u *Unmarshaler) Unmarshal() error {
	rawSegments, err := u.segmentExtractor.Extract()
	if err != nil {
		return err
	}
	for _, seg := range rawSegments {
		header, err := extractSegmentHeader(seg)
		if err != nil {
			return err
		}
		segmentID := segment.VersionedSegment{ID: header.ID.Val(), Version: header.Version.Val()}
		if segment.KnownSegments.IsUnmarshaler(segmentID) {
			unmarshaler, err := segment.KnownSegments.UnmarshalerForSegment(segmentID)
			if err != nil {
				return err
			}
			err = unmarshaler.UnmarshalHBCI(seg)
			if err != nil {
				return fmt.Errorf("error unmarshaling segment %q: %w", header, err)
			}
			segments, ok := u.segments[segmentID.ID]
			if !ok {
				segments = make([]segment.Segment, 0)
			}
			segments = append(segments, unmarshaler.(segment.Segment))
			u.segments[segmentID.ID] = segments
		}
	}
	return nil
}

// UnmarshalSegment unmarshals the segment with the given ID and version
func (u *Unmarshaler) UnmarshalSegment(segmentID string, version int) (segment.Segment, error) {
	unmarshaler, err := segment.KnownSegments.UnmarshalerForSegment(
		segment.VersionedSegment{ID: segmentID, Version: version},
	)
	if err != nil {
		return nil, err
	}
	segmentBytes, err := u.extractSegment(segmentID)
	if err != nil {
		return nil, err
	}
	err = unmarshaler.UnmarshalHBCI(segmentBytes)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling segment: %v", err)
	}
	return unmarshaler.(segment.Segment), nil
}

func (u *Unmarshaler) extractSegment(segmentID string) ([]byte, error) {
	_, err := u.segmentExtractor.Extract()
	if err != nil {
		return nil, err
	}
	segmentBytes := u.segmentExtractor.FindSegment(segmentID)
	if segmentBytes == nil {
		return nil, fmt.Errorf("Segment not found in message: %q", segmentID)
	}
	return segmentBytes, nil
}

// SegmentsByID returns all already unmarshaled segments for the given ID
func (u *Unmarshaler) SegmentsByID(segmentID string) []segment.Segment {
	return u.segments[segmentID]
}

// SegmentByID returns the first segment found for ID
func (u *Unmarshaler) SegmentByID(segmentID string) segment.Segment {
	segments, ok := u.segments[segmentID]
	if ok {
		return segments[0]
	}
	return nil
}

// MarshaledSegmentsByID returns all segments for a given ID as array of bytes
func (u *Unmarshaler) MarshaledSegmentsByID(segmentID string) [][]byte {
	return u.segmentExtractor.FindSegments(segmentID)
}

// MarshaledSegmentByID returns the first segment for the given ID as array of bytes
func (u *Unmarshaler) MarshaledSegmentByID(segmentID string) []byte {
	return u.segmentExtractor.FindSegment(segmentID)
}

// MarshaledSegments return all marshaled segments as array of bytes
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

func extractSegmentHeader(segmentBytes []byte) (*element.SegmentHeader, error) {
	elements, err := segment.ExtractElements(segmentBytes)
	if err != nil {
		return nil, err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return nil, err
	}
	return header, nil
}
