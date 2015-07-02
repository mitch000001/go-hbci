package element

import "github.com/mitch000001/go-hbci/domain"

func NewAcknowledgements(acknowledgements []domain.Acknowledgement) *AcknowledgementsDataElement {
	ackDEs := make([]DataElement, len(acknowledgements))
	for i, acknowledgement := range acknowledgements {
		ackDEs[i] = NewAcknowledgement(acknowledgement)
	}
	a := &AcknowledgementsDataElement{
		arrayElementGroup: NewArrayElementGroup(AcknowledgementDEG, 1, 99, ackDEs...),
	}
	return a
}

type AcknowledgementsDataElement struct {
	*arrayElementGroup
}

func NewAcknowledgement(acknowledgement domain.Acknowledgement) *AcknowledgementDataElement {
	a := &AcknowledgementDataElement{
		Code:                 NewDigit(acknowledgement.Code, 4),
		ReferenceDataElement: NewAlphaNumeric(acknowledgement.ReferenceDataElement, 7),
		Text:                 NewAlphaNumeric(acknowledgement.Text, 80),
		Params:               NewParams(10, 10, acknowledgement.Params...),
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
