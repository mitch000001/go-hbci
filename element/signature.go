package element

import (
	"fmt"
	"time"

	"github.com/mitch000001/go-hbci/domain"
)

const (
	// SecurityHolderMessageSender represents the MessageSender as security holder
	SecurityHolderMessageSender = "MS"
	// SecurityHolderMessageReceiver represents the MessageReceiver as security holder
	SecurityHolderMessageReceiver = "MR"
)

// NewRDHSecurityIdentification returns a new SecurityIdentificationDataElement
func NewRDHSecurityIdentification(securityHolder, clientSystemID string) *SecurityIdentificationDataElement {
	var holder string
	if securityHolder == SecurityHolderMessageSender {
		holder = "1"
	} else if securityHolder == SecurityHolderMessageReceiver {
		holder = "2"
	} else {
		panic(fmt.Errorf("SecurityHolder must be 'MS' or 'MR'"))
	}
	s := &SecurityIdentificationDataElement{
		SecurityHolder: NewAlphaNumeric(holder, 3),
		ClientSystemID: NewIdentification(clientSystemID),
	}
	s.DataElement = NewDataElementGroup(securityIdentificationDEG, 3, s)
	return s
}

// SecurityIdentificationDataElement represents a security method for wire transfer
type SecurityIdentificationDataElement struct {
	DataElement
	// Bezeichner für Sicherheitspartei
	SecurityHolder *AlphaNumericDataElement
	CID            *BinaryDataElement
	ClientSystemID *IdentificationDataElement
}

// UnmarshalHBCI unmarshals value into the DataElement
func (a *SecurityIdentificationDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 3 {
		return fmt.Errorf("Malformed marshaled value")
	}
	a.DataElement = NewDataElementGroup(securityIdentificationDEG, 3, a)
	if len(elements) > 0 && len(elements[0]) > 0 {
		a.SecurityHolder = &AlphaNumericDataElement{}
		err = a.SecurityHolder.UnmarshalHBCI(elements[0])
		if err != nil {
			return err
		}
	}
	if len(elements) > 1 && len(elements[1]) > 0 {
		a.CID = &BinaryDataElement{}
		err = a.CID.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		a.ClientSystemID = &IdentificationDataElement{}
		err = a.ClientSystemID.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	return nil
}

// GroupDataElements returns the grouped DataElements
func (s *SecurityIdentificationDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.SecurityHolder,
		s.CID,
		s.ClientSystemID,
	}
}

const (
	// SecurityTimestamp defines the type of the SecurityDate
	SecurityTimestamp = "STS"
	// CertificateRevocationTime defines the type of the SecurityDate
	CertificateRevocationTime = "CRT"
)

// NewSecurityDate creates a new SecurityDate for the given type
func NewSecurityDate(dateID string, date time.Time) *SecurityDateDataElement {
	var id string
	if dateID == SecurityTimestamp {
		id = "1"
	} else if dateID == CertificateRevocationTime {
		id = "6"
	} else {
		panic(fmt.Errorf("DateIdentifier must be 'STS' or 'CRT'"))
	}
	s := &SecurityDateDataElement{
		DateIdentifier: NewAlphaNumeric(id, 3),
		Date:           NewDate(date),
		Time:           NewTime(date),
	}
	s.DataElement = NewDataElementGroup(securityDateDEG, 3, s)
	return s
}

// SecurityDateDataElement represents a date with a context type
type SecurityDateDataElement struct {
	DataElement
	DateIdentifier *AlphaNumericDataElement
	Date           *DateDataElement
	Time           *TimeDataElement
}

// GroupDataElements returns the grouped DataElements
func (s *SecurityDateDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.DateIdentifier,
		s.Date,
		s.Time,
	}
}

// NewDefaultHashAlgorithm creates a default HashAlgorithmDataElement with
// values ready to use for initial dialog comm
func NewDefaultHashAlgorithm() *HashAlgorithmDataElement {
	h := &HashAlgorithmDataElement{
		Usage:            NewAlphaNumeric("1", 3),
		Algorithm:        NewAlphaNumeric("999", 3),
		AlgorithmParamID: NewAlphaNumeric("1", 3),
	}
	h.DataElement = NewDataElementGroup(hashAlgorithmDEG, 4, h)
	return h
}

