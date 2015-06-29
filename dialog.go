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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/dataelement"
)

const initialDialogID = "0"
const initialClientSystemID = "0"
const anonymousClientID = "9999999999"

type Dialog interface {
	Init() (string, error)
}

func newDialog(bankId dataelement.BankId, hbciUrl string, clientId string, signatureProvider SignatureProvider, encryptionProvider EncryptionProvider) *dialog {
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
	BankID             dataelement.BankId
	ClientID           string
	ClientSystemID     string
	messageCount       int
	signatureProvider  SignatureProvider
	encryptionProvider EncryptionProvider
}

func (d *dialog) nextMessageNumber() int {
	d.messageCount += 1
	return d.messageCount
}

func (d *dialog) dialogEnd(dialogId string) *DialogFinishingMessage {
	dialogEnd := new(DialogFinishingMessage)
	messageNum := d.nextMessageNumber()
	dialogEnd.basicClientMessage = newBasicClientMessage(dialogEnd)
	dialogEnd.Header = NewMessageHeaderSegment(0, 220, dialogId, messageNum)
	dialogEnd.End = NewMessageEndSegment(8, messageNum)
	dialogEnd.DialogEnd = NewDialogEndSegment(dialogId)
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

func NewPinTanDialog(bankId dataelement.BankId, hbciUrl string, clientId string) *pinTanDialog {
	signatureProvider := NewPinTanSignatureProvider(nil, "")
	encryptionProvider := NewPinTanEncryptionProvider(nil, "")
	d := &pinTanDialog{
		dialog: newDialog(bankId, hbciUrl, clientId, signatureProvider, encryptionProvider),
	}
	return d
}

type pinTanDialog struct {
	*dialog
	pin           string
	signingKey    Key
	pinTanKeyName dataelement.KeyName
}

func (d *pinTanDialog) SetPin(pin string) {
	d.pin = pin
	pinKey := NewPinKey(pin, dataelement.NewPinTanKeyName(d.BankID, d.ClientID, "S"))
	d.signingKey = pinKey
	d.signatureProvider = NewPinTanSignatureProvider(pinKey, d.ClientSystemID)
	pinKey = NewPinKey(pin, dataelement.NewPinTanKeyName(d.BankID, d.ClientID, "V"))
	d.encryptionProvider = NewPinTanEncryptionProvider(pinKey, d.ClientSystemID)
}

func (d *pinTanDialog) Init() (string, error) {
	initMessage := NewDialogInitializationClientMessage()
	messageNum := d.nextMessageNumber()
	initMessage.Header = NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	initMessage.End = NewMessageEndSegment(8, messageNum)
	initMessage.Identification = NewIdentificationSegment(d.BankID, d.ClientID, initialClientSystemID, true)
	initMessage.ProcessingPreparation = NewProcessingPreparationSegment(0, 0, 1)
	controlRef := "1"
	initMessage.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	initMessage.SignatureEnd = NewSignatureEndSegment(7, controlRef)
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
	dialogEnd.SignatureEnd = NewSignatureEndSegment(7, controlRef)
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
	syncMessage := new(SynchronisationMessage)
	messageNum := d.nextMessageNumber()
	syncMessage.basicClientMessage = newBasicClientMessage(syncMessage)
	syncMessage.Header = NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	syncMessage.End = NewMessageEndSegment(8, messageNum)
	syncMessage.Identification = NewIdentificationSegment(d.BankID, d.ClientID, initialClientSystemID, true)
	syncMessage.ProcessingPreparation = NewProcessingPreparationSegment(0, 0, 1)
	syncMessage.Sync = NewSynchronisationSegment(0)
	controlRef := "1"
	syncMessage.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	syncMessage.SignatureEnd = NewSignatureEndSegment(7, controlRef)
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

	result := bytes.Split(response, []byte("'"))
	i := sort.Search(len(result), func(i int) bool {
		return bytes.HasPrefix(result[i], []byte("HISYN"))
	})
	if i < len(result) {
		syncResponse := result[i]
		dataElements := bytes.Split(syncResponse, []byte("+"))
		newClientSystemId := dataElements[1]
		fmt.Printf("SyncResponse: %s\n", syncResponse)
		d.ClientSystemID = string(newClientSystemId)
		d.signatureProvider.SetClientSystemID(d.ClientSystemID)
		d.encryptionProvider.SetClientSystemID(d.ClientSystemID)
	}
	messageHeader := result[0]
	dataElements := bytes.Split(messageHeader, []byte("+"))
	newDialogId := string(dataElements[3])
	fmt.Printf("New dialogID: %s\n", newDialogId)

	dialogEnd := d.dialogEnd(newDialogId)
	dialogEnd.SignatureBegin = d.signatureProvider.NewSignatureHeader(controlRef, 0)
	dialogEnd.SignatureEnd = NewSignatureEndSegment(7, controlRef)
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
	comm := NewCommunicationAccessMessage(d.BankID, d.BankID, 5, "")
	comm.Header = NewMessageHeaderSegment(0, 220, initialDialogID, 1)
	comm.End = NewMessageEndSegment(3, 1)
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
	initMessage := NewDialogInitializationClientMessage()
	messageNum := d.nextMessageNumber()
	initMessage.Header = NewMessageHeaderSegment(-1, 220, initialDialogID, messageNum)
	initMessage.End = NewMessageEndSegment(8, messageNum)
	initMessage.Identification = NewIdentificationSegment(d.BankID, d.ClientID, "0", false)
	initMessage.ProcessingPreparation = NewProcessingPreparationSegment(0, 0, 1)
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

func NewRDHDialog(bankId dataelement.BankId, hbciUrl string, clientId string) *rdhDialog {
	key, err := dataelement.GenerateSigningKey()
	if err != nil {
		panic(err)
	}
	signingKey := dataelement.NewRSAKey(key, dataelement.NewInitialKeyName(bankId.CountryCode, bankId.ID, clientId, "S"))
	provider := NewRDHSignatureProvider(signingKey)
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
	SigningKey  Key
}

func NewDialogInitializationClientMessage() *DialogInitializationClientMessage {
	d := &DialogInitializationClientMessage{}
	d.basicClientMessage = newBasicClientMessage(d)
	return d
}

type DialogInitializationClientMessage struct {
	*basicClientMessage
	Identification             *IdentificationSegment
	ProcessingPreparation      *ProcessingPreparationSegment
	PublicSigningKeyRequest    *PublicKeyRequestSegment
	PublicEncryptionKeyRequest *PublicKeyRequestSegment
}

func (d *DialogInitializationClientMessage) Jobs() SegmentSequence {
	return SegmentSequence{
		d.Identification,
		d.ProcessingPreparation,
		d.PublicSigningKeyRequest,
		d.PublicEncryptionKeyRequest,
	}
}

type DialogInitializationBankMessage struct {
	*basicBankMessage
	BankParams            SegmentSequence
	UserParams            SegmentSequence
	PublicKeyTransmission *PublicKeyTransmissionSegment
	Announcement          *BankAnnouncementSegment
}

type DialogFinishingMessage struct {
	*basicClientMessage
	DialogEnd *DialogEndSegment
}

func (d *DialogFinishingMessage) Jobs() SegmentSequence {
	return SegmentSequence{
		d.DialogEnd,
	}
}

func NewDialogCancellationMessage(messageAcknowledgement *MessageAcknowledgement) *DialogCancellationMessage {
	d := &DialogCancellationMessage{
		MessageAcknowledgements: messageAcknowledgement,
	}
	return d
}

type DialogCancellationMessage struct {
	*basicMessage
	MessageAcknowledgements *MessageAcknowledgement
}

type AnonymousDialogMessage struct {
	*basicMessage
	Identification        *IdentificationSegment
	ProcessingPreparation *ProcessingPreparationSegment
}

func NewDialogEndSegment(dialogId string) *DialogEndSegment {
	d := &DialogEndSegment{
		DialogID: dataelement.NewIdentificationDataElement(dialogId),
	}
	d.Segment = NewBasicSegment("HKEND", 3, 1, d)
	return d
}

type DialogEndSegment struct {
	Segment
	DialogID *dataelement.IdentificationDataElement
}

func (d *DialogEndSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		d.DialogID,
	}
}

func NewProcessingPreparationSegment(bdpVersion int, udpVersion int, language int) *ProcessingPreparationSegment {
	p := &ProcessingPreparationSegment{
		BPDVersion:     dataelement.NewNumberDataElement(bdpVersion, 3),
		UPDVersion:     dataelement.NewNumberDataElement(udpVersion, 3),
		DialogLanguage: dataelement.NewNumberDataElement(language, 3),
		ProductName:    dataelement.NewAlphaNumericDataElement(productName, 25),
		ProductVersion: dataelement.NewAlphaNumericDataElement(productVersion, 5),
	}
	p.Segment = NewBasicSegment("HKVVB", 4, 2, p)
	return p
}

type ProcessingPreparationSegment struct {
	Segment
	BPDVersion *dataelement.NumberDataElement
	UPDVersion *dataelement.NumberDataElement
	// 0 for undefined
	// Sprachkennzeichen | Bedeutung   | Sprachencode ISO 639 | ISO 8859 Subset | ISO 8859- Codeset
	// --------------------------------------------------------------------------------------------
	// 1				 | Deutsch	   | de (German) ￼	      | Deutsch ￼ ￼		| 1 (Latin 1)
	// 2				 | Englisch	   | en (English)		  | Englisch		| 1 (Latin 1)
	// 3 				 | Französisch | fr (French)  		  | Französisch ￼	| 1 (Latin 1)
	DialogLanguage *dataelement.NumberDataElement
	ProductName    *dataelement.AlphaNumericDataElement
	ProductVersion *dataelement.AlphaNumericDataElement
}

func (p *ProcessingPreparationSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		p.BPDVersion,
		p.UPDVersion,
		p.DialogLanguage,
		p.ProductName,
		p.ProductVersion,
	}
}

func NewBankAnnouncementSegment(subject, body string) *BankAnnouncementSegment {
	b := &BankAnnouncementSegment{
		Subject: dataelement.NewAlphaNumericDataElement(subject, 35),
		Body:    dataelement.NewTextDataElement(body, 2048),
	}
	b.Segment = NewBasicSegment("HIKIM", 8, 2, b)
	return b
}

type BankAnnouncementSegment struct {
	Segment
	Subject *dataelement.AlphaNumericDataElement
	Body    *dataelement.TextDataElement
}

func (b *BankAnnouncementSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		b.Subject,
		b.Body,
	}
}
