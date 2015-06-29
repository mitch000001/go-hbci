package hbci

import (
	"fmt"
	"sort"
	"time"

	"github.com/mitch000001/go-hbci/dataelement"
	"github.com/mitch000001/go-hbci/domain"
)

type Key interface {
	KeyName() domain.KeyName
	SetKeyNumber(number int)
	SetKeyVersion(version int)
	Sign(message []byte) (signature []byte, err error)
	Encrypt(message []byte) (encrypted []byte, err error)
	CanSign() bool
	CanEncrypt() bool
}

func NewRDHSignatureProvider(signingKey *domain.RSAKey) SignatureProvider {
	return &RDHSignatureProvider{signingKey: signingKey}
}

type RDHSignatureProvider struct {
	signingKey     *domain.RSAKey
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

const initialKeyVersion = 999

type InitialPublicKeyRenewalMessage struct {
	*basicMessage
	Identification             *IdentificationSegment
	PublicSigningKeyRequest    *PublicKeyRequestSegment
	PublicEncryptionKeyRequest *PublicKeyRequestSegment
}

func NewPublicKeyRenewalSegment(number int, keyName domain.KeyName, pubKey *domain.PublicKey) *PublicKeyRenewalSegment {
	if keyName.KeyType == "B" {
		panic(fmt.Errorf("KeyType may not be 'B'"))
	}
	p := &PublicKeyRenewalSegment{
		MessageID:  dataelement.NewNumberDataElement(2, 1),
		FunctionID: dataelement.NewNumberDataElement(112, 3),
		KeyName:    dataelement.NewKeyNameDataElement(keyName),
		PublicKey:  dataelement.NewPublicKeyDataElement(pubKey),
	}
	p.Segment = NewBasicSegment("HKSAK", number, 2, p)
	return p
}

type PublicKeyRenewalSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *dataelement.NumberDataElement
	// "112" für ‘Certificate Replacement’ (Ersatz des Zertifikats))
	FunctionID *dataelement.NumberDataElement
	// Key type may not equal 'B'
	KeyName     *dataelement.KeyNameDataElement
	PublicKey   *dataelement.PublicKeyDataElement
	Certificate *dataelement.CertificateDataElement
}

func (p *PublicKeyRenewalSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.PublicKey,
		p.Certificate,
	}
}

func NewPublicKeyRequestSegment(number int, keyName domain.KeyName) *PublicKeyRequestSegment {
	p := &PublicKeyRequestSegment{
		MessageID:  dataelement.NewNumberDataElement(2, 1),
		FunctionID: dataelement.NewNumberDataElement(124, 3),
		KeyName:    dataelement.NewKeyNameDataElement(keyName),
	}
	p.Segment = NewBasicSegment("HKISA", number, 2, p)
	return p
}

type PublicKeyRequestSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *dataelement.NumberDataElement
	// "124" für ‘Certificate Status Request’
	FunctionID  *dataelement.NumberDataElement
	KeyName     *dataelement.KeyNameDataElement
	Certificate *dataelement.CertificateDataElement
}

func (p *PublicKeyRequestSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.Certificate,
	}
}

func NewPublicKeyTransmissionSegment(dialogId string, number int, messageReference int, keyName domain.KeyName, pubKey *domain.PublicKey, refSegment *PublicKeyRequestSegment) *PublicKeyTransmissionSegment {
	if messageReference <= 0 {
		panic(fmt.Errorf("Message Reference number must be greater 0"))
	}
	p := &PublicKeyTransmissionSegment{
		MessageID:  dataelement.NewNumberDataElement(1, 1),
		DialogID:   dataelement.NewIdentificationDataElement(dialogId),
		MessageRef: dataelement.NewNumberDataElement(messageReference, 4),
		FunctionID: dataelement.NewNumberDataElement(224, 3),
		KeyName:    dataelement.NewKeyNameDataElement(keyName),
		PublicKey:  dataelement.NewPublicKeyDataElement(pubKey),
	}
	header := dataelement.NewReferencingSegmentHeader("HIISA", number, 2, refSegment.Header().Number.Val())
	p.Segment = NewBasicSegmentWithHeader(header, p)
	return p
}

