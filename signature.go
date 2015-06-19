package hbci

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"reflect"
	"time"

	"golang.org/x/crypto/ripemd160"
)

const initializationVector = "\x01\x23\x45\x67\x89\xAB\xCD\xEF\xFE\xDC\xBA\x98\x76\x54\x32\x10\xF0\xE1\xD2\xC3"

func GenerateSigningKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 1024)
}

func MessageHashSum(message fmt.Stringer) []byte {
	h := ripemd160.New()
	io.WriteString(h, initializationVector)
	io.WriteString(h, message.String())
	return h.Sum(nil)
}

func SignMessageHash(messageHash []byte, key *rsa.PrivateKey) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, key, 0, messageHash)
}

func EncryptMessage(message fmt.Stringer) (string, error) {
	block, err := des.NewCipher([]byte(initializationVector))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, des.BlockSize+len(message.String()))
	iv := ciphertext[:des.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[des.BlockSize:], []byte(message.String()))

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return fmt.Sprintf("%x", ciphertext), nil
}

func NewSignatureHeaderSegment(controlReference string, signatureId int, securityHolder, holderId string, keyName KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderSegment{
		SecurityFunction:         NewAlphaNumericDataElement("1", 3),
		SecurityControlRef:       NewAlphaNumericDataElement(controlReference, 14),
		SecurityApplicationRange: NewAlphaNumericDataElement("1", 3),
		SecuritySupplierRole:     NewAlphaNumericDataElement("1", 3),
		SecurityID:               NewRDHSecurityIdentificationDataElement(securityHolder, holderId),
		SecurityRefNumber:        NewNumberDataElement(signatureId, 16),
		SecurityDate:             NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		HashAlgorithm:            NewDefaultHashAlgorithmDataElement(),
		SignatureAlgorithm:       NewRDHSignatureAlgorithmDataElement(),
		KeyName:                  NewKeyNameDataElement(keyName),
	}
	s.basicSegment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

type SignatureHeaderSegment struct {
	*basicSegment
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	SecurityFunction   *AlphaNumericDataElement
	SecurityControlRef *AlphaNumericDataElement
	// "1" for SHM (SignatureHeader and HBCI-Data)
	// "2" for SHT (SignatureHeader to SignatureEnd)
	SecurityApplicationRange *AlphaNumericDataElement
	// "1" for ISS, Herausgeber der signierten Nachricht (z.B. Erfasser oder Erstsignatur)
	// "3" for CON, der Unterzeichnete unterstützt den Inhalt der Nachricht (z.B. bei Zweitsignatur)
	// "4" for WIT, der Unterzeichnete ist Zeuge (z.B. Übermittler), aber für den Inhalt der Nachricht nicht verantwortlich)
	SecuritySupplierRole *AlphaNumericDataElement
	SecurityID           *SecurityIdentificationDataElement
	SecurityRefNumber    *NumberDataElement
	SecurityDate         *SecurityDateDataElement
	HashAlgorithm        *HashAlgorithmDataElement
	SignatureAlgorithm   *SignatureAlgorithmDataElement
	KeyName              *KeyNameDataElement
	Certificate          *CertificateDataElement
}

func (s *SignatureHeaderSegment) elements() []DataElement {
	return []DataElement{
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

func NewSignatureEndSegment(number int, controlReference string, signature []byte) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: NewAlphaNumericDataElement(controlReference, 14),
		Signature:          NewBinaryDataElement(signature, 512),
	}
	s.basicSegment = NewBasicSegment("HNSHA", number, 1, s)
	return s
}

type SignatureEndSegment struct {
	*basicSegment
	SecurityControlRef *AlphaNumericDataElement
	Signature          *BinaryDataElement
}

func (s *SignatureEndSegment) elements() []DataElement {
	return []DataElement{
		s.SecurityControlRef,
		s.Signature,
	}
}

const (
	SecurityHolderMessageSender   = "MS"
	SecurityHolderMessageReceiver = "MR"
)

func NewRDHSecurityIdentificationDataElement(securityHolder, holderId string) *SecurityIdentificationDataElement {
	var holder string
	if securityHolder == SecurityHolderMessageSender {
		holder = "1"
	} else if securityHolder == SecurityHolderMessageReceiver {
		holder = "2"
	} else {
		panic(fmt.Errorf("SecurityHolder must be 'MS' or 'MR'"))
	}
	s := &SecurityIdentificationDataElement{
		SecurityHolder: NewAlphaNumericDataElement(holder, 3),
		HolderID:       NewIdentificationDataElement(holderId),
	}
	s.elementGroup = NewDataElementGroup(SecurityIdentificationDEG, 3, s)
	return s
}

type SecurityIdentificationDataElement struct {
	*elementGroup
	// Bezeichner für Sicherheitspartei
	SecurityHolder *AlphaNumericDataElement
	CID            *BinaryDataElement
	HolderID       *IdentificationDataElement
}

func (s *SecurityIdentificationDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.SecurityHolder,
		s.CID,
		s.HolderID,
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
		DateIdentifier: NewAlphaNumericDataElement(id, 3),
		Date:           NewDateDataElement(date),
		Time:           NewTimeDataElement(date),
	}
	s.elementGroup = NewDataElementGroup(SecurityDateDEG, 3, s)
	return s
}

