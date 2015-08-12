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
	Init() error
	SyncClientSystemID() (string, error)
	End() error
}

func newDialog(bankId domain.BankId, hbciUrl string, userId string, signatureProvider message.SignatureProvider, cryptoProvider message.CryptoProvider) *dialog {
	return &dialog{
		httpClient:        http.DefaultClient,
		hbciUrl:           hbciUrl,
		BankID:            bankId,
		UserID:            userId,
		clientID:          userId,
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
	UserID            string
	clientID          string
	ClientSystemID    string
	Language          domain.Language
	UserParameterData domain.UserParameterData
	Accounts          []domain.AccountInformation
	messageCount      int
	dialogID          string
	securityFn        string
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

func (d *dialog) SetSecurityFunction(securityFn string) {
	d.securityFn = securityFn
	d.signatureProvider.SetSecurityFunction(d.securityFn)
}

func (d *dialog) AccountInformation(allAccounts bool) error {
	err := d.Init()
	if err != nil {
		return err
	}
	defer func() {
		d.End()
	}()
	account := *d.Accounts[len(d.Accounts)-1].AccountConnection
	fmt.Printf("Account: %#v\n", account)
	accountInformationRequest := segment.NewAccountInformationRequestSegment(account, allAccounts)
	clientMessage := d.newBasicMessage(message.NewHBCIMessage(accountInformationRequest))
	signedMessage, err := clientMessage.Sign(d.signatureProvider)
	if err != nil {
		return err
	}
	encMessage, err := signedMessage.Encrypt(d.cryptoProvider)
	if err != nil {
		return err
	}
	decryptedMessage, err := d.request(encMessage)
	if err != nil {
		return err
	}
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
		return fmt.Errorf("Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	accountInfoResponse := decryptedMessage.FindSegment("HIKIF")
	if accountInfoResponse != nil {
		fmt.Printf("Account Info: %s\n", accountInfoResponse)
		return nil
	} else {
		return fmt.Errorf("Malformed response: expected HIKIF segment")
	}
	return nil
}

func (d *dialog) Balances(allAccounts bool) ([]domain.AccountBalance, error) {
	err := d.Init()
	if err != nil {
		return nil, err
	}
	defer func() {
		d.End()
	}()
	account := *d.Accounts[len(d.Accounts)-1].AccountConnection
	fmt.Printf("Account: %#v\n", account)
	accountBalanceRequest := segment.NewAccountBalanceRequestSegment(account, allAccounts)
	clientMessage := d.newBasicMessage(message.NewHBCIMessage(accountBalanceRequest))
	signedMessage, err := clientMessage.Sign(d.signatureProvider)
	if err != nil {
		return nil, err
	}
	encMessage, err := signedMessage.Encrypt(d.cryptoProvider)
	if err != nil {
		return nil, err
	}
	decryptedMessage, err := d.request(encMessage)
	if err != nil {
		return nil, err
	}
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
		return nil, fmt.Errorf("Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	var balances []domain.AccountBalance
	balanceResponses := decryptedMessage.FindSegments("HISAL")
	if balanceResponses != nil {
		for _, marshaledSegment := range balanceResponses {
			balanceSegment := &segment.AccountBalanceResponseSegment{}
			err = balanceSegment.UnmarshalHBCI(marshaledSegment)
			if err != nil {
				return nil, fmt.Errorf("Error while parsing account balance: %v", err)
			}
			balances = append(balances, balanceSegment.AccountBalance())
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("Malformed response: expected HISAL segment")
	}

	return balances, nil
}

func (d *dialog) SyncClientSystemID() (string, error) {
	syncMessage := &message.SynchronisationMessage{
		Identification:        segment.NewIdentificationSegment(d.BankID, d.clientID, initialClientSystemID, true),
		ProcessingPreparation: segment.NewProcessingPreparationSegment(0, 0, 1),
		Sync: segment.NewSynchronisationSegment(0),
	}
	syncMessage.BasicMessage = d.newBasicMessage(syncMessage)
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

	err = d.parseBankParameterData(decryptedMessage)
	if err != nil {
		return "", err
	}

	err = d.parseUserParameterData(decryptedMessage)
	if err != nil {
		return "", err
	}

	err = d.End()
	if err != nil {
		return "", err
	}

	return d.ClientSystemID, nil
}

func (d *dialog) Init() error {
	err := d.init()
	if err != nil {
		return err
	}
	d.dialogID = initialDialogID
	d.messageCount = 0
	initMessage := &message.DialogInitializationClientMessage{
		Identification:        segment.NewIdentificationSegment(d.BankID, d.clientID, d.ClientSystemID, true),
		ProcessingPreparation: segment.NewProcessingPreparationSegment(d.BankParameterDataVersion(), d.UserParameterDataVersion(), d.Language),
	}
	initMessage.BasicMessage = d.newBasicMessage(initMessage)
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

	err = d.parseBankParameterData(decryptedMessage)
	if err != nil {
		return err
	}

	err = d.parseUserParameterData(decryptedMessage)
	if err != nil {
		return err
	}

	bankInfoMessageBytes := decryptedMessage.FindSegment("HIKIM")
	fmt.Printf("INFO:\n%q\n", bankInfoMessageBytes)

	newSecurityFn := d.securityFn
	errors := make([]string, 0)
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.Code == 3920 {
			supportedSecurityFns := ack.Params
			if len(supportedSecurityFns) != 0 {
				fmt.Printf("Supported securityFunctions: %q\n", supportedSecurityFns)
				// TODO: proper handling of each case, see FINTS3.0 docu
				newSecurityFn = supportedSecurityFns[0]
			}
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("DialogEnd: Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	if d.securityFn != newSecurityFn {
		err = d.End()
		if err != nil {
			return err
		}
		d.SetSecurityFunction(newSecurityFn)
		err = d.Init()
		if err != nil {
			return err
		}
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

func (d *dialog) init() error {
	if d.ClientSystemID == initialClientSystemID {
		id, err := d.SyncClientSystemID()
		if err != nil {
			return err
		}
		d.ClientSystemID = id
	}
	return nil
}

func (d *dialog) newClientMessage(hbciMessage message.HBCIMessage) message.ClientMessage {
	return d.newBasicMessage(hbciMessage)
}

func (d *dialog) newBasicMessage(hbciMessage message.HBCIMessage) *message.BasicMessage {
	messageNum := d.nextMessageNumber()
	clientMessage := message.NewBasicMessage(hbciMessage)
	clientMessage.Header = segment.NewMessageHeaderSegment(-1, 220, d.dialogID, messageNum)
	clientMessage.End = segment.NewMessageEndSegment(-1, messageNum)
	return clientMessage
}

func (d *dialog) parseBankParameterData(bankMessage message.BankMessage) error {
	bankParamData := bankMessage.FindSegment("HIBPA")
	if bankParamData != nil {
		paramSegment := &segment.CommonBankParameterSegment{}
		err := paramSegment.UnmarshalHBCI(bankParamData)
		if err != nil {
			return fmt.Errorf("Error while unmarshaling Bank Parameter Data: %v", err)
		}
		d.BankParameterData = paramSegment.BankParameterData()
	}
	pinTanTransactions := bankMessage.FindSegment("DIPINS")
	if pinTanTransactions != nil {
		pinTanTransactionSegment := &segment.PinTanBusinessTransactionParamsSegment{}
		err := pinTanTransactionSegment.UnmarshalHBCI(pinTanTransactions)
		if err != nil {
			return fmt.Errorf("Error while unmarshaling PinTan Segment Parameter Data: %v", err)
		}
		pinTransactions := make(map[string]bool)
		for _, transaction := range pinTanTransactionSegment.PinTanBusinessTransactions() {
			pinTransactions[transaction.SegmentID] = transaction.NeedsTan
		}
		d.BankParameterData.PinTanBusinessTransactions = pinTransactions
	}
	return nil
}

func (d *dialog) parseUserParameterData(bankMessage message.BankMessage) error {
	userParamData := bankMessage.FindSegment("HIUPA")
	if userParamData != nil {
		paramSegment := &segment.CommonUserParameterDataSegment{}
		err := paramSegment.UnmarshalHBCI(userParamData)
		if err != nil {
			return fmt.Errorf("Error while unmarshaling User Parameter Data: %v", err)
		}
		d.UserParameterData = paramSegment.UserParameterData()
		d.clientID = d.UserParameterData.UserID
	}

	accountData := bankMessage.FindSegments("HIUPD")
	if accountData != nil {
		for _, acc := range accountData {
			infoSegment := &segment.AccountInformationSegment{}
			err := infoSegment.UnmarshalHBCI(acc)
			if err != nil {
				return fmt.Errorf("Error while unmarshaling Accounts: %v", err)
			}
			d.Accounts = append(d.Accounts, infoSegment.Account())
		}
	}

	return nil
}

func (d *dialog) request(clientMessage message.ClientMessage) (message.BankMessage, error) {
	marshaledMessage, err := clientMessage.MarshalHBCI()
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

	var bankMessage message.BankMessage
	if response.IsEncrypted() {
		encMessage, err := extractEncryptedMessage(response)
		if err != nil {
			return nil, err
		}

		decryptedMessage, err := encMessage.Decrypt(d.cryptoProvider)
		if err != nil {
			return nil, fmt.Errorf("Error while decrypting message: %v", err)
		}
		bankMessage = decryptedMessage
	} else {
		decryptedMessage, err := extractUnencryptedMessage(response)
		if err != nil {
			return nil, err
		}
		bankMessage = decryptedMessage
	}

	return bankMessage, err
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
		return nil, fmt.Errorf("Malformed response: missing encrypted data: \n%q", response)
	}
	return encMessage, nil
}

func extractUnencryptedMessage(response *transport.Response) (*message.DecryptedMessage, error) {
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
	decryptedMessage, err := message.NewDecryptedMessage(header, nil, response.MarshaledResponse)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling unencrypted data: %v", err)
	}
	return decryptedMessage, nil
}

func (d *dialog) nextMessageNumber() int {
	d.messageCount += 1
	return d.messageCount
}

func (d *dialog) dialogEnd() *message.DialogFinishingMessage {
	dialogEnd := &message.DialogFinishingMessage{
		DialogEnd: segment.NewDialogEndSegment(d.dialogID),
	}
	dialogEnd.BasicMessage = d.newBasicMessage(dialogEnd)
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
