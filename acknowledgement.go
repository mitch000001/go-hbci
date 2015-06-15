package hbci

type MessageAcknowledgement struct {
	*segment
	Acknowledgements []*AcknowledgementDataElement
}
type SegmentAcknowledgement struct{}

func NewAcknowledgementDataElement(code int, referenceDataElement, text string, params []string) *AcknowledgementDataElement {
	paramDataElements := make([]*AlphaNumericDataElement, len(params))
	if params != nil {
		for i, p := range params {
			paramDataElements[i] = NewAlphaNumericDataElement(p, 35)
		}
	}
	a := &AcknowledgementDataElement{
		Code:                 NewDigitDataElement(code, 4),
		ReferenceDataElement: NewAlphaNumericDataElement(referenceDataElement, 7),
		Text:                 NewAlphaNumericDataElement(text, 80),
		Params:               paramDataElements,
	}
	a.elementGroup = NewDataElementGroup(AcknowledgementDEG, 3+len(paramDataElements), a)
	return a
}

type AcknowledgementDataElement struct {
	*elementGroup
	Code                 *DigitDataElement
	ReferenceDataElement *AlphaNumericDataElement
	Text                 *AlphaNumericDataElement
	Params               []*AlphaNumericDataElement
}

func (a *AcknowledgementDataElement) Valid() bool {
	if len(a.Params) != 0 || len(a.Params) != 10 {
		return false
	} else {
		if a.Code == nil || a.Text == nil {
			return false
		} else {
			return a.elementGroup.Valid()
		}
	}
}

func (a *AcknowledgementDataElement) GroupDataElements() []DataElement {
	elements := []DataElement{
		a.Code,
		a.ReferenceDataElement,
		a.Text,
	}
	var paramElements []DataElement
	for _, p := range a.Params {
		paramElements = append(paramElements, p)
	}
	elements = append(elements, paramElements...)
	return elements
}
