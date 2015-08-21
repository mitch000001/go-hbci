package segment

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitch000001/go-hbci/charset"
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
	String() string
}

type ClientSegment interface {
	Segment
	MarshalHBCI() ([]byte, error)
}

type BankSegment interface {
	Segment
	Unmarshaler
}

type basicSegment interface {
	Version() int
	ID() string
	referencedId() string
	sender() string
	elements() []element.DataElement
}

type Unmarshaler interface {
	UnmarshalHBCI([]byte) error
}

func SegmentFromHeaderBytes(headerBytes []byte, seg basicSegment) (*segment, error) {
	elements, err := element.ExtractElements(headerBytes)
	var header *element.SegmentHeader
	id := charset.ToUtf8(elements[0])
	numStr := elements[1]
	number, err := strconv.Atoi(charset.ToUtf8(numStr))
	if err != nil {
		return nil, fmt.Errorf("Malformed segment header number")
	}
	version, err := strconv.Atoi(charset.ToUtf8(elements[2]))
	if err != nil {
		return nil, fmt.Errorf("Malformed segment header version")
	}
	if len(elements) == 4 && len(elements[3]) > 0 {
		ref, err := strconv.Atoi(charset.ToUtf8(elements[3]))
		if err != nil {
			return nil, fmt.Errorf("Malformed segment header reference")
		}
		header = element.NewReferencingSegmentHeader(id, number, version, ref)
	} else {
		header = element.NewSegmentHeader(id, number, version)
	}
	return NewBasicSegmentWithHeader(header, seg), nil
}

func NewReferencingBasicSegment(number int, ref int, seg basicSegment) *segment {
	header := element.NewReferencingSegmentHeader(seg.ID(), number, seg.Version(), ref)
	return NewBasicSegmentWithHeader(header, seg)
}

func NewBasicSegment(number int, seg basicSegment) *segment {
	header := element.NewSegmentHeader(seg.ID(), number, seg.Version())
	return NewBasicSegmentWithHeader(header, seg)
}

func NewBasicSegmentWithHeader(header *element.SegmentHeader, seg basicSegment) *segment {
	return &segment{header: header, segment: seg}
}

type segment struct {
	segment basicSegment
	header  *element.SegmentHeader
}

func (s *segment) String() string {
	elementStrings := make([]string, len(s.segment.elements())+1)
	elementStrings[0] = s.header.String()
	for i, de := range s.segment.elements() {
		val := reflect.ValueOf(de)
		if val.IsValid() && !val.IsNil() {
			elementStrings[i+1] = de.String()
		}
	}
	return strings.Join(elementStrings, "+") + "'"
}

func (s *segment) MarshalHBCI() ([]byte, error) {
	elementBytes := make([][]byte, len(s.segment.elements())+1)
	headerBytes, err := s.header.MarshalHBCI()
	if err != nil {
		return nil, err
	}
	elementBytes[0] = headerBytes
	for i, de := range s.segment.elements() {
		val := reflect.ValueOf(de)
		if val.IsValid() && !val.IsNil() {
			marshaled, err := de.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			elementBytes[i+1] = marshaled
		}
	}
	marshaled := bytes.Join(elementBytes, []byte("+"))
	marshaled = append(marshaled, '\'')
	return marshaled, nil
}

func (s *segment) DataElements() []element.DataElement {
	var dataElements []element.DataElement
	dataElements = append(dataElements, s.header)
	dataElements = append(dataElements, s.segment.elements()...)
	return dataElements
}

func (s *segment) Header() *element.SegmentHeader {
	return s.header
}

func (s *segment) ID() string {
	return s.header.ID.Val()
}

func (s *segment) Version() int {
	return s.header.Version.Val()
}

func (s *segment) SetNumber(numberFn func() int) {
	s.header.SetNumber(numberFn())
}

func (s *segment) SetReference(ref int) {
	s.header.SetReference(ref)
}
