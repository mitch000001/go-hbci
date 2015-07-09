package hbci

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
)

const initialDialogID = "0"
const initialClientSystemID = "0"
const anonymousClientID = "9999999999"

type Dialog interface {
	Init() (string, error)
}

func newDialog(bankId domain.BankId, hbciUrl string, clientId string, signatureProvider message.SignatureProvider, encryptionProvider message.EncryptionProvider) *dialog {
	return &dialog{
		hbciUrl:            hbciUrl,
		BankID:             bankId,
		ClientID:           clientId,
		ClientSystemID:     initialClientSystemID,
		signatureProvider:  signatureProvider,
		encryptionProvider: encryptionProvider,
	}
}

type dialog struct {
	hbciUrl            string
	BankID             domain.BankId
	ClientID           string
	ClientSystemID     string
	messageCount       int
	signatureProvider  message.SignatureProvider
	encryptionProvider message.EncryptionProvider
}

func (d *dialog) nextMessageNumber() int {
	d.messageCount += 1
	return d.messageCount
}

func (d *dialog) dialogEnd(dialogId string) *message.DialogFinishingMessage {
	dialogEnd := new(message.DialogFinishingMessage)
	messageNum := d.nextMessageNumber()
	dialogEnd.BasicClientMessage = message.NewBasicClientMessage(dialogEnd)
	dialogEnd.Header = segment.NewMessageHeaderSegment(0, 220, dialogId, messageNum)
	dialogEnd.End = segment.NewMessageEndSegment(8, messageNum)
	dialogEnd.DialogEnd = segment.NewDialogEndSegment(dialogId)
	return dialogEnd
}

func (d *dialog) post(message []byte) ([]byte, error) {
	encodedMessage := base64.StdEncoding.EncodeToString(message)
	response, err := http.Post(d.hbciUrl, "application/vnd.hbci", strings.NewReader(encodedMessage))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		decodedReader := base64.NewDecoder(base64.StdEncoding, response.Body)
		bodyBytes, err := ioutil.ReadAll(decodedReader)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	} else {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	}
}

func (d *dialog) dial(message []byte) ([]byte, error) {
	conn, err := net.Dial("tcp4", d.hbciUrl)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	fmt.Fprintf(conn, "%s\r\n\r\n", string(message))
	buf := bufio.NewReader(conn)
	// read answer header
	header, err := buf.ReadString('\'')
	if err != nil {
		return nil, err
	}
	headerItems := strings.Split(header, "+")
	if len(headerItems) < 2 {
		return nil, fmt.Errorf("Response header too short")
	}
	sizeString := headerItems[1]
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing message size: %T:%v\n")
	}
	messageBuf := make([]byte, size)
	buf.Read(messageBuf)
	var retBuf bytes.Buffer
	retBuf.WriteString(header)
	retBuf.Write(messageBuf)
	return retBuf.Bytes(), err
}

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
	if err != nil && err != io.EOF {
		return "", err
	}
	fmt.Printf("Response: \n%s\n", bytes.Join(bytes.Split(response, []byte("'")), []byte("'\n")))

	extractor := NewSegmentExtractor(response)
	_, err = extractor.Extract()
	if err != nil {
		return "", err
	}
	messageHeader := extractor.FindSegment("HNHBK")
	if messageHeader == nil {
		return "", fmt.Errorf("Malformed response: %q", response)
	}
	dataElements := bytes.Split(messageHeader, []byte("+"))
	newDialogId := string(dataElements[3])

	syncResponse := extractor.FindSegment("HISYN")
	if syncResponse != nil {
		dataElements := bytes.Split(syncResponse, []byte("+"))
		newClientSystemId := dataElements[1]
		d.ClientSystemID = string(newClientSystemId)
		d.signatureProvider.SetClientSystemID(d.ClientSystemID)
		d.encryptionProvider.SetClientSystemID(d.ClientSystemID)
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
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(response), nil
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

func NewRDHDialog(bankId domain.BankId, hbciUrl string, clientId string) *rdhDialog {
	key, err := domain.GenerateSigningKey()
	if err != nil {
		panic(err)
	}
	signingKey := domain.NewRSAKey(key, domain.NewInitialKeyName(bankId.CountryCode, bankId.ID, clientId, "S"))
	provider := message.NewRDHSignatureProvider(signingKey)
	d := &rdhDialog{
		dialog:      newDialog(bankId, hbciUrl, clientId, provider, nil),
		SigningKey:  signingKey,
		SignatureID: 12345,
	}
	return d
}

type rdhDialog struct {
	*dialog
	SignatureID int
	SigningKey  domain.Key
}
