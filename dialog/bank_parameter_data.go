package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

type BankParameterData struct {
	domain.BankParameterData   `yaml:",inline"`
	SupportedSegmentParameters []SegmentParameter `yaml:"supportedSegments"`
}

type SegmentParameter struct {
	segment.VersionedSegment `yaml:",inline"`
	Parameters               segment.Segment `yaml:",omitempty"`
}
