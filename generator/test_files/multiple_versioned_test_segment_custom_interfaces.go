package test_files

import "github.com/mitch000001/go-hbci/element"

type MultipleVersionedTestSegmentCustomInterfaces struct {
	BankSegment
}

type MultipleVersionedTestSegmentCustomInterfacesV1 struct {
	versionInterface1
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (m *MultipleVersionedTestSegmentCustomInterfacesV1) elements() []element.DataElement {
	return []element.DataElement{
		m.Abc,
		m.Def,
	}
}

type MultipleVersionedTestSegmentCustomInterfacesV2 struct {
	versionInterface2
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (m *MultipleVersionedTestSegmentCustomInterfacesV2) elements() []element.DataElement {
	return []element.DataElement{
		m.Abc,
		m.Def,
	}
}
