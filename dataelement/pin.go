package dataelement

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

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
