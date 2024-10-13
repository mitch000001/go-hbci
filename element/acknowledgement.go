package element

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

// These represent HBCI acknowledgement codes. Codes starting with 3 are meant
// to be warnings.
const (
	AcknowledgementAdditionalInformation     = 3040
	AcknowledgementSupportedSecurityFunction = 3920
)

// NewAcknowledgement returns a new acknowledgement DataElement
func NewAcknowledgement(acknowledgement domain.Acknowledgement) *AcknowledgementDataElement {
	a := &AcknowledgementDataElement{
		Code:                 NewDigit(acknowledgement.Code, 4),
		ReferenceDataElement: NewAlphaNumeric(acknowledgement.ReferenceDataElementPosition, 7),
		Text:                 NewAlphaNumeric(acknowledgement.Text, 80),
		Params:               NewParams(10, 10, acknowledgement.Params...),
	}
	a.DataElement = NewDataElementGroup(acknowledgementDEG, 4, a)
	return a
}

// An AcknowledgementDataElement defines a group of DataElements used to
// transmit information from the bank institute to the client.
type AcknowledgementDataElement struct {
	DataElement
	Code                       *DigitDataElement
	ReferenceDataElement       *AlphaNumericDataElement
	Text                       *AlphaNumericDataElement
	Params                     *ParamsDataElement
	referencingMessage         domain.MessageReference
	referencingSegmentPosition int
	typ                        string
}

// SetReferencingMessage is used by MessageAcknowledgements to set the reference
// to the previously sent message. This is needed to identify the message within
// an ongoing dialog.
func (a *AcknowledgementDataElement) SetReferencingMessage(reference domain.MessageReference) {
	a.referencingMessage = reference
}

// SetReferencingSegmentPosition is a setter for setting the Segment position this
// Acknowledgement is referring to.
func (a *AcknowledgementDataElement) SetReferencingSegmentPosition(position int) {
	a.referencingSegmentPosition = position
}

// SetType sets the type of the Acknowledgement. There are only two types of
// Acknowledgements, message related ones and segment related ones.
//
// See domain.MessageAcknowledgements or domain.SegmentAcknowledgement for
// details on semantics.
func (a *AcknowledgementDataElement) SetType(acknowledgementType string) {
	a.typ = acknowledgementType
}

// Val returns the underlying value of the acknowledgement.
func (a *AcknowledgementDataElement) Val() domain.Acknowledgement {
	return domain.Acknowledgement{
		Code:                         a.Code.Val(),
		ReferenceDataElementPosition: a.ReferenceDataElement.Val(),
		Text:                         a.Text.Val(),
		Params:                       a.Params.Val(),
		Type:                         a.typ,
		ReferencingMessage:           a.referencingMessage,
		ReferencingSegmentNumber:     a.referencingSegmentPosition,
	}
}

// IsValid returns true if the acknowledgment code and text is set and the
// underlying value is valid.
func (a *AcknowledgementDataElement) IsValid() bool {
	if a.Code == nil || a.Text == nil {
		return false
	}
	return a.DataElement.IsValid()
}

// GroupDataElements returns all grouped elements within the acknowledgment.
func (a *AcknowledgementDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		a.Code,
		a.ReferenceDataElement,
		a.Text,
		a.Params,
	}
}

// UnmarshalHBCI unmarshals the value into an acknowledgment.
func (a *AcknowledgementDataElement) UnmarshalHBCI(value []byte) error {
	acknowledgement := domain.Acknowledgement{}
	chunks, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(chunks) < 3 {
		return fmt.Errorf("malformed acknowledgment to unmarshal")
	}
	code, err := strconv.Atoi(charset.ToUTF8(chunks[0]))
	if err != nil {
		return fmt.Errorf("%T: Malformed code", a)
	}
	acknowledgement.Code = code
	acknowledgement.ReferenceDataElementPosition = charset.ToUTF8(chunks[1])
	acknowledgement.Text = charset.ToUTF8(chunks[2])
	if len(chunks) > 3 {
		params := make([]string, len(chunks[3:]))
		for i, chunk := range chunks[3:] {
			params[i] = charset.ToUTF8(chunk)
		}
		acknowledgement.Params = params
	}
	*a = *NewAcknowledgement(acknowledgement)
	return nil
}

// NewParams returns a new ParamsDataElement.
func NewParams(min, max int, params ...string) *ParamsDataElement {
	var paramDE []DataElement
	for _, p := range params {
		paramDE = append(paramDE, NewAlphaNumeric(p, 35))
	}
	return &ParamsDataElement{arrayElementGroup: newArrayElementGroup(acknowlegdementParamsGDEG, min, max, paramDE)}
}

// ParamsDataElement defines a DataElement describing generic parameters
type ParamsDataElement struct {
	*arrayElementGroup
}

// Val returns the underlying value of the DataElement.
func (p *ParamsDataElement) Val() []string {
	params := make([]string, len(p.array))
	for i, de := range p.array {
		params[i] = de.Value().(string)
	}
	return params
}

// UnmarshalHBCI unmarshals the value into a ParamsDataElement
func (p *ParamsDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return nil
	}
	if len(elements) > 10 {
		return fmt.Errorf("malformed params")
	}
	dataElements := make([]DataElement, len(elements))
	for i, elem := range elements {
		dataElements[i] = NewAlphaNumeric(charset.ToUTF8(elem), 35)
	}
	p.arrayElementGroup = newArrayElementGroup(acknowlegdementParamsGDEG, 10, 10, dataElements)
	return nil
}
