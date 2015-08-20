package test_files

import "github.com/mitch000001/go-hbci/element"

type TestSegment struct {
	Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
	Xyz *element.NumberDataElement
}

func (t *TestSegment) elements() []element.DataElement {
	return []element.DataElement{
		t.Abc,
		t.Xyz,
	}
}
