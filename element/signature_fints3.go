package element

func NewPinTanSecurityProfile(securityMethod int) *SecurityProfileDataElement {
	s := &SecurityProfileDataElement{
		SecurityMethod:        NewAlphaNumeric("PIN", 3),
		SecurityMethodVersion: NewNumber(securityMethod, 3),
	}
	s.DataElement = NewDataElementGroup(SecurityProfileDEG, 2, s)
	return s
}

type SecurityProfileDataElement struct {
	DataElement
	SecurityMethod        *AlphaNumericDataElement
	SecurityMethodVersion *NumberDataElement
}

func (s *SecurityProfileDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.SecurityMethod,
		s.SecurityMethodVersion,
	}
}
