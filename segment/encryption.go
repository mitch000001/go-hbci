package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPinTanEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     element.NewAlphaNumeric("998", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewPinTanEncryptionAlgorithm(),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e := &EncryptionHeaderSegment{
		version: v2,
	}
	e.Segment = NewBasicSegment("HNVSK", 998, 2, e)
	return e
}

func NewEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     element.NewAlphaNumeric("4", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewRDHEncryptionAlgorithm(key),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e := &EncryptionHeaderSegment{
		version: v2,
	}
	e.Segment = NewBasicSegment("HNVSK", 998, 2, e)
	return e
}

type EncryptionHeaderSegment struct {
	Segment
	version
}

func (e *EncryptionHeaderSegment) elements() []element.DataElement {
	return e.version.versionedElements()
}

type EncryptionHeaderVersion2 struct {
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

func (e *EncryptionHeaderVersion2) version() int {
	return 2
}

func (e *EncryptionHeaderVersion2) versionedElements() []element.DataElement {
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

func NewEncryptedDataSegment(encryptedData []byte) *EncryptedDataSegment {
	e := &EncryptedDataSegment{
		Data: element.NewBinary(encryptedData, -1),
	}
	e.Segment = NewBasicSegment("HNVSD", 999, 1, e)
	return e
}

type EncryptedDataSegment struct {
	Segment
	Data *element.BinaryDataElement
}

func (e *EncryptedDataSegment) elements() []element.DataElement {
	return []element.DataElement{
		e.Data,
	}
}
