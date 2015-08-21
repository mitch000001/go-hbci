package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
)

func NewRDHDialog(bankId domain.BankId, hbciUrl string, clientId string, hbciVersion segment.Version) *rdhDialog {
	key, err := domain.GenerateSigningKey()
	if err != nil {
		panic(err)
	}
	signingKey := domain.NewRSAKey(key, domain.NewInitialKeyName(bankId.CountryCode, bankId.ID, clientId, "S"))
	provider := message.NewRDHSignatureProvider(signingKey, 12345, hbciVersion)
	d := &rdhDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, hbciVersion, provider, nil),
	}
	return d
}

type rdhDialog struct {
	*dialog
}
