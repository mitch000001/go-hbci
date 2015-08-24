package dialog

import (
	"fmt"
	"io"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
)

func NewPinTanDialog(bankId domain.BankId, hbciUrl string, userId string, hbciVersion segment.HBCIVersion) *PinTanDialog {
	pinKey := domain.NewPinKey("", domain.NewPinTanKeyName(bankId, userId, "S"))
	signatureProvider := message.NewPinTanSignatureProvider(pinKey, "0", hbciVersion)
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
	d.transport = transport.NewHttpsTransport()
	return d
}

type PinTanDialog struct {
	*dialog
}

func (d *PinTanDialog) SetPin(pin string) {
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "S"))
	d.signatureProvider = message.NewPinTanSignatureProvider(pinKey, d.ClientSystemID, d.hbciVersion)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.UserID, "V"))
	d.cryptoProvider = message.NewPinTanCryptoProvider(pinKey, d.ClientSystemID)
}

func (d *PinTanDialog) CommunicationAccess() (string, error) {
	comm := message.NewCommunicationAccessMessage(d.BankID, d.BankID, 5, "")
	comm.Header = segment.NewMessageHeaderSegment(0, 220, initialDialogID, 1)
	comm.End = segment.NewMessageEndSegment(3, 1)
	comm.SetSize()

	response, err := d.request(comm)
	if err != nil && err != io.EOF {
		return "", err
	}

	responseBytes, err := response.(*message.DecryptedMessage).MarshalHBCI()
	if err != nil {
		return "", err
	}
	return string(responseBytes), nil
}

func (d *PinTanDialog) Anonymous(fn func() (string, error)) (string, error) {
	initMessage := message.NewDialogInitializationClientMessage(d.hbciVersion)
	messageNum := d.nextMessageNumber()
	initMessage.Header = segment.NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	initMessage.End = segment.NewMessageEndSegment(8, messageNum)
	initMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.clientID, "0", false)
	initMessage.ProcessingPreparation = segment.NewProcessingPreparationSegment(0, 0, 1)
	initMessage.SetNumbers()
	initMessage.SetSize()
	marshaledMessage, err := initMessage.MarshalHBCI()
	if err != nil {
		return "", err
	}

	response, err := d.post(marshaledMessage)
	if err != nil && err != io.EOF {
		return "", err
	}

	_, err = fn()
	if err != nil && err != io.EOF {
		return "", err
	}

	dialogEnd := &message.DialogFinishingMessage{
		DialogEnd: segment.NewDialogEndSegment(d.dialogID),
	}
	dialogEnd.BasicMessage = d.newBasicMessage(dialogEnd)
	dialogEnd.SetNumbers()
	dialogEnd.SetSize()
	marshaledEndMessage, err := dialogEnd.MarshalHBCI()
	if err != nil {
		return "", err
	}
	response, err = d.post(marshaledEndMessage)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("Error while ending dialog: %v", err)
	}

	return string(response), nil
}
