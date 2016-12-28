package transport

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
)

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
