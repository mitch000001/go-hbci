package test_files

import "github.com/mitch000001/go-hbci/element"

type TestSegmentUnknownElement struct {
	Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (t *TestSegmentUnknownElement) elements() []element.DataElement {
	return []element.DataElement{
		t.Abc,
		&element.NumberDataElement{},
	}
}
