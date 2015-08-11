package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci"
	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/token"
)

func RegisterUnmarshaler(segmentId string, unmarshalerFn func() hbci.Unmarshaler) {
	knownUnmarshalers[segmentId] = unmarshalerFn
}

type unmarshalerIndex map[string]func() hbci.Unmarshaler

func (u unmarshalerIndex) Unmarshaler(segmentId string) (hbci.Unmarshaler, bool) {
	unmarhsalerFn, ok := u[segmentId]
	if ok {
		return unmarhsalerFn(), ok
	} else {
		return nil, ok
	}
}

var knownUnmarshalers = unmarshalerIndex{
	"HNHBK": func() hbci.Unmarshaler { return &segment.MessageHeaderSegment{} },
	"HIRMG": func() hbci.Unmarshaler { return &segment.MessageAcknowledgement{} },
	"HIRMS": func() hbci.Unmarshaler { return &segment.SegmentAcknowledgement{} },
	"HNVSD": func() hbci.Unmarshaler { return &segment.EncryptedDataSegment{} },
	"HISYN": func() hbci.Unmarshaler { return &segment.SynchronisationResponseSegment{} },
	"HIUPD": func() hbci.Unmarshaler { return &segment.AccountInformationSegment{} },
}

func NewUnmarshaler(message []byte) *Unmarshaler {
	return &Unmarshaler{
		rawMessage:       message,
		segmentExtractor: segment.NewSegmentExtractor(message),
	}
}

type Unmarshaler struct {
	rawMessage       []byte
	segmentExtractor *segment.SegmentExtractor
	segments         map[string][]segment.Segment
}

func (u *Unmarshaler) CanUnmarshal(segmentId string) bool {
	_, ok := knownUnmarshalers[segmentId]
	return ok
}

func (u *Unmarshaler) Unmarshal() error {
	rawSegments, err := u.segmentExtractor.Extract()
	if err != nil {
		return err
	}
	for _, seg := range rawSegments {
		segmentId, err := extractSegmentID(seg)
		if err != nil {
			return err
		}
		unmarshaler, ok := knownUnmarshalers.Unmarshaler(segmentId)
		if ok {
			err = unmarshaler.UnmarshalHBCI(seg)
			if err != nil {
				return err
			}
			segments, ok := u.segments[segmentId]
			if !ok {
				segments = make([]segment.Segment, 0)
			}
			segments = append(segments, unmarshaler.(segment.Segment))
			u.segments[segmentId] = segments
		}
	}
	return nil
}

func (u *Unmarshaler) UnmarshalSegment(segmentId string) (segment.Segment, error) {
	unmarshaler, ok := knownUnmarshalers.Unmarshaler(segmentId)
	if !ok {
		return nil, fmt.Errorf("Unknown segment: %q", segmentId)
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

func extractSegmentID(segment []byte) (string, error) {
	lexer := token.NewStringLexer("SegmentIdExtractor", string(segment))
	if lexer.HasNext() {
		idToken := lexer.Next()
		if idToken.Type() != token.ALPHA_NUMERIC {
			return "", fmt.Errorf("Malformed segment: segment ID not alphanumeric")
		} else {
			return idToken.Value(), nil
		}
	} else {
		return "", fmt.Errorf("Malformed segment: empty")
	}
}
