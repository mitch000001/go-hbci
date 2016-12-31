package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
)

func NewPinTanDialog(bankId domain.BankId, hbciUrl string, userId string, hbciVersion segment.HBCIVersion) *PinTanDialog {
	pinKey := domain.NewPinKey("", domain.NewPinTanKeyName(bankId, userId, "S"))
	signatureProvider := message.NewPinTanSignatureProvider(pinKey, "0")
	pinKey = domain.NewPinKey("", domain.NewPinTanKeyName(bankId, userId, "V"))
	cryptoProvider := message.NewPinTanCryptoProvider(pinKey, "0")
	d := &PinTanDialog{
		dialog: newDialog(
			bankId,
			hbciUrl,
			userId,
			hbciVersion,
			signatureProvider,
			cryptoProvider,
		),
	}
	d.transport = transport.NewHttpsBase64Transport()
	return d
}

type PinTanDialog struct {
	*dialog
}

func (d *PinTanDialog) SetPin(pin string) {
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "S"))
	d.signatureProvider = message.NewPinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "V"))
	d.cryptoProvider = message.NewPinTanCryptoProvider(pinKey, d.ClientSystemID)
}
