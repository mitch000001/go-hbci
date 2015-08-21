package test_files

import "github.com/mitch000001/go-hbci/element"

type VersionedTestSegmentCustomInterface struct {
	versionedTestSegmentCustomInterface
}

type versionedTestSegmentCustomInterface interface {
	BankSegment
}

type VersionedTestSegmentCustomInterfaceV1 struct {
	Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (v *VersionedTestSegmentCustomInterfaceV1) elements() []element.DataElement {
	return []element.DataElement{
		v.Abc,
		v.Def,
	}
}
