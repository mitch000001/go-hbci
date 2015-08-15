package segment

import "github.com/mitch000001/go-hbci/element"

func NewTanRequestProcess2(jobReference string, anotherTANFollows bool) *TanRequestSegment {
	t := &TanRequestSegment{
		TANProcess:        element.NewAlphaNumeric("2", 1),
		JobReference:      element.NewAlphaNumeric(jobReference, 35),
		AnotherTanFollows: element.NewBoolean(anotherTANFollows),
	}
	t.Segment = NewBasicSegment(1, t)
	return t
}

func NewTanRequestProcess4() *TanRequestSegment {
	t := &TanRequestSegment{
		TANProcess: element.NewAlphaNumeric("4", 1),
	}
	t.Segment = NewBasicSegment(1, t)
	return t
}

type TanRequestSegment struct {
	Segment
	TANProcess        *element.AlphaNumericDataElement
	JobHash           *element.BinaryDataElement
	JobReference      *element.AlphaNumericDataElement
	TanListNumber     *element.AlphaNumericDataElement
	AnotherTanFollows *element.BooleanDataElement
	TANInformation    *element.AlphaNumericDataElement
}

func (t *TanRequestSegment) Version() int         { return 1 }
func (t *TanRequestSegment) ID() string           { return "HKTAN" }
func (t *TanRequestSegment) referencedId() string { return "" }
func (t *TanRequestSegment) sender() string       { return senderUser }

func (t *TanRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.JobHash,
		t.JobReference,
		t.TanListNumber,
		t.AnotherTanFollows,
		t.TANInformation,
	}
}
