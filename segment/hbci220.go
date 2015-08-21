package segment

var HBCI220 = HBCIVersion{
	version:                   220,
	PinTanEncryptionHeader:    NewPinTanEncryptionHeaderSegment,
	RDHEncryptionHeader:       NewEncryptionHeaderSegment,
	PinTanSignatureHeader:     NewPinTanSignatureHeaderSegment,
	RDHSignatureHeader:        NewRDHSignatureHeaderSegment,
	SignatureEnd:              NewSignatureEndSegment,
	SynchronisationRequest:    NewSynchronisationSegmentV2,
	AccountTransactionRequest: NewAccountTransactionRequestSegmentV5,
}
