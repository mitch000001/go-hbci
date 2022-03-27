package message

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
	"golang.org/x/crypto/ripemd160"
)

// A SignatureProvider represents a provider to sign a message
type SignatureProvider interface {
	SetSecurityFunction(securityFn string)
	SetClientSystemID(clientSystemID string)
	Sign(message []byte) ([]byte, error)
	WriteSignatureHeader(segment.SignatureHeader)
	WriteSignature(end segment.SignatureEnd, signature []byte)
}

// HashSum calculates the riemd160 hash sum of message
func HashSum(message string) []byte {
	h := ripemd160.New()
	//io.WriteString(h, initializationVector)
	io.WriteString(h, message)
	return h.Sum(nil)
}

// SignMessageHash signs the messageHash with key
func SignMessageHash(messageHash []byte, key *rsa.PrivateKey) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, key, 0, messageHash)
}

func generateControlReference(key domain.Key) string {
	h := crc32.NewIEEE()
	keyName := key.KeyName()
	io.WriteString(h, fmt.Sprintf("%d", keyName.BankID.CountryCode))
	io.WriteString(h, keyName.BankID.ID)
	io.WriteString(h, keyName.UserID)
	io.WriteString(h, string(keyName.KeyType))
	io.WriteString(h, fmt.Sprintf("%d", keyName.KeyNumber))
	io.WriteString(h, fmt.Sprintf("%d", keyName.KeyVersion))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// NewPinTanSignatureProvider creates a SignatureProvider for the pin key
func NewPinTanSignatureProvider(key *domain.PinKey, clientSystemID string) SignatureProvider {
	controlReference := generateControlReference(key)
	return &pinTanSignatureProvider{
		key:              key,
		clientSystemID:   clientSystemID,
		controlReference: controlReference,
		securityFn:       "999",
	}
}

type pinTanSignatureProvider struct {
	key              *domain.PinKey
	clientSystemID   string
	securityFn       string
	controlReference string
}

func (p *pinTanSignatureProvider) SetClientSystemID(clientSystemID string) {
	p.clientSystemID = clientSystemID
}

func (p *pinTanSignatureProvider) SetSecurityFunction(securityFn string) {
	p.securityFn = securityFn
}

func (p *pinTanSignatureProvider) Sign(message []byte) ([]byte, error) {
	return p.key.Sign(message)
}

func (p *pinTanSignatureProvider) WriteSignatureHeader(header segment.SignatureHeader) {
	header.SetSecurityFunction(p.securityFn)
	header.SetClientSystemID(p.clientSystemID)
	header.SetSigningKeyName(p.key.KeyName())
	header.SetControlReference(p.controlReference)
	header.SetSignatureID(0)
}

func (p *pinTanSignatureProvider) WriteSignature(end segment.SignatureEnd, signature []byte) {
	end.SetPinTan(p.key.Pin(), "")
	end.SetControlReference(p.controlReference)
}

// NewRDHSignatureProvider creates a new SignatureProvider for the given signingKey
func NewRDHSignatureProvider(signingKey *domain.RSAKey, signatureID int) SignatureProvider {
	controlReference := generateControlReference(signingKey)
	return &rdhSignatureProvider{
		signingKey:       signingKey,
		controlReference: controlReference,
		signatureID:      signatureID,
		securityFn:       "1",
	}
}

type rdhSignatureProvider struct {
	signingKey       *domain.RSAKey
	clientSystemID   string
	controlReference string
	securityFn       string
	signatureID      int
}

func (r *rdhSignatureProvider) SetClientSystemID(clientSystemID string) {
	r.clientSystemID = clientSystemID
}

func (r *rdhSignatureProvider) SetSecurityFunction(securityFn string) {
	r.securityFn = securityFn
}

func (r *rdhSignatureProvider) Sign(message []byte) ([]byte, error) {
	hashSum := HashSum(string(message))
	return r.signingKey.Sign(hashSum)
}

func (r *rdhSignatureProvider) WriteSignatureHeader(header segment.SignatureHeader) {
	header.SetSecurityFunction(r.securityFn)
	header.SetClientSystemID(r.clientSystemID)
	header.SetSigningKeyName(r.signingKey.KeyName())
	header.SetSignatureID(r.signatureID)
	header.SetControlReference(r.controlReference)
}

func (r *rdhSignatureProvider) WriteSignature(end segment.SignatureEnd, signature []byte) {
	end.SetSignature(signature)
	end.SetControlReference(r.controlReference)
}
