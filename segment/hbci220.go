package segment

var HBCI220 = HBCIVersion{
	version:                220,
	PinTanEncryptionHeader: NewPinTanEncryptionHeaderSegment,
	RDHEncryptionHeader:    NewEncryptionHeaderSegment,
	SignatureHeader:        NewSignatureHeaderSegmentV3,
	PinTanSignatureHeader:  NewPinTanSignatureHeaderSegment,
	RDHSignatureHeader:     NewRDHSignatureHeaderSegment,
	SignatureEnd:           NewSignatureEndSegmentV1,
	SynchronisationRequest: NewSynchronisationSegmentV2,
}
