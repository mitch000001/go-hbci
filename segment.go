package hbci

import (
	"reflect"
	"strings"

	"github.com/mitch000001/go-hbci/dataelement"
	"github.com/mitch000001/go-hbci/domain"
)

type Segment interface {
	Header() *dataelement.SegmentHeader
	SetNumber(int)
	DataElements() []dataelement.DataElement
	String() string
}

type segment interface {
	elements() []dataelement.DataElement
}

func NewBasicSegment(id string, number int, version int, seg segment) Segment {
	header := dataelement.NewSegmentHeader(id, number, version)
	return NewBasicSegmentWithHeader(header, seg)
}

func NewBasicSegmentWithHeader(header *dataelement.SegmentHeader, seg segment) Segment {
	return &basicSegment{header: header, segment: seg}
}

type basicSegment struct {
	segment segment
	header  *dataelement.SegmentHeader
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

func (s *basicSegment) DataElements() []dataelement.DataElement {
	var dataElements []dataelement.DataElement
	dataElements = append(dataElements, s.header)
	dataElements = append(dataElements, s.segment.elements()...)
	return dataElements
}

func (s *basicSegment) Header() *dataelement.SegmentHeader {
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
	var clientSystemStatus *dataelement.NumberDataElement
	if systemIdRequired {
		clientSystemStatus = dataelement.NewNumberDataElement(1, 1)
	} else {
		clientSystemStatus = dataelement.NewNumberDataElement(0, 1)
	}
	id := &IdentificationSegment{
		BankId:             dataelement.NewBankIndentificationDataElement(bankId),
		ClientId:           dataelement.NewIdentificationDataElement(clientId),
		ClientSystemId:     dataelement.NewIdentificationDataElement(clientSystemId),
		ClientSystemStatus: clientSystemStatus,
	}
	id.Segment = NewBasicSegment("HKIDN", 3, 2, id)
	return id
}

type IdentificationSegment struct {
	Segment
	BankId             *dataelement.BankIdentificationDataElement
	ClientId           *dataelement.IdentificationDataElement
	ClientSystemId     *dataelement.IdentificationDataElement
	ClientSystemStatus *dataelement.NumberDataElement
}

func (i *IdentificationSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		i.BankId,
		i.ClientId,
		i.ClientSystemId,
		i.ClientSystemStatus,
	}
}
