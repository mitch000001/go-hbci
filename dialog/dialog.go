package dialog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/internal"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
)

// Dialog represents the common interface to use when talking to bank institutes
type Dialog interface {
	SyncClientSystemID() (string, error)
	SendMessage(message.HBCIMessage) (message.BankMessage, error)
}

const (
	initialDialogID                 = "0"
	initialClientSystemID           = "0"
	initialBankParameterDataVersion = 0
	initialUserParameterDataVersion = 0
	anonymousClientID               = "9999999999"
)

func newDialog(
	bankID domain.BankID,
	hbciURL string,
	userID string,
	hbciVersion segment.HBCIVersion,
	productName string,
	productVersion string,
	signatureProvider message.SignatureProvider,
	cryptoProvider message.CryptoProvider,
) *dialog {
	return &dialog{
		hbciURL:           hbciURL,
		BankID:            bankID,
		UserID:            userID,
		clientID:          userID,
		ClientSystemID:    initialClientSystemID,
		Language:          domain.German,
		Accounts:          make([]domain.AccountInformation, 0),
		signatureProvider: signatureProvider,
		cryptoProvider:    cryptoProvider,
		dialogID:          initialDialogID,
		hbciVersion:       hbciVersion,
		productName:       productName,
		productVersion:    productVersion,
	}
}

type dialog struct {
	transport         transport.Transport
	hbciURL           string
	BankID            domain.BankID
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
	BankParameterData BankParameterData
	hbciVersion       segment.HBCIVersion
	productName       string
	productVersion    string
	supportedSegments []segment.VersionedSegment
}

func (d *dialog) UserParameterDataVersion() int {
	return d.UserParameterData.Version
}

func (d *dialog) BankParameterDataVersion() int {
	return d.BankParameterData.Version
}

func (d *dialog) SupportedSegments() []segment.VersionedSegment {
	return d.supportedSegments
}

func (d *dialog) SetClientSystemID(clientSystemID string) {
	d.ClientSystemID = clientSystemID
	d.signatureProvider.SetClientSystemID(d.ClientSystemID)
	d.cryptoProvider.SetClientSystemID(d.ClientSystemID)
}

func (d *dialog) SetSecurityFunction(securityFn string) {
	d.securityFn = securityFn
	d.signatureProvider.SetSecurityFunction(d.securityFn)
	d.cryptoProvider.SetSecurityFunction(d.securityFn)
}

func (d *dialog) SendMessage(clientMessage message.HBCIMessage) (message.BankMessage, error) {
	err := d.init()
	if err != nil {
		return nil, err
	}
	defer func() { logErr(d.end()) }()
	requestMessage := d.newBasicMessage(clientMessage)
	signedMessage, err := requestMessage.Sign(d.signatureProvider)
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
			internal.Info.Printf("%v\n", ack)
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return nil, fmt.Errorf("institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	return decryptedMessage, nil
}

func (d *dialog) SyncUserParameterData() error {
	internal.Info.Printf("Initializing dialog")
	err := d.init()
	if err != nil {
		return err
	}
	defer func() {
		internal.Info.Printf("Ending dialog")
		logErr(d.end())
	}()
	return nil
}

func (d *dialog) SyncClientSystemID() (string, error) {
	syncMessage := message.NewSynchronisationMessage(d.hbciVersion)
	syncMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.clientID, initialClientSystemID, true)
	syncMessage.ProcessingPreparation = segment.NewProcessingPreparationSegmentV3(
		initialBankParameterDataVersion, initialUserParameterDataVersion, domain.German, d.productName, d.productVersion,
	)
	syncMessage.TanRequest = segment.NewTanProcess4RequestSegmentV6(segment.IdentificationID)
	syncMessage.Sync = d.hbciVersion.SynchronisationRequest(segment.SyncModeAquireClientID)
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
		return "", fmt.Errorf("error while extracting encrypted message: %v", err)
	}

	messageHeader := decryptedMessage.MessageHeader()
	if messageHeader == nil {
		return "", fmt.Errorf("malformed response message: %q", decryptedMessage)
	}
	d.dialogID = messageHeader.DialogID.Val()
	var errors []string
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsSuccess() {
			internal.Debug.Printf("%v\n", ack)
		}
		if ack.IsWarning() {
			internal.Info.Printf("%v\n", ack)
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return "", fmt.Errorf("institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	if err := d.updateSecurityFunctionIfNeeded(decryptedMessage); err != nil {
		return "", fmt.Errorf("error updating security function: %w", err)
	}

	syncResponse := decryptedMessage.FindSegment("HISYN")
	if syncResponse != nil {
		syncSegment := syncResponse.(segment.SynchronisationResponse)
		d.SetClientSystemID(syncSegment.ClientSystemID())
	} else {
		return "", fmt.Errorf("malformed message: missing unmarshaler for SynchronisationResponse")
	}

	err = d.parseBankParameterData(decryptedMessage)
	if err != nil {
		return "", err
	}

	err = d.parseUserParameterData(decryptedMessage)
	if err != nil {
		return "", err
	}

	err = d.end()
	if err != nil {
		return "", err
	}

	return d.ClientSystemID, nil
}

