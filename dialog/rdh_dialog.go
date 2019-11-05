package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
)

// NewRDHDialog creates a dialog to use with cardreader flow
func NewRDHDialog(bankID domain.BankID, hbciURL string, clientID string, hbciVersion segment.HBCIVersion, productName string) Dialog {
	key, err := domain.GenerateSigningKey()
	if err != nil {
		panic(err)
	}
	signingKey := domain.NewRSAKey(key, domain.NewInitialKeyName(bankID.CountryCode, bankID.ID, clientID, "S"))
	provider := message.NewRDHSignatureProvider(signingKey, 12345)
	d := &rdhDialog{
		dialog: newDialog(bankID, hbciURL, clientID, hbciVersion, productName, provider, nil),
	}
	return d
}

type rdhDialog struct {
	*dialog
}
