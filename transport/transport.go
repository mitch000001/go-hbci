package transport

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mitch000001/go-hbci/segment"
)

type Transport interface {
	Do(*Request) (*Response, error)
}

func NewHttpsTransport() *HttpsTransport {
	return &HttpsTransport{
		httpClient: http.DefaultClient,
	}
}

type HttpsTransport struct {
	httpClient *http.Client
}

func (h *HttpsTransport) Do(request *Request) (*Response, error) {
	encodedMessage := base64.StdEncoding.EncodeToString(request.MarshaledMessage)
	httpResponse, err := h.httpClient.Post(request.URL, "application/vnd.hbci", strings.NewReader(encodedMessage))
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	var marshaledResponse []byte
	if httpResponse.StatusCode == http.StatusOK {
		decodedReader := base64.NewDecoder(base64.StdEncoding, httpResponse.Body)
		marshaledResponse, err = ioutil.ReadAll(decodedReader)
		if err != nil {
			return nil, err
		}
	} else {
		marshaledResponse, err = ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
	}
	return ReadResponse(marshaledResponse, request)
}

func NewRequest() *Request {
	return &Request{}
}

type Request struct {
	URL              string
	MarshaledMessage []byte
}

func ReadResponse(marshaledMessage []byte, request *Request) (*Response, error) {
	extractor := segment.NewSegmentExtractor(marshaledMessage)
	_, err := extractor.Extract()
	if err != nil {
		return nil, err
	}
	response := &Response{
		Request:           request,
		marshaledResponse: marshaledMessage,
		SegmentExtractor:  extractor,
	}
	return response, nil
}

type Response struct {
	*segment.SegmentExtractor
	Request           *Request
	marshaledResponse []byte
}

func (h *Response) IsEncrypted() bool {
	return false
}
