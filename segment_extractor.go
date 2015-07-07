package hbci

import "bytes"

func NewSegmentExtractor(messageBytes []byte) *SegmentExtractor {
	return &SegmentExtractor{rawMessage: messageBytes}
}

type SegmentExtractor struct {
	rawMessage []byte
}

func (s *SegmentExtractor) Extract() ([][]byte, error) {
	// TODO: Fix naive extract method
	result := bytes.Split(s.rawMessage, []byte("'"))
	result = result[:len(result)-1]
	return result, nil
}
