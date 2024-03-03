package segment

import (
	"fmt"
	"sort"

	"github.com/mitch000001/go-hbci/element"
)

type tanProcess4Constructor func(referencingSegmentID string) *TanRequestSegment

var tanProcess4RequestSegmentConstructors = map[int](tanProcess4Constructor){
	7: NewTanProcess4RequestSegmentV7,
	6: NewTanProcess4RequestSegmentV6,
	1: NewTanProcess4RequestSegmentV1,
}

func TanProcess4RequestBuilder(versions []int) (tanProcess4Constructor, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		builder, ok := tanProcess4RequestSegmentConstructors[version]
		if ok {
			return builder, nil
		}
	}
	return nil, fmt.Errorf("unsupported versions %v", versions)
}

type TanRequestSegment struct {
	tanRequestSegment
}

type tanRequestSegment interface {
	ClientSegment
}

func NewTanRequestProcess2(jobReference string, anotherTANFollows bool) *TanRequestSegmentV1 {
	t := &TanRequestSegmentV1{
		TANProcess:        element.NewAlphaNumeric("2", 1),
		JobReference:      element.NewAlphaNumeric(jobReference, 35),
		AnotherTanFollows: element.NewBoolean(anotherTANFollows),
	}
	t.ClientSegment = NewBasicSegment(1, t)
	return t
}

func NewTanProcess4RequestSegmentV1(referencingSegmentID string) *TanRequestSegment {
	t := &TanRequestSegmentV1{
		TANProcess: element.NewAlphaNumeric("4", 1),
	}
	t.ClientSegment = NewBasicSegment(1, t)

	segment := &TanRequestSegment{
		tanRequestSegment: t,
	}
	return segment
}

type TanRequestSegmentV1 struct {
	ClientSegment
	TANProcess        *element.AlphaNumericDataElement
	JobHash           *element.BinaryDataElement
	JobReference      *element.AlphaNumericDataElement
	TanListNumber     *element.AlphaNumericDataElement
	AnotherTanFollows *element.BooleanDataElement
	TANInformation    *element.AlphaNumericDataElement
}

func (t *TanRequestSegmentV1) Version() int         { return 1 }
func (t *TanRequestSegmentV1) ID() string           { return "HKTAN" }
func (t *TanRequestSegmentV1) referencedId() string { return "" }
func (t *TanRequestSegmentV1) sender() string       { return senderUser }

func (t *TanRequestSegmentV1) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.JobHash,
		t.JobReference,
		t.TanListNumber,
		t.AnotherTanFollows,
		t.TANInformation,
	}
}

func NewTanProcess4RequestSegmentV6(referencingSegmentID string) *TanRequestSegment {
	t := &TanRequestSegmentV6{
		TANProcess:           element.NewAlphaNumeric("4", 1),
		ReferencingSegmentID: element.NewAlphaNumeric(referencingSegmentID, 6),
	}
	t.ClientSegment = NewBasicSegment(1, t)

	segment := &TanRequestSegment{
		tanRequestSegment: t,
	}
	return segment
}

type TanRequestSegmentV6 struct {
	ClientSegment
	TANProcess           *element.AlphaNumericDataElement
	ReferencingSegmentID *element.AlphaNumericDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	TanListNumber        *element.AlphaNumericDataElement
	AnotherTanFollows    *element.BooleanDataElement
	TANInformation       *element.AlphaNumericDataElement
}

func (t *TanRequestSegmentV6) Version() int         { return 6 }
func (t *TanRequestSegmentV6) ID() string           { return "HKTAN" }
func (t *TanRequestSegmentV6) referencedId() string { return "" }
func (t *TanRequestSegmentV6) sender() string       { return senderUser }

func (t *TanRequestSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.ReferencingSegmentID,
		t.JobHash,
		t.JobReference,
		t.TanListNumber,
		t.AnotherTanFollows,
		t.TANInformation,
	}
}

func NewTanProcess4RequestSegmentV7(referencingSegmentID string) *TanRequestSegment {
	t := &TanRequestSegmentV7{
		TANProcess:           element.NewAlphaNumeric("4", 1),
		ReferencingSegmentID: element.NewAlphaNumeric(referencingSegmentID, 6),
	}
	t.ClientSegment = NewBasicSegment(1, t)

	segment := &TanRequestSegment{
		tanRequestSegment: t,
	}
	return segment
}

type TanRequestSegmentV7 struct {
	ClientSegment
	TANProcess           *element.AlphaNumericDataElement
	ReferencingSegmentID *element.AlphaNumericDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	TanListNumber        *element.AlphaNumericDataElement
	AnotherTanFollows    *element.BooleanDataElement
	TANInformation       *element.AlphaNumericDataElement
}

func (t *TanRequestSegmentV7) Version() int         { return 7 }
func (t *TanRequestSegmentV7) ID() string           { return "HKTAN" }
func (t *TanRequestSegmentV7) referencedId() string { return "" }
func (t *TanRequestSegmentV7) sender() string       { return senderUser }

func (t *TanRequestSegmentV7) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.ReferencingSegmentID,
		t.JobHash,
		t.JobReference,
		t.TanListNumber,
		t.AnotherTanFollows,
		t.TANInformation,
	}
}

type TanResponse interface {
	BankSegment
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment TanResponseSegment -segment_interface TanResponse -segment_versions="TanResponseSegmentV6:6:Segment,TanResponseSegmentV7:7:Segment"

type TanResponseSegment struct {
	TanResponse
}

type TanResponseSegmentV6 struct {
	Segment
	TANProcess           *element.AlphaNumericDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	Challenge            *element.AlphaNumericDataElement
	ChallengeHHD_UC      *element.BinaryDataElement
	TANMediumDescription *element.AlphaNumericDataElement
	ChallengeExpiryDate  *element.TanChallengeExpiryDate
}

func (t *TanResponseSegmentV6) Version() int         { return 6 }
func (t *TanResponseSegmentV6) ID() string           { return "HITAN" }
func (t *TanResponseSegmentV6) referencedId() string { return "" }
func (t *TanResponseSegmentV6) sender() string       { return senderBank }

func (t *TanResponseSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.JobHash,
		t.JobReference,
	}
}

type TanResponseSegmentV7 struct {
	Segment
	TANProcess           *element.AlphaNumericDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	Challenge            *element.AlphaNumericDataElement
	ChallengeHHD_UC      *element.BinaryDataElement
	TANMediumDescription *element.AlphaNumericDataElement
	ChallengeExpiryDate  *element.TanChallengeExpiryDate
}

func (t *TanResponseSegmentV7) Version() int         { return 7 }
func (t *TanResponseSegmentV7) ID() string           { return "HITAN" }
func (t *TanResponseSegmentV7) referencedId() string { return "" }
func (t *TanResponseSegmentV7) sender() string       { return senderBank }

func (t *TanResponseSegmentV7) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.JobHash,
		t.JobReference,
	}
}
