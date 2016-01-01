package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

type CryptoProvider interface {
	SetClientSystemID(clientSystemId string)
	SetSecurityFunction(securityFn string)
	Encrypt(message []byte) ([]byte, error)
	Decrypt(encryptedMessage []byte) ([]byte, error)
	WriteEncryptionHeader(header segment.EncryptionHeader)
}

func NewPinTanCryptoProvider(key *domain.PinKey, clientSystemId string) *PinTanCryptoProvider {
	return &PinTanCryptoProvider{
		key:            key,
		clientSystemId: clientSystemId,
		securityFn:     "999",
	}
}

type PinTanCryptoProvider struct {
	key            *domain.PinKey
	clientSystemId string
	securityFn     string
}

func (p *PinTanCryptoProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanCryptoProvider) SetSecurityFunction(securityFn string) {
	p.securityFn = securityFn
}

func (p *PinTanCryptoProvider) Encrypt(message []byte) ([]byte, error) {
	if p.key.Pin() == "" {
		return nil, fmt.Errorf("Malformed PIN")
	}
	return p.key.Encrypt(message)
}

func (p *PinTanCryptoProvider) Decrypt(encryptedMessage []byte) ([]byte, error) {
	return p.key.Decrypt(encryptedMessage)
}

func (p *PinTanCryptoProvider) WriteEncryptionHeader(header segment.EncryptionHeader) {
	header.SetClientSystemID(p.clientSystemId)
	header.SetSecurityProfile(p.securityFn)
	header.SetEncryptionKeyName(p.key.KeyName())
	header.SetEncryptionAlgorithm(element.NewPinTanEncryptionAlgorithm())
}
