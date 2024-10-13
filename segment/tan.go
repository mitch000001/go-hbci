package segment

import (
	"fmt"
	"sort"

	"github.com/mitch000001/go-hbci/domain"
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
	SetTANProcess(string)
	SetAnotherTanFollows(bool)
	SetTANParams(domain.TanParams)
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
func (t *TanRequestSegmentV1) SetTANProcess(process string) {
	t.TANProcess = element.NewAlphaNumeric(process, 1)
}
func (t *TanRequestSegmentV1) SetAnotherTanFollows(another bool) {
	t.AnotherTanFollows = element.NewBoolean(another)
}
func (t *TanRequestSegmentV1) SetTANParams(params domain.TanParams) {
	t.JobReference = element.NewAlphaNumeric(params.JobReference, 35)
	if params.JobHash != nil {
		t.JobHash = element.NewBinary(params.JobHash, 256)
	}
}

func NewTanProcess4RequestSegmentV6(referencingSegmentID string) *TanRequestSegment {
	t := &TanRequestSegmentV6{
		TANProcess:           element.NewCode("4", 1, []string{"1", "2", "3", "4"}),
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
	TANProcess           *element.CodeDataElement
	ReferencingSegmentID *element.AlphaNumericDataElement
	AccountConnection    *element.InternationalAccountConnectionDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	TanListNumber        *element.AlphaNumericDataElement
	AnotherTanFollows    *element.BooleanDataElement
	CancelJob            *element.BooleanDataElement
	SMSDebitAccount      *element.BooleanDataElement
	ChallengeClass       *element.NumberDataElement
	ChallengeClassParams *element.AlphaNumericDataElement
	TANMediumDescription *element.AlphaNumericDataElement
	ResponseHHD_UC       *element.AlphaNumericDataElement
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
		t.AccountConnection,
		t.JobHash,
		t.JobReference,
		t.TanListNumber,
		t.AnotherTanFollows,
		t.CancelJob,
		t.SMSDebitAccount,
		t.ChallengeClass,
		t.ChallengeClassParams,
		t.TANMediumDescription,
		t.ResponseHHD_UC,
	}
}

func (t *TanRequestSegmentV6) SetTANProcess(process string) {
	t.TANProcess = element.NewCode(process, 1, []string{"1", "2", "3", "4"})
}
func (t *TanRequestSegmentV6) SetAnotherTanFollows(another bool) {
	t.AnotherTanFollows = element.NewBoolean(another)
}
func (t *TanRequestSegmentV6) SetTANParams(params domain.TanParams) {
	t.JobReference = element.NewAlphaNumeric(params.JobReference, 35)
	if params.JobHash != nil {
		t.JobHash = element.NewBinary(params.JobHash, 256)
	}
}

func NewTanProcess4RequestSegmentV7(referencingSegmentID string) *TanRequestSegment {
	t := &TanRequestSegmentV7{
		TANProcess:           element.NewCode("4", 1, []string{"1", "2", "3", "4", "S"}),
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
	TANProcess           *element.CodeDataElement
	ReferencingSegmentID *element.AlphaNumericDataElement
	AccountConnection    *element.InternationalAccountConnectionDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	TanListNumber        *element.AlphaNumericDataElement
	AnotherTanFollows    *element.BooleanDataElement
	CancelJob            *element.BooleanDataElement
	SMSDebitAccount      *element.BooleanDataElement
	ChallengeClass       *element.NumberDataElement
	ChallengeClassParams *element.AlphaNumericDataElement
	TANMediumDescription *element.AlphaNumericDataElement
	ResponseHHD_UC       *element.AlphaNumericDataElement
}

func (t *TanRequestSegmentV7) Version() int         { return 7 }
func (t *TanRequestSegmentV7) ID() string           { return "HKTAN" }
func (t *TanRequestSegmentV7) referencedId() string { return "" }
func (t *TanRequestSegmentV7) sender() string       { return senderUser }

func (t *TanRequestSegmentV7) elements() []element.DataElement {
	return []element.DataElement{
		t.TANProcess,
		t.ReferencingSegmentID,
		t.AccountConnection,
		t.JobHash,
		t.JobReference,
		t.TanListNumber,
		t.AnotherTanFollows,
		t.CancelJob,
		t.SMSDebitAccount,
		t.ChallengeClass,
		t.ChallengeClassParams,
		t.TANMediumDescription,
		t.ResponseHHD_UC,
	}
}

func (t *TanRequestSegmentV7) SetTANProcess(process string) {
	t.TANProcess = element.NewCode(process, 1, []string{"1", "2", "3", "4", "S"})
}
func (t *TanRequestSegmentV7) SetAnotherTanFollows(another bool) {
	t.AnotherTanFollows = element.NewBoolean(another)
}
func (t *TanRequestSegmentV7) SetTANParams(params domain.TanParams) {
	t.JobReference = element.NewAlphaNumeric(params.JobReference, 35)
	if params.JobHash != nil {
		t.JobHash = element.NewBinary(params.JobHash, 256)
	}
}

type TanResponse interface {
	BankSegment
	TanParams() domain.TanParams
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
	ChallengeExpiryDate  *element.TanChallengeExpiryDate
	TANMediumDescription *element.AlphaNumericDataElement
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
		t.Challenge,
		t.ChallengeHHD_UC,
		t.ChallengeExpiryDate,
		t.TANMediumDescription,
	}
}

func (t *TanResponseSegmentV6) TanParams() domain.TanParams {
	params := domain.TanParams{}
	params.TANProcess = t.TANProcess.Val()
	if t.JobHash != nil {
		params.JobHash = t.JobHash.Val()
	}
	if t.JobReference != nil {
		params.JobReference = t.JobReference.Val()
	}
	if t.Challenge != nil {
		params.Challenge = t.Challenge.Val()
	}
	if t.ChallengeHHD_UC != nil {
		params.ChallengeHHD_UC = t.ChallengeHHD_UC.Val()
	}
	if t.TANMediumDescription != nil {
		params.TANMediumDescription = t.TANMediumDescription.Val()
	}
	if t.ChallengeExpiryDate != nil {
		params.ChallengeExpiryDate = t.ChallengeExpiryDate.Val()
	}
	return params
}

type TanResponseSegmentV7 struct {
	Segment
	TANProcess           *element.AlphaNumericDataElement
	JobHash              *element.BinaryDataElement
	JobReference         *element.AlphaNumericDataElement
	Challenge            *element.AlphaNumericDataElement
	ChallengeHHD_UC      *element.BinaryDataElement
	ChallengeExpiryDate  *element.TanChallengeExpiryDate
	TANMediumDescription *element.AlphaNumericDataElement
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
		t.Challenge,
		t.ChallengeHHD_UC,
		t.ChallengeExpiryDate,
		t.TANMediumDescription,
	}
}

func (t *TanResponseSegmentV7) TanParams() domain.TanParams {
	params := domain.TanParams{}
	params.TANProcess = t.TANProcess.Val()
	if t.JobHash != nil {
		params.JobHash = t.JobHash.Val()
	}
	if t.JobReference != nil {
		params.JobReference = t.JobReference.Val()
	}
	if t.Challenge != nil {
		params.Challenge = t.Challenge.Val()
	}
	if t.ChallengeHHD_UC != nil {
		params.ChallengeHHD_UC = t.ChallengeHHD_UC.Val()
	}
	if t.TANMediumDescription != nil {
		params.TANMediumDescription = t.TANMediumDescription.Val()
	}
	if t.ChallengeExpiryDate != nil {
		params.ChallengeExpiryDate = t.ChallengeExpiryDate.Val()
	}
	return params
}
