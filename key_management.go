package hbci

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"
)

type Key interface {
	KeyName() KeyName
	SetKeyNumber(number int)
	SetKeyVersion(version int)
	Sign(message []byte) (signature []byte, err error)
	Encrypt(message []byte) (encrypted []byte, err error)
	CanSign() bool
	CanEncrypt() bool
}

func NewRSAKey(pubKey *PublicKey, keyName *KeyName) *RSAKey {
	return &RSAKey{PublicKey: pubKey, keyName: keyName}
}

type RSAKey struct {
	*PublicKey
	keyName *KeyName
}

func (r *RSAKey) KeyName() KeyName {
	return *r.keyName
}

func (r *RSAKey) SetKeyNumber(number int) {
	r.keyName.KeyNumber = number
}

func (r *RSAKey) SetKeyVersion(version int) {
	r.keyName.KeyVersion = version
}

func (r *RSAKey) CanSign() bool {
	return r.PublicKey.rsaPrivateKey != nil
}

func (r *RSAKey) CanEncrypt() bool {
	return r.PublicKey.rsaPublicKey != nil
}

func NewRDHSignatureProvider(signingKey *RSAKey) SignatureProvider {
	return &RDHSignatureProvider{signingKey: signingKey}
}

type RDHSignatureProvider struct {
	signingKey     *RSAKey
	clientSystemId string
}

func (r *RDHSignatureProvider) SetClientSystemID(clientSystemId string) {
	r.clientSystemId = clientSystemId
}

func (r *RDHSignatureProvider) SignMessage(message SignedHBCIMessage) error {
	marshaledMessage := message.SignatureBeginSegment().String()
	for _, segment := range message.HBCISegments() {
		marshaledMessage += segment.String()
	}
	hashSum := MessageHashSum(marshaledMessage)
	sig, err := r.signingKey.Sign(hashSum)
	if err != nil {
		return err
	}
	message.SignatureEndSegment().SetSignature(sig)
	return nil
}

func (p *RDHSignatureProvider) NewSignatureHeader(controlReference string, signatureId int) *SignatureHeaderSegment {
	return NewRDHSignatureHeaderSegment(controlReference, signatureId, p.clientSystemId, p.signingKey.KeyName())
}

func NewEncryptionKey(modulus, exponent []byte) *PublicKey {
	p := &PublicKey{
		Type: "V",
	}
	copy(p.Modulus, modulus)
	copy(p.Exponent, exponent)
	mod := new(big.Int).SetBytes(modulus)
	exp := new(big.Int).SetBytes(exponent)
	pubKey := rsa.PublicKey{
		N: mod,
		E: int(exp.Int64()),
	}
	p.rsaPublicKey = &pubKey
	return p
}

type PublicKey struct {
	Type          string
	Modulus       []byte
	Exponent      []byte
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

func (p *PublicKey) SigningKey() *rsa.PrivateKey {
	return p.rsaPrivateKey
}

func (p *PublicKey) Sign(message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, p.rsaPrivateKey, 0, message)
}

func (p *PublicKey) Encrypt(message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, p.rsaPublicKey, message)
}

const initialKeyVersion = 999

type InitialPublicKeyRenewalMessage struct {
	*basicMessage
	Identification             *IdentificationSegment
	PublicSigningKeyRequest    *PublicKeyRequestSegment
	PublicEncryptionKeyRequest *PublicKeyRequestSegment
}

func NewPublicKeyRenewalSegment(number int, keyName KeyName, pubKey *PublicKey) *PublicKeyRenewalSegment {
	if keyName.KeyType == "B" {
		panic(fmt.Errorf("KeyType may not be 'B'"))
	}
	p := &PublicKeyRenewalSegment{
		MessageID:  NewNumberDataElement(2, 1),
		FunctionID: NewNumberDataElement(112, 3),
		KeyName:    NewKeyNameDataElement(keyName),
		PublicKey:  NewPublicKeyDataElement(pubKey),
	}
	p.Segment = NewBasicSegment("HKSAK", number, 2, p)
	return p
}

type PublicKeyRenewalSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *NumberDataElement
	// "112" für ‘Certificate Replacement’ (Ersatz des Zertifikats))
	FunctionID *NumberDataElement
	// Key type may not equal 'B'
	KeyName     *KeyNameDataElement
	PublicKey   *PublicKeyDataElement
	Certificate *CertificateDataElement
}

func (p *PublicKeyRenewalSegment) elements() []DataElement {
	return []DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.PublicKey,
		p.Certificate,
	}
}

func NewPublicKeyRequestSegment(number int, keyName KeyName) *PublicKeyRequestSegment {
	p := &PublicKeyRequestSegment{
		MessageID:  NewNumberDataElement(2, 1),
		FunctionID: NewNumberDataElement(124, 3),
		KeyName:    NewKeyNameDataElement(keyName),
	}
	p.Segment = NewBasicSegment("HKISA", number, 2, p)
	return p
}

type PublicKeyRequestSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *NumberDataElement
	// "124" für ‘Certificate Status Request’
	FunctionID  *NumberDataElement
	KeyName     *KeyNameDataElement
	Certificate *CertificateDataElement
}

func (p *PublicKeyRequestSegment) elements() []DataElement {
	return []DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.Certificate,
	}
}

