package segment

import (
	"github.com/mitch000001/go-hbci/dataelement"
	"github.com/mitch000001/go-hbci/domain"
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
		BPDVersion:               dataelement.NewNumber(bpdVersion, 3),
		BankID:                   dataelement.NewBankIndentification(bankId),
		BankName:                 dataelement.NewAlphaNumeric(bankName, 60),
		BusinessTransactionCount: dataelement.NewNumber(businessTransactionCount, 3),
		SupportedLanguages:       dataelement.NewSupportedLanguages(supportedLanguages...),
		SupportedHBCIVersions:    dataelement.NewSupportedHBCIVersions(supportedHBCIVersions...),
		MaxMessageSize:           dataelement.NewNumber(maxMessageSize, 4),
	}
	header := dataelement.NewReferencingSegmentHeader("HIBPA", 1, 2, HKVVBSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

type CommonBankParameterSegment struct {
	Segment
	BPDVersion               *dataelement.NumberDataElement
	BankID                   *dataelement.BankIdentificationDataElement
	BankName                 *dataelement.AlphaNumericDataElement
	BusinessTransactionCount *dataelement.NumberDataElement
	SupportedLanguages       *dataelement.SupportedLanguagesDataElement
	SupportedHBCIVersions    *dataelement.SupportedHBCIVersionsDataElement
	MaxMessageSize           *dataelement.NumberDataElement
}

func (c *CommonBankParameterSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
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
	MixAllowed       *dataelement.BooleanDataElement
	SupportedMethods *dataelement.SupportedSecurityMethodDataElement
}

func (s *SecurityMethodSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		s.MixAllowed,
		s.SupportedMethods,
	}
}

type CompressionMethodSegment struct {
	Segment
	SupportedCompressionMethods *dataelement.SupportedCompressionMethodsDataElement
}

func (c *CompressionMethodSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		c.SupportedCompressionMethods,
	}
}
