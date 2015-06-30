package dataelement

import (
	"fmt"
	"time"

	"github.com/mitch000001/go-hbci/domain"
)

const (
	SecurityHolderMessageSender   = "MS"
	SecurityHolderMessageReceiver = "MR"
)

func NewRDHSecurityIdentificationDataElement(securityHolder, clientSystemId string) *SecurityIdentificationDataElement {
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
		ClientSystemID: NewIdentification(clientSystemId),
	}
	s.DataElement = NewDataElementGroup(SecurityIdentificationDEG, 3, s)
	return s
}

type SecurityIdentificationDataElement struct {
	DataElement
	// Bezeichner für Sicherheitspartei
	SecurityHolder *AlphaNumericDataElement
	CID            *BinaryDataElement
	ClientSystemID *IdentificationDataElement
}

func (s *SecurityIdentificationDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.SecurityHolder,
		s.CID,
		s.ClientSystemID,
	}
}

const (
	SecurityTimestamp         = "STS"
	CertificateRevocationTime = "CRT"
)

func NewSecurityDateDataElement(dateId string, date time.Time) *SecurityDateDataElement {
	var id string
	if dateId == SecurityTimestamp {
		id = "1"
	} else if dateId == CertificateRevocationTime {
		id = "6"
	} else {
		panic(fmt.Errorf("DateIdentifier must be 'STS' or 'CRT'"))
	}
	s := &SecurityDateDataElement{
		DateIdentifier: NewAlphaNumeric(id, 3),
		Date:           NewDate(date),
		Time:           NewTime(date),
	}
	s.DataElement = NewDataElementGroup(SecurityDateDEG, 3, s)
	return s
}

type SecurityDateDataElement struct {
	DataElement
	DateIdentifier *AlphaNumericDataElement
	Date           *DateDataElement
	Time           *TimeDataElement
}

func (s *SecurityDateDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.DateIdentifier,
		s.Date,
		s.Time,
	}
}

func NewDefaultHashAlgorithmDataElement() *HashAlgorithmDataElement {
	h := &HashAlgorithmDataElement{
		Usage:            NewAlphaNumeric("1", 3),
		Algorithm:        NewAlphaNumeric("999", 3),
		AlgorithmParamId: NewAlphaNumeric("1", 3),
	}
	h.DataElement = NewDataElementGroup(HashAlgorithmDEG, 4, h)
	return h
}

type HashAlgorithmDataElement struct {
	DataElement
	// "1" for OHA, Owner Hashing
	Usage *AlphaNumericDataElement
	// "999" for ZZZ (RIPEMD-160)
	Algorithm *AlphaNumericDataElement
	// "1" for IVC, Initialization value, clear text
	AlgorithmParamId *AlphaNumericDataElement
	// may not be used in versions 2.20 and below
	AlgorithmParamValue *BinaryDataElement
}

func (h *HashAlgorithmDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		h.Usage,
		h.Algorithm,
		h.AlgorithmParamId,
		h.AlgorithmParamValue,
	}
}

func NewRDHSignatureAlgorithmDataElement() *SignatureAlgorithmDataElement {
	s := &SignatureAlgorithmDataElement{
		Usage:         NewAlphaNumeric("6", 3),
		Algorithm:     NewAlphaNumeric("10", 3),
		OperationMode: NewAlphaNumeric("16", 3),
	}
	s.DataElement = NewDataElementGroup(SignatureAlgorithmDEG, 3, s)
	return s
}

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

func (s *SignatureAlgorithmDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.Usage,
		s.Algorithm,
		s.OperationMode,
	}
}

func NewKeyNameDataElement(keyName domain.KeyName) *KeyNameDataElement {
	a := &KeyNameDataElement{
		Bank:       NewBankIndentification(keyName.BankID),
		UserID:     NewIdentification(keyName.UserID),
		KeyType:    NewAlphaNumeric(keyName.KeyType, 1),
		KeyNumber:  NewNumber(keyName.KeyNumber, 3),
		KeyVersion: NewNumber(keyName.KeyVersion, 3),
	}
	a.DataElement = NewDataElementGroup(KeyNameDEG, 5, a)
	return a
}

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

func (k *KeyNameDataElement) Val() domain.KeyName {
	return domain.KeyName{
		BankID:     domain.BankId{k.Bank.CountryCode.Val(), k.Bank.BankID.Val()},
		UserID:     k.UserID.Val(),
		KeyType:    k.KeyType.Val(),
		KeyNumber:  k.KeyNumber.Val(),
		KeyVersion: k.KeyVersion.Val(),
	}
}

func (k *KeyNameDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		k.Bank,
		k.UserID,
		k.KeyType,
		k.KeyNumber,
		k.KeyVersion,
	}
}

func NewCertificateDataElement(typ int, certificate []byte) *CertificateDataElement {
	c := &CertificateDataElement{
		CertificateType: NewNumber(typ, 1),
		Content:         NewBinary(certificate, 2048),
	}
	c.DataElement = NewDataElementGroup(CertificateDEG, 2, c)
	return c
}

type CertificateDataElement struct {
	DataElement
	// "1" for ZKA
	// "2" for UN/EDIFACT
	// "3" for X.509
	CertificateType *NumberDataElement
	Content         *BinaryDataElement
}

func (c *CertificateDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		c.CertificateType,
		c.Content,
	}
}
