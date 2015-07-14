package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
)

func NewRDHDialog(bankId domain.BankId, hbciUrl string, clientId string) *rdhDialog {
	key, err := domain.GenerateSigningKey()
	if err != nil {
		panic(err)
	}
	signingKey := domain.NewRSAKey(key, domain.NewInitialKeyName(bankId.CountryCode, bankId.ID, clientId, "S"))
	provider := message.NewRDHSignatureProvider(signingKey)
	d := &rdhDialog{
		dialog:      newDialog(bankId, hbciUrl, clientId, provider, nil),
		SigningKey:  signingKey,
		SignatureID: 12345,
	}
	return d
}

type rdhDialog struct {
	*dialog
	SignatureID int
	SigningKey  domain.Key
}
