package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewFINTS3PinTanEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment {
	e := &EncryptionHeaderSegmentV3{
		SecurityProfile:      element.NewPinTanSecurityProfile(),
		SecurityFunction:     element.NewAlphaNumeric("998", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewPinTanEncryptionAlgorithm(),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e.ClientSegment = NewBasicSegment(998, e)

	segment := &EncryptionHeaderSegment{
		Segment: e,
	}
	return segment
}

type EncryptionHeaderSegmentV3 struct {
	ClientSegment
	SecurityProfile *element.SecurityProfileDataElement
	// "4" for ENC, Encryption (encryption and eventually compression)
	// "998" for Cleartext
	SecurityFunction *element.AlphaNumericDataElement
	// "1" for ISS,  Herausgeber der chiffrierten Nachricht (Erfasser)
	// "4" for WIT, der Unterzeichnete ist Zeuge, aber für den Inhalt der
	// Nachricht nicht verantwortlich (Übermittler, welcher nicht Erfasser ist)
	SecuritySupplierRole *element.AlphaNumericDataElement
	SecurityID           *element.SecurityIdentificationDataElement
	SecurityDate         *element.SecurityDateDataElement
	EncryptionAlgorithm  *element.EncryptionAlgorithmDataElement
	KeyName              *element.KeyNameDataElement
	CompressionFunction  *element.AlphaNumericDataElement
	Certificate          *element.CertificateDataElement
}

func (e *EncryptionHeaderSegmentV3) Version() int         { return 3 }
func (e *EncryptionHeaderSegmentV3) ID() string           { return "HNVSK" }
func (e *EncryptionHeaderSegmentV3) referencedId() string { return "" }
func (e *EncryptionHeaderSegmentV3) sender() string       { return senderBoth }

func (e *EncryptionHeaderSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		e.SecurityProfile,
		e.SecurityFunction,
		e.SecuritySupplierRole,
		e.SecurityID,
		e.SecurityDate,
		e.EncryptionAlgorithm,
		e.KeyName,
		e.CompressionFunction,
		e.Certificate,
	}
}
