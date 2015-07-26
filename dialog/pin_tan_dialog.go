package dialog

import (
	"fmt"
	"io"
	"strings"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
)

func NewPinTanDialog(bankId domain.BankId, hbciUrl string, clientId string) *PinTanDialog {
	d := &PinTanDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, nil, nil),
	}
	d.transport = transport.NewHttpsTransport()
	return d
}

type PinTanDialog struct {
	*dialog
	pin string
}

func (d *PinTanDialog) SetPin(pin string) {
	d.pin = pin
	pinKey := domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.ClientID, "S"))
	d.signatureProvider = message.NewPinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = domain.NewPinKey(pin, domain.NewPinTanKeyName(d.BankID, d.ClientID, "V"))
	d.cryptoProvider = message.NewPinTanCryptoProvider(pinKey, d.ClientSystemID)
}

func (d *PinTanDialog) Init() error {
	d.dialogID = initialDialogID
	d.messageCount = 0
	initMessage := message.NewDialogInitializationClientMessage()
	messageNum := d.nextMessageNumber()
	initMessage.Header = segment.NewMessageHeaderSegment(-1, 220, d.dialogID, messageNum)
	initMessage.End = segment.NewMessageEndSegment(8, messageNum)
	initMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.ClientID, d.ClientSystemID, true)
	initMessage.ProcessingPreparation = segment.NewProcessingPreparationSegment(d.BankParameterDataVersion(), d.UserParameterDataVersion(), d.Language)
	signedInitMessage, err := initMessage.Sign(d.signatureProvider)
	if err != nil {
		return err
	}
	encryptedInitMessage, err := signedInitMessage.Encrypt(d.cryptoProvider)
	if err != nil {
		return err
	}

	decryptedMessage, err := d.request(encryptedInitMessage)
	if err != nil {
		return fmt.Errorf("Error while initializing dialog: %v", err)
	}
	messageHeader := decryptedMessage.MessageHeader()
	if messageHeader == nil {
		return fmt.Errorf("Malformed response message: %q", decryptedMessage)
	}
	d.dialogID = messageHeader.DialogID.Val()

	bankInfoMessageBytes := decryptedMessage.FindSegment("HIKIM")
	fmt.Printf("INFO:\n%q\n", bankInfoMessageBytes)

	errors := make([]string, 0)
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("DialogEnd: Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func (d *PinTanDialog) End() error {
	dialogEnd := d.dialogEnd()
	signedDialogEnd, err := dialogEnd.Sign(d.signatureProvider)
	if err != nil {
		return err
	}
	encryptedDialogEnd, err := signedDialogEnd.Encrypt(d.cryptoProvider)
	if err != nil {
		return err
	}

	decryptedMessage, err := d.request(encryptedDialogEnd)
	if err != nil {
		return fmt.Errorf("Error while ending dialog: %v", err)
	}

	errors := make([]string, 0)
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("DialogEnd: Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func (d *PinTanDialog) SyncClientSystemID() (string, error) {
	syncMessage := new(message.SynchronisationMessage)
	messageNum := d.nextMessageNumber()
	syncMessage.BasicClientMessage = message.NewBasicClientMessage(syncMessage)
	syncMessage.Header = segment.NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	syncMessage.End = segment.NewMessageEndSegment(8, messageNum)
	syncMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.ClientID, initialClientSystemID, true)
	syncMessage.ProcessingPreparation = segment.NewProcessingPreparationSegment(0, 0, 1)
	syncMessage.Sync = segment.NewSynchronisationSegment(0)
	signedSyncMessage, err := syncMessage.Sign(d.signatureProvider)
	if err != nil {
		return "", err
	}
	d.cryptoProvider.SetClientSystemID(initialClientSystemID)
	encryptedSyncMessage, err := signedSyncMessage.Encrypt(d.cryptoProvider)
	if err != nil {
		return "", err
	}

	decryptedMessage, err := d.request(encryptedSyncMessage)
	if err != nil {
		return "", fmt.Errorf("Error while extracting encrypted message: %v", err)
	}

	messageHeader := decryptedMessage.MessageHeader()
	if messageHeader == nil {
		return "", fmt.Errorf("Malformed response message: %q", decryptedMessage)
	}
	d.dialogID = messageHeader.DialogID.Val()

	var errors []string
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsWarning() {
			fmt.Printf("%v\n", ack)
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return "", fmt.Errorf("Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}

	syncResponse := decryptedMessage.FindSegment("HISYN")
	if syncResponse != nil {
		syncSegment := &segment.SynchronisationResponseSegment{}
		err = syncSegment.UnmarshalHBCI(syncResponse)
		if err != nil {
			return "", fmt.Errorf("Error while unmarshaling sync response: %v", err)
		}
		d.SetClientSystemID(syncSegment.ClientSystemID.Val())
	} else {
		return "", fmt.Errorf("Malformed message: missing SynchronisationResponse")
	}

	bankParamData := decryptedMessage.FindSegment("HIBPA")
	if bankParamData != nil {
		paramSegment := &segment.CommonBankParameterSegment{}
		err = paramSegment.UnmarshalHBCI(bankParamData)
		if err != nil {
			return "", fmt.Errorf("Error while unmarshaling Bank Parameter Data: %v", err)
		}
		d.BankParameterData = paramSegment.BankParameterData()
	}

	userParamData := decryptedMessage.FindSegment("HIUPA")
	if userParamData != nil {
		paramSegment := &segment.CommonUserParameterDataSegment{}
		err = paramSegment.UnmarshalHBCI(userParamData)
		if err != nil {
			return "", fmt.Errorf("Error while unmarshaling User Parameter Data: %v", err)
		}
		d.UserParameterData = paramSegment.UserParameterData()
	}

	accountData := decryptedMessage.FindSegments("HIUPD")
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

	err = d.End()
	if err != nil {
		return "", err
	}

	return d.ClientSystemID, nil
}

func (d *PinTanDialog) request(message message.ClientMessage) (message.BankMessage, error) {
	marshaledMessage, err := message.MarshalHBCI()
	if err != nil {
		return nil, err
	}

	request := &transport.Request{
		URL:              d.hbciUrl,
		MarshaledMessage: marshaledMessage,
	}

	response, err := d.transport.Do(request)
	if err != nil {
		return nil, err
	}

	encMessage, err := extractEncryptedMessage(response)
	if err != nil {
		return nil, err
	}

	decryptedMessage, err := encMessage.Decrypt(d.cryptoProvider)
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting message: %v", err)
	}

	return decryptedMessage, err
}

func extractEncryptedMessage(response *transport.Response) (*message.EncryptedMessage, error) {
	messageHeader := response.FindSegment("HNHBK")
	if messageHeader == nil {
		return nil, fmt.Errorf("Malformed response: missing Message Header")
	}
	header := &segment.MessageHeaderSegment{}
	err := header.UnmarshalHBCI(messageHeader)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling message header: %v", err)
	}
	// TODO: parse messageEnd
	// TODO: parse encryptionHeader

	encMessage := message.NewEncryptedMessage(header, nil)

	encryptedData := response.FindSegment("HNVSD")
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

	dialogEnd := d.dialogEnd()
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
