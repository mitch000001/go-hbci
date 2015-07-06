package element

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/domain"
)

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

func (a *AcknowledgementDataElement) UnmarshalHBCI(value []byte) error {
	acknowledgement := domain.Acknowledgement{}
	chunks := bytes.Split(value, []byte(":"))
	if len(chunks) < 3 {
		return fmt.Errorf("Malformed acknowledgment to unmarshal")
	}
	code, err := strconv.Atoi(string(chunks[0]))
	if err != nil {
		return fmt.Errorf("%T: Malformed code", a)
	}
	acknowledgement.Code = code
	acknowledgement.ReferenceDataElement = string(chunks[1])
	acknowledgement.Text = string(chunks[2])
	if len(chunks) > 3 {
		params := make([]string, len(chunks[3:]))
		for i, chunk := range chunks[3:] {
			params[i] = string(chunk)
		}
		acknowledgement.Params = params
	}
	*a = *NewAcknowledgement(acknowledgement)
	return nil
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

func (p *ParamsDataElement) UnmarshalHBCI(value []byte) error {
	elements := bytes.Split(value, []byte(":"))
	if len(elements) > 10 {
		return fmt.Errorf("Malformed params")
	}
	dataElements := make([]DataElement, len(elements))
	for i, elem := range elements {
		dataElements[i] = NewAlphaNumeric(string(elem), 35)
	}
	p.arrayElementGroup = NewArrayElementGroup(AcknowlegdementParamsGDEG, 10, 10, dataElements...)
	return nil
}
