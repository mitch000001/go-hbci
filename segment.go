package hbci

import "strings"

type Segment interface {
	DataElements() []DataElement
}

func NewSegment(header *SegmentHeader, seg Segment) *segment {
	return &segment{Header: header, segment: seg}
}

type segment struct {
	segment Segment
	Header  *SegmentHeader
}

func (s *segment) String() string {
	elementStrings := make([]string, len(s.segment.DataElements())+1)
	elementStrings[0] = s.Header.String()
	for i, de := range s.segment.DataElements() {
		elementStrings[i+1] = de.String()
	}
	return strings.Join(elementStrings, "+") + "'"
}

func (s *segment) SetNumber(number int) {
	s.Header.SetNumber(number)
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
	header := NewSegmentHeader("HKIDN", 3, 2)
	id.segment = NewSegment(header, id)
	return id
}

type IdentificationSegment struct {
	*segment
	BankId             *BankIdentificationDataElement
	ClientId           *IdentificationDataElement
	ClientSystemId     *IdentificationDataElement
	ClientSystemStatus *NumberDataElement
}

func (i *IdentificationSegment) DataElements() []DataElement {
	return []DataElement{
		i.BankId,
		i.ClientId,
		i.ClientSystemId,
		i.ClientSystemStatus,
	}
}

func NewReferencingSegmentHeader(id string, number, version, reference int) *SegmentHeader {
	return &SegmentHeader{
		ID:      NewAlphaNumericDataElement(id, 6),
		Number:  NewNumberDataElement(number, 3),
		Version: NewNumberDataElement(version, 3),
		Ref:     NewNumberDataElement(reference, 3),
	}
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

func (s *SegmentHeader) Valid() bool {
	if s.ID == nil || s.Number == nil || s.Version == nil {
		return false
	} else {
		return s.elementGroup.Valid()
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
