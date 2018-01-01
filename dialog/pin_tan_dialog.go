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

// NewPinTanDialog creates a new dialog to use for pin/tan transport
func NewPinTanDialog(bankID domain.BankID, hbciURL string, userID string, hbciVersion segment.HBCIVersion) *PinTanDialog {
	pinKey := domain.NewPinKey("", domain.NewPinTanKeyName(bankID, userID, "S"))
	signatureProvider := message.NewPinTanSignatureProvider(pinKey, "0")
	pinKey = domain.NewPinKey("", domain.NewPinTanKeyName(bankID, userID, "V"))
	cryptoProvider := message.NewPinTanCryptoProvider(pinKey, "0")
	d := &PinTanDialog{
		dialog: newDialog(
			bankID,
			hbciURL,
			userID,
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

// PinTanDialog represents a dialog to use in pin/tan flow with HTTPS transport
type PinTanDialog struct {
	*dialog
}

// SetPin lets the user reset the pin after creation
func (d *PinTanDialog) SetPin(pin string) {
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "S"))
	d.signatureProvider = message.NewPinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "V"))
	d.cryptoProvider = message.NewPinTanCryptoProvider(pinKey, d.ClientSystemID)
}
