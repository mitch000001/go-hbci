package segment

import "github.com/mitch000001/go-hbci/domain"

var SupportedHBCIVersions = map[int]Version{
	220: HBCI220,
	300: FINTS300,
}

type Version struct {
	version                int
	PinTanEncryptionHeader func(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment
	RDHEncryptionHeader    func(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment
	PinTanSignatureHeader  func(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment
	RDHSignatureHeader     func(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment
	SignatureEnd           func(number int, controlReference string) *SignatureEndSegment
	SynchronisationRequest func(modus int) *SynchronisationRequestSegment
}

func (v Version) Version() int {
	return v.version
}
