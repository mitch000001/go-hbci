package hbci

import "time"

func NewFINTS3PinTanSignatureHeaderSegment(controlReference string, holderId string, keyName KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         NewAlphaNumericDataElement("998", 3),
		SecurityControlRef:       NewAlphaNumericDataElement(controlReference, 14),
		SecurityApplicationRange: NewAlphaNumericDataElement("1", 3),
		SecuritySupplierRole:     NewAlphaNumericDataElement("1", 3),
		SecurityID:               NewRDHSecurityIdentificationDataElement(SecurityHolderMessageSender, holderId),
		SecurityDate:             NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		HashAlgorithm:            NewDefaultHashAlgorithmDataElement(),
		SignatureAlgorithm:       NewRDHSignatureAlgorithmDataElement(),
		KeyName:                  NewKeyNameDataElement(keyName),
	}
	v4 := &SignatureHeaderVersion4{
		SignatureHeaderVersion3: v3,
		SecurityProfile:         NewPinTanSecurityProfile(),
	}
	s := &SignatureHeaderSegment{
		version: v4,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 4, s)
	return s
}

func NewFINTS3SignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: NewAlphaNumericDataElement(controlReference, 14),
	}
	s.Segment = NewBasicSegment("HNSHA", number, 2, s)
	return s
}

type SignatureHeaderVersion4 struct {
	*SignatureHeaderVersion3
	SecurityProfile *SecurityProfileDataElement
}

func (s *SignatureHeaderVersion4) versionedElements() []DataElement {
	return []DataElement{
		s.SecurityProfile,
		s.SecurityFunction,
		s.SecurityControlRef,
		s.SecurityApplicationRange,
		s.SecuritySupplierRole,
		s.SecurityID,
		s.SecurityRefNumber,
		s.SecurityDate,
		s.HashAlgorithm,
		s.SignatureAlgorithm,
		s.KeyName,
		s.Certificate,
	}
}

func NewPinTanSecurityProfile() *SecurityProfileDataElement {
	s := &SecurityProfileDataElement{
		SecurityMethod:        NewAlphaNumericDataElement("PIN", 3),
		SecurityMethodVersion: NewNumberDataElement(1, 3),
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
