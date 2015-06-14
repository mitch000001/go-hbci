package hbci

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"

	"golang.org/x/crypto/ripemd160"
)

const initializationVector = "\x01\x23\x45\x67\x89\xAB\xCD\xEF\xFE\xDC\xBA\x98\x76\x54\x32\x10\xF0\xE1\xD2\xC3"

func GenerateSigningKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 1024)
}

func MessageHashSum(message fmt.Stringer) []byte {
	h := ripemd160.New()
	io.WriteString(h, initializationVector)
	io.WriteString(h, message.String())
	return h.Sum(nil)
}

func SignMessageHash(messageHash []byte, key *rsa.PrivateKey) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, key, 0, messageHash)
}

func EncryptMessage(message fmt.Stringer) (string, error) {
	block, err := des.NewCipher([]byte(initializationVector))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, des.BlockSize+len(message.String()))
	iv := ciphertext[:des.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[des.BlockSize:], []byte(message.String()))

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return fmt.Sprintf("%x", ciphertext), nil
}

func NewKeyNameDataElement(countryCode int, bankId string, userId string, keyType string, keyNumber, keyVersion int) *KeyNameDataElement {
	a := &KeyNameDataElement{
		Bank:       NewBankIndentificationDataElementWithBankId(countryCode, bankId),
		UserID:     NewIdentificationDataElement(userId),
		KeyType:    NewAlphaNumericDataElement(keyType, 1),
		KeyNumber:  NewNumberDataElement(keyNumber, 3),
		KeyVersion: NewNumberDataElement(keyVersion, 3),
	}
	a.elementGroup = NewDataElementGroup(KeyNameDEG, 5, a)
	return a
}

type KeyNameDataElement struct {
	*elementGroup
	Bank       *BankIdentificationDataElement
	UserID     *IdentificationDataElement
	KeyType    *AlphaNumericDataElement
	KeyNumber  *NumberDataElement
	KeyVersion *NumberDataElement
}

func (k *KeyNameDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		k.Bank,
		k.UserID,
		k.KeyType,
		k.KeyNumber,
		k.KeyVersion,
	}
}
