package segment

import (
	"fmt"
	"sort"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPublicKeyRenewalSegment(number int, keyName domain.KeyName, pubKey *domain.PublicKey) *PublicKeyRenewalSegment {
	if keyName.KeyType == "B" {
		panic(fmt.Errorf("KeyType may not be 'B'"))
	}
	p := &PublicKeyRenewalSegment{
		MessageID:  element.NewNumber(2, 1),
		FunctionID: element.NewNumber(112, 3),
		KeyName:    element.NewKeyName(keyName),
		PublicKey:  element.NewPublicKey(pubKey),
	}
	p.Segment = NewBasicSegment("HKSAK", number, 2, p)
	return p
}

type PublicKeyRenewalSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *element.NumberDataElement
	// "112" für ‘Certificate Replacement’ (Ersatz des Zertifikats))
	FunctionID *element.NumberDataElement
	// Key type may not equal 'B'
	KeyName     *element.KeyNameDataElement
	PublicKey   *element.PublicKeyDataElement
	Certificate *element.CertificateDataElement
}

func (p *PublicKeyRenewalSegment) elements() []element.DataElement {
	return []element.DataElement{
		p.MessageID,
		p.FunctionID,
		p.KeyName,
		p.PublicKey,
		p.Certificate,
	}
}

func NewPublicKeyRequestSegment(number int, keyName domain.KeyName) *PublicKeyRequestSegment {
	p := &PublicKeyRequestSegment{
		MessageID:  element.NewNumber(2, 1),
		FunctionID: element.NewNumber(124, 3),
		KeyName:    element.NewKeyName(keyName),
	}
	p.Segment = NewBasicSegment("HKISA", number, 2, p)
	return p
}

type PublicKeyRequestSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *element.NumberDataElement
	// "124" für ‘Certificate Status Request’
	FunctionID  *element.NumberDataElement
	KeyName     *element.KeyNameDataElement
	Certificate *element.CertificateDataElement
}

func (p *PublicKeyRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
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
		MessageID:  element.NewNumber(1, 1),
		DialogID:   element.NewIdentification(dialogId),
		MessageRef: element.NewNumber(messageReference, 4),
		FunctionID: element.NewNumber(224, 3),
		KeyName:    element.NewKeyName(keyName),
		PublicKey:  element.NewPublicKey(pubKey),
	}
	header := element.NewReferencingSegmentHeader("HIISA", number, 2, refSegment.Header().Number.Val())
	p.Segment = NewBasicSegmentWithHeader(header, p)
	return p
}

type PublicKeyTransmissionSegment struct {
	Segment
	// "1" für ‘Key-Management-Nachricht ist Antwort’
	MessageID  *element.NumberDataElement
	DialogID   *element.IdentificationDataElement
	MessageRef *element.NumberDataElement
	// "224" für ‘Certificate Status Notice’
	FunctionID  *element.NumberDataElement
	KeyName     *element.KeyNameDataElement
	PublicKey   *element.PublicKeyDataElement
	Certificate *element.CertificateDataElement
}

func (p *PublicKeyTransmissionSegment) elements() []element.DataElement {
	return []element.DataElement{
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
		MessageID:        element.NewNumber(2, 1),
		FunctionID:       element.NewNumber(130, 3),
		KeyName:          element.NewKeyName(keyName),
		RevocationReason: element.NewAlphaNumeric(reason, 3),
		Date:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
	}
	p.Segment = NewBasicSegment("HKSSP", number, 2, p)
	return p
}

type PublicKeyRevocationSegment struct {
	Segment
	// "2" für ‘Key-Management-Nachricht erwartet Antwort’
	MessageID *element.NumberDataElement
	// "130" für ‘Certificate Revocation’ (Zertifikatswiderruf)
	FunctionID *element.NumberDataElement
	KeyName    *element.KeyNameDataElement
	// "1" für ‘Schlüssel des Zertifikatseigentümers kompromittiert’
	// "501" für ‘Zertifikat ungültig wegen Verdacht auf Kompromittierung’
	// "999" für ‘gesperrt aus sonstigen Gründen’
	RevocationReason *element.AlphaNumericDataElement
	Date             *element.SecurityDateDataElement
	Certificate      *element.CertificateDataElement
}

func (p *PublicKeyRevocationSegment) elements() []element.DataElement {
	return []element.DataElement{
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
		MessageID:        element.NewNumber(1, 1),
		DialogID:         element.NewIdentification(dialogId),
		MessageRef:       element.NewNumber(messageReference, 4),
		FunctionID:       element.NewNumber(231, 3),
		KeyName:          element.NewKeyName(keyName),
		RevocationReason: element.NewAlphaNumeric(reason, 3),
		Date:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
	}
	header := element.NewReferencingSegmentHeader("HISSP", number, 2, refSegment.Header().Number.Val())
	p.Segment = NewBasicSegmentWithHeader(header, p)
	return p
}

type PublicKeyRevocationConfirmationSegment struct {
	Segment
	// "1" für ‘Key-Management-Nachricht ist Antwort’
	MessageID  *element.NumberDataElement
	DialogID   *element.IdentificationDataElement
	MessageRef *element.NumberDataElement
	// "231" für ‘Revocation Confirmation’ (Bestätigung des Zertifikatswiderrufs)
	FunctionID *element.NumberDataElement
	KeyName    *element.KeyNameDataElement
	// "1" für ‘Schlüssel des Zertifikatseigentümers kompromittiert’
	// "501" für ‘Zertifikat ungültig wegen Verdacht auf Kompromittierung’
	// "999" für ‘gesperrt aus sonstigen Gründen’
	RevocationReason *element.AlphaNumericDataElement
	Date             *element.SecurityDateDataElement
	Certificate      *element.CertificateDataElement
}

func (p *PublicKeyRevocationConfirmationSegment) elements() []element.DataElement {
	return []element.DataElement{
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
