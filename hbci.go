package hbci

import "crypto/rsa"

func MakeCall() string {
	return ""
}

func InitializeDialog() {}

type Client struct {
	rsaKey *rsa.PublicKey
}
