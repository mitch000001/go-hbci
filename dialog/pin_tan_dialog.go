package dialog

import (
	"encoding/base64"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
	https "github.com/mitch000001/go-hbci/transport/https"
	middleware "github.com/mitch000001/go-hbci/transport/middleware"
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
	var dialogTransport transport.Transport
	dialogTransport = https.New()
	dialogTransport = middleware.Base64Encoding(base64.StdEncoding)(dialogTransport)
	d.transport = dialogTransport
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
