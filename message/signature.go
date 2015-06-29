package message

import (
	"crypto/rand"
	"crypto/rsa"
	"io"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
	"golang.org/x/crypto/ripemd160"
)

type SignatureProvider interface {
	SetClientSystemID(clientSystemId string)
	SignMessage(SignedHBCIMessage) error
	NewSignatureHeader(controlReference string, signatureId int) *segment.SignatureHeaderSegment
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

func NewPinTanSignatureProvider(key *domain.PinKey, clientSystemId string) SignatureProvider {
	return &PinTanSignatureProvider{key: key, clientSystemId: clientSystemId}
}

type PinTanSignatureProvider struct {
	key            *domain.PinKey
	clientSystemId string
}

func (p *PinTanSignatureProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanSignatureProvider) SignMessage(signedMessage SignedHBCIMessage) error {
	signedMessage.SignatureEndSegment().SetPinTan(p.key.Pin(), "")
	return nil
}

func (p *PinTanSignatureProvider) NewSignatureHeader(controlReference string, signatureId int) *segment.SignatureHeaderSegment {
	return segment.NewPinTanSignatureHeaderSegment(controlReference, p.clientSystemId, p.key.KeyName())
}

func NewRDHSignatureProvider(signingKey *domain.RSAKey) SignatureProvider {
	return &RDHSignatureProvider{signingKey: signingKey}
}

type RDHSignatureProvider struct {
	signingKey     *domain.RSAKey
	clientSystemId string
}

func (r *RDHSignatureProvider) SetClientSystemID(clientSystemId string) {
	r.clientSystemId = clientSystemId
}

func (r *RDHSignatureProvider) SignMessage(message SignedHBCIMessage) error {
	marshaledMessage := message.SignatureBeginSegment().String()
	for _, segment := range message.HBCISegments() {
		marshaledMessage += segment.String()
	}
	hashSum := MessageHashSum(marshaledMessage)
	sig, err := r.signingKey.Sign(hashSum)
	if err != nil {
		return err
	}
	message.SignatureEndSegment().SetSignature(sig)
	return nil
}

func (p *RDHSignatureProvider) NewSignatureHeader(controlReference string, signatureId int) *segment.SignatureHeaderSegment {
	return segment.NewRDHSignatureHeaderSegment(controlReference, signatureId, p.clientSystemId, p.signingKey.KeyName())
}
