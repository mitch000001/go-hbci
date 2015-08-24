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

type SignatureProvider interface {
	SetSecurityFunction(securityFn string)
	SetClientSystemID(clientSystemId string)
	Sign(message []byte) ([]byte, error)
	WriteSignatureHeader(segment.SignatureHeader)
	WriteSignature(end segment.SignatureEnd, signature []byte)
}

func MessageHashSum(message string) []byte {
	h := ripemd160.New()
	//io.WriteString(h, initializationVector)
	io.WriteString(h, message)
	return h.Sum(nil)
}

func SignMessageHash(messageHash []byte, key *rsa.PrivateKey) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, key, 0, messageHash)
}

func generateControlReference(key domain.Key) string {
	h := crc32.NewIEEE()
	keyName := key.KeyName()
	io.WriteString(h, fmt.Sprintf("%d", keyName.BankID.CountryCode))
	io.WriteString(h, keyName.BankID.ID)
	io.WriteString(h, keyName.UserID)
	io.WriteString(h, keyName.KeyType)
	io.WriteString(h, fmt.Sprintf("%d", keyName.KeyNumber))
	io.WriteString(h, fmt.Sprintf("%d", keyName.KeyVersion))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewPinTanSignatureProvider(key *domain.PinKey, clientSystemId string) SignatureProvider {
	controlReference := generateControlReference(key)
	return &PinTanSignatureProvider{
		key:              key,
		clientSystemId:   clientSystemId,
		controlReference: controlReference,
		securityFn:       "999",
	}
}

type PinTanSignatureProvider struct {
	key              *domain.PinKey
	clientSystemId   string
	securityFn       string
	controlReference string
}

func (p *PinTanSignatureProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanSignatureProvider) SetSecurityFunction(securityFn string) {
	p.securityFn = securityFn
}

func (p *PinTanSignatureProvider) Sign(message []byte) ([]byte, error) {
	return p.key.Sign(message)
}

func (p *PinTanSignatureProvider) WriteSignatureHeader(header segment.SignatureHeader) {
	header.SetSecurityFunction(p.securityFn)
	header.SetClientSystemID(p.clientSystemId)
	header.SetSigningKeyName(p.key.KeyName())
}

func (p *PinTanSignatureProvider) WriteSignature(end segment.SignatureEnd, signature []byte) {
	end.SetPinTan(p.key.Pin(), "")
}

func NewRDHSignatureProvider(signingKey *domain.RSAKey, signatureId int) SignatureProvider {
	controlReference := generateControlReference(signingKey)
	return &RDHSignatureProvider{
		signingKey:       signingKey,
		controlReference: controlReference,
		signatureId:      signatureId,
		securityFn:       "1",
	}
}

type RDHSignatureProvider struct {
	signingKey       *domain.RSAKey
	clientSystemId   string
	controlReference string
	securityFn       string
	signatureId      int
}

func (r *RDHSignatureProvider) SetClientSystemID(clientSystemId string) {
	r.clientSystemId = clientSystemId
}

func (r *RDHSignatureProvider) SetSecurityFunction(securityFn string) {
	r.securityFn = securityFn
}

func (r *RDHSignatureProvider) Sign(message []byte) ([]byte, error) {
	hashSum := MessageHashSum(string(message))
	return r.signingKey.Sign(hashSum)
}

func (r *RDHSignatureProvider) WriteSignatureHeader(header segment.SignatureHeader) {
	header.SetSecurityFunction(r.securityFn)
	header.SetClientSystemID(r.clientSystemId)
	header.SetSigningKeyName(r.signingKey.KeyName())
}

func (r *RDHSignatureProvider) WriteSignature(end segment.SignatureEnd, signature []byte) {
	end.SetSignature(signature)
}
