package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewFINTS3CommunicationAccessRequestSegment(fromBank domain.BankId, toBank domain.BankId, maxEntries int) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID: element.NewBankIndentification(fromBank),
		ToBankID:   element.NewBankIndentification(toBank),
		MaxEntries: element.NewNumber(maxEntries, 4),
	}
	c.Segment = NewBasicSegment(2, c)
	return c
}

func NewCommunicationAccessRequestSegment(fromBank domain.BankId, toBank domain.BankId, maxEntries int, aufsetzpunkt string) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID:   element.NewBankIndentification(fromBank),
		ToBankID:     element.NewBankIndentification(toBank),
		MaxEntries:   element.NewNumber(maxEntries, 4),
		Aufsetzpunkt: element.NewAlphaNumeric(aufsetzpunkt, 35),
	}
	c.Segment = NewBasicSegment(2, c)
	return c
}

type CommunicationAccessRequestSegment struct {
	Segment
	FromBankID *element.BankIdentificationDataElement
	ToBankID   *element.BankIdentificationDataElement
	MaxEntries *element.NumberDataElement
	// TODO: find a fitting name
	Aufsetzpunkt *element.AlphaNumericDataElement
}

func (c *CommunicationAccessRequestSegment) Version() int         { return 3 }
func (c *CommunicationAccessRequestSegment) ID() string           { return "HKKOM" }
func (c *CommunicationAccessRequestSegment) referencedId() string { return "" }
func (c *CommunicationAccessRequestSegment) sender() string       { return senderUser }

func (c *CommunicationAccessRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.FromBankID,
		c.ToBankID,
		c.MaxEntries,
		c.Aufsetzpunkt,
	}
}

const HKKOMSegmentNumber = -1

func NewCommunicationAccessResponseSegment(bankId domain.BankId, language int, params domain.CommunicationParameter) *CommunicationAccessResponseSegment {
	c := &CommunicationAccessResponseSegment{
		BankID:              element.NewBankIndentification(bankId),
		StandardLanguage:    element.NewNumber(language, 3),
		CommunicationParams: element.NewCommunicationParameter(params),
	}
	header := element.NewReferencingSegmentHeader("HIKOM", 4, 3, HKKOMSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

type CommunicationAccessResponseSegment struct {
	Segment
	BankID              *element.BankIdentificationDataElement
	StandardLanguage    *element.NumberDataElement
	CommunicationParams *element.CommunicationParameterDataElement
}

func (c *CommunicationAccessResponseSegment) Version() int         { return 3 }
func (c *CommunicationAccessResponseSegment) ID() string           { return "HIKOM" }
func (c *CommunicationAccessResponseSegment) referencedId() string { return "HKKOM" }
func (c *CommunicationAccessResponseSegment) sender() string       { return senderBank }

func (c *CommunicationAccessResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.BankID,
		c.StandardLanguage,
		c.CommunicationParams,
	}
}
