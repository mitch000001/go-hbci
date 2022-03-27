package dialog

import (
	"encoding/base64"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/internal"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
	https "github.com/mitch000001/go-hbci/transport/https"
	middleware "github.com/mitch000001/go-hbci/transport/middleware"
)

// Config contains the configuration of a PinTanDialog
type Config struct {
	BankID      domain.BankID
	HBCIURL     string
	UserID      string
	HBCIVersion segment.HBCIVersion
	Transport   transport.Transport
}

// NewPinTanDialog creates a new dialog to use for pin/tan transport
func NewPinTanDialog(config Config) *PinTanDialog {
	pinKey := domain.NewPinKey("", domain.NewPinTanKeyName(config.BankID, config.UserID, domain.KeyTypeSigning))
	signatureProvider := message.NewPinTanSignatureProvider(pinKey, initialClientSystemID)
	pinKey = domain.NewPinKey("", domain.NewPinTanKeyName(config.BankID, config.UserID, domain.KeyTypeEncryption))
	cryptoProvider := message.NewPinTanCryptoProvider(pinKey, initialClientSystemID)
	d := &PinTanDialog{
		dialog: newDialog(
			config.BankID,
			config.HBCIURL,
			config.UserID,
			config.HBCIVersion,
			signatureProvider,
			cryptoProvider,
		),
	}

	var dialogTransport transport.Transport
	if config.Transport == nil {
		dialogTransport = https.New()
	} else {
		dialogTransport = config.Transport
	}
	dialogTransport = middleware.Base64Encoding(base64.StdEncoding)(dialogTransport)
	dialogTransport = middleware.Logging(internal.Debug)(dialogTransport)
	d.transport = dialogTransport
	return d
}

// PinTanDialog represents a dialog to use in pin/tan flow with HTTPS transport
type PinTanDialog struct {
	*dialog
}

// SetPin lets the user reset the pin after creation
func (d *PinTanDialog) SetPin(pin string) {
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, domain.KeyTypeSigning))
	d.signatureProvider = message.NewPinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, domain.KeyTypeEncryption))
	d.cryptoProvider = message.NewPinTanCryptoProvider(pinKey, d.ClientSystemID)
}
