package hbci

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"time"

	"github.com/mitch000001/go-hbci/dataelement"
	"golang.org/x/crypto/ripemd160"
)

type SignatureProvider interface {
	SetClientSystemID(clientSystemId string)
	SignMessage(SignedHBCIMessage) error
	NewSignatureHeader(controlReference string, signatureId int) *SignatureHeaderSegment
}

const initializationVector = "\x01\x23\x45\x67\x89\xAB\xCD\xEF\xFE\xDC\xBA\x98\x76\x54\x32\x10\xF0\xE1\xD2\xC3"

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

func NewPinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName dataelement.KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         dataelement.NewAlphaNumericDataElement("999", 3),
		SecurityControlRef:       dataelement.NewAlphaNumericDataElement(controlReference, 14),
		SecurityApplicationRange: dataelement.NewAlphaNumericDataElement("1", 3),
		SecuritySupplierRole:     dataelement.NewAlphaNumericDataElement("1", 3),
		SecurityID:               dataelement.NewRDHSecurityIdentificationDataElement(dataelement.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        dataelement.NewNumberDataElement(0, 16),
		SecurityDate:             dataelement.NewSecurityDateDataElement(dataelement.SecurityTimestamp, time.Now()),
		HashAlgorithm:            dataelement.NewDefaultHashAlgorithmDataElement(),
		SignatureAlgorithm:       dataelement.NewRDHSignatureAlgorithmDataElement(),
		KeyName:                  dataelement.NewKeyNameDataElement(keyName),
	}
	s := &SignatureHeaderSegment{
		version: v3,
	}
	s.Segment = NewBasicSegment("HNSHK", 2, 3, s)
	return s
}

func NewRDHSignatureHeaderSegment(controlReference string, signatureId int, clientSystemId string, keyName dataelement.KeyName) *SignatureHeaderSegment {
	v3 := &SignatureHeaderVersion3{
		SecurityFunction:         dataelement.NewAlphaNumericDataElement("1", 3),
		SecurityControlRef:       dataelement.NewAlphaNumericDataElement(controlReference, 14),
		SecurityApplicationRange: dataelement.NewAlphaNumericDataElement("1", 3),
		SecuritySupplierRole:     dataelement.NewAlphaNumericDataElement("1", 3),
		SecurityID:               dataelement.NewRDHSecurityIdentificationDataElement(dataelement.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        dataelement.NewNumberDataElement(signatureId, 16),
		SecurityDate:             dataelement.NewSecurityDateDataElement(dataelement.SecurityTimestamp, time.Now()),
		HashAlgorithm:            dataelement.NewDefaultHashAlgorithmDataElement(),
		SignatureAlgorithm:       dataelement.NewRDHSignatureAlgorithmDataElement(),
		KeyName:                  dataelement.NewKeyNameDataElement(keyName),
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

func (s *SignatureHeaderSegment) elements() []dataelement.DataElement {
	return s.version.versionedElements()
}

type SignatureHeaderVersion3 struct {
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	// "999" for PIN/TAN
	SecurityFunction   *dataelement.AlphaNumericDataElement
	SecurityControlRef *dataelement.AlphaNumericDataElement
	// "1" for SHM (SignatureHeader and HBCI-Data)
	// "2" for SHT (SignatureHeader to SignatureEnd)
	SecurityApplicationRange *dataelement.AlphaNumericDataElement
	// "1" for ISS, Herausgeber der signierten Nachricht (z.B. Erfasser oder Erstsignatur)
	// "3" for CON, der Unterzeichnete unterstützt den Inhalt der Nachricht (z.B. bei Zweitsignatur)
	// "4" for WIT, der Unterzeichnete ist Zeuge (z.B. Übermittler), aber für den Inhalt der Nachricht nicht verantwortlich)
	SecuritySupplierRole *dataelement.AlphaNumericDataElement
	SecurityID           *dataelement.SecurityIdentificationDataElement
	SecurityRefNumber    *dataelement.NumberDataElement
	SecurityDate         *dataelement.SecurityDateDataElement
	HashAlgorithm        *dataelement.HashAlgorithmDataElement
	SignatureAlgorithm   *dataelement.SignatureAlgorithmDataElement
	KeyName              *dataelement.KeyNameDataElement
	Certificate          *dataelement.CertificateDataElement
}

func (s SignatureHeaderVersion3) version() int {
	return 3
}

func (s *SignatureHeaderVersion3) versionedElements() []dataelement.DataElement {
	return []dataelement.DataElement{
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
		SecurityControlRef: dataelement.NewAlphaNumericDataElement(controlReference, 14),
	}
	s.Segment = NewBasicSegment("HNSHA", number, 1, s)
	return s
}

type SignatureEndSegment struct {
	Segment
	SecurityControlRef *dataelement.AlphaNumericDataElement
	Signature          *dataelement.BinaryDataElement
	PinTan             *dataelement.PinTanDataElement
}

func (s *SignatureEndSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.Signature = dataelement.NewBinaryDataElement(signature, 512)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.PinTan = dataelement.NewPinTanDataElement(pin, tan)
}
