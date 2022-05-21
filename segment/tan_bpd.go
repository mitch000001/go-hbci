package segment

import (
	"github.com/mitch000001/go-hbci/element"
	"gopkg.in/yaml.v3"
)

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
	Segment
	MaxJobs                     *element.NumberDataElement
	MinSignatures               *element.NumberDataElement
	SecurityClass               *element.CodeDataElement
	Tan2StepSubmissionParameter *element.Tan2StepSubmissionParameterV6
}

func (t *TanBankParameterV6) Version() int         { return 6 }
func (t *TanBankParameterV6) ID() string           { return "HITANS" }
func (t *TanBankParameterV6) referencedId() string { return "HKVVB" }
func (t *TanBankParameterV6) sender() string       { return senderBank }

func (t *TanBankParameterV6) elements() []element.DataElement {
	return []element.DataElement{
		t.MaxJobs,
		t.MinSignatures,
		t.SecurityClass,
		t.Tan2StepSubmissionParameter,
	}
}

func (t *TanBankParameterV6) MarshalYAML() (interface{}, error) {
	return map[string]yaml.Marshaler{
		"MaxJobs":                     t.MaxJobs,
		"MinSignatures":               t.MinSignatures,
		"SecurityClass":               t.SecurityClass,
		"Tan2StepSubmissionParameter": t.Tan2StepSubmissionParameter,
	}, nil
}
