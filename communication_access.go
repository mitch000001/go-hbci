package hbci

import (
	"github.com/mitch000001/go-hbci/dataelement"
	"github.com/mitch000001/go-hbci/domain"
)

func NewCommunicationAccessMessage(fromBank domain.BankId, toBank domain.BankId, maxEntries int, aufsetzpunkt string) *CommunicationAccessMessage {
	c := &CommunicationAccessMessage{
		Request: NewCommunicationAccessRequestSegment(fromBank, toBank, maxEntries, aufsetzpunkt),
	}
	c.basicMessage = newBasicMessage(c)
	return c
}

type CommunicationAccessMessage struct {
	*basicMessage
	Request *CommunicationAccessRequestSegment
}

func (c *CommunicationAccessMessage) HBCISegments() []Segment {
	return []Segment{
		c.Request,
	}
}

func NewFINTS3CommunicationAccessRequestSegment(fromBank domain.BankId, toBank domain.BankId, maxEntries int) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID: dataelement.NewBankIndentificationDataElement(fromBank),
		ToBankID:   dataelement.NewBankIndentificationDataElement(toBank),
		MaxEntries: dataelement.NewNumberDataElement(maxEntries, 4),
	}
	c.Segment = NewBasicSegment("HKKOM", 2, 4, c)
	return c
}

func NewCommunicationAccessRequestSegment(fromBank domain.BankId, toBank domain.BankId, maxEntries int, aufsetzpunkt string) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID:   dataelement.NewBankIndentificationDataElement(fromBank),
		ToBankID:     dataelement.NewBankIndentificationDataElement(toBank),
		MaxEntries:   dataelement.NewNumberDataElement(maxEntries, 4),
		Aufsetzpunkt: dataelement.NewAlphaNumericDataElement(aufsetzpunkt, 35),
	}
	c.Segment = NewBasicSegment("HKKOM", 2, 3, c)
	return c
}

type CommunicationAccessRequestSegment struct {
	Segment
	FromBankID *dataelement.BankIdentificationDataElement
	ToBankID   *dataelement.BankIdentificationDataElement
	MaxEntries *dataelement.NumberDataElement
	// TODO: find a fitting name
	Aufsetzpunkt *dataelement.AlphaNumericDataElement
}

func (c *CommunicationAccessRequestSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		c.FromBankID,
		c.ToBankID,
		c.MaxEntries,
		c.Aufsetzpunkt,
	}
}

const HKKOMSegmentNumber = -1

func NewCommunicationAccessResponseSegment(bankId domain.BankId, language int, params domain.CommunicationParameter) *CommunicationAccessResponseSegment {
	c := &CommunicationAccessResponseSegment{
		BankID:              dataelement.NewBankIndentificationDataElement(bankId),
		StandardLanguage:    dataelement.NewNumberDataElement(language, 3),
		CommunicationParams: dataelement.NewCommunicationParameterDataElement(params),
	}
	header := dataelement.NewReferencingSegmentHeader("HIKOM", 4, 3, HKKOMSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

type CommunicationAccessResponseSegment struct {
	Segment
	BankID              *dataelement.BankIdentificationDataElement
	StandardLanguage    *dataelement.NumberDataElement
	CommunicationParams *dataelement.CommunicationParameterDataElement
}

func (c *CommunicationAccessResponseSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		c.BankID,
		c.StandardLanguage,
		c.CommunicationParams,
	}
}
