package test_files

import "github.com/mitch000001/go-hbci/element"

// These methods are for testing purpose only, just to make the compiler happy

func ExtractElements([]byte) ([][]byte, error)                { return nil, nil }
func SegmentFromHeaderBytes([]byte, Segment) (Segment, error) { return nil, nil }

type Segment interface {
	elements() []element.DataElement
}
