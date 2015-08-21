package element

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

func NewCustomSignature(pin, tan string) *CustomSignatureDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 99),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 99)
	}
	p.DataElement = NewDataElementGroup(PinTanDEG, 2, p)
	cust := &CustomSignatureDataElement{
		PinTanDataElement: p,
	}
	return cust
}

type CustomSignatureDataElement struct {
	*PinTanDataElement
}

func NewPinTan(pin, tan string) *PinTanDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 6),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 35)
	}
	p.DataElement = NewDataElementGroup(PinTanDEG, 2, p)
	return p
}

type PinTanDataElement struct {
	DataElement
	PIN *AlphaNumericDataElement
	TAN *AlphaNumericDataElement
}

func (p *PinTanDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.PIN,
		p.TAN,
	}
}
