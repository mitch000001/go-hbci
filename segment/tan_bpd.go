package segment

import (
	"github.com/mitch000001/go-hbci/element"
)

const TanBankParameterID = "HITANS"

type TanBankParameter interface {
	BankSegment
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment TanBankParameterSegment -segment_interface TanBankParameter -segment_versions="TanBankParameterV6:6:Segment"

type TanBankParameterSegment struct {
	TanBankParameter
}

// TanBankParameterV6
//
// Zwei-Schritt-TAN-Einreichung, Parameter
type TanBankParameterV6 struct {
	Segment                     `yaml:"-"`
	MaxJobs                     *element.NumberDataElement             `yaml:"MaxJobs"`
	MinSignatures               *element.NumberDataElement             `yaml:"MinSignatures"`
	SecurityClass               *element.CodeDataElement               `yaml:"SecurityClass"`
	Tan2StepSubmissionParameter *element.Tan2StepSubmissionParameterV6 `yaml:"Tan2StepSubmissionParameter"`
}

func (t *TanBankParameterV6) Version() int         { return 6 }
func (t *TanBankParameterV6) ID() string           { return TanBankParameterID }
func (t *TanBankParameterV6) referencedId() string { return ProcessingPreparationID }
func (t *TanBankParameterV6) sender() string       { return senderBank }

func (t *TanBankParameterV6) elements() []element.DataElement {
	return []element.DataElement{
		t.MaxJobs,
		t.MinSignatures,
		t.SecurityClass,
		t.Tan2StepSubmissionParameter,
	}
}

