package hbci

func NewFINTS3PinTanEncryptionProvider(key *PinKey, clientSystemId string) *FINTS3PinTanEncryptionProvider {
	return &FINTS3PinTanEncryptionProvider{
		&PinTanEncryptionProvider{
			key:            key,
			clientSystemId: clientSystemId,
		},
	}
}

type FINTS3PinTanEncryptionProvider struct {
	*PinTanEncryptionProvider
}

func (p *FINTS3PinTanEncryptionProvider) Encrypt(message []byte) (*EncryptedMessage, error) {
	encryptedBytes, _ := p.key.Encrypt(message)
	encryptedMessage := NewFINTS3EncryptedPinTanMessage(p.clientSystemId, p.key.KeyName(), encryptedBytes)
	return encryptedMessage, nil
}

func NewFINTS3PinTanSignatureProvider(key *PinKey) SignatureProvider {
	return &FINTS3PinTanSignatureProvider{&PinTanSignatureProvider{key: key}}
}

type FINTS3PinTanSignatureProvider struct {
	*PinTanSignatureProvider
}

func (f *FINTS3PinTanSignatureProvider) NewSignatureHeader(controlReference string, signatureId int) *SignatureHeaderSegment {
	return NewFINTS3PinTanSignatureHeaderSegment(controlReference, f.clientSystemId, f.key.KeyName())
}
