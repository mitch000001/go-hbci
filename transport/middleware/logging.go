package transport

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mitch000001/go-hbci/internal"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/transport"
)

// Logging creates a middleware that logs every request and response sent over
// the transport
func Logging(logger *log.Logger, cryptoProvider message.CryptoProvider) transport.Middleware {
	if logger == nil {
		logger = internal.Debug
	}
	return func(t transport.Transport) transport.Transport {
		return transport.Func(func(req *transport.Request) (*transport.Response, error) {
			var buf bytes.Buffer
			marshaledRequest, err := io.ReadAll(io.TeeReader(req.Body, &buf))
			if err != nil {
				logger.Printf("Error reading request body: %v", err)
			}
			req.Body = io.NopCloser(&buf)
			logger.Println("Decrypted Request:")
			if err := writeMessage(cryptoProvider, logger.Writer(), marshaledRequest); err != nil {
				logger.Printf("%s\n", marshaledRequest)
			}

			res, err := t.Do(req)
			if err != nil {
				logger.Printf("Error executing request:\n%v", err)
				return nil, err
			}
			var responseBuf bytes.Buffer
			marshaledResponse, err := io.ReadAll(io.TeeReader(res.Body, &responseBuf))
			if err != nil {
				logger.Printf("Error reading response body: %v", err)
			}
			res.Body = io.NopCloser(&responseBuf)
			logger.Println("Decrypted Response:")
			if err := writeMessage(cryptoProvider, logger.Writer(), marshaledResponse); err != nil {
				logger.Printf("%s\n", marshaledResponse)
			}
			return res, nil
		})
	}
}

func writeMessage(cryptoProvider message.CryptoProvider, w io.Writer, marshaledMessage []byte) error {
	bankMessage, err := readMessageData(cryptoProvider, marshaledMessage)
	if err != nil {
		return writeRawSegments(w, marshaledMessage)
	}
	segmentProvider, ok := bankMessage.(marshaledSegmentsProvider)
	if !ok {
		return writeRawSegments(w, marshaledMessage)
	}
	var errs errors
	if _, err := fmt.Fprintf(w, "\t%v\n", bankMessage.MessageHeader()); err != nil {
		errs = append(errs, err)
	}
	for _, s := range segmentProvider.MarshaledSegments() {
		_, err := fmt.Fprintf(w, "\t%s\n", string(s))
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}

func writeRawSegments(w io.Writer, marshaledMessage []byte) error {
	extractor := message.NewSegmentExtractor(marshaledMessage)
	segments, err := extractor.Extract()
	if err != nil {
		return fmt.Errorf("error extracting segments from message: %v", err)
	}
	var errs errors
	for _, s := range segments {
		_, err := fmt.Fprintf(w, ">>%s\n", string(s))
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}

func readMessageData(cryptoProvider message.CryptoProvider, rawMessage []byte) (message.Message, error) {
	response, err := transport.ReadResponse(bufio.NewReader(bytes.NewReader(rawMessage)), nil)
	if err != nil {
		return nil, fmt.Errorf("transport#ReadRequest: %v", err)
	}

	var bankMessage message.BankMessage
	if response.IsEncrypted() {
		encMessage, err := extractEncryptedMessage(response)
		if err != nil {
			return nil, err
		}

		decryptedMessage, err := encMessage.Decrypt(cryptoProvider)
		if err != nil {
			return nil, fmt.Errorf("error while decrypting message: %v", err)
		}
		bankMessage = decryptedMessage
	} else {
		decryptedMessage, err := extractUnencryptedMessage(response)
		if err != nil {
			return nil, err
		}
		bankMessage = decryptedMessage
	}
	return bankMessage, nil
}

func extractEncryptedMessage(response *transport.Response) (*message.EncryptedMessage, error) {
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

	encMessage := message.NewEncryptedMessage(header, nil, segment.FINTS300)

	encryptedData := response.FindSegment("HNVSD")
	if encryptedData != nil {
		encSegment := &segment.EncryptedDataSegment{}
		err = encSegment.UnmarshalHBCI(encryptedData)
		if err != nil {
			return nil, fmt.Errorf("Error while unmarshaling encrypted data: %v", err)
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
		return nil, fmt.Errorf("error while unmarshaling message header: %v", err)
	}
	// TODO: parse messageEnd
	decryptedMessage, err := message.NewDecryptedMessage(header, nil, response.MarshaledResponse)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling unencrypted data: %v", err)
	}
	return decryptedMessage, nil
}

type errors []error

func (e errors) Error() string {
	errs := make([]string, len(e))
	for i, err := range e {
		errs[i] = err.Error()
	}
	return strings.Join(errs, ",")
}

type marshaledSegmentsProvider interface {
	MarshaledSegments() [][]byte
}
