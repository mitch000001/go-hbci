package hbci

import "crypto/rsa"

const productName = "go-hbci library"
const productVersion = "0.0.1"

func MakeCall() string {
	return ""
}

func InitializeDialog() {}

type Client struct {
	rsaKey *rsa.PublicKey
}
