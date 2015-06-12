package hbci

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/ripemd160"
)

const initializationVector = "\x01\x23\x45\x67\x89\xAB\xCD\xEF\xFE\xDC\xBA\x98\x76\x54\x32\x10\xF0\xE1\xD2\xC3"
const hashPadding = "\x00\x00\x00\x00"

func MessageHashSum(message fmt.Stringer) []byte {
	h := ripemd160.New()
	io.WriteString(h, initializationVector)
	io.WriteString(h, message.String())
	hash := h.Sum(nil)
	hash = append(hash, []byte(hashPadding)...)
	return hash
}

func SignMessageHash(messageHash []byte) ([]byte, error) {
	return nil, nil

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
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(message.String()))

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
