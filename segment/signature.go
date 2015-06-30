package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/dataelement"
	"github.com/mitch000001/go-hbci/domain"
)

func NewPinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         dataelement.NewAlphaNumeric("999", 3),
		SecurityControlRef:       dataelement.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: dataelement.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     dataelement.NewAlphaNumeric("1", 3),
		SecurityID:               dataelement.NewRDHSecurityIdentification(dataelement.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        dataelement.NewNumber(0, 16),
		SecurityDate:             dataelement.NewSecurityDate(dataelement.SecurityTimestamp, time.Now()),
		HashAlgorithm:            dataelement.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       dataelement.NewRDHSignatureAlgorithm(),
		KeyName:                  dataelement.NewKeyName(keyName),
	}
	s := &SignatureHeaderSegment{
		version: v3,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

func NewRDHSignatureHeaderSegment(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         dataelement.NewAlphaNumeric("1", 3),
		SecurityControlRef:       dataelement.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: dataelement.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     dataelement.NewAlphaNumeric("1", 3),
		SecurityID:               dataelement.NewRDHSecurityIdentification(dataelement.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        dataelement.NewNumber(signatureId, 16),
		SecurityDate:             dataelement.NewSecurityDate(dataelement.SecurityTimestamp, time.Now()),
		HashAlgorithm:            dataelement.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       dataelement.NewRDHSignatureAlgorithm(),
		KeyName:                  dataelement.NewKeyName(keyName),
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

func (s *SignatureHeaderSegment) elements() []dataelement.DataElement {
	return s.version.versionedElements()
}

type SignatureHeaderVersion3 struct {
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	// "999" for PIN/TAN
	SecurityFunction   *dataelement.AlphaNumericDataElement
	SecurityControlRef *dataelement.AlphaNumericDataElement
	// "1" for SHM (SignatureHeader and HBCI-Data)
	// "2" for SHT (SignatureHeader to SignatureEnd)
	SecurityApplicationRange *dataelement.AlphaNumericDataElement
	// "1" for ISS, Herausgeber der signierten Nachricht (z.B. Erfasser oder Erstsignatur)
	// "3" for CON, der Unterzeichnete unterstützt den Inhalt der Nachricht (z.B. bei Zweitsignatur)
	// "4" for WIT, der Unterzeichnete ist Zeuge (z.B. Übermittler), aber für den Inhalt der Nachricht nicht verantwortlich)
	SecuritySupplierRole *dataelement.AlphaNumericDataElement
	SecurityID           *dataelement.SecurityIdentificationDataElement
	SecurityRefNumber    *dataelement.NumberDataElement
	SecurityDate         *dataelement.SecurityDateDataElement
	HashAlgorithm        *dataelement.HashAlgorithmDataElement
	SignatureAlgorithm   *dataelement.SignatureAlgorithmDataElement
	KeyName              *dataelement.KeyNameDataElement
	Certificate          *dataelement.CertificateDataElement
}

func (s SignatureHeaderVersion3) version() int {
	return 3
}

func (s *SignatureHeaderVersion3) versionedElements() []dataelement.DataElement {
	return []dataelement.DataElement{
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
		SecurityControlRef: dataelement.NewAlphaNumeric(controlReference, 14),
	}
	s.Segment = NewBasicSegment("HNSHA", number, 1, s)
	return s
}

type SignatureEndSegment struct {
	Segment
	SecurityControlRef *dataelement.AlphaNumericDataElement
	Signature          *dataelement.BinaryDataElement
	PinTan             *dataelement.PinTanDataElement
}

func (s *SignatureEndSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.Signature = dataelement.NewBinary(signature, 512)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.PinTan = dataelement.NewPinTan(pin, tan)
}
