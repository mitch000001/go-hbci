package dataelement

func NewAcknowledgement(code int, referenceDataElement, text string, params []string) *AcknowledgementDataElement {
	paramDataElements := make([]*AlphaNumericDataElement, len(params))
	if params != nil {
		for i, p := range params {
			paramDataElements[i] = NewAlphaNumeric(p, 35)
		}
	}
	a := &AcknowledgementDataElement{
		Code:                 NewDigit(code, 4),
		ReferenceDataElement: NewAlphaNumeric(referenceDataElement, 7),
		Text:                 NewAlphaNumeric(text, 80),
		Params:               NewParams(10, 10, params...),
	}
	a.DataElement = NewDataElementGroup(AcknowledgementDEG, 4, a)
	return a
}

type AcknowledgementDataElement struct {
	DataElement
	Code                 *DigitDataElement
	ReferenceDataElement *AlphaNumericDataElement
	Text                 *AlphaNumericDataElement
	Params               *ParamsDataElement
}

func (a *AcknowledgementDataElement) IsValid() bool {
	if a.Code == nil || a.Text == nil {
		return false
	} else {
		return a.DataElement.IsValid()
	}
}

func (a *AcknowledgementDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.Code,
		a.ReferenceDataElement,
		a.Text,
		a.Params,
	}
}

func NewParams(min, max int, params ...string) *ParamsDataElement {
	var paramDE []DataElement
	for _, p := range params {
		paramDE = append(paramDE, NewAlphaNumeric(p, 35))
	}
	return &ParamsDataElement{arrayElementGroup: NewArrayElementGroup(AcknowlegdementParamsGDEG, min, max, paramDE...)}
}

type ParamsDataElement struct {
	*arrayElementGroup
}