func NewPublicKeyTransmissionSegment(dialogId string, number int, messageReference int, keyName KeyName, pubKey *PublicKey, refSegment *PublicKeyRequestSegment) *PublicKeyTransmissionSegment {
	if messageReference <= 0 {
		panic(fmt.Errorf("Message Reference number must be greater 0"))
	}
	p := &PublicKeyTransmissionSegment{
		MessageID:  NewNumberDataElement(1, 1),
		DialogID:   NewIdentificationDataElement(dialogId),
		MessageRef: NewNumberDataElement(messageReference, 4),
		FunctionID: NewNumberDataElement(224, 3),
		KeyName:    NewKeyNameDataElement(keyName),
		PublicKey:  NewPublicKeyDataElement(pubKey),
	}
	header := NewReferencingSegmentHeader("HIISA", number, 2, refSegment.Header().Number.Val())
	p.Segment = NewBasicSegmentWithHeader(header, p)
	return p
}

type PublicKeyTransmissionSegment struct {
	Segment
	// "1" für ‘Key-Management-Nachricht ist Antwort’
	MessageID  *NumberDataElement
	DialogID   *IdentificationDataElement
	MessageRef *NumberDataElement
	// "224" für ‘Certificate Status Notice’
	FunctionID  *NumberDataElement
	KeyName     *KeyNameDataElement
	PublicKey   *PublicKeyDataElement
	Certificate *CertificateDataElement
}

func (p *PublicKeyTransmissionSegment) elements() []DataElement {
	return []DataElement{
		p.MessageID,
		p.DialogID,
		p.MessageRef,
		p.FunctionID,
		p.KeyName,
		p.PublicKey,
		p.Certificate,
	}
}

const (
	KeyCompromitted      = "1"
	KeyMaybeCompromitted = "501"
	KeyRevocationMisc    = "999"
)

var validRevocationReasons = []string{
	KeyCompromitted,
	KeyMaybeCompromitted,
	KeyRevocationMisc,
}

func NewPublicKeyRevocationSegment(number int, keyName KeyName, reason string) *PublicKeyRevocationSegment {
	if sort.SearchStrings(validRevocationReasons, reason) > len(validRevocationReasons) {
		panic(fmt.Errorf("Reason must be one of %v", validRevocationReasons))
	}
	p := &PublicKeyRevocationSegment{
		MessageID:        NewNumberDataElement(2, 1),
		FunctionID:       NewNumberDataElement(130, 3),
		KeyName:          NewKeyNameDataElement(keyName),
		RevocationReason: NewAlphaNumericDataElement(reason, 3),
		Date:             NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
	}
	p.Segment = NewBasicSegment("HKSSP", number, 2, p)
	return p
}

type PublicKeyRevocationSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *NumberDataElement
	// "130" für ‘Certificate Revocation’ (Zertifikatswiderruf)
	FunctionID *NumberDataElement
	KeyName    *KeyNameDataElement
	// "1" für ‘Schlüssel des Zertifikatseigentümers kompromittiert’
	// "501" für ‘Zertifikat ungültig wegen Verdacht auf Kompromittierung’
	// "999" für ‘gesperrt aus sonstigen Gründen’
	RevocationReason *AlphaNumericDataElement
	Date             *SecurityDateDataElement
	Certificate      *CertificateDataElement
}

func (p *PublicKeyRevocationSegment) elements() []DataElement {
	return []DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.RevocationReason,
		p.Date,
		p.Certificate,
	}
}

func NewPublicKeyRevocationConfirmationSegment(dialogId string, number int, messageReference int, keyName KeyName, reason string, refSegment *PublicKeyRevocationSegment) *PublicKeyRevocationConfirmationSegment {
	if messageReference <= 0 {
		panic(fmt.Errorf("Message Reference number must be greater 0"))
	}
	if sort.SearchStrings(validRevocationReasons, reason) > len(validRevocationReasons) {
		panic(fmt.Errorf("Reason must be one of %v", validRevocationReasons))
	}
	p := &PublicKeyRevocationConfirmationSegment{
		MessageID:        NewNumberDataElement(1, 1),
		DialogID:         NewIdentificationDataElement(dialogId),
		MessageRef:       NewNumberDataElement(messageReference, 4),
		FunctionID:       NewNumberDataElement(231, 3),
		KeyName:          NewKeyNameDataElement(keyName),
		RevocationReason: NewAlphaNumericDataElement(reason, 3),
		Date:             NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
	}
	header := NewReferencingSegmentHeader("HISSP", number, 2, refSegment.Header().Number.Val())
	p.Segment = NewBasicSegmentWithHeader(header, p)
	return p
}

type PublicKeyRevocationConfirmationSegment struct {
	Segment
	// "1" für ‘Key-Management-Nachricht ist Antwort’
	MessageID  *NumberDataElement
	DialogID   *IdentificationDataElement
	MessageRef *NumberDataElement
	// "231" für ‘Revocation Confirmation’ (Bestätigung des Zertifikatswiderrufs)
	FunctionID *NumberDataElement
	KeyName    *KeyNameDataElement
	// "1" für ‘Schlüssel des Zertifikatseigentümers kompromittiert’
	// "501" für ‘Zertifikat ungültig wegen Verdacht auf Kompromittierung’
	// "999" für ‘gesperrt aus sonstigen Gründen’
	RevocationReason *AlphaNumericDataElement
	Date             *SecurityDateDataElement
	Certificate      *CertificateDataElement
}

func (p *PublicKeyRevocationConfirmationSegment) elements() []DataElement {
	return []DataElement{
		p.MessageID,
		p.DialogID,
		p.MessageRef,
		p.FunctionID,
		p.KeyName,
		p.RevocationReason,
		p.Date,
		p.Certificate,
	}
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
	p.DataElement = NewDataElementGroup(PublicKeyDEG, 7, p)
	return p
}

type PublicKeyDataElement struct {
	DataElement
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
