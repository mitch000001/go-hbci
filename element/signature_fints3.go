package element

import "fmt"

// NewPinTanSecurityProfile returns a new SecurityProfile for the provided
// securityMehtod
func NewPinTanSecurityProfile(securityMethod int) *SecurityProfileDataElement {
	s := &SecurityProfileDataElement{
		SecurityMethod:        NewAlphaNumeric("PIN", 3),
		SecurityMethodVersion: NewNumber(securityMethod, 3),
	}
	s.DataElement = NewDataElementGroup(securityProfileDEG, 2, s)
	return s
}

// SecurityProfileDataElement defines a security method for the dialog flow
type SecurityProfileDataElement struct {
	DataElement
	SecurityMethod        *AlphaNumericDataElement
	SecurityMethodVersion *NumberDataElement
}

// UnmarshalHBCI unmarshals value into the DataElement
func (s *SecurityProfileDataElement) UnmarshalHBCI(value []byte) error {
	s = &SecurityProfileDataElement{}
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 2 {
		return fmt.Errorf("malformed marshaled value")
	}
	s.DataElement = NewDataElementGroup(securityProfileDEG, 5, s)
	s.SecurityMethod = &AlphaNumericDataElement{}
	err = s.SecurityMethod.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	s.SecurityMethodVersion = &NumberDataElement{}
	err = s.SecurityMethodVersion.UnmarshalHBCI(elements[1])
	return err
}

// GroupDataElements returns the grouped DataElements
func (s *SecurityProfileDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.SecurityMethod,
		s.SecurityMethodVersion,
	}
}
