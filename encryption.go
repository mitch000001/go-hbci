package hbci

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/mitch000001/go-hbci/dataelement"
)

type EncryptionProvider interface {
	SetClientSystemID(clientSystemId string)
	Encrypt(message []byte) (*EncryptedMessage, error)
	EncryptWithInitialKeyName(message []byte) (*EncryptedMessage, error)
}

const encryptionInitializationVector = "\x00\x00\x00\x00\x00\x00\x00\x00"

func GenerateMessageKey() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewEncryptedPinTanMessage(clientSystemId string, keyName dataelement.KeyName, encryptedMessage []byte) *EncryptedMessage {
	e := &EncryptedMessage{
		EncryptionHeader: NewPinTanEncryptionHeaderSegment(clientSystemId, keyName),
		EncryptedData:    NewEncryptedDataSegment(encryptedMessage),
	}
	e.basicMessage = newBasicMessage(e)
	return e
}

type EncryptedMessage struct {
	*basicMessage
	EncryptionHeader *EncryptionHeaderSegment
	EncryptedData    *EncryptedDataSegment
}

func (e *EncryptedMessage) HBCISegments() []Segment {
	return []Segment{
		e.EncryptionHeader,
		e.EncryptedData,
	}
}

func (e *EncryptedMessage) SetNumbers() {
	panic(fmt.Errorf("SetNumbers: Operation not allowed on encrypted messages"))
}

func NewPinTanEncryptionHeaderSegment(clientSystemId string, keyName dataelement.KeyName) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     dataelement.NewAlphaNumericDataElement("998", 3),
		SecuritySupplierRole: dataelement.NewAlphaNumericDataElement("1", 3),
		SecurityID:           dataelement.NewRDHSecurityIdentificationDataElement(dataelement.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         dataelement.NewSecurityDateDataElement(dataelement.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  dataelement.NewPinTanEncryptionAlgorithmDataElement(),
		KeyName:              dataelement.NewKeyNameDataElement(keyName),
		CompressionFunction:  dataelement.NewAlphaNumericDataElement("0", 3),
	}
	e := &EncryptionHeaderSegment{
		version: v2,
	}
	e.Segment = NewBasicSegment("HNVSK", 998, 2, e)
	return e
}

func NewEncryptionHeaderSegment(clientSystemId string, keyName dataelement.KeyName, key []byte) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     dataelement.NewAlphaNumericDataElement("4", 3),
		SecuritySupplierRole: dataelement.NewAlphaNumericDataElement("1", 3),
		SecurityID:           dataelement.NewRDHSecurityIdentificationDataElement(dataelement.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         dataelement.NewSecurityDateDataElement(dataelement.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  dataelement.NewRDHEncryptionAlgorithmDataElement(key),
		KeyName:              dataelement.NewKeyNameDataElement(keyName),
		CompressionFunction:  dataelement.NewAlphaNumericDataElement("0", 3),
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

func (e *EncryptionHeaderSegment) elements() []dataelement.DataElement {
	return e.version.versionedElements()
}

type EncryptionHeaderVersion2 struct {
	// "4" for ENC, Encryption (encryption and eventually compression)
	// "998" for Cleartext
	SecurityFunction *dataelement.AlphaNumericDataElement
	// "1" for ISS,  Herausgeber der chiffrierten Nachricht (Erfasser)
	// "4" for WIT, der Unterzeichnete ist Zeuge, aber für den Inhalt der
	// Nachricht nicht verantwortlich (Übermittler, welcher nicht Erfasser ist)
	SecuritySupplierRole *dataelement.AlphaNumericDataElement
	SecurityID           *dataelement.SecurityIdentificationDataElement
	SecurityDate         *dataelement.SecurityDateDataElement
	EncryptionAlgorithm  *dataelement.EncryptionAlgorithmDataElement
	KeyName              *dataelement.KeyNameDataElement
	CompressionFunction  *dataelement.AlphaNumericDataElement
	Certificate          *dataelement.CertificateDataElement
}

func (e *EncryptionHeaderVersion2) version() int {
	return 2
}

func (e *EncryptionHeaderVersion2) versionedElements() []dataelement.DataElement {
	return []dataelement.DataElement{
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
		Data: dataelement.NewBinaryDataElement(encryptedData, -1),
	}
	e.Segment = NewBasicSegment("HNVSD", 999, 1, e)
	return e
}

type EncryptedDataSegment struct {
	Segment
	Data *dataelement.BinaryDataElement
}

func (e *EncryptedDataSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		e.Data,
	}
}
