package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/internal"
)

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

// NewCustomSignature returns a new CustomSignatureDataElement for the pin and
// tan
func NewCustomSignature(pin, tan string) *CustomSignatureDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 99),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 99)
	}
	p.DataElement = NewDataElementGroup(pinTanDEG, 2, p)
	cust := &CustomSignatureDataElement{
		PinTanDataElement: p,
	}
	return cust
}

// CustomSignatureDataElement represents a custom signature
type CustomSignatureDataElement struct {
	*PinTanDataElement
}

func (c *CustomSignatureDataElement) UnmarshalHBCI(value []byte) error {
	p := &PinTanDataElement{}
	if err := p.UnmarshalHBCI(value); err != nil {
		return err
	}
	c.PinTanDataElement = p
	return nil
}

// NewPinTan returns a new PinTanDataElement for pin and tan
func NewPinTan(pin, tan string) *PinTanDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumeric(pin, 6),
	}
	if tan != "" {
		p.TAN = NewAlphaNumeric(tan, 35)
	}
	p.DataElement = NewDataElementGroup(pinTanDEG, 2, p)
	return p
}

// PinTanDataElement represents a DataElement which contains the PIN and the
// TAN for a transaction
type PinTanDataElement struct {
	DataElement
	PIN *AlphaNumericDataElement
	TAN *AlphaNumericDataElement
}

// GroupDataElements returns the grouped DataElements
func (p *PinTanDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.PIN,
		p.TAN,
	}
}

// UnmarshalHBCI unmarshals value into the DataElement
func (p *PinTanDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) > 2 {
		return fmt.Errorf("malformed marshaled value")
	}
	p.DataElement = NewDataElementGroup(securityIdentificationDEG, 3, p)
	if len(elements) > 0 && len(elements[0]) > 0 {
		p.PIN = &AlphaNumericDataElement{}
		err = p.PIN.UnmarshalHBCI(elements[0])
		if err != nil {
			return err
		}
	}
	if len(elements) > 1 && len(elements[1]) > 0 {
		p.TAN = &AlphaNumericDataElement{}
		err = p.TAN.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	return nil
}

type PinTanSpecificParamDataElement struct {
	DataElement                  `yaml:"-"`
	PinMinLength                 *NumberDataElement                   `yaml:"PinMinLength"`
	PinMaxLength                 *NumberDataElement                   `yaml:"PinMaxLength"`
	TanMaxLength                 *NumberDataElement                   `yaml:"TanMaxLength"`
	UserIDText                   *AlphaNumericDataElement             `yaml:"UserIDText"`
	CustomerIDText               *AlphaNumericDataElement             `yaml:"CustomerIDText"`
	JobSpecificPinTanInformation *PinTanBusinessTransactionParameters `yaml:"JobSpecificPinTanInformation"`
}

// Elements returns the grouped DataElements
func (p *PinTanSpecificParamDataElement) Elements() []DataElement {
	return []DataElement{
		p.PinMinLength,
		p.PinMaxLength,
		p.TanMaxLength,
		p.UserIDText,
		p.CustomerIDText,
		p.JobSpecificPinTanInformation,
	}
}

// UnmarshalHBCI unmarshals value
func (t *PinTanSpecificParamDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	iter := internal.NewIterator(elements)
	var PinMinLength NumberDataElement
	if err := PinMinLength.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling PinMinLength: %v", err)
	}
	t.PinMinLength = &PinMinLength
	var PinMaxLength NumberDataElement
	if err := PinMaxLength.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling PinMaxLength: %v", err)
	}
	t.PinMaxLength = &PinMaxLength
	var TanMaxLength NumberDataElement
	if err := TanMaxLength.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling TanMaxLength: %v", err)
	}
	t.TanMaxLength = &TanMaxLength
	var UserIDText AlphaNumericDataElement
	if err := UserIDText.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling UserIDText: %v", err)
	}
	t.UserIDText = &UserIDText
	var CustomerIDText AlphaNumericDataElement
	if err := CustomerIDText.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling CustomerIDText: %v", err)
	}
	t.CustomerIDText = &CustomerIDText
	var JobSpecificPinTanInformation PinTanBusinessTransactionParameters
	if err := JobSpecificPinTanInformation.UnmarshalHBCI(bytes.Join(iter.Remainder(), []byte(":"))); err != nil {
		return fmt.Errorf("error unmarshaling JobSpecificPinTanInformation: %v", err)
	}
	t.JobSpecificPinTanInformation = &JobSpecificPinTanInformation
	t.DataElement = NewGroupDataElementGroup(pinTanSpecificParamDataElementDEG, 2, t)
	return nil
}
