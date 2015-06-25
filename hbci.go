package hbci

import "crypto/rsa"

const productName = "go-hbci library"
const productVersion = "0.0.1"

type Marshaler interface {
	MarshalHBCI() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalHBCI([]byte) error
}

// supportedUnmarshaler maps segment IDs to Unmarshalers
var supportedUnmarshaler = map[string]Unmarshaler{}

func MakeCall() string {
	return ""
}

type Client struct {
	rsaKey *rsa.PublicKey
}
