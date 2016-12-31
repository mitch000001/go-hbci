package charset

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/charmap"
)

func mustConvert(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}

func ToUtf8(iso8859_1_buf []byte) string {
	decoder := charmap.ISO8859_1.NewDecoder()
	// TODO: propagate errors
	return string(mustConvert(decoder.Bytes(iso8859_1_buf)))
}

func ToISO8859_1(utf8String string) []byte {
	encoder := charmap.ISO8859_1.NewEncoder()
	// TODO: propagate errors
	return mustConvert(encoder.Bytes([]byte(utf8String)))
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
