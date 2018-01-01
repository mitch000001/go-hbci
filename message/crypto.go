package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

// CryptoProvider represents a provider to encrypt and decrypt messages
type CryptoProvider interface {
	SetClientSystemID(clientSystemID string)
	SetSecurityFunction(securityFn string)
	Encrypt(message []byte) ([]byte, error)
	Decrypt(encryptedMessage []byte) ([]byte, error)
	WriteEncryptionHeader(header segment.EncryptionHeader)
}

// NewPinTanCryptoProvider creates a new CryptoProvider for the pin key
func NewPinTanCryptoProvider(key *domain.PinKey, clientSystemID string) CryptoProvider {
	return &pinTanCryptoProvider{
		key:            key,
		clientSystemID: clientSystemID,
		securityFn:     "999",
	}
}

type pinTanCryptoProvider struct {
	key            *domain.PinKey
	clientSystemID string
	securityFn     string
}

func (p *pinTanCryptoProvider) SetClientSystemID(clientSystemID string) {
	p.clientSystemID = clientSystemID
}

func (p *pinTanCryptoProvider) SetSecurityFunction(securityFn string) {
	p.securityFn = securityFn
}

func (p *pinTanCryptoProvider) Encrypt(message []byte) ([]byte, error) {
	if p.key.Pin() == "" {
		return nil, fmt.Errorf("Malformed PIN")
	}
	return p.key.Encrypt(message)
}

func (p *pinTanCryptoProvider) Decrypt(encryptedMessage []byte) ([]byte, error) {
	return p.key.Decrypt(encryptedMessage)
}

func (p *pinTanCryptoProvider) WriteEncryptionHeader(header segment.EncryptionHeader) {
	header.SetClientSystemID(p.clientSystemID)
	header.SetSecurityProfile(p.securityFn)
	header.SetEncryptionKeyName(p.key.KeyName())
	header.SetEncryptionAlgorithm(element.NewPinTanEncryptionAlgorithm())
}
