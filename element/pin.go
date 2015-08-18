package element

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

func NewFINTSPinTan(pin, tan string) *PinTanDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 99),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 99)
	}
	p.DataElement = NewDataElementGroup(PinTanDEG, 2, p)
	return p
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
