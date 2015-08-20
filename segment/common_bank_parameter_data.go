package segment

import (
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

func (c *CommonBankParameterSegment) BankParameterData() domain.BankParameterData {
	return domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
	}
}
