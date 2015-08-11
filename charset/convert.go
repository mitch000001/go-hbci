package charset

import (
	"bytes"
	"io"
)

func ToUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func ToISO8859_1(utf8String string) []byte {
	buf := make([]byte, 0)
	runes := bytes.Runes([]byte(utf8String))
	for _, r := range runes {
		buf = append(buf, byte(r))
	}
	return buf
}

func NewISO8859_1Writer() io.Writer {
	var b bytes.Buffer
	return &iso8859_1Writer{&b}
}

type iso8859_1Writer struct {
	*bytes.Buffer
}

func (i *iso8859_1Writer) Write(p []byte) (int, error) {
	iso8859_1Bytes := ToISO8859_1(string(p))
	return i.Buffer.Write(iso8859_1Bytes)
}
