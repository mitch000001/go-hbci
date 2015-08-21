package test_files

import "github.com/mitch000001/go-hbci/element"

type VersionedTestSegment struct {
	BankSegment
}

type VersionedTestSegmentV1 struct {
	Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (v *VersionedTestSegmentV1) elements() []element.DataElement {
	return []element.DataElement{
		v.Abc,
		v.Def,
	}
}
