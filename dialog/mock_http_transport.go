package dialog

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type MockHttpTransport struct {
	requests  []*http.Request
	responses []*http.Response
	errors    []error
	callCount int
}

func (m *MockHttpTransport) init() {
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

func (m *MockHttpTransport) checkAndAdaptBoundaries(req *http.Request) {
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
			bodyBytes, _ := ioutil.ReadAll(req.Body)
			m.errors = append(m.errors, fmt.Errorf("Unexpected request: %+#v\nBody: %q", req, bodyBytes))
		}
	}
}

func (m *MockHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.checkAndAdaptBoundaries(req)
	m.requests = append(m.requests, req)
	response, err := m.responses[m.callCount], m.errors[m.callCount]
	m.callCount += 1
	return response, err
}

func (m *MockHttpTransport) SetResponsePayload(payload []byte) {
	m.init()
	encodedMessage := base64.StdEncoding.EncodeToString(payload)
	reader := strings.NewReader(encodedMessage)
	m.responses = append(m.responses, &http.Response{
		Body:          ioutil.NopCloser(reader),
		ContentLength: int64(len(payload)),
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
	})
}

func (m *MockHttpTransport) SetResponsePayloads(payloads [][]byte) {
	m.init()
	m.responses = make([]*http.Response, len(payloads))
	m.errors = make([]error, len(payloads))
	for i, payload := range payloads {
		encodedMessage := base64.StdEncoding.EncodeToString(payload)
		reader := strings.NewReader(encodedMessage)
		m.responses[i] = &http.Response{
			Body:          ioutil.NopCloser(reader),
			ContentLength: int64(len(payload)),
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.0",
			ProtoMajor:    1,
			ProtoMinor:    0,
		}
	}
}

func (m *MockHttpTransport) CallCount() int {
	return m.callCount
}

func (m *MockHttpTransport) Request(index int) *http.Request {
	if len(m.requests) < index {
		return nil
	} else {
		return m.requests[index]
	}
}

func (m *MockHttpTransport) Requests() []*http.Request {
	if m.requests == nil {
		return make([]*http.Request, 0)
	} else {
		return m.requests
	}
}

func (m *MockHttpTransport) Error(index int) error {
	if len(m.errors) < index {
		return nil
	} else {
		return m.errors[index]
	}
}

func (m *MockHttpTransport) Errors() []error {
	if m.errors == nil {
		return make([]error, 0)
	} else {
		return m.errors
	}
}

func (m *MockHttpTransport) Reset() {
	m.requests = make([]*http.Request, 0)
	m.responses = make([]*http.Response, 0)
	m.errors = make([]error, 0)
	m.callCount = 0
}
