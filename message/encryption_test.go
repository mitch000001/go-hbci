package message

import (
	"fmt"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

func TestEncryptedPinTanMessageDecrypt(t *testing.T) {
	keyName := domain.NewPinTanKeyName(domain.BankId{CountryCode: 280, ID: "1"}, "userID", "V")
	pinKey := domain.NewPinKey("abcde", keyName)

	provider := NewPinTanCryptoProvider(pinKey, "clientSystemID")

	syncSegment := "HISYN:2:3:8+newClientSystemID'"
	acknowledgement := "HIRMG:2:2:1+0100::Dialog beendet'"

	body := fmt.Sprintf("%s%s", acknowledgement, syncSegment)

	header := segment.NewMessageHeaderSegment(1, 220, "abcde", 1)
	end := segment.NewMessageEndSegment(4, 1)
	encryptedMessage := NewEncryptedMessage(header, end, segment.HBCI220)
	encryptedMessage.EncryptionHeader = segment.HBCI220.PinTanEncryptionHeader("0", *keyName)
	provider.WriteEncryptionHeader(encryptedMessage.EncryptionHeader)
	encryptedMessage.EncryptedData = segment.NewEncryptedDataSegment([]byte(body))

	decryptedMessage, err := encryptedMessage.Decrypt(provider)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actualSyncSegment := decryptedMessage.FindMarshaledSegment("HISYN")

	if syncSegment != string(actualSyncSegment) {
		t.Logf("Expected decrypted message to include SynchronisationResponse, but had not\n")
		t.Fail()
	}
}
