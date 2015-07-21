package segment

import (
	"reflect"
	"strings"

	"github.com/mitch000001/go-hbci/element"
)

const (
	senderBank = "I"
	senderUser = "K"
	senderBoth = "K/I"
)

type Segment interface {
	Header() *element.SegmentHeader
	SetNumber(func() int)
	DataElements() []element.DataElement
	String() string
}

type Segments map[string]Segment

type segment interface {
	version() int
	id() string
	referencedId() string
	sender() string
	elements() []element.DataElement
}

func NewReferencingBasicSegment(number int, ref int, seg segment) Segment {
	header := element.NewReferencingSegmentHeader(seg.id(), number, seg.version(), ref)
	return NewBasicSegmentWithHeader(header, seg)
}

func NewBasicSegment(number int, seg segment) Segment {
	header := element.NewSegmentHeader(seg.id(), number, seg.version())
	return NewBasicSegmentWithHeader(header, seg)
}

func NewBasicSegmentWithHeader(header *element.SegmentHeader, seg segment) Segment {
	return &basicSegment{header: header, segment: seg}
}

type basicSegment struct {
	segment segment
	header  *element.SegmentHeader
}

func (s *basicSegment) String() string {
	elementStrings := make([]string, len(s.segment.elements())+1)
	elementStrings[0] = s.header.String()
	for i, de := range s.segment.elements() {
		if !reflect.ValueOf(de).IsNil() {
			elementStrings[i+1] = de.String()
		}
	}
	return strings.Join(elementStrings, "+") + "'"
}

func (s *basicSegment) DataElements() []element.DataElement {
	var dataElements []element.DataElement
	dataElements = append(dataElements, s.header)
	dataElements = append(dataElements, s.segment.elements()...)
	return dataElements
}

func (s *basicSegment) Header() *element.SegmentHeader {
	return s.header
}

func (s *basicSegment) ID() string {
	return s.header.ID.Val()
}

func (s *basicSegment) SetNumber(numberFn func() int) {
	s.header.SetNumber(numberFn())
}

func (s *basicSegment) SetReference(ref int) {
	s.header.SetReference(ref)
}