type SecurityDateDataElement struct {
	*elementGroup
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
		Usage:            NewAlphaNumericDataElement("1", 3),
		Algorithm:        NewAlphaNumericDataElement("999", 3),
		AlgorithmParamId: NewAlphaNumericDataElement("1", 3),
	}
	h.elementGroup = NewDataElementGroup(HashAlgorithmDEG, 4, h)
	return h
}

type HashAlgorithmDataElement struct {
	*elementGroup
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
		Usage:         NewAlphaNumericDataElement("6", 3),
		Algorithm:     NewAlphaNumericDataElement("10", 3),
		OperationMode: NewAlphaNumericDataElement("16", 3),
	}
	s.elementGroup = NewDataElementGroup(SignatureAlgorithmDEG, 3, s)
	return s
}

type SignatureAlgorithmDataElement struct {
	*elementGroup
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

func NewCertificateDataElement(typ int, certificate []byte) *CertificateDataElement {
	c := &CertificateDataElement{
		CertificateType: NewNumberDataElement(typ, 1),
		Content:         NewBinaryDataElement(certificate, 2048),
	}
	c.elementGroup = NewDataElementGroup(CertificateDEG, 2, c)
	return c
}

func NewInitialKeyName(countryCode int, bankId, userId string, keyType string) KeyName {
	return KeyName{
		CountryCode: countryCode,
		BankID:      bankId,
		UserID:      userId,
		KeyType:     keyType,
		KeyNumber:   999,
		KeyVersion:  999,
	}
}

type KeyName struct {
	CountryCode int
	BankID      string
	UserID      string
	KeyType     string
	KeyNumber   int
	KeyVersion  int
}

func (k KeyName) IsInitial() bool {
	return k.KeyNumber == 999 && k.KeyVersion == 999
}

func NewKeyNameDataElement(keyName KeyName) *KeyNameDataElement {
	a := &KeyNameDataElement{
		Bank:       NewBankIndentificationDataElementWithBankId(keyName.CountryCode, keyName.BankID),
		UserID:     NewIdentificationDataElement(keyName.UserID),
		KeyType:    NewAlphaNumericDataElement(keyName.KeyType, 1),
		KeyNumber:  NewNumberDataElement(keyName.KeyNumber, 3),
		KeyVersion: NewNumberDataElement(keyName.KeyVersion, 3),
	}
	a.elementGroup = NewDataElementGroup(KeyNameDEG, 5, a)
	return a
}

type KeyNameDataElement struct {
	*elementGroup
	Bank       *BankIdentificationDataElement
	UserID     *IdentificationDataElement
	KeyType    *AlphaNumericDataElement
	KeyNumber  *NumberDataElement
	KeyVersion *NumberDataElement
}

func (k *KeyNameDataElement) Val() KeyName {
	return KeyName{
		CountryCode: k.Bank.CountryCode.Val(),
		BankID:      k.Bank.BankID.Val(),
		UserID:      k.UserID.Val(),
		KeyType:     k.KeyType.Val(),
		KeyNumber:   k.KeyNumber.Val(),
		KeyVersion:  k.KeyVersion.Val(),
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

type CertificateDataElement struct {
	*elementGroup
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

type PublicKey struct {
	Type     string
	Modulus  []byte
	Exponent []byte
}

func NewPublicKeyDataElement(pubKey *PublicKey) *PublicKeyDataElement {
	if !reflect.DeepEqual(pubKey.Exponent, []byte("65537")) {
		panic(fmt.Errorf("Exponent must equal 65537 (% X)", "65537"))
	}
	p := &PublicKeyDataElement{
		Usage:         NewAlphaNumericDataElement(pubKey.Type, 3),
		OperationMode: NewAlphaNumericDataElement("16", 3),
		Cipher:        NewAlphaNumericDataElement("10", 3),
		Modulus:       NewBinaryDataElement(pubKey.Modulus, 512),
		ModulusID:     NewAlphaNumericDataElement("12", 3),
		Exponent:      NewBinaryDataElement(pubKey.Exponent, 512),
		ExponentID:    NewAlphaNumericDataElement("13", 3),
	}
	p.elementGroup = NewDataElementGroup(PublicKeyDEG, 7, p)
	return p
}

type PublicKeyDataElement struct {
	*elementGroup
	// "5" for OCF, Owner Ciphering (Encryption key)
	// "6" for OSG, Owner Signing (Signing key)
	Usage *AlphaNumericDataElement
	// "16" for DSMR (ISO 9796)
	OperationMode *AlphaNumericDataElement
	// "10" for RSA
	Cipher  *AlphaNumericDataElement
	Modulus *BinaryDataElement
	// "12" for MOD, Modulus
	ModulusID *AlphaNumericDataElement
	// "65537"
	Exponent *BinaryDataElement
	// "13" for EXP, Exponent
	ExponentID *AlphaNumericDataElement
}

func (p *PublicKeyDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.Usage,
		p.OperationMode,
		p.Cipher,
		p.Modulus,
		p.ModulusID,
		p.Exponent,
		p.ExponentID,
	}
}

func (p *PublicKeyDataElement) Val() *PublicKey {
	return &PublicKey{
		Type:     p.Usage.Val(),
		Modulus:  p.Modulus.Val(),
		Exponent: p.Exponent.Val(),
	}
}
