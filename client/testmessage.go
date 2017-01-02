package client

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mitch000001/go-hbci/charset"
)

func encryptedTestMessage(dialogID string, segments ...string) []byte {
	encodedSegments := charset.ToISO8859_1(strings.Join(segments, ""))
	encryptionHeader := "HNVSK:998:2:+998+1+1::0+1:20150713:173634+2:2:13:@8@\x00\x00\x00\x00\x00\x00\x00\x00:5:1:+280:10000000:12345:V:0:0+0+'"
	encryptionData := fmt.Sprintf("HNVSD:999:1:+@%d@%s'", len(encodedSegments), encodedSegments)
	messageEnd := fmt.Sprintf("HNHBS:%d:1:+1'", len(segments)+1)
	messageHeader := fmt.Sprintf("HNHBK:1:3+%012d+220+%s+1+'", 31+len(dialogID)+len(encryptionHeader)+len(encryptionData)+len(messageEnd), dialogID)
	encryptedMessage := [][]byte{
		charset.ToISO8859_1(messageHeader),
		charset.ToISO8859_1(encryptionHeader),
		[]byte(encryptionData),
		charset.ToISO8859_1(messageEnd),
	}
	return bytes.Join(encryptedMessage, []byte(""))
}
