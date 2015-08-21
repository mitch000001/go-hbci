package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderV3{
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
	s.ClientSegment = NewBasicSegment(2, s)

	segment := &SignatureHeaderSegment{
		signatureHeaderSegment: s,
	}
	return segment
}

func NewRDHSignatureHeaderSegment(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderV3{
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
	s.ClientSegment = NewBasicSegment(2, s)

	segment := &SignatureHeaderSegment{
		signatureHeaderSegment: s,
	}
	return segment
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment SignatureHeaderSegment -segment_interface signatureHeaderSegment -segment_versions="SignatureHeaderV3:3:ClientSegment,SignatureHeaderSegmentV4:4:ClientSegment"

type SignatureHeaderSegment struct {
	signatureHeaderSegment
}

type signatureHeaderSegment interface {
	CommonSegment
	SetSecurityFunction(string)
}

func (s *SignatureHeaderSegment) SetSecurityFunction(securityFn string) {
	s.signatureHeaderSegment.SetSecurityFunction(securityFn)
}

type SignatureHeaderV3 struct {
	ClientSegment
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

func (s *SignatureHeaderV3) SetSecurityFunction(securityFn string) {
	s.SecurityFunction = element.NewAlphaNumeric(securityFn, 3)
}

func (s *SignatureHeaderV3) Version() int         { return 3 }
func (s *SignatureHeaderV3) ID() string           { return "HNSHK" }
func (s *SignatureHeaderV3) referencedId() string { return "" }
func (s *SignatureHeaderV3) sender() string       { return senderBoth }

func (s *SignatureHeaderV3) elements() []element.DataElement {
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
	s := &SignatureEndV1{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.ClientSegment = NewBasicSegment(number, s)

	segment := &SignatureEndSegment{
		signatureEndSegment: s,
	}
	return segment
}

type signatureEndSegment interface {
	ClientSegment
	SetSignature(signature []byte)
	SetPinTan(pin, tan string)
}

type SignatureEndSegment struct {
	signatureEndSegment
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.signatureEndSegment.SetSignature(signature)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.signatureEndSegment.SetPinTan(pin, tan)
}

type SignatureEndV1 struct {
	ClientSegment
	SecurityControlRef *element.AlphaNumericDataElement
	Signature          *element.BinaryDataElement
	PinTan             *element.PinTanDataElement
}

func (s *SignatureEndV1) Version() int         { return 1 }
func (s *SignatureEndV1) ID() string           { return "HNSHA" }
func (s *SignatureEndV1) referencedId() string { return "" }
func (s *SignatureEndV1) sender() string       { return senderBoth }

func (s *SignatureEndV1) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndV1) SetSignature(signature []byte) {
	s.Signature = element.NewBinary(signature, 512)
}

func (s *SignatureEndV1) SetPinTan(pin, tan string) {
	s.PinTan = element.NewPinTan(pin, tan)
}

func NewFINTS3PinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderSegmentV4{
		SecurityProfile:          element.NewPinTanSecurityProfile(1),
		SecurityFunction:         element.NewCode("999", 3, []string{"1", "2", "999"}),
		SecurityControlRef:       element.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: element.NewCode("1", 3, []string{"1", "2"}),
		SecuritySupplierRole:     element.NewCode("1", 3, []string{"1", "3", "4"}),
		SecurityID:               element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
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

type SignatureHeaderSegmentV4 struct {
	ClientSegment
	SecurityProfile *element.SecurityProfileDataElement
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	// "999" for PIN/TAN
	SecurityFunction   *element.CodeDataElement
	SecurityControlRef *element.AlphaNumericDataElement
	// "1" for SHM (SignatureHeader and HBCI-Data)
	// "2" for SHT (SignatureHeader to SignatureEnd)
	SecurityApplicationRange *element.CodeDataElement
	// "1" for ISS, Herausgeber der signierten Nachricht (z.B. Erfasser oder Erstsignatur)
	// "3" for CON, der Unterzeichnete unterstützt den Inhalt der Nachricht (z.B. bei Zweitsignatur)
	// "4" for WIT, der Unterzeichnete ist Zeuge (z.B. Übermittler), aber für den Inhalt der Nachricht nicht verantwortlich)
	SecuritySupplierRole *element.CodeDataElement
	SecurityID           *element.SecurityIdentificationDataElement
	SecurityRefNumber    *element.NumberDataElement
	SecurityDate         *element.SecurityDateDataElement
	HashAlgorithm        *element.HashAlgorithmDataElement
	SignatureAlgorithm   *element.SignatureAlgorithmDataElement
	KeyName              *element.KeyNameDataElement
	Certificate          *element.CertificateDataElement
}

func (s *SignatureHeaderSegmentV4) SetSecurityFunction(securityFn string) {
	s.SecurityFunction = element.NewCode(securityFn, 3, []string{"1", "2", "999", securityFn})
	if securityFn == "999" {
		s.SecurityProfile = element.NewPinTanSecurityProfile(1)
	} else {
		s.SecurityProfile = element.NewPinTanSecurityProfile(2)
	}
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

func NewFINTS3SignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndV2{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.ClientSegment = NewBasicSegment(number, s)

	segment := &SignatureEndSegment{
		signatureEndSegment: s,
	}
	return segment
}

type SignatureEndV2 struct {
	ClientSegment
	SecurityControlRef *element.AlphaNumericDataElement
	Signature          *element.BinaryDataElement
	PinTan             *element.PinTanDataElement
}

func (s *SignatureEndV2) Version() int         { return 1 }
func (s *SignatureEndV2) ID() string           { return "HNSHA" }
func (s *SignatureEndV2) referencedId() string { return "" }
func (s *SignatureEndV2) sender() string       { return senderBoth }

func (s *SignatureEndV2) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndV2) SetSignature(signature []byte) {
	s.Signature = element.NewBinary(signature, 512)
}

func (s *SignatureEndV2) SetPinTan(pin, tan string) {
	s.PinTan = element.NewFINTSPinTan(pin, tan)
}
