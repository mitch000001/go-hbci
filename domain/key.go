package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

type Key interface {
	KeyName() KeyName
	SetKeyNumber(number int)
	SetKeyVersion(version int)
	Sign(message []byte) (signature []byte, err error)
	Encrypt(message []byte) (encrypted []byte, err error)
	CanSign() bool
	CanEncrypt() bool
}

const initialKeyVersion = 999

func NewPinTanKeyName(bankId BankId, userId string, keyType string) *KeyName {
	return &KeyName{
		BankID:     bankId,
		UserID:     userId,
		KeyType:    keyType,
		KeyNumber:  0,
		KeyVersion: 0,
	}
}

func NewInitialKeyName(countryCode int, bankId, userId string, keyType string) *KeyName {
	return &KeyName{
		BankID:     BankId{CountryCode: countryCode, ID: bankId},
		UserID:     userId,
		KeyType:    keyType,
		KeyNumber:  999,
		KeyVersion: 999,
	}
}

type KeyName struct {
	BankID     BankId
	UserID     string
	KeyType    string
	KeyNumber  int
	KeyVersion int
}

func (k *KeyName) IsInitial() bool {
	return k.KeyNumber == 999 && k.KeyVersion == 999
}

func (k *KeyName) SetInitial() {
	k.KeyNumber = 999
	k.KeyVersion = 999
}

func NewPinKey(pin string, keyName *KeyName) *PinKey {
	return &PinKey{pin: pin, keyName: keyName}
}

type PinKey struct {
	pin     string
	keyName *KeyName
}

func (p *PinKey) KeyName() KeyName {
	return *p.keyName
}

func (p *PinKey) SetKeyNumber(number int) {
	p.keyName.KeyNumber = number
}

func (p *PinKey) SetKeyVersion(version int) {
	p.keyName.KeyVersion = version
}

func (p *PinKey) CanSign() bool {
	return true
}

func (p *PinKey) CanEncrypt() bool {
	return true
}

func (p *PinKey) Pin() string {
	return p.pin
}

func (p *PinKey) Sign(message []byte) ([]byte, error) {
	return []byte(p.pin), nil
}

func (p *PinKey) Encrypt(message []byte) ([]byte, error) {
	encMessage := make([]byte, len(message))
	// Make a deep copy, just in case
	copy(encMessage, message)
	return encMessage, nil
}

func (p *PinKey) Decrypt(encryptedMessage []byte) ([]byte, error) {
	decMessage := make([]byte, len(encryptedMessage))
	// Make a deep copy, just in case
	copy(decMessage, encryptedMessage)
	return decMessage, nil
}

func GenerateSigningKey() (*PublicKey, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 768)
	if err != nil {
		return nil, err
	}
	p := PublicKey{
		Type:          "S",
		Modulus:       rsaKey.N.Bytes(),
		Exponent:      big.NewInt(int64(rsaKey.E)).Bytes(),
		rsaPrivateKey: rsaKey,
	}
	return &p, nil
}

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
