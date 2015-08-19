package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

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
