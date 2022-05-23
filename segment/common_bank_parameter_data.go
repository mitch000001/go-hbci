package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const CommonBankParameterID string = "HIBPA"

type CommonBankParameter interface {
	BankSegment
	BankParameterData() domain.BankParameterData
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment CommonBankParameterSegment -segment_interface commonBankParameterSegment -segment_versions="CommonBankParameterV2:2,CommonBankParameterV3:3"

type CommonBankParameterSegment struct {
	commonBankParameterSegment
}

type commonBankParameterSegment interface {
	BankSegment
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
func (c *CommonBankParameterV2) ID() string           { return CommonBankParameterID }
func (c *CommonBankParameterV2) referencedId() string { return ProcessingPreparationID }
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

func (c *CommonBankParameterV2) BankParameterData() domain.BankParameterData {
	bpd := domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
	}
	if c.MaxMessageSize != nil {
		bpd.MaxMessageSize = c.MaxMessageSize.Val()
	}
	return bpd
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
func (c *CommonBankParameterV3) referencedId() string { return ProcessingPreparationID }
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
		c.MinTimeoutValue,
		c.MaxTimeoutValue,
	}
}

func (c *CommonBankParameterV3) BankParameterData() domain.BankParameterData {
	bpd := domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
	}
	if c.MaxMessageSize != nil {
		bpd.MaxMessageSize = c.MaxMessageSize.Val()
	}
	if c.MinTimeoutValue != nil {
		bpd.MinTimeout = c.MinTimeoutValue.Val()
	}
	if c.MaxTimeoutValue != nil {
		bpd.MaxTimeout = c.MaxTimeoutValue.Val()
	}
	return bpd
}
