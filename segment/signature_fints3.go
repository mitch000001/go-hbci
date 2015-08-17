package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewFINTS3PinTanSignatureHeaderSegment(controlReference string, holderId string, keyName domain.KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderSegmentV4{
		SecurityProfile:          element.NewPinTanSecurityProfile(),
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
	s.ClientSegment = NewBasicSegment(2, s)
	segment := &SignatureHeaderSegment{
		signatureHeaderSegment: s,
	}
	return segment
}

func NewFINTS3SignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.ClientSegment = NewBasicSegment(number, s)
	return s
}

type SignatureHeaderSegmentV4 struct {
	ClientSegment
	SecurityProfile *element.SecurityProfileDataElement
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	// "999" for PIN/TAN
	SecurityFunction   *element.AlphaNumericDataElement
	SecurityControlRef *element.AlphaNumericDataElement
	// "1" for SHM (SignatureHeader and HBCI-Data)
	// "2" for SHT (SignatureHeader to SignatureEnd)
	SecurityApplicationRange *element.AlphaNumericDataElement
	// "1" for ISS, Herausgeber der signierten Nachricht (z.B. Erfasser oder Erstsignatur)
	// "3" for CON, der Unterzeichnete unterstützt den Inhalt der Nachricht (z.B. bei Zweitsignatur)
	// "4" for WIT, der Unterzeichnete ist Zeuge (z.B. Übermittler), aber für den Inhalt der Nachricht nicht verantwortlich)
	SecuritySupplierRole *element.AlphaNumericDataElement
	SecurityID           *element.SecurityIdentificationDataElement
	SecurityRefNumber    *element.NumberDataElement
	SecurityDate         *element.SecurityDateDataElement
	HashAlgorithm        *element.HashAlgorithmDataElement
	SignatureAlgorithm   *element.SignatureAlgorithmDataElement
	KeyName              *element.KeyNameDataElement
	Certificate          *element.CertificateDataElement
}

func (s *SignatureHeaderSegmentV4) SetSecurityFunction(securityFn string) {
	s.SecurityFunction = element.NewAlphaNumeric(securityFn, 3)
}

func (s *SignatureHeaderSegmentV4) Version() int         { return 3 }
func (s *SignatureHeaderSegmentV4) ID() string           { return "HNSHK" }
func (s *SignatureHeaderSegmentV4) referencedId() string { return "" }
func (s *SignatureHeaderSegmentV4) sender() string       { return senderBoth }

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
