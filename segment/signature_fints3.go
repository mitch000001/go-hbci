package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewFINTS3PinTanSignatureHeaderSegment(controlReference string, holderId string, keyName domain.KeyName) *SignatureHeaderSegmentV4 {
	v3 := &SignatureHeaderSegment{
		SecurityFunction:         element.NewAlphaNumeric("998", 3),
		SecurityControlRef:       element.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: element.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     element.NewAlphaNumeric("1", 3),
		SecurityID:               element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, holderId),
		SecurityDate:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		HashAlgorithm:            element.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       element.NewRDHSignatureAlgorithm(),
		KeyName:                  element.NewKeyName(keyName),
	}
	s := &SignatureHeaderSegmentV4{
		SignatureHeaderSegment: v3,
		SecurityProfile:        element.NewPinTanSecurityProfile(),
	}
	s.Segment = NewBasicSegment(2, s)
	return s
}

func NewFINTS3SignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.Segment = NewBasicSegment(number, s)
	return s
}

type SignatureHeaderSegmentV4 struct {
	*SignatureHeaderSegment
	SecurityProfile *element.SecurityProfileDataElement
}

func (s *SignatureHeaderSegmentV4) elements() []element.DataElement {
	return []element.DataElement{
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
