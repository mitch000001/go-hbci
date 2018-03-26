package hbci

import "github.com/mitch000001/go-hbci/internal"

const Version = "0.1.3"

func SetDebugMode(debug bool) {
	internal.SetDebugMode(debug)
}

func SetInfoLog(info bool) {
	internal.SetInfoLog(info)
}

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
