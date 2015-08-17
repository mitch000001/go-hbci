package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPinTanEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment {
	e := &EncryptionHeaderV2{
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

func NewEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment {
	e := &EncryptionHeaderV2{
		SecurityFunction:     element.NewAlphaNumeric("4", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewRDHEncryptionAlgorithm(key),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e.ClientSegment = NewBasicSegment(998, e)

	segment := &EncryptionHeaderSegment{
		Segment: e,
	}
	return segment
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment EncryptionHeaderSegment

type EncryptionHeaderSegment struct {
	Segment
}

type EncryptionHeaderV2 struct {
	ClientSegment
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

func (e *EncryptionHeaderV2) Version() int         { return 2 }
func (e *EncryptionHeaderV2) ID() string           { return "HNVSK" }
func (e *EncryptionHeaderV2) referencedId() string { return "" }
func (e *EncryptionHeaderV2) sender() string       { return senderBoth }

func (e *EncryptionHeaderV2) elements() []element.DataElement {
	return []element.DataElement{
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
