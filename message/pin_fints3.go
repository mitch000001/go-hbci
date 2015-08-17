package message

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

func NewFINTS3PinTanEncryptionProvider(key *domain.PinKey, clientSystemId string) *FINTS3PinTanEncryptionProvider {
	return &FINTS3PinTanEncryptionProvider{
		&PinTanCryptoProvider{
			key:            key,
			clientSystemId: clientSystemId,
		},
	}
}

type FINTS3PinTanEncryptionProvider struct {
	*PinTanCryptoProvider
}

func (p *FINTS3PinTanEncryptionProvider) WriteEncryptionHeader(message *EncryptedMessage) {
	message.EncryptionHeader = segment.NewFINTS3PinTanEncryptionHeaderSegment(p.clientSystemId, p.key.KeyName())
}

func NewFINTS3PinTanSignatureProvider(key *domain.PinKey, clientSystemId string) SignatureProvider {
	controlReference := generateControlReference(key)
	return &FINTS3PinTanSignatureProvider{
		&PinTanSignatureProvider{
			key:              key,
			clientSystemId:   clientSystemId,
			controlReference: controlReference,
			securityFn:       "999",
		},
	}
}

type FINTS3PinTanSignatureProvider struct {
	*PinTanSignatureProvider
}

func (p *FINTS3PinTanSignatureProvider) SignMessage(signedMessage SignedHBCIMessage) error {
	signatureHeader := segment.NewFINTS3PinTanSignatureHeaderSegment(p.controlReference, p.clientSystemId, p.key.KeyName())
	signatureHeader.SetSecurityFunction(p.securityFn)
	signatureEnd := segment.NewSignatureEndSegment(0, p.controlReference)
	signatureEnd.SetPinTan(p.key.Pin(), "")
	signedMessage.SetSignatureHeader(signatureHeader)
	signedMessage.SetSignatureEnd(signatureEnd)
	signedMessage.SetNumbers()
	return nil
}
