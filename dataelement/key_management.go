package dataelement

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
	"reflect"
)

func NewRSAKey(pubKey *PublicKey, keyName *KeyName) *RSAKey {
	return &RSAKey{PublicKey: pubKey, keyName: keyName}
}

type RSAKey struct {
	*PublicKey
	keyName *KeyName
}

func (r *RSAKey) KeyName() KeyName {
	return *r.keyName
}

func (r *RSAKey) SetKeyNumber(number int) {
	r.keyName.KeyNumber = number
}

func (r *RSAKey) SetKeyVersion(version int) {
	r.keyName.KeyVersion = version
}

func (r *RSAKey) CanSign() bool {
	return r.PublicKey.rsaPrivateKey != nil
}

func (r *RSAKey) CanEncrypt() bool {
	return r.PublicKey.rsaPublicKey != nil
}

func NewEncryptionKey(modulus, exponent []byte) *PublicKey {
	p := &PublicKey{
		Type: "V",
	}
	copy(p.Modulus, modulus)
	copy(p.Exponent, exponent)
	mod := new(big.Int).SetBytes(modulus)
	exp := new(big.Int).SetBytes(exponent)
	pubKey := rsa.PublicKey{
		N: mod,
		E: int(exp.Int64()),
	}
	p.rsaPublicKey = &pubKey
	return p
}

type PublicKey struct {
	Type          string
	Modulus       []byte
	Exponent      []byte
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

func (p *PublicKey) SigningKey() *rsa.PrivateKey {
	return p.rsaPrivateKey
}

func (p *PublicKey) Sign(message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, p.rsaPrivateKey, 0, message)
}

func (p *PublicKey) Encrypt(message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, p.rsaPublicKey, message)
}

func NewPublicKeyDataElement(pubKey *PublicKey) *PublicKeyDataElement {
	if !reflect.DeepEqual(pubKey.Exponent, []byte("65537")) {
		panic(fmt.Errorf("Exponent must equal 65537 (% X)", "65537"))
	}
	p := &PublicKeyDataElement{
		Usage:         NewAlphaNumericDataElement(pubKey.Type, 3),
		OperationMode: NewAlphaNumericDataElement("16", 3),
		Cipher:        NewAlphaNumericDataElement("10", 3),
		Modulus:       NewBinaryDataElement(pubKey.Modulus, 512),
		ModulusID:     NewAlphaNumericDataElement("12", 3),
		Exponent:      NewBinaryDataElement(pubKey.Exponent, 512),
		ExponentID:    NewAlphaNumericDataElement("13", 3),
	}
	p.DataElement = NewDataElementGroup(PublicKeyDEG, 7, p)
	return p
}

type PublicKeyDataElement struct {
	DataElement
	// "5" for OCF, Owner Ciphering (Encryption key)
	// "6" for OSG, Owner Signing (Signing key)
	Usage *AlphaNumericDataElement
	// "16" for DSMR (ISO 9796)
	OperationMode *AlphaNumericDataElement
	// "10" for RSA
	Cipher  *AlphaNumericDataElement
	Modulus *BinaryDataElement
	// "12" for MOD, Modulus
	ModulusID *AlphaNumericDataElement
	// "65537"
	Exponent *BinaryDataElement
	// "13" for EXP, Exponent
	ExponentID *AlphaNumericDataElement
}

func (p *PublicKeyDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.Usage,
		p.OperationMode,
		p.Cipher,
		p.Modulus,
		p.ModulusID,
		p.Exponent,
		p.ExponentID,
	}
}

func (p *PublicKeyDataElement) Val() *PublicKey {
	return &PublicKey{
		Type:     p.Usage.Val(),
		Modulus:  p.Modulus.Val(),
		Exponent: p.Exponent.Val(),
	}
}
