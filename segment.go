package hbci

import "strings"

type Segment interface {
	Header() *SegmentHeader
	DataElements() []DataElement
	String() string
}

type segment interface {
	elements() []DataElement
}

func NewBasicSegment(id string, number int, version int, seg segment) *basicSegment {
	header := NewSegmentHeader(id, number, version)
	return NewBasicSegmentWithHeader(header, seg)
}

func NewBasicSegmentWithHeader(header *SegmentHeader, seg segment) *basicSegment {
	return &basicSegment{header: header, segment: seg}
}

type basicSegment struct {
	segment segment
	header  *SegmentHeader
}

func (s *basicSegment) String() string {
	elementStrings := make([]string, len(s.segment.elements())+1)
	elementStrings[0] = s.header.String()
	for i, de := range s.segment.elements() {
		elementStrings[i+1] = de.String()
	}
	return strings.Join(elementStrings, "+") + "'"
}

func (s *basicSegment) DataElements() []DataElement {
	var dataElements []DataElement
	dataElements = append(dataElements, s.header)
	dataElements = append(dataElements, s.segment.elements()...)
	return dataElements
}

func (s *basicSegment) Header() *SegmentHeader {
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

func NewIdentificationSegment(countryCode int, bankId string, clientId string, clientSystemId string, systemIdRequired bool) *IdentificationSegment {
	var clientSystemStatus *NumberDataElement
	if systemIdRequired {
		clientSystemStatus = NewNumberDataElement(1, 1)
	} else {
		clientSystemStatus = NewNumberDataElement(0, 1)
	}
	id := &IdentificationSegment{
		BankId:             NewBankIndentificationDataElementWithBankId(countryCode, bankId),
		ClientId:           NewIdentificationDataElement(clientId),
		ClientSystemId:     NewIdentificationDataElement(clientSystemId),
		ClientSystemStatus: clientSystemStatus,
	}
	id.basicSegment = NewBasicSegment("HKIDN", 3, 2, id)
	return id
}

type IdentificationSegment struct {
	*basicSegment
	BankId             *BankIdentificationDataElement
	ClientId           *IdentificationDataElement
	ClientSystemId     *IdentificationDataElement
	ClientSystemStatus *NumberDataElement
}

func (i *IdentificationSegment) elements() []DataElement {
	return []DataElement{
		i.BankId,
		i.ClientId,
		i.ClientSystemId,
		i.ClientSystemStatus,
	}
}

func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	header := NewSegmentHeader(id, number, version)
	header.Ref = NewNumberDataElement(reference, 3)
	return header
}

func NewSegmentHeader(id string, number, version int) *SegmentHeader {
	header := &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
	}
	header.elementGroup = NewDataElementGroup(SegmentHeaderDEG, 4, header)
	return header
}

type SegmentHeader struct {
	*elementGroup
	ID      *AlphaNumericDataElement
	Number  *NumberDataElement
	Version *NumberDataElement
	Ref     *NumberDataElement
}

func (s *SegmentHeader) SetNumber(number int) {
	s.Number = NewNumberDataElement(number, 3)
}

func (s *SegmentHeader) SetReference(ref int) {
	s.Ref = NewNumberDataElement(ref, 3)
}

func (s *SegmentHeader) IsValid() bool {
	if s.ID == nil || s.Number == nil || s.Version == nil {
		return false
	} else {
		return s.elementGroup.IsValid()
	}
}

func (s *SegmentHeader) Value() interface{} {
	return s
}

func (s *SegmentHeader) GroupDataElements() []DataElement {
	return []DataElement{
		s.ID,
		s.Number,
		s.Version,
		s.Ref,
	}
}
