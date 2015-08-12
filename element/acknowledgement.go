package element

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
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
	Code                     *DigitDataElement
	ReferenceDataElement     *AlphaNumericDataElement
	Text                     *AlphaNumericDataElement
	Params                   *ParamsDataElement
	referencingMessage       domain.ReferencingMessage
	referencingSegmentNumber int
	typ                      string
}

func (a *AcknowledgementDataElement) SetReferencingMessage(reference domain.ReferencingMessage) {
	a.referencingMessage = reference
}

func (a *AcknowledgementDataElement) SetReferencingSegmentNumber(number int) {
	a.referencingSegmentNumber = number
}

func (a *AcknowledgementDataElement) SetType(acknowledgementType string) {
	a.typ = acknowledgementType
}

func (a *AcknowledgementDataElement) Val() domain.Acknowledgement {
	return domain.Acknowledgement{
		Code:                     a.Code.Val(),
		ReferenceDataElement:     a.ReferenceDataElement.Val(),
		Text:                     a.Text.Val(),
		Params:                   a.Params.Val(),
		Type:                     a.typ,
		ReferencingMessage:       a.referencingMessage,
		ReferencingSegmentNumber: a.referencingSegmentNumber,
	}
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
	chunks, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(chunks) < 3 {
		return fmt.Errorf("Malformed acknowledgment to unmarshal")
	}
	code, err := strconv.Atoi(charset.ToUtf8(chunks[0]))
	if err != nil {
		return fmt.Errorf("%T: Malformed code", a)
	}
	acknowledgement.Code = code
	acknowledgement.ReferenceDataElement = charset.ToUtf8(chunks[1])
	acknowledgement.Text = charset.ToUtf8(chunks[2])
	if len(chunks) > 3 {
		params := make([]string, len(chunks[3:]))
		for i, chunk := range chunks[3:] {
			params[i] = charset.ToUtf8(chunk)
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

func (p *ParamsDataElement) Val() []string {
	params := make([]string, len(p.array))
	for i, de := range p.array {
		params[i] = de.Value().(string)
	}
	return params
}

func (p *ParamsDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return nil
	}
	if len(elements) > 10 {
		return fmt.Errorf("Malformed params")
	}
	dataElements := make([]DataElement, len(elements))
	for i, elem := range elements {
		dataElements[i] = NewAlphaNumeric(charset.ToUtf8(elem), 35)
	}
	p.arrayElementGroup = NewArrayElementGroup(AcknowlegdementParamsGDEG, 10, 10, dataElements...)
	return nil
}
