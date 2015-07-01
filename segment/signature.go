package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         element.NewAlphaNumeric("999", 3),
		SecurityControlRef:       element.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: element.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     element.NewAlphaNumeric("1", 3),
		SecurityID:               element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        element.NewNumber(0, 16),
		SecurityDate:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		HashAlgorithm:            element.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       element.NewRDHSignatureAlgorithm(),
		KeyName:                  element.NewKeyName(keyName),
	}
	s := &SignatureHeaderSegment{
		version: v3,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

func NewRDHSignatureHeaderSegment(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         element.NewAlphaNumeric("1", 3),
		SecurityControlRef:       element.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: element.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     element.NewAlphaNumeric("1", 3),
		SecurityID:               element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        element.NewNumber(signatureId, 16),
		SecurityDate:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		HashAlgorithm:            element.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       element.NewRDHSignatureAlgorithm(),
		KeyName:                  element.NewKeyName(keyName),
	}
	s := &SignatureHeaderSegment{
		version: v3,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

type SignatureHeaderSegment struct {
	Segment
	version
}

func (s *SignatureHeaderSegment) elements() []element.DataElement {
	return s.version.versionedElements()
}

type SignatureHeaderVersion3 struct {
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

func (s SignatureHeaderVersion3) version() int {
	return 3
}

func (s *SignatureHeaderVersion3) versionedElements() []element.DataElement {
	return []element.DataElement{
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

func NewSignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.Segment = NewBasicSegment("HNSHA", number, 1, s)
	return s
}

type SignatureEndSegment struct {
	Segment
	SecurityControlRef *element.AlphaNumericDataElement
	Signature          *element.BinaryDataElement
	PinTan             *element.PinTanDataElement
}

func (s *SignatureEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.Signature = element.NewBinary(signature, 512)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.PinTan = element.NewPinTan(pin, tan)
}
