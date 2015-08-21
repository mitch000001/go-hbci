package segment

var FINTS300 = Version{
	version:                300,
	PinTanEncryptionHeader: NewPinTanEncryptionHeaderSegmentV3,
	SynchronisationRequest: NewSynchronisationSegmentV3,
	PinTanSignatureHeader:  NewPinTanSignatureHeaderSegmentV4,
	SignatureEnd:           NewSignatureEndSegmentV2,
}
