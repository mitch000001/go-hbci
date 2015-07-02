package segment

import (
	"reflect"
	"strings"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const (
	senderBank = "I"
	senderUser = "K"
	senderBoth = "K/I"
)

type Segment interface {
	Header() *element.SegmentHeader
	SetNumber(int)
	DataElements() []element.DataElement
	String() string
}

type segment interface {
	init()
	version() int
	id() string
	referencedId() string
	sender() string
	elements() []element.DataElement
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

func (s *basicSegment) SetNumber(number int) {
	s.header.SetNumber(number)
}

func (s *basicSegment) SetReference(ref int) {
	s.header.SetReference(ref)
}

func NewIdentificationSegment(bankId domain.BankId, clientId string, clientSystemId string, systemIdRequired bool) *IdentificationSegment {
	var clientSystemStatus *element.NumberDataElement
	if systemIdRequired {
		clientSystemStatus = element.NewNumber(1, 1)
	} else {
		clientSystemStatus = element.NewNumber(0, 1)
	}
	id := &IdentificationSegment{
		BankId:             element.NewBankIndentification(bankId),
		ClientId:           element.NewIdentification(clientId),
		ClientSystemId:     element.NewIdentification(clientSystemId),
		ClientSystemStatus: clientSystemStatus,
	}
	id.Segment = NewBasicSegment(3, id)
	return id
}

type IdentificationSegment struct {
	Segment
	BankId             *element.BankIdentificationDataElement
	ClientId           *element.IdentificationDataElement
	ClientSystemId     *element.IdentificationDataElement
	ClientSystemStatus *element.NumberDataElement
}

func (i *IdentificationSegment) init() {
	*i.BankId = *new(element.BankIdentificationDataElement)
	*i.ClientId = *new(element.IdentificationDataElement)
	*i.ClientSystemId = *new(element.IdentificationDataElement)
	*i.ClientSystemStatus = *new(element.NumberDataElement)
}
func (i *IdentificationSegment) version() int         { return 2 }
func (i *IdentificationSegment) id() string           { return "HKIDN" }
func (i *IdentificationSegment) referencedId() string { return "" }
func (i *IdentificationSegment) sender() string       { return senderUser }

func (i *IdentificationSegment) elements() []element.DataElement {
	return []element.DataElement{
		i.BankId,
		i.ClientId,
		i.ClientSystemId,
		i.ClientSystemStatus,
	}
}
