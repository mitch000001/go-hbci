package test_files

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

// These methods are for testing purpose only, just to make the compiler happy

func ExtractElements([]byte) ([][]byte, error)                { return nil, nil }
func SegmentFromHeaderBytes([]byte, Segment) (Segment, error) { return nil, nil }

type Segment interface {
	ID() string
	Version() int
	elements() []element.DataElement
}

type versionInterface1 interface {
}

type versionInterface2 interface {
}

type BankSegment interface {
	Segment
	UnmarshalHBCI([]byte) error
}

type VersionedSegment struct {
	ID      string
	Version int
}

func (v VersionedSegment) String() string {
	return fmt.Sprintf("%s:%d", v.ID, v.Version)
}

type SegmentIndex struct {
	segmentMap map[VersionedSegment]func() Segment
}

func (s SegmentIndex) IsIndexed(segmentId VersionedSegment) bool {
	_, ok := s.segmentMap[segmentId]
	return ok
}

func (s *SegmentIndex) addToIndex(segmentIdentifier VersionedSegment, segmentProviderFn func() Segment) error {
	if s.IsIndexed(segmentIdentifier) {
		return fmt.Errorf("Segment already in index: %s", segmentIdentifier)
	}
	s.segmentMap[segmentIdentifier] = segmentProviderFn
	return nil
}

func (s *SegmentIndex) mustAddToIndex(segmentIdentifier VersionedSegment, segmentProviderFn func() Segment) {
	err := s.addToIndex(segmentIdentifier, segmentProviderFn)
	if err != nil {
		panic(err)
	}
}

var KnownSegments = SegmentIndex{segmentMap: make(map[VersionedSegment]func() Segment)}
