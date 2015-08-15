package segment

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

var HKVVBSegmentNumber = -1

func NewCommonBankParameterSegment(
	bpdVersion int,
	bankId domain.BankId,
	bankName string,
	businessTransactionCount int,
	supportedLanguages []int,
	supportedHBCIVersions []int,
	maxMessageSize int) *CommonBankParameterSegment {
	c := &CommonBankParameterSegment{
		BPDVersion:               element.NewNumber(bpdVersion, 3),
		BankID:                   element.NewBankIndentification(bankId),
		BankName:                 element.NewAlphaNumeric(bankName, 60),
		BusinessTransactionCount: element.NewNumber(businessTransactionCount, 3),
		SupportedLanguages:       element.NewSupportedLanguages(supportedLanguages...),
		SupportedHBCIVersions:    element.NewSupportedHBCIVersions(supportedHBCIVersions...),
		MaxMessageSize:           element.NewNumber(maxMessageSize, 4),
	}
	header := element.NewReferencingSegmentHeader("HIBPA", 1, 2, HKVVBSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

type CommonBankParameterSegment struct {
	Segment
	BPDVersion               *element.NumberDataElement
	BankID                   *element.BankIdentificationDataElement
	BankName                 *element.AlphaNumericDataElement
	BusinessTransactionCount *element.NumberDataElement
	SupportedLanguages       *element.SupportedLanguagesDataElement
	SupportedHBCIVersions    *element.SupportedHBCIVersionsDataElement
	MaxMessageSize           *element.NumberDataElement
}

func (c *CommonBankParameterSegment) Version() int         { return 2 }
func (c *CommonBankParameterSegment) ID() string           { return "HIBPA" }
func (c *CommonBankParameterSegment) referencedId() string { return "HKVVB" }
func (c *CommonBankParameterSegment) sender() string       { return senderBank }

func (c *CommonBankParameterSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.BPDVersion,
		c.BankID,
		c.BankName,
		c.BusinessTransactionCount,
		c.SupportedLanguages,
		c.SupportedHBCIVersions,
		c.MaxMessageSize,
	}
}

func (c *CommonBankParameterSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 || len(elements) < 7 {
		return fmt.Errorf("Malformed marshaled value")
	}
	segment, err := SegmentFromHeaderBytes(elements[0], c)
	if err != nil {
		return err
	}
	c.Segment = segment
	version, err := strconv.Atoi(charset.ToUtf8(elements[1]))
	if err != nil {
		return err
	}
	c.BPDVersion = element.NewNumber(version, 3)
	bankId := &element.BankIdentificationDataElement{}
	err = bankId.UnmarshalHBCI(elements[2])
	if err != nil {
		return err
	}
	c.BankID = bankId
	c.BankName = element.NewAlphaNumeric(charset.ToUtf8(elements[3]), 60)
	transactionCount, err := strconv.Atoi(charset.ToUtf8(elements[4]))
	if err != nil {
		return err
	}
	c.BusinessTransactionCount = element.NewNumber(transactionCount, 3)
	languages := &element.SupportedLanguagesDataElement{}
	err = languages.UnmarshalHBCI(elements[5])
	if err != nil {
		return err
	}
	c.SupportedLanguages = languages
	versions := &element.SupportedHBCIVersionsDataElement{}
	err = versions.UnmarshalHBCI(elements[6])
	if err != nil {
		return err
	}
	c.SupportedHBCIVersions = versions
	if len(elements) == 8 {
		maxSize, err := strconv.Atoi(charset.ToUtf8(elements[7]))
		if err != nil {
			return err
		}
		c.MaxMessageSize = element.NewNumber(maxSize, 4)
	}
	return nil
}

func (c *CommonBankParameterSegment) BankParameterData() domain.BankParameterData {
	return domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
	}
}

type SecurityMethodSegment struct {
	Segment
	MixAllowed       *element.BooleanDataElement
	SupportedMethods *element.SupportedSecurityMethodDataElement
}

func (s *SecurityMethodSegment) version() int         { return 2 }
func (s *SecurityMethodSegment) id() string           { return "HISHV" }
func (s *SecurityMethodSegment) referencedId() string { return "HKVVB" }
func (s *SecurityMethodSegment) sender() string       { return senderBank }

func (s *SecurityMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.MixAllowed,
		s.SupportedMethods,
	}
}

type CompressionMethodSegment struct {
	Segment
	SupportedCompressionMethods *element.SupportedCompressionMethodsDataElement
}

func (c *CompressionMethodSegment) version() int         { return 1 }
func (c *CompressionMethodSegment) id() string           { return "HIKPV" }
func (c *CompressionMethodSegment) referencedId() string { return "HKVVB" }
func (c *CompressionMethodSegment) sender() string       { return senderBank }

func (c *CompressionMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.SupportedCompressionMethods,
	}
}

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
	maxJobs, err := strconv.Atoi(charset.ToUtf8(elements[1]))
	if err != nil {
		return fmt.Errorf("%T: Malformed max jobs: %v", b, err)
	}
	b.MaxJobs = element.NewNumber(maxJobs, 4)
	minSignatures, err := strconv.Atoi(charset.ToUtf8(elements[2]))
	if err != nil {
		return fmt.Errorf("%T: Malformed min signatures: %v", b, err)
	}
	b.MinSignatures = element.NewNumber(minSignatures, 2)
	return nil
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
