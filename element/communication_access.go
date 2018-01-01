package element

import "github.com/mitch000001/go-hbci/domain"

// NewCommunicationParameter returns a new CommunicationParameterDataElement
func NewCommunicationParameter(params domain.CommunicationParameter) *CommunicationParameterDataElement {
	c := &CommunicationParameterDataElement{
		Protocol:              NewNumber(params.Protocol, 2),
		Address:               NewAlphaNumeric(params.Address, 512),
		AddressAddition:       NewAlphaNumeric(params.AddressAddition, 512),
		FilterFunction:        NewAlphaNumeric(params.FilterFunction, 3),
		FilterFunctionVersion: NewNumber(params.FilterFunctionVersion, 3),
	}
	c.DataElement = NewDataElementGroup(communicationParameterDEG, 5, c)
	return c
}

// CommunicationParameterDataElement represents the supported way of
// communicating with a HBCI server.
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

// GroupDataElements returns the grouped DataElements
func (c *CommunicationParameterDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		c.Protocol,
		c.Address,
		c.AddressAddition,
		c.FilterFunction,
		c.FilterFunctionVersion,
	}
}
