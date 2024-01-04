package transport

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MockHTTPTransport implements a http.RoundTripper.
// It can be used to mock http.Client details.
type MockHTTPTransport struct {
	requests  []*http.Request
	responses []*http.Response
	errors    []error
	callCount int
}

// RoundTrip satisties the http.RoundTripper interface.
//
// it will store the request to later retrospects and return the stored response
// or the stored error respectively.
// callCount will be incremented on each call.
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.checkAndAdaptBoundaries(req)
	m.requests = append(m.requests, req)
	response, err := m.responses[m.callCount], m.errors[m.callCount]
	m.callCount++
	return response, err
}

// SetResponsePayload registers a given payload in the transport.
// The payload will be encoded as Base64 and wrapped into a io.Reader. This also
// sets a status code of 200 and the error for the given call will be nil.
// The index of the response will be returned.
func (m *MockHTTPTransport) SetResponsePayload(payload []byte) int {
	m.init()
	encodedMessage := base64.StdEncoding.EncodeToString(payload)
	reader := strings.NewReader(encodedMessage)
	m.responses = append(m.responses, &http.Response{
		Body:          io.NopCloser(reader),
		ContentLength: int64(len(payload)),
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
	})
	m.errors = append(m.errors, nil)
	return len(m.responses) - 1
}

// SetResponsePayloads registers multiple payloads in the transport.
// The payloads will be encoded as Base64 and wrapped into a io.Reader. This also
// sets a status code of 200 and the errors for the given calls will be nil.
// This method will reset any responses or errors previously registered.
func (m *MockHTTPTransport) SetResponsePayloads(payloads [][]byte) {
	m.init()
	m.responses = make([]*http.Response, len(payloads))
	m.errors = make([]error, len(payloads))
	for i, payload := range payloads {
		encodedMessage := base64.StdEncoding.EncodeToString(payload)
		reader := strings.NewReader(encodedMessage)
		m.responses[i] = &http.Response{
			Body:          io.NopCloser(reader),
			ContentLength: int64(len(payload)),
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.0",
			ProtoMajor:    1,
			ProtoMinor:    0,
		}
	}
}

// CallCount returns the number of calls received by the transport.
func (m *MockHTTPTransport) CallCount() int {
	return m.callCount
}

// Request returns the http.Request for the current index. If there are less
// requests than the index nil is returned.
func (m *MockHTTPTransport) Request(index int) *http.Request {
	if len(m.requests) < index {
		return nil
	}
	return m.requests[index]
}

// Requests returns all http.Requests the transport has received. If there are
// no requests, an empty slice is returned.
func (m *MockHTTPTransport) Requests() []*http.Request {
	if m.requests == nil {
		return make([]*http.Request, 0)
	}
	return m.requests
}

// Error returns the error for the given index. If there are less errors than
// the index nil is returned.
func (m *MockHTTPTransport) Error(index int) error {
	if len(m.errors) < index {
		return nil
	}
	return m.errors[index]
}

// Errors returns all errors registered in the transport. If there is no
// registered error, an empty slice is returned.
func (m *MockHTTPTransport) Errors() []error {
	if m.errors == nil {
		return make([]error, 0)
	}
	return m.errors
}

// Reset resets the registered responses and the received requests. It also sets
// the callCount to zero.
func (m *MockHTTPTransport) Reset() {
	m.requests = make([]*http.Request, 0)
	m.responses = make([]*http.Response, 0)
	m.errors = make([]error, 0)
	m.callCount = 0
}

func (m *MockHTTPTransport) init() {
	if m.requests == nil {
		m.requests = make([]*http.Request, 0)
	}
	if m.responses == nil {
		m.responses = make([]*http.Response, 0)
	}
	if m.errors == nil {
		m.errors = make([]error, 0)
	}
	m.callCount = 0
}

func (m *MockHTTPTransport) checkAndAdaptBoundaries(req *http.Request) {
	if m.requests == nil {
		m.requests = make([]*http.Request, m.callCount)
	}
	if len(m.responses) <= m.callCount {
		if m.responses == nil {
			m.responses = make([]*http.Response, m.callCount)
		} else {
			m.responses = append(m.responses, nil)
		}
	}
	if len(m.errors) <= m.callCount {
		if m.errors == nil {
			m.errors = make([]error, m.callCount)
		} else {
			bodyBytes, _ := io.ReadAll(req.Body)
			m.errors = append(m.errors, fmt.Errorf("Unexpected request: %+#v\nBody: %q", req, bodyBytes))
		}
	}
}
