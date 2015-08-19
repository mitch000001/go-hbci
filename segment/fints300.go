package segment

var FINTS300 = Version{
	version:                300,
	PinTanEncryptionHeader: NewFINTS3PinTanEncryptionHeaderSegment,
	SynchronisationRequest: NewSynchronisationSegmentV3,
	PinTanSignatureHeader:  NewFINTS3PinTanSignatureHeaderSegment,
	SignatureEnd:           NewFINTS3SignatureEndSegment,
}
