package hbci

func NewCommunicationAccessMessage(fromBank BankId, toBank BankId, maxEntries int, aufsetzpunkt string) *CommunicationAccessMessage {
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

func NewFINTS3CommunicationAccessRequestSegment(fromBank BankId, toBank BankId, maxEntries int) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID: NewBankIndentificationDataElement(fromBank),
		ToBankID:   NewBankIndentificationDataElement(toBank),
		MaxEntries: NewNumberDataElement(maxEntries, 4),
	}
	c.Segment = NewBasicSegment("HKKOM", 2, 4, c)
	return c
}

func NewCommunicationAccessRequestSegment(fromBank BankId, toBank BankId, maxEntries int, aufsetzpunkt string) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID:   NewBankIndentificationDataElement(fromBank),
		ToBankID:     NewBankIndentificationDataElement(toBank),
		MaxEntries:   NewNumberDataElement(maxEntries, 4),
		Aufsetzpunkt: NewAlphaNumericDataElement(aufsetzpunkt, 35),
	}
	c.Segment = NewBasicSegment("HKKOM", 2, 3, c)
	return c
}

type CommunicationAccessRequestSegment struct {
	Segment
	FromBankID *BankIdentificationDataElement
	ToBankID   *BankIdentificationDataElement
	MaxEntries *NumberDataElement
	// TODO: find a fitting name
	Aufsetzpunkt *AlphaNumericDataElement
}

func (c *CommunicationAccessRequestSegment) elements() []DataElement {
	return []DataElement{
		c.FromBankID,
		c.ToBankID,
		c.MaxEntries,
		c.Aufsetzpunkt,
	}
}

const HKKOMSegmentNumber = -1

func NewCommunicationAccessResponseSegment(bankId BankId, language int, params CommunicationParameter) *CommunicationAccessResponseSegment {
	c := &CommunicationAccessResponseSegment{
		BankID:              NewBankIndentificationDataElement(bankId),
		StandardLanguage:    NewNumberDataElement(language, 3),
		CommunicationParams: NewCommunicationParameterDataElement(params),
	}
	header := NewReferencingSegmentHeader("HIKOM", 4, 3, HKKOMSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

type CommunicationAccessResponseSegment struct {
	Segment
	BankID              *BankIdentificationDataElement
	StandardLanguage    *NumberDataElement
	CommunicationParams *CommunicationParameterDataElement
}

func (c *CommunicationAccessResponseSegment) elements() []DataElement {
	return []DataElement{
		c.BankID,
		c.StandardLanguage,
		c.CommunicationParams,
	}
}

type CommunicationParameter struct {
	Protocol              int
	Address               string
	AddressAddition       string
	FilterFunction        string
	FilterFunctionVersion int
}

func NewCommunicationParameterDataElement(params CommunicationParameter) *CommunicationParameterDataElement {
	c := &CommunicationParameterDataElement{
		Protocol:              NewNumberDataElement(params.Protocol, 2),
		Address:               NewAlphaNumericDataElement(params.Address, 512),
		AddressAddition:       NewAlphaNumericDataElement(params.AddressAddition, 512),
		FilterFunction:        NewAlphaNumericDataElement(params.FilterFunction, 3),
		FilterFunctionVersion: NewNumberDataElement(params.FilterFunctionVersion, 3),
	}
	c.DataElement = NewDataElementGroup(CommunicationParameterDEG, 5, c)
	return c
}

type CommunicationParameterDataElement struct {
	DataElement
	// Code | Zugang   | Protokollstack
	// ---------------------------------------------------￼
	// 1	| T-Online | ETSI 300 072 (CEPT), EHKP, BtxFIF
	// 2 	| TCP/IP ￼ | SLIP/PPP
	// 3	| HTTPS	   | (für PIN/TAN-Verfahren)
	Protocol *NumberDataElement
	// Zugang ￼ |  Adresse ￼    |￼ Anmerkungen
	// ---------------------------------------------------------------------------------------------------
	// T-Online | Gateway-Seite | als numerischer Wert (ohne die Steuerzeichen * und #) einzustellen.
	// TCP/IP	| IP-Adresse ￼	| als alphanumerischer Wert (z.B. ‘123.123.123.123’)
	// HTTPS    | Adresse       | als alphanumerischer Wert (z.B. ‚https://www.xyz.de:7000/PinTanServlet‘)
	Address *AlphaNumericDataElement
	// Zugang ￼ | Adressenzusatz ￼| Anmerkungen
	// ----------------------------------------------------------------------------------
	// T-Online | Regionalbereich | Für ein bundesweites Angebot ist ‘00’ ein- zustellen’
	// TCP/IP ￼ | nicht belegt	  |
	// HTTPS  ￼ | nicht belegt	  |
	AddressAddition *AlphaNumericDataElement
	// Code | Bedeutung ￼
	// ------------------------
	// MIM  | MIME Base 64
	// UUE ￼| Uuencode/Uudecode
	FilterFunction        *AlphaNumericDataElement
	FilterFunctionVersion *NumberDataElement
}

func (c *CommunicationParameterDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		c.Protocol,
		c.Address,
		c.AddressAddition,
		c.FilterFunction,
		c.FilterFunctionVersion,
	}
}
