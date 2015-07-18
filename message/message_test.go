package message

import (
	"fmt"
	"strings"
	"testing"
)

func EncryptedTestMessage(encryptedData ...string) []byte {
	encryptionHeader := "HNVSK:998:2:+998+1+1::0+1:20150713:173634+2:2:13:@8@\x00\x00\x00\x00\x00\x00\x00\x00:5:1:+280:10000000:12345:V:0:0+0+'"
	encryptionData := fmt.Sprintf("HNVSD:999:1:+@%d@%s'", len(strings.Join(encryptedData, "")), strings.Join(encryptedData, ""))
	messageEnd := fmt.Sprintf("HNHBS:%d:1:+1'", len(encryptedData)+1)
	messageHeader := fmt.Sprintf("HNHBK:1:3+%012d+220+abcde+1+'", 36+len(encryptionHeader)+len(encryptionData)+len(messageEnd))
	encryptedMessage := []string{
		messageHeader,
		encryptionHeader,
		encryptionData,
		messageEnd,
	}
	return []byte(strings.Join(encryptedMessage, ""))
}

func TestBasicMessageUnmarshalHBCI(t *testing.T) {
	//testMessage := "HNHBK:1:3+000000000273+220+9631542215616260+1+9631542215616260:1'HIRMG:2:2+9010::Nachricht ist komplett nicht bearbeitet (HBMSG=10319)+9110::Unbekannter Aufbau (HBMSG=10311)+9800::Dialog abgebrochen (HBMSG=10321)'HIRMS:3:2:998+9160:2:Syntaxfehler (HBMSG=10001)'HNHBS:4:1+1'"

	//var b basicMessage

	//err := b.UnmarshalHBCI([]byte(testMessage))

	//if err != nil {
	//t.Logf("Expected no error, got %T:%v\n", err, err)
	//t.Fail()
	//}
}
