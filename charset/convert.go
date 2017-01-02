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

// ToUTF8 decodes the provided buffer from the ISO8859_1 encoding to UTF8.
// The result is returned as string.
// Any error from the decoder will cause a panic.
//
// This function hasn't change its signature as it is used in almost any package.
// To address the issues with panicing errors this package will most likely get
// a new decoding function which return possible decoding errors.
func ToUTF8(iso8859_1Buf []byte) string {
	decoder := charmap.ISO8859_1.NewDecoder()
	// TODO: propagate errors
	return string(mustConvert(decoder.Bytes(iso8859_1Buf)))
}

// ToISO8859_1 encodes the provided utf8String to the ISO8859_1 encoding.
// The result is returned as byte slice.
// Any error from the encoder will cause a panic.
//
// This function hasn't change its signature as it is used in almost any package.
// To address the issues with panicing errors this package will most likely get
// a new encoding function which return possible encoding errors.
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
