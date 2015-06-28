package hbci

import (
	"crypto/rand"
	"fmt"
	"time"
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

func NewEncryptedPinTanMessage(clientSystemId string, keyName KeyName, encryptedMessage []byte) *EncryptedMessage {
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

func NewPinTanEncryptionHeaderSegment(clientSystemId string, keyName KeyName) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     NewAlphaNumericDataElement("998", 3),
		SecuritySupplierRole: NewAlphaNumericDataElement("1", 3),
		SecurityID:           NewRDHSecurityIdentificationDataElement(SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  NewPinTanEncryptionAlgorithmDataElement(),
		KeyName:              NewKeyNameDataElement(keyName),
		CompressionFunction:  NewAlphaNumericDataElement("0", 3),
	}
	e := &EncryptionHeaderSegment{
		version: v2,
	}
	e.Segment = NewBasicSegment("HNVSK", 998, 2, e)
	return e
}

func NewEncryptionHeaderSegment(clientSystemId string, keyName KeyName, key []byte) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     NewAlphaNumericDataElement("4", 3),
		SecuritySupplierRole: NewAlphaNumericDataElement("1", 3),
		SecurityID:           NewRDHSecurityIdentificationDataElement(SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  NewRDHEncryptionAlgorithmDataElement(key),
		KeyName:              NewKeyNameDataElement(keyName),
		CompressionFunction:  NewAlphaNumericDataElement("0", 3),
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

func (e *EncryptionHeaderSegment) elements() []DataElement {
	return e.version.versionedElements()
}

type EncryptionHeaderVersion2 struct {
	// "4" for ENC, Encryption (encryption and eventually compression)
	// "998" for Cleartext
	SecurityFunction *AlphaNumericDataElement
	// "1" for ISS,  Herausgeber der chiffrierten Nachricht (Erfasser)
	// "4" for WIT, der Unterzeichnete ist Zeuge, aber für den Inhalt der
	// Nachricht nicht verantwortlich (Übermittler, welcher nicht Erfasser ist)
	SecuritySupplierRole *AlphaNumericDataElement
	SecurityID           *SecurityIdentificationDataElement
	SecurityDate         *SecurityDateDataElement
	EncryptionAlgorithm  *EncryptionAlgorithmDataElement
	KeyName              *KeyNameDataElement
	CompressionFunction  *AlphaNumericDataElement
	Certificate          *CertificateDataElement
}

func (e *EncryptionHeaderVersion2) version() int {
	return 2
}

func (e *EncryptionHeaderVersion2) versionedElements() []DataElement {
	return []DataElement{
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
		Data: NewBinaryDataElement(encryptedData, -1),
	}
	e.Segment = NewBasicSegment("HNVSD", 999, 1, e)
	return e
}

type EncryptedDataSegment struct {
	Segment
	Data *BinaryDataElement
}

func (e *EncryptedDataSegment) elements() []DataElement {
	return []DataElement{
		e.Data,
	}
}

func NewPinTanEncryptionAlgorithmDataElement() *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumericDataElement("2", 3),
		OperationMode:              NewAlphaNumericDataElement("2", 3),
		Algorithm:                  NewAlphaNumericDataElement("13", 3),
		Key:                        NewBinaryDataElement([]byte(defaultPinTan), 512),
		KeyParamID:                 NewAlphaNumericDataElement("5", 3),
		InitializationValueParamID: NewAlphaNumericDataElement("1", 3),
	}
	e.DataElement = NewDataElementGroup(EncryptionAlgorithmDEG, 7, e)
	return e
}

func NewRDHEncryptionAlgorithmDataElement(pubKey []byte) *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumericDataElement("2", 3),
		OperationMode:              NewAlphaNumericDataElement("2", 3),
		Algorithm:                  NewAlphaNumericDataElement("13", 3),
		Key:                        NewBinaryDataElement(pubKey, 512),
		KeyParamID:                 NewAlphaNumericDataElement("6", 3),
		InitializationValueParamID: NewAlphaNumericDataElement("1", 3),
	}
	e.DataElement = NewDataElementGroup(EncryptionAlgorithmDEG, 7, e)
	return e
}

type EncryptionAlgorithmDataElement struct {
	DataElement
	// "2" for OSY, Owner Symmetric
	Usage *AlphaNumericDataElement
	// "2" for CBC, Cipher Block Chaining.
	OperationMode *AlphaNumericDataElement
	// "13" for 2-Key-Triple-DES
	Algorithm *AlphaNumericDataElement
	Key       *BinaryDataElement
	// "5" for KYE, Symmetric key, en-/decryption with a symmetric key (DDV)
	// "6" for KYP, Symmetric key, encryption with a public key (RDH).
	KeyParamID                 *AlphaNumericDataElement
	InitializationValueParamID *AlphaNumericDataElement
	InitializationValue        *BinaryDataElement
}

func (e *EncryptionAlgorithmDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		e.Usage,
		e.OperationMode,
		e.Algorithm,
		e.Key,
		e.KeyParamID,
		e.InitializationValueParamID,
		e.InitializationValue,
	}
}
