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

type SecurityMethodSegment struct {
	Segment
	MixAllowed       *element.BooleanDataElement
	SupportedMethods *element.SupportedSecurityMethodDataElement
}

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

func (c *CompressionMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.SupportedCompressionMethods,
	}
}
