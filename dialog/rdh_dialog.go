package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
)

func NewRDHDialog(bankId domain.BankId, hbciUrl string, clientId string, hbciVersion int) *rdhDialog {
	key, err := domain.GenerateSigningKey()
	if err != nil {
		panic(err)
	}
	signingKey := domain.NewRSAKey(key, domain.NewInitialKeyName(bankId.CountryCode, bankId.ID, clientId, "S"))
	provider := message.NewRDHSignatureProvider(signingKey, 12345)
	d := &rdhDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, hbciVersion, provider, nil),
	}
	return d
}

type rdhDialog struct {
	*dialog
}
