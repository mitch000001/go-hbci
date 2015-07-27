package dialog

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
)

const initialDialogID = "0"
const initialClientSystemID = "0"
const anonymousClientID = "9999999999"

type Dialog interface {
	Init() (string, error)
	SyncClientSystemID() (string, error)
	End() error
}

func newDialog(bankId domain.BankId, hbciUrl string, clientId string, signatureProvider message.SignatureProvider, cryptoProvider message.CryptoProvider) *dialog {
	return &dialog{
		httpClient:        http.DefaultClient,
		hbciUrl:           hbciUrl,
		BankID:            bankId,
		ClientID:          clientId,
		ClientSystemID:    initialClientSystemID,
		Language:          domain.German,
		Accounts:          make([]domain.AccountInformation, 0),
		signatureProvider: signatureProvider,
		cryptoProvider:    cryptoProvider,
		dialogID:          initialDialogID,
	}
}

type dialog struct {
	transport         transport.Transport
	httpClient        *http.Client
	hbciUrl           string
	BankID            domain.BankId
	ClientID          string
	ClientSystemID    string
	Language          domain.Language
	UserParameterData domain.UserParameterData
	Accounts          []domain.AccountInformation
	messageCount      int
	dialogID          string
	signatureProvider message.SignatureProvider
	cryptoProvider    message.CryptoProvider
	BankParameterData domain.BankParameterData
}

func (d *dialog) UserParameterDataVersion() int {
	return d.UserParameterData.Version
}

func (d *dialog) BankParameterDataVersion() int {
	return d.BankParameterData.Version
}

func (d *dialog) SetClientSystemID(clientSystemID string) {
	d.ClientSystemID = clientSystemID
	d.signatureProvider.SetClientSystemID(d.ClientSystemID)
	d.cryptoProvider.SetClientSystemID(d.ClientSystemID)
}

func (d *dialog) SyncClientSystemID() (string, error) {
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

func (d *dialog) Init() error {
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

func (d *dialog) End() error {
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

func (d *dialog) request(message message.ClientMessage) (message.BankMessage, error) {
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

func (d *dialog) nextMessageNumber() int {
	d.messageCount += 1
	return d.messageCount
}

func (d *dialog) dialogEnd() *message.DialogFinishingMessage {
	dialogEnd := new(message.DialogFinishingMessage)
	messageNum := d.nextMessageNumber()
	dialogEnd.BasicClientMessage = message.NewBasicClientMessage(dialogEnd)
	dialogEnd.Header = segment.NewMessageHeaderSegment(0, 220, d.dialogID, messageNum)
	dialogEnd.End = segment.NewMessageEndSegment(8, messageNum)
	dialogEnd.DialogEnd = segment.NewDialogEndSegment(d.dialogID)
	return dialogEnd
}

func (d *dialog) post(message []byte) ([]byte, error) {
	encodedMessage := base64.StdEncoding.EncodeToString(message)
	response, err := d.httpClient.Post(d.hbciUrl, "application/vnd.hbci", strings.NewReader(encodedMessage))
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
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
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
