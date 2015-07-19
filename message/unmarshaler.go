package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci"
	"github.com/mitch000001/go-hbci/segment"
)

func RegisterUnmarshaler(segmentId string, unmarshalerFn func() hbci.Unmarshaler) {
	knownUnmarshaler[segmentId] = unmarshalerFn
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

var knownUnmarshaler = unmarshalerIndex{
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
}

func (u *Unmarshaler) Unmarshal(segmentId string) (segment.Segment, error) {
	unmarshaler, ok := knownUnmarshaler.Unmarshaler(segmentId)
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
