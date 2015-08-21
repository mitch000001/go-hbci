package test_files

import "github.com/mitch000001/go-hbci/element"

type MultipleVersionedTestSegment struct {
	BankSegment
}

type MultipleVersionedTestSegmentV1 struct {
	Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (m *MultipleVersionedTestSegmentV1) elements() []element.DataElement {
	return []element.DataElement{
		m.Abc,
		m.Def,
	}
}

type MultipleVersionedTestSegmentV2 struct {
	Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (m *MultipleVersionedTestSegmentV2) elements() []element.DataElement {
	return []element.DataElement{
		m.Abc,
		m.Def,
	}
}
