package segment

var FINTS300 = HBCIVersion{
	version:                   300,
	PinTanEncryptionHeader:    NewPinTanEncryptionHeaderSegmentV3,
	SynchronisationRequest:    NewSynchronisationSegmentV3,
	PinTanSignatureHeader:     NewPinTanSignatureHeaderSegmentV4,
	SignatureEnd:              NewSignatureEndSegmentV2,
	AccountTransactionRequest: NewAccountTransactionRequestSegmentV6,
}
