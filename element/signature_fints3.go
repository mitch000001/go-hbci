package element

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

// GroupDataElements returns the grouped DataElements
func (s *SecurityProfileDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.SecurityMethod,
		s.SecurityMethodVersion,
	}
}
