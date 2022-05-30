package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const TanBankParameterID = "HITANS"

type TanBankParameter interface {
	BankSegment
	TanProcessParameters() []domain.TanProcessParameter
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment TanBankParameterSegment -segment_interface TanBankParameter -segment_versions="TanBankParameterV6:6:Segment,TanBankParameterV7:7:Segment"

type TanBankParameterSegment struct {
	TanBankParameter
}

// TanBankParameterV6
//
// Zwei-Schritt-TAN-Einreichung, Parameter
type TanBankParameterV6 struct {
	Segment       `yaml:"-"`
	MaxJobs       *element.NumberDataElement `yaml:"MaxJobs"`
	MinSignatures *element.NumberDataElement `yaml:"MinSignatures"`
	// Codierung:
	// 0: kein Sicherheitsdienst erforderlich
	// 1: Authentikation
	// 2: Non-Repudiation mit fortgeschrittener elektronischer Signatur gemäß §2, SigG
	// 3: Non-Repudiation mit fortgeschrittener elektronischer Signatur gemäß §2, SigG und zwingender Zertifikatsprüfung
	// 4: Non-Repudiation mit qualifizierter elektronischer Signatur gemäß §2, SigG und zwingender Zertifikatsprüfung
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

func (t *TanBankParameterV6) TanProcessParameters() []domain.TanProcessParameter {
	var processParams []domain.TanProcessParameter
	for _, de := range t.Tan2StepSubmissionParameter.ProcessParameters.GroupDataElements() {
		pp, ok := de.(*element.Tan2StepSubmissionProcessParameterV6)
		if !ok {
			continue
		}
		processParams = append(processParams, domain.TanProcessParameter{
			Steps:                 2,
			SecurityFunction:      pp.SecurityFunction.Val(),
			TanProcess:            pp.TanProcess.Val(),
			TanProcessTechnicalID: pp.TechnicalIDTanProcess.Val(),
			TanProcessName:        pp.ZKATanProcess.Val(),
			TanProcessVersion:     pp.ZKATanProcessVersion.Val(),
			TwoStepProcessName:    pp.TwoStepProcessName.Val(),
		})
	}
	return processParams
}

// TanBankParameterV7
//
// Zwei-Schritt-TAN-Einreichung, Parameter
type TanBankParameterV7 struct {
	Segment                     `yaml:"-"`
	MaxJobs                     *element.NumberDataElement             `yaml:"MaxJobs"`
	MinSignatures               *element.NumberDataElement             `yaml:"MinSignatures"`
	SecurityClass               *element.CodeDataElement               `yaml:"SecurityClass"`
	Tan2StepSubmissionParameter *element.Tan2StepSubmissionParameterV7 `yaml:"Tan2StepSubmissionParameter"`
}

func (t *TanBankParameterV7) Version() int         { return 7 }
func (t *TanBankParameterV7) ID() string           { return TanBankParameterID }
func (t *TanBankParameterV7) referencedId() string { return ProcessingPreparationID }
func (t *TanBankParameterV7) sender() string       { return senderBank }

func (t *TanBankParameterV7) elements() []element.DataElement {
	return []element.DataElement{
		t.MaxJobs,
		t.MinSignatures,
		t.SecurityClass,
		t.Tan2StepSubmissionParameter,
	}
}

func (t *TanBankParameterV7) TanProcessParameters() []domain.TanProcessParameter {
	var processParams []domain.TanProcessParameter
	for _, de := range t.Tan2StepSubmissionParameter.ProcessParameters.GroupDataElements() {
		pp, ok := de.(*element.Tan2StepSubmissionProcessParameterV7)
		if !ok {
			continue
		}
		processParams = append(processParams, domain.TanProcessParameter{
			Steps:                 2,
			SecurityFunction:      pp.SecurityFunction.Val(),
			TanProcess:            pp.TanProcess.Val(),
			TanProcessTechnicalID: pp.TechnicalIDTanProcess.Val(),
			TanProcessName:        pp.DKTanProcess.Val(),
			TanProcessVersion:     pp.DKTanProcessVersion.Val(),
			TwoStepProcessName:    pp.TwoStepProcessName.Val(),
		})
	}
	return processParams
}
