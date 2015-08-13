package message

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash/crc32"
	"io"
	"reflect"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
	"golang.org/x/crypto/ripemd160"
)

type SignatureProvider interface {
	SetSecurityFunction(securityFn string)
	SetClientSystemID(clientSystemId string)
	SignMessage(SignedHBCIMessage) error
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

func (p *PinTanSignatureProvider) SignMessage(signedMessage SignedHBCIMessage) error {
	if p.key.Pin() == "" {
		return fmt.Errorf("Malformed PIN")
	}
	signatureHeader := segment.NewPinTanSignatureHeaderSegment(p.controlReference, p.clientSystemId, p.key.KeyName())
	signatureHeader.SetSecurityFunction(p.securityFn)
	signatureEnd := segment.NewSignatureEndSegment(0, p.controlReference)
	signatureEnd.SetPinTan(p.key.Pin(), "")
	signedMessage.SetSignatureHeader(signatureHeader)
	signedMessage.SetSignatureEnd(signatureEnd)
	signedMessage.SetNumbers()
	return nil
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

func (r *RDHSignatureProvider) SignMessage(message SignedHBCIMessage) error {
	signatureHeader := segment.NewRDHSignatureHeaderSegment(r.controlReference, r.signatureId, r.clientSystemId, r.signingKey.KeyName())
	signatureEnd := segment.NewSignatureEndSegment(0, r.controlReference)
	message.SetSignatureHeader(signatureHeader)
	message.SetSignatureEnd(signatureEnd)
	message.SetNumbers()
	var buffer bytes.Buffer
	buffer.WriteString(signatureHeader.String())
	for _, segment := range message.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			buffer.WriteString(segment.String())
		}
	}
	hashSum := MessageHashSum(buffer.String())
	sig, err := r.signingKey.Sign(hashSum)
	if err != nil {
		return err
	}
	// TODO: clean up the pointer chaos
	signatureEnd.SetSignature(sig)
	return nil
}
