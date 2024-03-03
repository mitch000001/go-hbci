package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const PinTanBankParameterID = "HIPINS"

type PinTanBankParameter interface {
	BankSegment
	PinTanBusinessTransactions() []domain.PinTanBusinessTransaction
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment PinTanBankParameterSegment -segment_interface PinTanBankParameter -segment_versions="PinTanBankParameterV1:1:Segment"

type PinTanBankParameterSegment struct {
	PinTanBankParameter
}

// PinTanBankParameterV1
//
// PIN/TAN-spezifische Informationen
type PinTanBankParameterV1 struct {
	Segment       `yaml:"-"`
	MaxJobs       *element.NumberDataElement `yaml:"MaxJobs"`
	MinSignatures *element.NumberDataElement `yaml:"MinSignatures"`
	// TODO/FIXME: find out which parameters are here actually
	XXX_Unknown          *element.NumberDataElement              `yaml:"-"`
	PinTanSpecificParams *element.PinTanSpecificParamDataElement `yaml:"PinTanSpecificParams"`
}

func (t *PinTanBankParameterV1) Version() int         { return 1 }
func (t *PinTanBankParameterV1) ID() string           { return PinTanBankParameterID }
func (t *PinTanBankParameterV1) referencedId() string { return ProcessingPreparationID }
func (t *PinTanBankParameterV1) sender() string       { return senderBank }

func (t *PinTanBankParameterV1) elements() []element.DataElement {
	return []element.DataElement{
		t.MaxJobs,
		t.MinSignatures,
		t.XXX_Unknown,
		t.PinTanSpecificParams,
	}
}

func (t *PinTanBankParameterV1) PinTanBusinessTransactions() []domain.PinTanBusinessTransaction {
	var transactions []domain.PinTanBusinessTransaction
	for _, transactionDe := range t.PinTanSpecificParams.JobSpecificPinTanInformation.GroupDataElements() {
		transactions = append(transactions, transactionDe.(*element.PinTanBusinessTransactionParameter).Val())
	}
	return transactions
}
