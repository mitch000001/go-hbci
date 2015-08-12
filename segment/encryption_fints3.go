package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewFINTS3PinTanEncryptionHeaderSegment(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegmentV3 {
	v2 := &EncryptionHeaderSegment{
		SecurityFunction:     element.NewAlphaNumeric("998", 3),
		SecuritySupplierRole: element.NewAlphaNumeric("1", 3),
		SecurityID:           element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityDate:         element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		EncryptionAlgorithm:  element.NewPinTanEncryptionAlgorithm(),
		KeyName:              element.NewKeyName(keyName),
		CompressionFunction:  element.NewAlphaNumeric("0", 3),
	}
	e := &EncryptionHeaderSegmentV3{
		EncryptionHeaderSegment: v2,
		SecurityProfile:         element.NewPinTanSecurityProfile(),
	}
	e.Segment = NewBasicSegment(998, e)
	return e
}

type EncryptionHeaderSegmentV3 struct {
	*EncryptionHeaderSegment
	SecurityProfile *element.SecurityProfileDataElement
}

func (e *EncryptionHeaderSegmentV3) version() int {
	return 3
}

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
