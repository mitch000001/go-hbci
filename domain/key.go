package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

// Key provides an interface to an encryption/signing key
type Key interface {
	// KeyName returns the KeyName
	KeyName() KeyName
	// SetKeyNumber sets the key number in the KeyName
	SetKeyNumber(number int)
	// SetKeyVersion sets the key version in the KeyName
	SetKeyVersion(version int)
	// Sign signs message
	Sign(message []byte) (signature []byte, err error)
	// Encrypt encrypts message
	Encrypt(message []byte) (encrypted []byte, err error)
	// CanSign returns true if the key can be used for signing
	CanSign() bool
	// CanEncrypt returns true if the key can be used for encryption
	CanEncrypt() bool
}

const initialKeyVersion = 999

// NewPinTanKeyName returns a new KeyName for the pin/tan flow
func NewPinTanKeyName(bankID BankID, userID string, keyType string) *KeyName {
	return &KeyName{
		BankID:     bankID,
		UserID:     userID,
		KeyType:    keyType,
		KeyNumber:  0,
		KeyVersion: 0,
	}
}

// NewInitialKeyName represents a KeyName ready to use for initial communication
func NewInitialKeyName(countryCode int, bankID, userID string, keyType string) *KeyName {
	return &KeyName{
		BankID:     BankID{CountryCode: countryCode, ID: bankID},
		UserID:     userID,
		KeyType:    keyType,
		KeyNumber:  999,
		KeyVersion: 999,
	}
}

// KeyName provides data about a given key
type KeyName struct {
	BankID     BankID
	UserID     string
	KeyType    string
	KeyNumber  int
	KeyVersion int
}

// IsInitial returns true if the KeyName represents an initial KeyName, false otherwise
func (k *KeyName) IsInitial() bool {
	return k.KeyNumber == 999 && k.KeyVersion == 999
}

// SetInitial resets the KeyName to reflect an initial KeyName
func (k *KeyName) SetInitial() {
	k.KeyNumber = 999
	k.KeyVersion = 999
}

// NewPinKey returns a new PinKey
func NewPinKey(pin string, keyName *KeyName) *PinKey {
	return &PinKey{pin: pin, keyName: keyName}
}

// PinKey represents a Key used for pin/tan flow
type PinKey struct {
	pin     string
	keyName *KeyName
}

// KeyName returns the KeyName
func (p *PinKey) KeyName() KeyName {
	return *p.keyName
}

// SetKeyNumber sets the key number in the KeyName
func (p *PinKey) SetKeyNumber(number int) {
	p.keyName.KeyNumber = number
}

// SetKeyVersion sets the key version in the KeyName
func (p *PinKey) SetKeyVersion(version int) {
	p.keyName.KeyVersion = version
}

// CanSign returns true if the key can be used for signing
func (p *PinKey) CanSign() bool {
	return true
}

// CanEncrypt returns true if the key can be used for encryption
func (p *PinKey) CanEncrypt() bool {
	return true
}

// Pin returns the pin within this key
func (p *PinKey) Pin() string {
	return p.pin
}

// Sign signs message
func (p *PinKey) Sign(message []byte) ([]byte, error) {
	return []byte(p.pin), nil
}

// Encrypt encryptes the message
func (p *PinKey) Encrypt(message []byte) ([]byte, error) {
	encMessage := make([]byte, len(message))
	// Make a deep copy, just in case
	copy(encMessage, message)
	return encMessage, nil
}

// Decrypt decrypts the encryptedMessage
func (p *PinKey) Decrypt(encryptedMessage []byte) ([]byte, error) {
	decMessage := make([]byte, len(encryptedMessage))
	// Make a deep copy, just in case
	copy(decMessage, encryptedMessage)
	return decMessage, nil
}

// GenerateSigningKey generates a new signing key
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

// NewRSAKey returns a new RSA key
func NewRSAKey(pubKey *PublicKey, keyName *KeyName) *RSAKey {
	return &RSAKey{PublicKey: pubKey, keyName: keyName}
}

// RSAKey represents a public RSA key which implements the Key interface
type RSAKey struct {
	*PublicKey
	keyName *KeyName
}

// KeyName returns the KeyName
func (r *RSAKey) KeyName() KeyName {
	return *r.keyName
}

// SetKeyNumber sets the key number in the KeyName
func (r *RSAKey) SetKeyNumber(number int) {
	r.keyName.KeyNumber = number
}

// SetKeyVersion sets the key version in the KeyName
func (r *RSAKey) SetKeyVersion(version int) {
	r.keyName.KeyVersion = version
}

// CanSign returns true if the key can be used for signing
func (r *RSAKey) CanSign() bool {
	return r.PublicKey.rsaPrivateKey != nil
}

// CanEncrypt returns true if the key can be used for encryption
func (r *RSAKey) CanEncrypt() bool {
	return r.PublicKey.rsaPublicKey != nil
}

// NewEncryptionKey creates a new RSA encryption key
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

// PublicKey represents a key which can either embed a private or public RSA key
type PublicKey struct {
	Type          string
	Modulus       []byte
	Exponent      []byte
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

// SigningKey returns the RSA private key to sign with, or nil when not set
func (p *PublicKey) SigningKey() *rsa.PrivateKey {
	return p.rsaPrivateKey
}

// Sign signs message with the private key
func (p *PublicKey) Sign(message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, p.rsaPrivateKey, 0, message)
}

// Encrypt encryptes the message with the public key
func (p *PublicKey) Encrypt(message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, p.rsaPublicKey, message)
}
