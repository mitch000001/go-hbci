package dialog

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
)

func NewPinTanDialog(bankId domain.BankId, hbciUrl string, clientId string) *pinTanDialog {
	signatureProvider := message.NewPinTanSignatureProvider(nil, "")
	encryptionProvider := message.NewPinTanEncryptionProvider(nil, "")
	d := &pinTanDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, signatureProvider, encryptionProvider),
	}
	return d
}

type pinTanDialog struct {
	*dialog
	pin           string
	pinTanKeyName domain.KeyName
}

func (d *pinTanDialog) SetPin(pin string) {
	d.pin = pin
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.ClientID, "S"))
	d.signatureProvider = message.NewPinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.ClientID, "V"))
	d.encryptionProvider = message.NewPinTanEncryptionProvider(pinKey, d.ClientSystemID)
}

func (d *pinTanDialog) Init() (string, error) {
	initMessage := message.NewDialogInitializationClientMessage()
	messageNum := d.nextMessageNumber()
	initMessage.Header = segment.NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	initMessage.End = segment.NewMessageEndSegment(8, messageNum)
	initMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.ClientID, initialClientSystemID, true)
	initMessage.ProcessingPreparation = segment.NewProcessingPreparationSegment(0, 0, 1)
	controlRef := "1"
	initMessage.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	initMessage.SignatureEnd = segment.NewSignatureEndSegment(7, controlRef)
	initMessage.SetNumbers()
	err := initMessage.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	encryptedInitMessage, err := initMessage.Encrypt(d.encryptionProvider)
	if err != nil {
		return "", err
	}
	encryptedInitMessage.SetSize()
	marshaledMessage, err := encryptedInitMessage.MarshalHBCI()
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
	dialogEnd.SignatureEnd = segment.NewSignatureEndSegment(7, controlRef)
	dialogEnd.SetNumbers()
	err = dialogEnd.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	dialogEnd.SetSize()
	encryptedDialogEnd, err := dialogEnd.Encrypt(d.encryptionProvider)
	if err != nil {
		return "", err
	}
	marshaledEndMessage, err := encryptedDialogEnd.MarshalHBCI()
	if err != nil {
		return "", err
	}
	response, err = d.post(marshaledEndMessage)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(response), nil
}

