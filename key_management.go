package hbci

import (
	"fmt"
	"sort"
	"time"
)

type PublicKeyRenewalMessage struct {
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
	header := NewSegmentHeader("HKSAK", number, 2)
	p.basicSegment = NewBasicSegment(header, p)
	return p
}

type PublicKeyRenewalSegment struct {
	*basicSegment
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
	header := NewSegmentHeader("HKISA", number, 2)
	p.basicSegment = NewBasicSegment(header, p)
	return p
}

type PublicKeyRequestSegment struct {
	*basicSegment
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
	header := NewReferencingSegmentHeader("HIISA", number, 2, refSegment.header.Number.Val())
	p.basicSegment = NewBasicSegment(header, p)
	return p
}

type PublicKeyTransmissionSegment struct {
	*basicSegment
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
	header := NewSegmentHeader("HKSSP", number, 2)
	p.basicSegment = NewBasicSegment(header, p)
	return p
}

type PublicKeyRevocationSegment struct {
	*basicSegment
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
	header := NewReferencingSegmentHeader("HISSP", number, 2, refSegment.header.Number.Val())
	p.basicSegment = NewBasicSegment(header, p)
	return p
}

type PublicKeyRevocationConfirmationSegment struct {
	*basicSegment
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
