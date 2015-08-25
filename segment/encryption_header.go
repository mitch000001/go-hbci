package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

type EncryptionHeader interface {
	ClientSegment
	SetClientSystemID(clientSystemID string)
	SetSecurityProfile(securityFn string)
	SetEncryptionKeyName(keyName domain.KeyName)
	SetEncryptionAlgorithm(algorithm *element.EncryptionAlgorithmDataElement)
}

func NewPinTanEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment {
	e := &EncryptionHeaderV2{
		SecurityFunction:     element.NewAlphaNumeric("998", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewPinTanEncryptionAlgorithm(),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e.ClientSegment = NewBasicSegment(998, e)

	segment := &EncryptionHeaderSegment{
		encryptionHeaderSegment: e,
	}
	return segment
}

func NewEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment {
	e := &EncryptionHeaderV2{
		SecurityFunction:     element.NewAlphaNumeric("4", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewRDHEncryptionAlgorithm(key),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e.ClientSegment = NewBasicSegment(998, e)

	segment := &EncryptionHeaderSegment{
		encryptionHeaderSegment: e,
	}
	return segment
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment EncryptionHeaderSegment -segment_interface encryptionHeaderSegment -segment_versions="EncryptionHeaderV2:2:ClientSegment,EncryptionHeaderSegmentV3:3:ClientSegment"

type EncryptionHeaderSegment struct {
	encryptionHeaderSegment
}

type encryptionHeaderSegment interface {
	EncryptionHeader
	Unmarshaler
}

type EncryptionHeaderV2 struct {
	ClientSegment
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

func (e *EncryptionHeaderV2) Version() int         { return 2 }
func (e *EncryptionHeaderV2) ID() string           { return "HNVSK" }
func (e *EncryptionHeaderV2) referencedId() string { return "" }
func (e *EncryptionHeaderV2) sender() string       { return senderBoth }

func (e *EncryptionHeaderV2) elements() []element.DataElement {
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

func (e *EncryptionHeaderV2) SetEncryptionKeyName(keyName domain.KeyName) {
	e.KeyName = element.NewKeyName(keyName)
}

func (e *EncryptionHeaderV2) SetEncryptionAlgorithm(algorithm *element.EncryptionAlgorithmDataElement) {
	e.EncryptionAlgorithm = algorithm
}

func (e *EncryptionHeaderV2) SetClientSystemID(clientSystemId string) {
	e.SecurityID = element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId)
}

func (e *EncryptionHeaderV2) SetSecurityProfile(securityFn string) {
	// NO OP
}

func NewPinTanEncryptionHeaderSegmentV3(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment {
	e := &EncryptionHeaderSegmentV3{
		SecurityProfile:      element.NewPinTanSecurityProfile(1),
		SecurityFunction:     element.NewCode("998", 3, []string{"4", "998"}),
		SecuritySupplierRole: element.NewCode("1", 3, []string{"1", "4"}),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewPinTanEncryptionAlgorithm(),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewCode("0", 3, []string{"0", "1", "2", "3", "4", "5", "6", "7", "999"}),
	}
	e.ClientSegment = NewBasicSegment(998, e)

	segment := &EncryptionHeaderSegment{
		encryptionHeaderSegment: e,
	}
	return segment
}

type EncryptionHeaderSegmentV3 struct {
	ClientSegment
	SecurityProfile *element.SecurityProfileDataElement
	// "4" for ENC, Encryption (encryption and eventually compression)
	// "998" for Cleartext
	SecurityFunction *element.CodeDataElement
	// "1" for ISS,  Herausgeber der chiffrierten Nachricht (Erfasser)
	// "4" for WIT, der Unterzeichnete ist Zeuge, aber für den Inhalt der
	// Nachricht nicht verantwortlich (Übermittler, welcher nicht Erfasser ist)
	SecuritySupplierRole *element.CodeDataElement
	SecurityID           *element.SecurityIdentificationDataElement
	SecurityDate         *element.SecurityDateDataElement
	EncryptionAlgorithm  *element.EncryptionAlgorithmDataElement
	KeyName              *element.KeyNameDataElement
	// 0: no compression (NULL)
	// 1: Lempel, Ziv, Welch (LZW)
	// 2: Optimized LZW (COM)
	// 3: Lempel, Ziv (LZSS)
	// 4: LZ + Huffman Coding (LZHuf)
	// 5: PKZIP (ZIP)
	// 6: deflate (GZIP) (http://www.gzip.org/zlib)
	// 7: bzip2 (http://sourceware.cygnus.com/bzip2/)
	// 999: Gegenseitig vereinbart (ZZZ)
	CompressionFunction *element.CodeDataElement
	Certificate         *element.CertificateDataElement
}

func (e *EncryptionHeaderSegmentV3) Version() int         { return 3 }
func (e *EncryptionHeaderSegmentV3) ID() string           { return "HNVSK" }
func (e *EncryptionHeaderSegmentV3) referencedId() string { return "" }
func (e *EncryptionHeaderSegmentV3) sender() string       { return senderBoth }

func (e *EncryptionHeaderSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		e.SecurityProfile,
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

func (e *EncryptionHeaderSegmentV3) SetEncryptionKeyName(keyName domain.KeyName) {
	e.KeyName = element.NewKeyName(keyName)
}

func (e *EncryptionHeaderSegmentV3) SetEncryptionAlgorithm(algorithm *element.EncryptionAlgorithmDataElement) {
	e.EncryptionAlgorithm = algorithm
}

func (e *EncryptionHeaderSegmentV3) SetClientSystemID(clientSystemId string) {
	e.SecurityID = element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId)
}

func (e *EncryptionHeaderSegmentV3) SetSecurityProfile(securityFn string) {
	if securityFn == "999" {
		e.SecurityProfile = element.NewPinTanSecurityProfile(1)
	} else {
		e.SecurityProfile = element.NewPinTanSecurityProfile(2)
	}
}
