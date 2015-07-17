package message

import (
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

func TestEncryptedPinTanMessageDecrypt(t *testing.T) {
	keyName := domain.NewPinTanKeyName(domain.BankId{CountryCode: 280, ID: "1"}, "userID", "V")
	pinKey := domain.NewPinKey("abcde", keyName)

	provider := NewPinTanCryptoProvider(pinKey, "clientSystemID")

	encMessage := "HISYN:2:3:8+newClientSystemID'"

	encryptedMessage := NewEncryptedPinTanMessage("clientSystemID", *keyName, []byte(encMessage))

	decryptedMessage, err := encryptedMessage.Decrypt(provider)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	syncSegment := decryptedMessage.FindSegment("HISYN")

	if string(syncSegment) != encMessage {
		t.Logf("Expected decrypted message to include SynchronisationResponse, but had not\n")
		t.Fail()
	}
}
