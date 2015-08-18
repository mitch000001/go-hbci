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

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment CommonBankParameterSegment

type CommonBankParameterSegment struct {
	commonBankParameterSegment
}

type commonBankParameterSegment interface {
	Segment
	BankParameterData() domain.BankParameterData
}

type CommonBankParameterV2 struct {
	Segment
	BPDVersion               *element.NumberDataElement
	BankID                   *element.BankIdentificationDataElement
	BankName                 *element.AlphaNumericDataElement
	BusinessTransactionCount *element.NumberDataElement
	SupportedLanguages       *element.SupportedLanguagesDataElement
	SupportedHBCIVersions    *element.SupportedHBCIVersionsDataElement
	MaxMessageSize           *element.NumberDataElement
}

func (c *CommonBankParameterV2) Version() int         { return 2 }
func (c *CommonBankParameterV2) ID() string           { return "HIBPA" }
func (c *CommonBankParameterV2) referencedId() string { return "HKVVB" }
func (c *CommonBankParameterV2) sender() string       { return senderBank }

func (c *CommonBankParameterV2) elements() []element.DataElement {
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

func (c *CommonBankParameterV2) UnmarshalHBCI(value []byte) error {
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

func (c *CommonBankParameterV2) BankParameterData() domain.BankParameterData {
	return domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
		MaxMessageSize:            c.MaxMessageSize.Val(),
	}
}

type CommonBankParameterV3 struct {
	Segment
	BPDVersion               *element.NumberDataElement
	BankID                   *element.BankIdentificationDataElement
	BankName                 *element.AlphaNumericDataElement
	BusinessTransactionCount *element.NumberDataElement
	SupportedLanguages       *element.SupportedLanguagesDataElement
	SupportedHBCIVersions    *element.SupportedHBCIVersionsDataElement
	MaxMessageSize           *element.NumberDataElement
	MinTimeoutValue          *element.NumberDataElement
	MaxTimeoutValue          *element.NumberDataElement
}

func (c *CommonBankParameterV3) Version() int         { return 3 }
func (c *CommonBankParameterV3) ID() string           { return "HIBPA" }
func (c *CommonBankParameterV3) referencedId() string { return "HKVVB" }
func (c *CommonBankParameterV3) sender() string       { return senderBank }

func (c *CommonBankParameterV3) elements() []element.DataElement {
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

func (c *CommonBankParameterV3) BankParameterData() domain.BankParameterData {
	return domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
		MaxMessageSize:            c.MaxMessageSize.Val(),
		MinTimeout:                c.MinTimeoutValue.Val(),
		MaxTimeout:                c.MaxTimeoutValue.Val(),
	}
}