func (d *dialog) GetAnonymousBankParameterData() (*BankParameterData, error) {
	err := d.anonymousInit()
	if err != nil {
		return nil, fmt.Errorf("Error while initating anonymous dialog: %v", err)
	}
	defer func() { logErr(d.anonymousEnd()) }()
	return &d.BankParameterData, nil
}

func (d *dialog) SendAnonymousMessage(clientMessage message.HBCIMessage) (message.BankMessage, error) {
	err := d.anonymousInit()
	if err != nil {
		return nil, fmt.Errorf("Error while initating anonymous dialog: %v", err)
	}
	defer func() { logErr(d.anonymousEnd()) }()
	// TODO: add checks if job needs signature or not
	requestMessage := d.newBasicMessage(clientMessage)
	requestMessage.SetSegmentPositions()
	bankMessage, err := d.request(requestMessage)
	if err != nil {
		return nil, err
	}
	var errors []string
	acknowledgements := bankMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsWarning() {
			fmt.Printf("%v\n", ack)
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return nil, fmt.Errorf("institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	return bankMessage, nil
}

func (d *dialog) anonymousInit() error {
	d.dialogID = initialDialogID
	d.messageCount = 0
	initMessage := message.NewDialogInitializationClientMessage(d.hbciVersion)
	initMessage.Identification = segment.NewIdentificationSegment(d.BankID, anonymousClientID, initialClientSystemID, false)
	initMessage.ProcessingPreparation = segment.NewProcessingPreparationSegmentV3(
		d.BankParameterDataVersion(), d.UserParameterDataVersion(), d.Language, d.productName, d.productVersion,
	)
	initMessage.BasicMessage = d.newBasicMessage(initMessage)
	initMessage.SetSegmentPositions()
	bankMessage, err := d.request(initMessage)
	if err != nil {
		return err
	}
	messageHeader := bankMessage.MessageHeader()
	if messageHeader == nil {
		return fmt.Errorf("malformed response message: %q", bankMessage)
	}
	d.dialogID = messageHeader.DialogID.Val()

	err = d.parseBankParameterData(bankMessage)
	if err != nil {
		return err
	}

	err = d.parseUserParameterData(bankMessage)
	if err != nil {
		return err
	}

	bankInfoMessage := bankMessage.FindSegment(segment.BankAnnouncementID)
	if bankInfoMessage != nil {
		bankInfoSegment := bankInfoMessage.(*segment.BankAnnouncementSegment)
		internal.Info.Printf("INFO:\n%s\n%s\n", bankInfoSegment.Subject.Val(), bankInfoSegment.Body.Val())
	}

	errors := make([]string, 0)
	acknowledgements := bankMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsWarning() {
			fmt.Printf("%v\n", ack)
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("DialogEnd: Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

func (d *dialog) anonymousEnd() error {
	dialogEnd := message.NewDialogFinishingMessage(d.hbciVersion, d.dialogID)
	dialogEnd.BasicMessage = d.newBasicMessage(dialogEnd)
	dialogEnd.SetSegmentPositions()

	decryptedMessage, err := d.request(dialogEnd)
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
	d.dialogID = initialDialogID
	d.messageCount = 0
	initMessage := message.NewDialogInitializationClientMessage(d.hbciVersion)
	initMessage.Identification = segment.NewIdentificationSegment(d.BankID, d.clientID, d.ClientSystemID, true)
	initMessage.ProcessingPreparation = segment.NewProcessingPreparationSegmentV3(
		d.BankParameterDataVersion(), d.UserParameterDataVersion(), d.Language, d.productName, d.productVersion,
	)
	builder := segment.NewBuilder(d.supportedSegments)
	tanRequest, err := builder.TanProcessV4Request(segment.IdentificationID)
	if err != nil {
		return fmt.Errorf("error building TAN V4 Process segment: %w", err)
	}
	initMessage.TanRequest = tanRequest
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
		return fmt.Errorf("error while initializing dialog: %v", err)
	}
	messageHeader := decryptedMessage.MessageHeader()
	if messageHeader == nil {
		return fmt.Errorf("malformed response message: %q", decryptedMessage)
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

	bankInfoMessage := decryptedMessage.FindSegment(segment.BankAnnouncementID)
	if bankInfoMessage != nil {
		bankInfoSegment := bankInfoMessage.(*segment.BankAnnouncementSegment)
		internal.Info.Printf("INFO:\n%s\n%s\n", bankInfoSegment.Subject.Val(), bankInfoSegment.Body.Val())
	}

	var errors []string
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.IsSuccess() {
			internal.Debug.Printf("%v\n", ack)
		}
		if ack.IsWarning() {
			internal.Info.Printf("%v\n", ack)
		}
		if ack.IsError() {
			errors = append(errors, ack.String())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("DialogEnd: Institute returned errors:\n%s", strings.Join(errors, "\n"))
	}

	if err := d.updateSecurityFunctionIfNeeded(decryptedMessage); err != nil {
		return fmt.Errorf("error updating security function: %w", err)
	}

	return nil
}

func (d *dialog) end() error {
	dialogEnd := message.NewDialogFinishingMessage(d.hbciVersion, d.dialogID)
	dialogEnd.BasicMessage = d.newBasicMessage(dialogEnd)
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

func (d *dialog) newClientMessage(hbciMessage message.HBCIMessage) message.ClientMessage {
	return d.newBasicMessage(hbciMessage)
}

func (d *dialog) newBasicMessage(hbciMessage message.HBCIMessage) *message.BasicMessage {
	messageNum := d.nextMessageNumber()
	clientMessage := message.NewBasicMessage(hbciMessage)
	clientMessage.Header = segment.NewMessageHeaderSegment(-1, d.hbciVersion.Version(), d.dialogID, messageNum)
	clientMessage.End = segment.NewMessageEndSegment(-1, messageNum)
	return clientMessage
}

func (d *dialog) updateSecurityFunctionIfNeeded(message message.BankMessage) error {
	if !d.hasSupportedSecurityAcknowledgement(message) {
		return nil
	}
	supportedSecurityFns, ok := d.supportedSecurityFunctionsFromBankMessage(message)
	if !ok {
		return fmt.Errorf("no supported security function implemented")
	}
	oldSecurityFn := d.securityFn
	newSecurityFn := ""
	newSecurityFnName := ""
	// TODO: proper handling of each case, see FINTS3.0 docu -> Ask user which to use
	for fn, name := range supportedSecurityFns {
		newSecurityFn = fn
		newSecurityFnName = name
		break
	}
	if oldSecurityFn != newSecurityFn {
		internal.Info.Printf(
			"New supported security function found. Setting new security function %q (%s).", newSecurityFnName, newSecurityFn,
		)
		d.SetSecurityFunction(newSecurityFn)

	}
	return nil
}

func (d *dialog) hasSupportedSecurityAcknowledgement(message message.BankMessage) bool {
	rawTanBPDSegment := message.FindSegment(segment.TanBankParameterID)
	if rawTanBPDSegment == nil {
		return false
	}
	for _, ack := range message.Acknowledgements() {
		if ack.Code == element.AcknowledgementSupportedSecurityFunction {
			return true
		}
	}
	return false
}

func (d *dialog) supportedSecurityFunctionsFromBankMessage(message message.BankMessage) (map[string]string, bool) {
	acknowledgements := message.Acknowledgements()
	var supportedSecurityFns []string
	for _, ack := range acknowledgements {
		if ack.Code == element.AcknowledgementSupportedSecurityFunction {
			supportedSecurityFns = ack.Params
			break
		}
	}
	if len(supportedSecurityFns) == 0 {
		return nil, false
	}
	securityFunctions := map[string]string{}
	availableSecurityFns := extractAvailableSecurityFnsFromMessage(message)
	for _, sf := range supportedSecurityFns {
		if pp, ok := availableSecurityFns[sf]; ok {
			securityFunctions[sf] = pp.TwoStepProcessName
		}
	}
	if len(securityFunctions) == 0 {
		return nil, false
	}
	return securityFunctions, true
}

func extractAvailableSecurityFnsFromMessage(message message.BankMessage) map[string]domain.TanProcessParameter {
	rawSegments := message.FindSegments(segment.TanBankParameterID)
	availableSecurityFns := map[string]domain.TanProcessParameter{}
	for _, rawSegment := range rawSegments {
		if rawSegment == nil {
			return nil
		}
		tanBankParams, ok := rawSegment.(segment.TanBankParameter)
		if !ok {
			return nil
		}
		for _, pp := range tanBankParams.TanProcessParameters() {
			availableSecurityFns[pp.SecurityFunction] = pp
		}
	}
	return availableSecurityFns
}

func (d *dialog) parseBankParameterData(bankMessage message.BankMessage) error {
	bankParamData := bankMessage.FindSegment(segment.CommonBankParameterID)
	if bankParamData == nil {
		return nil
	}
	paramSegment, ok := bankParamData.(segment.CommonBankParameter)
	if !ok {
		return fmt.Errorf("error converting common bank parameter data")
	}
	d.supportedSegments = bankMessage.SupportedSegments()
	d.BankParameterData = BankParameterData{
		BankParameterData:          paramSegment.BankParameterData(),
		SupportedSegmentParameters: make([]SegmentParameter, len(d.supportedSegments)),
	}
	pinTanTransactions := bankMessage.FindSegment(segment.PinTanBankParameterID)
	if pinTanTransactions != nil {
		pinTanTransactionSegment := pinTanTransactions.(segment.PinTanBusinessTransactionParams)
		pinTransactions := make(map[string]bool)
		for _, transaction := range pinTanTransactionSegment.PinTanBusinessTransactions() {
			pinTransactions[transaction.SegmentID] = transaction.NeedsTan
		}
		d.BankParameterData.PinTanBusinessTransactions = pinTransactions
	}
	for i, s := range d.supportedSegments {
		param := SegmentParameter{
			VersionedSegment: s,
		}
		parameterData := bankMessage.FindSegment(s.ID)
		if parameterData != nil && parameterData.Header().Version.Val() == s.Version {
			param.Parameters = parameterData
		}
		d.BankParameterData.SupportedSegmentParameters[i] = param
	}
	return nil
}

func (d *dialog) parseUserParameterData(bankMessage message.BankMessage) error {
	userParamData := bankMessage.FindSegment(segment.CommonUserParameterDataID)
	if userParamData != nil {
		paramSegment := userParamData.(segment.CommonUserParameterData)
		d.UserParameterData = paramSegment.UserParameterData()
		d.clientID = d.UserParameterData.UserID
	}

	accountData := bankMessage.FindSegments(segment.AccountInformationID)
	for _, acc := range accountData {
		infoSegment := acc.(segment.AccountInformation)
		d.Accounts = append(d.Accounts, infoSegment.Account())
	}
	return nil
}

func (d *dialog) request(clientMessage message.ClientMessage) (message.BankMessage, error) {
	marshaledMessage, err := clientMessage.MarshalHBCI()
	if err != nil {
		return nil, err
	}
	internal.Debug.Printf("Request:\n")
	for _, seg := range bytes.Split(marshaledMessage, []byte("'")) {
		internal.Debug.Printf("%q\n", seg)
	}

	reqBody := bytes.NewReader(marshaledMessage)

	request := &transport.Request{
		URL:  d.hbciURL,
		Body: io.NopCloser(reqBody),
	}

	response, err := d.transport.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error executing Transport request: %v", err)
	}
	response, err = transport.ReadResponse(bufio.NewReader(response.Body), response.Request)
	if err != nil {
		return nil, fmt.Errorf("error reading response from Transport: %v", err)
	}

	var bankMessage message.BankMessage
	if response.IsEncrypted() {
		encMessage, err := d.extractEncryptedMessage(response)
		if err != nil {
			return nil, err
		}

		decryptedMessage, err := encMessage.Decrypt(d.cryptoProvider)
		if err != nil {
			return nil, fmt.Errorf("error while decrypting message: %v", err)
		}
		internal.Debug.Printf("Response:\n %s\n", decryptedMessage.MessageHeader())
		bankMessage = decryptedMessage
	} else {
		decryptedMessage, err := extractUnencryptedMessage(response)
		if err != nil {
			return nil, err
		}
		internal.Debug.Printf("Response:\n %s\n", decryptedMessage.MessageHeader())
		bankMessage = decryptedMessage
	}

	return bankMessage, err
}

func (d *dialog) extractEncryptedMessage(response *transport.Response) (*message.EncryptedMessage, error) {
	messageHeader := response.FindSegment(segment.MessageHeaderID)
	if messageHeader == nil {
		return nil, fmt.Errorf("malformed response: missing Message Header")
	}
	header := &segment.MessageHeaderSegment{}
	err := header.UnmarshalHBCI(messageHeader)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling message header: %v", err)
	}
	// TODO: parse messageEnd
	// TODO: parse encryptionHeader

	encMessage := message.NewEncryptedMessage(header, nil, d.hbciVersion)

	encryptedData := response.FindSegment("HNVSD")
	if encryptedData != nil {
		encSegment := &segment.EncryptedDataSegment{}
		err = encSegment.UnmarshalHBCI(encryptedData)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshaling encrypted data: %v", err)
		}
		encMessage.EncryptedData = encSegment
	} else {
		return nil, fmt.Errorf("malformed response: missing encrypted data: \n%v", response)
	}
	return encMessage, nil
}

func extractUnencryptedMessage(response *transport.Response) (message.BankMessage, error) {
	messageHeader := response.FindSegment(segment.MessageHeaderID)
	if messageHeader == nil {
		return nil, fmt.Errorf("malformed response: missing Message Header")
	}
	header := &segment.MessageHeaderSegment{}
	err := header.UnmarshalHBCI(messageHeader)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling message header: %v", err)
	}
	// TODO: parse messageEnd
	decryptedMessage, err := message.NewDecryptedMessage(header, nil, response.MarshaledResponse)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling unencrypted data: %v", err)
	}
	return decryptedMessage, nil
}

func (d *dialog) nextMessageNumber() int {
	d.messageCount++
	return d.messageCount
}

func (d *dialog) dial(message []byte) ([]byte, error) {
	conn, err := net.Dial("tcp4", d.hbciURL)
	if err != nil {
		return nil, err
	}
	defer func() { logErr(conn.Close()) }()
	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(conn, "%s\r\n\r\n", string(message))
	buf := bufio.NewReader(conn)
	// read answer header
	header, err := buf.ReadString('\'')
	if err != nil {
		return nil, err
	}
	headerItems := strings.Split(header, "+")
	if len(headerItems) < 2 {
		return nil, fmt.Errorf("response header too short")
	}
	sizeString := headerItems[1]
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		return nil, fmt.Errorf("error while parsing message size: %T:%v", err, err)
	}
	messageBuf := make([]byte, size)
	_, err = buf.Read(messageBuf)
	if err != nil {
		return nil, fmt.Errorf("error while reading message: %T:%v", err, err)
	}
	var retBuf bytes.Buffer
	retBuf.WriteString(header)
	retBuf.Write(messageBuf)
	return retBuf.Bytes(), err
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
