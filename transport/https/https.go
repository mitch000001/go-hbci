package transport

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mitch000001/go-hbci/transport"
)

func NewBase64() transport.Transport {
	return &HttpsBase64Transport{
		httpClient: http.DefaultClient,
	}
}

type HttpsBase64Transport struct {
	httpClient *http.Client
}

func (h *HttpsBase64Transport) Do(request *transport.Request) (*transport.Response, error) {
	var buf bytes.Buffer
	encodingWriter := base64.NewEncoder(base64.StdEncoding, &buf)
	_, err := io.Copy(encodingWriter, request.Body)
	if err != nil {
		return nil, err
	}
	encodingWriter.Close()
	httpResponse, err := h.httpClient.Post(request.URL, "application/vnd.hbci", &buf)
	if err != nil {
		return nil, err
	}
	var reader io.Reader
	if httpResponse.StatusCode == http.StatusOK {
		reader = base64.NewDecoder(base64.StdEncoding, httpResponse.Body)
	} else {
		reader = httpResponse.Body
	}
	return &transport.Response{Body: ioutil.NopCloser(reader), Request: request}, nil
}

func New() transport.Transport {
	return &HttpsTransport{
		httpClient: http.DefaultClient,
	}
}

type HttpsTransport struct {
	httpClient *http.Client
}

func (h *HttpsTransport) Do(request *transport.Request) (*transport.Response, error) {
	httpResponse, err := h.httpClient.Post(request.URL, "application/vnd.hbci", request.Body)
	if err != nil {
		return nil, err
	}
	return &transport.Response{Body: httpResponse.Body, Request: request}, nil
}