func (d *pinTanDialog) SyncClientSystemID() (string, error) {
	syncMessage := new(message.SynchronisationMessage)
	messageNum := d.nextMessageNumber()
	syncMessage.BasicClientMessage = message.NewBasicClientMessage(syncMessage)
	syncMessage.Header = segment.NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	syncMessage.End = segment.NewMessageEndSegment(8, messageNum)
	syncMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.ClientID, initialClientSystemID, true)
	syncMessage.ProcessingPreparation = segment.NewProcessingPreparationSegment(0, 0, 1)
	syncMessage.Sync = segment.NewSynchronisationSegment(0)
	controlRef := "1"
	syncMessage.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	syncMessage.SignatureEnd = segment.NewSignatureEndSegment(7, controlRef)
	syncMessage.SetNumbers()
	err := syncMessage.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	d.encryptionProvider.SetClientSystemID(initialClientSystemID)
	encryptedSyncMessage, err := syncMessage.Encrypt(d.encryptionProvider)
	if err != nil {
		return "", err
	}
	encryptedSyncMessage.SetSize()
	marshaledMessage, err := encryptedSyncMessage.MarshalHBCI()
	if err != nil {
		return "", err
	}

	response, err := d.post(marshaledMessage)
	if err != nil {
		return "", err
	}

	encMessage, err := extractEncryptedMessage(response)
	if err != nil {
		return "", fmt.Errorf("Error while extracting encrypted message: %v", err)
	}

	messageHeader := encMessage.MessageHeader()
	if messageHeader == nil {
		return "", fmt.Errorf("Malformed response: %q", response)
	}
	newDialogId := messageHeader.DialogID.Val()

	encryptedData := encMessage.EncryptedData
	if encryptedData != nil {
		encExtractor := segment.NewSegmentExtractor(encryptedData.Data.Val())
		data, err := encExtractor.Extract()
		if err != nil {
			return "", fmt.Errorf("Error while decrypting message: %v", err)
		}
		fmt.Printf("Encrypted data: \n%s\n", bytes.Join(data, []byte("\n")))
		syncResponse := encExtractor.FindSegment("HISYN")
		if syncResponse != nil {
			syncSegment := &segment.SynchronisationResponseSegment{}
			err = syncSegment.UnmarshalHBCI(syncResponse)
			if err != nil {
				return "", fmt.Errorf("Error while unmarshaling sync response: %v", err)
			}
			d.ClientSystemID = syncSegment.ClientSystemID.Val()
			d.signatureProvider.SetClientSystemID(d.ClientSystemID)
			d.encryptionProvider.SetClientSystemID(d.ClientSystemID)
		}

		accountData := encExtractor.FindSegments("HIUPD")
		if accountData != nil {
			for _, acc := range accountData {
				infoSegment := &segment.AccountInformationSegment{}
				err = infoSegment.UnmarshalHBCI(acc)
				if err != nil {
					return "", fmt.Errorf("Error while unmarshaling Accounts: %v", err)
				}
				d.Accounts = append(d.Accounts, infoSegment.Account())
			}
		}
	} else {
		return "", fmt.Errorf("Expected encrypted message, but was not:\n%q", response)
	}

	dialogEnd := d.dialogEnd(newDialogId)
	dialogEnd.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	dialogEnd.SignatureEnd = segment.NewSignatureEndSegment(7, controlRef)
	dialogEnd.SetNumbers()
	err = dialogEnd.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	dialogEnd.SetSize()
	encryptedDialogEnd, err := dialogEnd.Encrypt(d.encryptionProvider)
	if err != nil {
		return "", err
	}
	marshaledEndMessage, err := encryptedDialogEnd.MarshalHBCI()
	if err != nil {
		return "", err
	}
	response, err = d.post(marshaledEndMessage)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func extractEncryptedMessage(response []byte) (*message.EncryptedMessage, error) {
	extractor := segment.NewSegmentExtractor(response)
	_, err := extractor.Extract()
	if err != nil {
		return nil, err
	}

	messageHeader := extractor.FindSegment("HNHBK")
	if messageHeader == nil {
		return nil, fmt.Errorf("Malformed response: missing Message Header")
	}
	header := &segment.MessageHeaderSegment{}
	err = header.UnmarshalHBCI(messageHeader)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling message header: %v", err)
	}
	// TODO: parse messageEnd

	encMessage := message.NewEncryptedMessage(header, nil)

	encryptedData := extractor.FindSegment("HNVSD")
	if encryptedData != nil {
		encSegment := &segment.EncryptedDataSegment{}
		err = encSegment.UnmarshalHBCI(encryptedData)
		if err != nil {
			return nil, fmt.Errorf("Error while unmarshaling encrypted data: %v", err)
		}
		encMessage.EncryptedData = encSegment
	} else {
		return nil, fmt.Errorf("Malformed response: missing encrypted data")
	}
	return encMessage, nil
}

func (d *pinTanDialog) CommunicationAccess() (string, error) {
	comm := message.NewCommunicationAccessMessage(d.BankID, d.BankID, 5, "")
	comm.Header = segment.NewMessageHeaderSegment(0, 220, initialDialogID, 1)
	comm.End = segment.NewMessageEndSegment(3, 1)
	comm.SetSize()
	marshaled, err := comm.MarshalHBCI()
	if err != nil {
		return "", err
	}
	fmt.Printf("Marshaled: %q\n", string(marshaled))
	response, err := d.post(marshaled)
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(response), nil
}

func (d *pinTanDialog) Anonymous(fn func() (string, error)) (string, error) {
	initMessage := message.NewDialogInitializationClientMessage()
	messageNum := d.nextMessageNumber()
	initMessage.Header = segment.NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	initMessage.End = segment.NewMessageEndSegment(8, messageNum)
	initMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.ClientID, "0", false)
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

	fmt.Printf("Response: %q\n", strings.Split(string(response), "'"))

	res, err := fn()
	if err != nil && err != io.EOF {
		return "", err
	}
	fmt.Printf("Response: %q\n", strings.Split(res, "'"))

	dialogEnd := d.dialogEnd("0")
	dialogEnd.SetNumbers()
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
