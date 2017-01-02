package transport

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mitch000001/go-hbci/transport"
)

// New returns a HttpsBase64Transport. It sets http.DefaultClient as http.Client
// to perform requests to the HBCI server.
//
// Each request will be encoded and each response will be decoded with Base64
// encoding.
func NewBase64() transport.Transport {
	return &HttpsBase64Transport{
		httpClient: http.DefaultClient,
	}
}

// A HttpsTransport inplements transport.Transport and performs request over
// HTTPS with also encoding requests and responses with Base64 encoding.
type HttpsBase64Transport struct {
	httpClient *http.Client
}

// Do performs the request to the HBCI server. If successful, it returns a
// populated transport.Response with the HTTP Response Body as Body and the
// request as Request.
//
// Before sending the request it will be encoded with Base64 encoding.
// When receiving the response with a status code 200 it will decode the response
// with Base64 encoding. A non 200 status code will be returned as is, without
// decoding it from Base64.
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

// New returns a HttpsTransport. It sets http.DefaultClient as http.Client to
// perform requests to the HBCI server
func New() transport.Transport {
	return &HttpsTransport{
		httpClient: http.DefaultClient,
	}
}

// A HttpsTransport inplements transport.Transport and performs request over HTTPS
type HttpsTransport struct {
	httpClient *http.Client
}

// Do performs the request to the HBCI server. If successful, it returns a
// populated transport.Response with the HTTP Response Body as Body and the
// request as Request
func (h *HttpsTransport) Do(request *transport.Request) (*transport.Response, error) {
	httpResponse, err := h.httpClient.Post(request.URL, "application/vnd.hbci", request.Body)
	if err != nil {
		return nil, err
	}
	return &transport.Response{Body: httpResponse.Body, Request: request}, nil
}
