package segment

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPinTanEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment {
	e := &EncryptionHeaderSegment{
		SecurityFunction:     element.NewAlphaNumeric("998", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewPinTanEncryptionAlgorithm(),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e.Segment = NewBasicSegment(998, e)
	return e
}

func NewEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment {
	e := &EncryptionHeaderSegment{
		SecurityFunction:     element.NewAlphaNumeric("4", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewRDHEncryptionAlgorithm(key),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e.Segment = NewBasicSegment(998, e)
	return e
}

type EncryptionHeaderSegment struct {
	Segment
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

func (e *EncryptionHeaderSegment) init() {
	*e.SecurityFunction = *new(element.AlphaNumericDataElement)
	*e.SecuritySupplierRole = *new(element.AlphaNumericDataElement)
	*e.SecurityID = *new(element.SecurityIdentificationDataElement)
	*e.SecurityDate = *new(element.SecurityDateDataElement)
	*e.EncryptionAlgorithm = *new(element.EncryptionAlgorithmDataElement)
	*e.KeyName = *new(element.KeyNameDataElement)
	*e.CompressionFunction = *new(element.AlphaNumericDataElement)
	*e.Certificate = *new(element.CertificateDataElement)
}
func (e *EncryptionHeaderSegment) version() int         { return 2 }
func (e *EncryptionHeaderSegment) id() string           { return "HNVSK" }
func (e *EncryptionHeaderSegment) referencedId() string { return "" }
func (e *EncryptionHeaderSegment) sender() string       { return senderBoth }

func (e *EncryptionHeaderSegment) elements() []element.DataElement {
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
	e.Segment = NewBasicSegment(999, e)
	return e
}

type EncryptedDataSegment struct {
	Segment
	Data *element.BinaryDataElement
}

func (e *EncryptedDataSegment) init() {
	*e.Data = *new(element.BinaryDataElement)
}
func (e *EncryptedDataSegment) version() int         { return 1 }
func (e *EncryptedDataSegment) id() string           { return "HNVSD" }
func (e *EncryptedDataSegment) referencedId() string { return "" }
func (e *EncryptedDataSegment) sender() string       { return senderBoth }

func (e *EncryptedDataSegment) UnmarshalHBCI(value []byte) error {
	elements := bytes.SplitN(value, []byte("+"), 2)
	header := elements[0]
	headerElems := bytes.Split(header, []byte(":"))
	num, err := strconv.Atoi(string(headerElems[1]))
	if err != nil {
		return fmt.Errorf("Malformed segment header")
	}
	e.Segment = NewBasicSegment(num, e)
	encryptedData := elements[1]
	e.Data = &element.BinaryDataElement{}
	err = e.Data.UnmarshalHBCI(encryptedData)
	if err != nil {
		return fmt.Errorf("Error while unmarshaling encrypted data: %v", err)
	}
	return nil
}

func (e *EncryptedDataSegment) elements() []element.DataElement {
	return []element.DataElement{
		e.Data,
	}
}
