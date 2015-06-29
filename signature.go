package hbci

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"
	"time"

	"golang.org/x/crypto/ripemd160"
)

type SignatureProvider interface {
	SetClientSystemID(clientSystemId string)
	SignMessage(SignedHBCIMessage) error
	NewSignatureHeader(controlReference string, signatureId int) *SignatureHeaderSegment
}

const initializationVector = "\x01\x23\x45\x67\x89\xAB\xCD\xEF\xFE\xDC\xBA\x98\x76\x54\x32\x10\xF0\xE1\xD2\xC3"

func GenerateSigningKey() (*PublicKey, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 768)
	if err != nil {
		return nil, err
	}
	p := PublicKey{
		Type:          "S",
		Modulus:       rsaKey.N.Bytes(),
		Exponent:      big.NewInt(int64(rsaKey.E)).Bytes(),
		rsaPrivateKey: rsaKey,
	}
	return &p, nil
}

func MessageHashSum(message string) []byte {
	h := ripemd160.New()
	//io.WriteString(h, initializationVector)
	io.WriteString(h, message)
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

func NewPinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         NewAlphaNumericDataElement("999", 3),
		SecurityControlRef:       NewAlphaNumericDataElement(controlReference, 14),
		SecurityApplicationRange: NewAlphaNumericDataElement("1", 3),
		SecuritySupplierRole:     NewAlphaNumericDataElement("1", 3),
		SecurityID:               NewRDHSecurityIdentificationDataElement(SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        NewNumberDataElement(0, 16),
		SecurityDate:             NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		HashAlgorithm:            NewDefaultHashAlgorithmDataElement(),
		SignatureAlgorithm:       NewRDHSignatureAlgorithmDataElement(),
		KeyName:                  NewKeyNameDataElement(keyName),
	}
	s := &SignatureHeaderSegment{
		version: v3,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

func NewRDHSignatureHeaderSegment(controlReference string, signatureId int, clientSystemId string, keyName KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         NewAlphaNumericDataElement("1", 3),
		SecurityControlRef:       NewAlphaNumericDataElement(controlReference, 14),
		SecurityApplicationRange: NewAlphaNumericDataElement("1", 3),
		SecuritySupplierRole:     NewAlphaNumericDataElement("1", 3),
		SecurityID:               NewRDHSecurityIdentificationDataElement(SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        NewNumberDataElement(signatureId, 16),
		SecurityDate:             NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		HashAlgorithm:            NewDefaultHashAlgorithmDataElement(),
		SignatureAlgorithm:       NewRDHSignatureAlgorithmDataElement(),
		KeyName:                  NewKeyNameDataElement(keyName),
	}
	s := &SignatureHeaderSegment{
		version: v3,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

type SignatureHeaderSegment struct {
	Segment
	version
}

func (s *SignatureHeaderSegment) elements() []DataElement {
	return s.version.versionedElements()
}

type SignatureHeaderVersion3 struct {
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	// "999" for PIN/TAN
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

func (s SignatureHeaderVersion3) version() int {
	return 3
}

func (s *SignatureHeaderVersion3) versionedElements() []DataElement {
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

func NewSignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: NewAlphaNumericDataElement(controlReference, 14),
	}
	s.Segment = NewBasicSegment("HNSHA", number, 1, s)
	return s
}

type SignatureEndSegment struct {
	Segment
	SecurityControlRef *AlphaNumericDataElement
	Signature          *BinaryDataElement
	PinTan             *PinTanDataElement
}

func (s *SignatureEndSegment) elements() []DataElement {
	return []DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.Signature = NewBinaryDataElement(signature, 512)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.PinTan = NewPinTanDataElement(pin, tan)
}

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
		SecurityHolder: NewAlphaNumericDataElement(holder, 3),
		ClientSystemID: NewIdentificationDataElement(clientSystemId),
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
		DateIdentifier: NewAlphaNumericDataElement(id, 3),
		Date:           NewDateDataElement(date),
		Time:           NewTimeDataElement(date),
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
		Usage:            NewAlphaNumericDataElement("1", 3),
		Algorithm:        NewAlphaNumericDataElement("999", 3),
		AlgorithmParamId: NewAlphaNumericDataElement("1", 3),
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
		Usage:         NewAlphaNumericDataElement("6", 3),
		Algorithm:     NewAlphaNumericDataElement("10", 3),
		OperationMode: NewAlphaNumericDataElement("16", 3),
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

func NewPinTanKeyName(bankId BankId, userId string, keyType string) *KeyName {
	return &KeyName{
		BankID:     bankId,
		UserID:     userId,
		KeyType:    keyType,
		KeyNumber:  0,
		KeyVersion: 0,
	}
}

func NewInitialKeyName(countryCode int, bankId, userId string, keyType string) *KeyName {
	return &KeyName{
		BankID:     BankId{CountryCode: countryCode, ID: bankId},
		UserID:     userId,
		KeyType:    keyType,
		KeyNumber:  999,
		KeyVersion: 999,
	}
}

type KeyName struct {
	BankID     BankId
	UserID     string
	KeyType    string
	KeyNumber  int
	KeyVersion int
}

func (k *KeyName) IsInitial() bool {
	return k.KeyNumber == 999 && k.KeyVersion == 999
}

func (k *KeyName) SetInitial() {
	k.KeyNumber = 999
	k.KeyVersion = 999
}

func NewKeyNameDataElement(keyName KeyName) *KeyNameDataElement {
	a := &KeyNameDataElement{
		Bank:       NewBankIndentificationDataElement(keyName.BankID),
		UserID:     NewIdentificationDataElement(keyName.UserID),
		KeyType:    NewAlphaNumericDataElement(keyName.KeyType, 1),
		KeyNumber:  NewNumberDataElement(keyName.KeyNumber, 3),
		KeyVersion: NewNumberDataElement(keyName.KeyVersion, 3),
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

func (k *KeyNameDataElement) Val() KeyName {
	return KeyName{
		BankID:     BankId{k.Bank.CountryCode.Val(), k.Bank.BankID.Val()},
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
		CertificateType: NewNumberDataElement(typ, 1),
		Content:         NewBinaryDataElement(certificate, 2048),
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
