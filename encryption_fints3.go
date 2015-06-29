package hbci

import "time"

func NewFINTS3EncryptedPinTanMessage(clientSystemId string, keyName KeyName, encryptedMessage []byte) *EncryptedMessage {
	e := &EncryptedMessage{
		EncryptionHeader: NewFINTS3PinTanEncryptionHeaderSegment(clientSystemId, keyName),
		EncryptedData:    NewEncryptedDataSegment(encryptedMessage),
	}
	e.basicMessage = newBasicMessage(e)
	return e
}

func NewFINTS3PinTanEncryptionHeaderSegment(clientSystemId string, keyName KeyName) *EncryptionHeaderSegment {
	v2 := &EncryptionHeaderVersion2{
		SecurityFunction:     NewAlphaNumericDataElement("998", 3),
		SecuritySupplierRole: NewAlphaNumericDataElement("1", 3),
		SecurityID:           NewRDHSecurityIdentificationDataElement(SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         NewSecurityDateDataElement(SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  NewPinTanEncryptionAlgorithmDataElement(),
		KeyName:              NewKeyNameDataElement(keyName),
		CompressionFunction:  NewAlphaNumericDataElement("0", 3),
	}
	v3 := &EncryptionHeaderVersion3{
		EncryptionHeaderVersion2: v2,
		SecurityProfile:          NewPinTanSecurityProfile(),
	}
	e := &EncryptionHeaderSegment{
		version: v3,
	}
	e.Segment = NewBasicSegment("HNVSK", 998, 3, e)
	return e
}

type EncryptionHeaderVersion3 struct {
	*EncryptionHeaderVersion2
	SecurityProfile *SecurityProfileDataElement
}

func (e *EncryptionHeaderVersion3) version() int {
	return 3
}

func (e *EncryptionHeaderVersion3) versionedElements() []DataElement {
	return []DataElement{
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
