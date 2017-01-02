package segment

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

type BusinessTransactionParamsSegment struct {
	Segment
	id            string
	version       int
	MaxJobs       *element.NumberDataElement
	MinSignatures *element.NumberDataElement
	Params        element.DataElementGroup
}

func (b *BusinessTransactionParamsSegment) Version() int         { return b.version }
func (b *BusinessTransactionParamsSegment) ID() string           { return b.id }
func (b *BusinessTransactionParamsSegment) referencedId() string { return "HKVVB" }
func (b *BusinessTransactionParamsSegment) sender() string       { return senderBank }

func (b *BusinessTransactionParamsSegment) elements() []element.DataElement {
	return []element.DataElement{
		b.MaxJobs,
		b.MinSignatures,
		b.Params,
	}
}

func (b *BusinessTransactionParamsSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	seg, err := SegmentFromHeaderBytes(elements[0], b)
	if err != nil {
		return err
	}
	b.Segment = seg
	if len(elements) < 4 {
		return fmt.Errorf("%T: Malformed marshaled value", b)
	}
	maxJobs, err := strconv.Atoi(charset.ToUTF8(elements[1]))
	if err != nil {
		return fmt.Errorf("%T: Malformed max jobs: %v", b, err)
	}
	b.MaxJobs = element.NewNumber(maxJobs, 4)
	minSignatures, err := strconv.Atoi(charset.ToUTF8(elements[2]))
	if err != nil {
		return fmt.Errorf("%T: Malformed min signatures: %v", b, err)
	}
	b.MinSignatures = element.NewNumber(minSignatures, 2)
	return nil
}

type PinTanBusinessTransactionParams interface {
	BankSegment
	PinTanBusinessTransactions() []domain.PinTanBusinessTransaction
}

type PinTanBusinessTransactionParamsSegment struct {
	*BusinessTransactionParamsSegment
}

func (p *PinTanBusinessTransactionParamsSegment) Version() int { return 1 }
func (p *PinTanBusinessTransactionParamsSegment) ID() string   { return "DIPINS" }

func (p *PinTanBusinessTransactionParamsSegment) UnmarshalHBCI(value []byte) error {
	businessTransactionSegment := &BusinessTransactionParamsSegment{}
	err := businessTransactionSegment.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	p.BusinessTransactionParamsSegment = businessTransactionSegment
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	pinTanParams := &element.PinTanBusinessTransactionParameters{}
	err = pinTanParams.UnmarshalHBCI(elements[3])
	if err != nil {
		return err
	}
	p.BusinessTransactionParamsSegment.Params = pinTanParams
	return nil
}

func (p *PinTanBusinessTransactionParamsSegment) PinTanBusinessTransactions() []domain.PinTanBusinessTransaction {
	var transactions []domain.PinTanBusinessTransaction
	for _, transactionDe := range p.Params.GroupDataElements() {
		transactions = append(transactions, transactionDe.(*element.PinTanBusinessTransactionParameter).Val())
	}
	return transactions
}
