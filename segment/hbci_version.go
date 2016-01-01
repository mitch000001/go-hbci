package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
)

var SupportedHBCIVersions = map[int]HBCIVersion{
	220: HBCI220,
	300: FINTS300,
}

type HBCIVersion struct {
	version                       int
	PinTanEncryptionHeader        func(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment
	RDHEncryptionHeader           func(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment
	SignatureHeader               func() *SignatureHeaderSegment
	PinTanSignatureHeader         func(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment
	RDHSignatureHeader            func(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment
	SignatureEnd                  func() *SignatureEndSegment
	SynchronisationRequest        func(modus int) *SynchronisationRequestSegment
	AccountBalanceRequest         func(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest
	AccountTransactionRequest     func(account domain.AccountConnection, allAccounts bool) *AccountTransactionRequestSegment
	SepaAccountTransactionRequest func(account domain.InternationalAccountConnection, allAccounts bool) *AccountTransactionRequestSegment
	StatusProtocolRequest         func(from, to time.Time, maxEntries int, continuationReference string) StatusProtocolRequest
}

func (v HBCIVersion) Version() int {
	return v.version
}
