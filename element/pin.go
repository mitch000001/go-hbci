package element

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

// NewCustomSignature returns a new CustomSignatureDataElement for the pin and
// tan
func NewCustomSignature(pin, tan string) *CustomSignatureDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 99),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 99)
	}
	p.DataElement = NewDataElementGroup(pinTanDEG, 2, p)
	cust := &CustomSignatureDataElement{
		PinTanDataElement: p,
	}
	return cust
}

// CustomSignatureDataElement represents a custom signature
type CustomSignatureDataElement struct {
	*PinTanDataElement
}

// NewPinTan returns a new PinTanDataElement for pin and tan
func NewPinTan(pin, tan string) *PinTanDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 6),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 35)
	}
	p.DataElement = NewDataElementGroup(pinTanDEG, 2, p)
	return p
}

// PinTanDataElement represents a DataElement which contains the PIN and the
// TAN for a transaction
type PinTanDataElement struct {
	DataElement
	PIN *AlphaNumericDataElement
	TAN *AlphaNumericDataElement
}

// GroupDataElements returns the grouped DataElements
func (p *PinTanDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.PIN,
		p.TAN,
	}
}
