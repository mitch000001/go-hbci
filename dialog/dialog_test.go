package dialog

import (
	"fmt"
	"strings"
)

func encryptedTestMessage(encryptedData ...string) []byte {
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
