package dialog

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/transport"
)

func NewFINTS3PinTanDialog(bankId domain.BankId, hbciUrl string, userId string) *PinTanDialog {
	pinKey := domain.NewPinKey("", domain.NewPinTanKeyName(bankId, userId, "S"))
	signatureProvider := message.NewFINTS3PinTanSignatureProvider(pinKey, "0")
	pinKey = domain.NewPinKey("", domain.NewPinTanKeyName(bankId, userId, "V"))
	cryptoProvider := message.NewFINTS3PinTanEncryptionProvider(pinKey, "0")
	d := NewPinTanDialog(
		bankId,
		hbciUrl,
		userId,
	)
	d.signatureProvider = signatureProvider
	d.cryptoProvider = cryptoProvider
	d.transport = transport.NewHttpsTransport()
	return d
}

func (d *PinTanDialog) SetFINTS3Pin(pin string) {
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "S"))
	d.signatureProvider = message.NewFINTS3PinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "V"))
	d.cryptoProvider = message.NewFINTS3PinTanEncryptionProvider(pinKey, d.ClientSystemID)
}
