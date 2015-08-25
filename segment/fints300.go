package segment

var FINTS300 = HBCIVersion{
	version:                   300,
	PinTanEncryptionHeader:    NewPinTanEncryptionHeaderSegmentV3,
	SynchronisationRequest:    NewSynchronisationSegmentV3,
	SignatureHeader:           NewSignatureHeaderSegmentV4,
	PinTanSignatureHeader:     NewPinTanSignatureHeaderSegmentV4,
	SignatureEnd:              NewSignatureEndSegmentV2,
	AccountBalanceRequest:     NewAccountBalanceRequestV6,
	AccountTransactionRequest: NewAccountTransactionRequestSegmentV6,
	StatusProtocolRequest:     NewStatusProtocolRequestV4,
}
