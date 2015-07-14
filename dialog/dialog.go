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
)

const initialDialogID = "0"
const initialClientSystemID = "0"
const anonymousClientID = "9999999999"

type Dialog interface {
	Init() (string, error)
}

func newDialog(bankId domain.BankId, hbciUrl string, clientId string, signatureProvider message.SignatureProvider, encryptionProvider message.EncryptionProvider) *dialog {
	return &dialog{
		httpClient:         http.DefaultClient,
		hbciUrl:            hbciUrl,
		BankID:             bankId,
		ClientID:           clientId,
		ClientSystemID:     initialClientSystemID,
		signatureProvider:  signatureProvider,
		encryptionProvider: encryptionProvider,
	}
}

type dialog struct {
	httpClient         *http.Client
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
