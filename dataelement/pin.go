package dataelement

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

func NewPinTanDataElement(pin, tan string) *PinTanDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumericDataElement(pin, 6),
	}
	if tan != "" {
		p.TAN = NewAlphaNumericDataElement(tan, 35)
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
