package hbci

import (
	"fmt"
	"io"
	"strings"
)

func NewFINTS3PinTanDialog(bankId BankId, hbciUrl string, clientId string) *pinTanDialog {
	signatureProvider := NewFINTS3PinTanSignatureProvider(nil)
	encryptionProvider := NewFINTS3PinTanEncryptionProvider(nil, "")
	d := &pinTanDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, signatureProvider, encryptionProvider),
	}
	return d
}

func (d *pinTanDialog) SetFINTS3Pin(pin string) {
	d.pin = pin
	pinKey := NewPinKey(pin, NewPinTanKeyName(d.BankID, d.ClientID, "S"))
	d.signingKey = pinKey
	d.signatureProvider = NewFINTS3PinTanSignatureProvider(pinKey)
	d.encryptionProvider = NewFINTS3PinTanEncryptionProvider(pinKey, d.ClientSystemID)
}

func (d *pinTanDialog) FINTS3SyncClientSystemID() (string, error) {
	syncMessage := new(SynchronisationMessage)
	messageNum := d.nextMessageNumber()
	syncMessage.basicClientMessage = newBasicClientMessage(syncMessage)
	syncMessage.Header = NewMessageHeaderSegment(-1, 300, initialDialogID, messageNum)
	syncMessage.End = NewMessageEndSegment(8, messageNum)
	syncMessage.Identification = NewIdentificationSegment(d.BankID, d.ClientID, initialClientSystemID, true)
	syncMessage.ProcessingPreparation = NewProcessingPreparationSegment(0, 0, 1)
	syncMessage.Sync = NewSynchronisationSegment(0)
	controlRef := "1"
	syncMessage.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	syncMessage.SignatureEnd = NewFINTS3SignatureEndSegment(7, controlRef)
	syncMessage.SetNumbers()
	err := syncMessage.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	syncMessage.SetSize()
	encryptedSyncMessage, err := syncMessage.Encrypt(d.encryptionProvider)
	if err != nil {
		return "", err
	}
	marshaledMessage, err := encryptedSyncMessage.MarshalHBCI()
	if err != nil {
		return "", err
	}

	response, err := d.post(marshaledMessage)
	if err != nil && err != io.EOF {
		return "", err
	}
	fmt.Printf("Response: %q\n", strings.Split(string(response), "'"))

	dialogEnd := d.dialogEnd(initialDialogID)
	dialogEnd.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	dialogEnd.SignatureEnd = NewSignatureEndSegment(7, controlRef)
	dialogEnd.SetNumbers()
	err = dialogEnd.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	dialogEnd.SetSize()
	marshaledEndMessage, err := dialogEnd.MarshalHBCI()
	if err != nil {
		return "", err
	}
	response, err = d.post(marshaledEndMessage)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(response), nil
}
