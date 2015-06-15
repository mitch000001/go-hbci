package hbci

import (
	"crypto/rand"
	"crypto/rsa"
	"time"
)

const encryptionInitializationVector = "\x00\x00\x00\x00\x00\x00\x00\x00"

func GenerateMessageKey() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateEncryptionKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 1024)
}

func NewEncryptionHeaderSegment(signatureId int, securityHolder, holderId string, keyName KeyName, key []byte) *EncryptionHeaderSegment {
	e := &EncryptionHeaderSegment{
		SecurityFunction:     NewAlphaNumericDataElement("4", 3),
		SecuritySupplierRole: NewAlphaNumericDataElement("1", 3),
		SecurityID:           NewRDHSecurityIdentificationDataElement(securityHolder, holderId),
		SecurityDate:         NewSecurityDateDataElement(SecurityDateIdentifierSecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  NewRDHEncryptionAlgorithmDataElement(key),
		CompressionFunction:  NewAlphaNumericDataElement("0", 3),
		KeyName:              NewKeyNameDataElement(keyName),
	}
	header := NewSegmentHeader("HNVSK", 2, 2)
	e.segment = NewSegment(header, e)
	return e
}

type EncryptionHeaderSegment struct {
	*segment
	// "4" for ENC, Encryption (encryption and eventually compression)
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

func (e *EncryptionHeaderSegment) DataElements() []DataElement {
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

func NewRDHEncryptionAlgorithmDataElement(pubKey []byte) *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumericDataElement("2", 3),
		OperationMode:              NewAlphaNumericDataElement("2", 3),
		Algorithm:                  NewAlphaNumericDataElement("13", 3),
		Key:                        NewBinaryDataElement(pubKey, 512),
		KeyParamID:                 NewAlphaNumericDataElement("6", 3),
		InitializationValueParamID: NewAlphaNumericDataElement("1", 3),
		InitializationValue:        NewBinaryDataElement([]byte(encryptionInitializationVector), 8),
	}
	return e
}

type EncryptionAlgorithmDataElement struct {
	*elementGroup
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

type EncryptedData struct{}
