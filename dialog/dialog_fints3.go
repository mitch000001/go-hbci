package hbci

import "github.com/mitch000001/go-hbci/message"

func NewFINTS3PinTanDialog(bankId BankId, hbciUrl string, clientId string) *pinTanDialog {
	signatureProvider := message.NewFINTS3PinTanSignatureProvider(nil)
	encryptionProvider := message.NewFINTS3PinTanEncryptionProvider(nil, "")
	d := &pinTanDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, signatureProvider, encryptionProvider),
	}
	return d
}

func (d *pinTanDialog) SetFINTS3Pin(pin string) {
	d.pin = pin
	pinKey := NewPinKey(pin, NewPinTanKeyName(d.BankID, d.ClientID, "S"))
	d.signingKey = pinKey
	d.signatureProvider = message.NewFINTS3PinTanSignatureProvider(pinKey)
	d.encryptionProvider = message.NewFINTS3PinTanEncryptionProvider(pinKey, d.ClientSystemID)
}
