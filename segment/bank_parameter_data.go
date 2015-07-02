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

func (c *CommonBankParameterSegment) init() {
	*c.BPDVersion = *new(element.NumberDataElement)
	*c.BankID = *new(element.BankIdentificationDataElement)
	*c.BankName = *new(element.AlphaNumericDataElement)
	*c.BusinessTransactionCount = *new(element.NumberDataElement)
	*c.SupportedLanguages = *new(element.SupportedLanguagesDataElement)
	*c.SupportedHBCIVersions = *new(element.SupportedHBCIVersionsDataElement)
	*c.MaxMessageSize = *new(element.NumberDataElement)
}
func (c *CommonBankParameterSegment) version() int         { return 2 }
func (c *CommonBankParameterSegment) id() string           { return "HIBPA" }
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

type SecurityMethodSegment struct {
	Segment
	MixAllowed       *element.BooleanDataElement
	SupportedMethods *element.SupportedSecurityMethodDataElement
}

func (s *SecurityMethodSegment) init() {
	*s.MixAllowed = *new(element.BooleanDataElement)
	*s.SupportedMethods = *new(element.SupportedSecurityMethodDataElement)
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

func (c *CompressionMethodSegment) init() {
	*c.SupportedCompressionMethods = *new(element.SupportedCompressionMethodsDataElement)
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