// HashAlgorithmDataElement defines a hash algorithm
type HashAlgorithmDataElement struct {
	DataElement
	// "1" for OHA, Owner Hashing
	Usage *AlphaNumericDataElement
	// "999" for ZZZ (RIPEMD-160)
	Algorithm *AlphaNumericDataElement
	// "1" for IVC, Initialization value, clear text
	AlgorithmParamID *AlphaNumericDataElement
	// may not be used in versions 2.20 and below
	AlgorithmParamValue *BinaryDataElement
}

// GroupDataElements returns the grouped DataElements
func (h *HashAlgorithmDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		h.Usage,
		h.Algorithm,
		h.AlgorithmParamID,
		h.AlgorithmParamValue,
	}
}

// NewRDHSignatureAlgorithm creates a SignatureAlgorithm ready to use for RDH
func NewRDHSignatureAlgorithm() *SignatureAlgorithmDataElement {
	s := &SignatureAlgorithmDataElement{
		Usage:         NewAlphaNumeric("6", 3),
		Algorithm:     NewAlphaNumeric("10", 3),
		OperationMode: NewAlphaNumeric("16", 3),
	}
	s.DataElement = NewDataElementGroup(signatureAlgorithmDEG, 3, s)
	return s
}

// A SignatureAlgorithmDataElement represents a signature algorithm
type SignatureAlgorithmDataElement struct {
	DataElement
	// "1" for OSG, Owner Signing
	Usage *AlphaNumericDataElement
	// "1" for DES (DDV)
	// "10" for RSA (RDH)
	Algorithm *AlphaNumericDataElement
	// "16" for DSMR, Digital Signature Scheme giving Message Recovery: ISO 9796 (RDH)
	// "999" for ZZZ (DDV)
	OperationMode *AlphaNumericDataElement
}

// GroupDataElements returns the grouped DataElements
func (s *SignatureAlgorithmDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.Usage,
		s.Algorithm,
		s.OperationMode,
	}
}

// NewKeyName creates a new KeyNameDataElement for keyName
func NewKeyName(keyName domain.KeyName) *KeyNameDataElement {
	a := &KeyNameDataElement{
		Bank:       NewBankIdentification(keyName.BankID),
		UserID:     NewIdentification(keyName.UserID),
		KeyType:    NewAlphaNumeric(keyName.KeyType, 1),
		KeyNumber:  NewNumber(keyName.KeyNumber, 3),
		KeyVersion: NewNumber(keyName.KeyVersion, 3),
	}
	a.DataElement = NewDataElementGroup(keyNameDEG, 5, a)
	return a
}

// KeyNameDataElement represents metadata for keys
type KeyNameDataElement struct {
	DataElement
	Bank   *BankIdentificationDataElement
	UserID *IdentificationDataElement
	// "S" for Signing key
	// "V" for Encryption key
	KeyType    *AlphaNumericDataElement
	KeyNumber  *NumberDataElement
	KeyVersion *NumberDataElement
}

// Val returns the KeyName as domain.KeyName
func (k *KeyNameDataElement) Val() domain.KeyName {
	return domain.KeyName{
		BankID: domain.BankID{
			CountryCode: k.Bank.CountryCode.Val(),
			ID:          k.Bank.BankID.Val()},
		UserID:     k.UserID.Val(),
		KeyType:    k.KeyType.Val(),
		KeyNumber:  k.KeyNumber.Val(),
		KeyVersion: k.KeyVersion.Val(),
	}
}

// GroupDataElements returns the grouped DataElements
func (k *KeyNameDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		k.Bank,
		k.UserID,
		k.KeyType,
		k.KeyNumber,
		k.KeyVersion,
	}
}

// NewCertificate embodies a certificate into a DataElement
func NewCertificate(typ int, certificate []byte) *CertificateDataElement {
	c := &CertificateDataElement{
		CertificateType: NewNumber(typ, 1),
		Content:         NewBinary(certificate, 2048),
	}
	c.DataElement = NewDataElementGroup(certificateDEG, 2, c)
	return c
}

// CertificateDataElement embodies certificate bytes into a DataElement
type CertificateDataElement struct {
	DataElement
	// "1" for ZKA
	// "2" for UN/EDIFACT
	// "3" for X.509
	CertificateType *NumberDataElement
	Content         *BinaryDataElement
}

// GroupDataElements returns the grouped DataElements
func (c *CertificateDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		c.CertificateType,
		c.Content,
	}
}