type PublicKeyTransmissionSegment struct {
	Segment
	// "1" für ‘Key-Management-Nachricht ist Antwort’
	MessageID  *dataelement.NumberDataElement
	DialogID   *dataelement.IdentificationDataElement
	MessageRef *dataelement.NumberDataElement
	// "224" für ‘Certificate Status Notice’
	FunctionID  *dataelement.NumberDataElement
	KeyName     *dataelement.KeyNameDataElement
	PublicKey   *dataelement.PublicKeyDataElement
	Certificate *dataelement.CertificateDataElement
}

func (p *PublicKeyTransmissionSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
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

func NewPublicKeyRevocationSegment(number int, keyName domain.KeyName, reason string) *PublicKeyRevocationSegment {
	if sort.SearchStrings(validRevocationReasons, reason) > len(validRevocationReasons) {
		panic(fmt.Errorf("Reason must be one of %v", validRevocationReasons))
	}
	p := &PublicKeyRevocationSegment{
		MessageID:        dataelement.NewNumberDataElement(2, 1),
		FunctionID:       dataelement.NewNumberDataElement(130, 3),
		KeyName:          dataelement.NewKeyNameDataElement(keyName),
		RevocationReason: dataelement.NewAlphaNumericDataElement(reason, 3),
		Date:             dataelement.NewSecurityDateDataElement(dataelement.SecurityTimestamp, time.Now()),
	}
	p.Segment = NewBasicSegment("HKSSP", number, 2, p)
	return p
}

type PublicKeyRevocationSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *dataelement.NumberDataElement
	// "130" für ‘Certificate Revocation’ (Zertifikatswiderruf)
	FunctionID *dataelement.NumberDataElement
	KeyName    *dataelement.KeyNameDataElement
	// "1" für ‘Schlüssel des Zertifikatseigentümers kompromittiert’
	// "501" für ‘Zertifikat ungültig wegen Verdacht auf Kompromittierung’
	// "999" für ‘gesperrt aus sonstigen Gründen’
	RevocationReason *dataelement.AlphaNumericDataElement
	Date             *dataelement.SecurityDateDataElement
	Certificate      *dataelement.CertificateDataElement
}

func (p *PublicKeyRevocationSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.RevocationReason,
		p.Date,
		p.Certificate,
	}
}

func NewPublicKeyRevocationConfirmationSegment(dialogId string, number int, messageReference int, keyName domain.KeyName, reason string, refSegment *PublicKeyRevocationSegment) *PublicKeyRevocationConfirmationSegment {
	if messageReference <= 0 {
		panic(fmt.Errorf("Message Reference number must be greater 0"))
	}
	if sort.SearchStrings(validRevocationReasons, reason) > len(validRevocationReasons) {
		panic(fmt.Errorf("Reason must be one of %v", validRevocationReasons))
	}
	p := &PublicKeyRevocationConfirmationSegment{
		MessageID:        dataelement.NewNumberDataElement(1, 1),
		DialogID:         dataelement.NewIdentificationDataElement(dialogId),
		MessageRef:       dataelement.NewNumberDataElement(messageReference, 4),
		FunctionID:       dataelement.NewNumberDataElement(231, 3),
		KeyName:          dataelement.NewKeyNameDataElement(keyName),
		RevocationReason: dataelement.NewAlphaNumericDataElement(reason, 3),
		Date:             dataelement.NewSecurityDateDataElement(dataelement.SecurityTimestamp, time.Now()),
	}
	header := dataelement.NewReferencingSegmentHeader("HISSP", number, 2, refSegment.Header().Number.Val())
	p.Segment = NewBasicSegmentWithHeader(header, p)
	return p
}

type PublicKeyRevocationConfirmationSegment struct {
	Segment
	// "1" für ‘Key-Management-Nachricht ist Antwort’
	MessageID  *dataelement.NumberDataElement
	DialogID   *dataelement.IdentificationDataElement
	MessageRef *dataelement.NumberDataElement
	// "231" für ‘Revocation Confirmation’ (Bestätigung des Zertifikatswiderrufs)
	FunctionID *dataelement.NumberDataElement
	KeyName    *dataelement.KeyNameDataElement
	// "1" für ‘Schlüssel des Zertifikatseigentümers kompromittiert’
	// "501" für ‘Zertifikat ungültig wegen Verdacht auf Kompromittierung’
	// "999" für ‘gesperrt aus sonstigen Gründen’
	RevocationReason *dataelement.AlphaNumericDataElement
	Date             *dataelement.SecurityDateDataElement
	Certificate      *dataelement.CertificateDataElement
}

func (p *PublicKeyRevocationConfirmationSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
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
