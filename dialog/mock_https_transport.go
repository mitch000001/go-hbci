package dialog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/mitch000001/go-hbci/transport"
)

type mockHTTPSTransport struct {
	requests  []*transport.Request
	responses []*transport.Response
	errors    []error
	callCount int
}

func (m *mockHTTPSTransport) Do(request *transport.Request) (*transport.Response, error) {
	m.checkAndAdaptBoundaries(request)
	m.requests = append(m.requests, request)
	response, err := m.responses[m.callCount], m.errors[m.callCount]
	m.callCount++
	return response, err
}

func (m *mockHTTPSTransport) SetResponseMessage(message []byte) {
	response, err := transport.ReadResponse(bufio.NewReader(bytes.NewReader(message)), &transport.Request{})
	if err != nil {
		panic(err)
	}
	m.init()
	m.responses = append(m.responses, response)
	m.errors = append(m.errors, nil)
}

func (m *mockHTTPSTransport) SetResponseMessages(responses [][]byte) {
	m.init()
	m.responses = make([]*transport.Response, len(responses))
	m.errors = make([]error, len(responses))
	for i, response := range responses {
		res, err := transport.ReadResponse(bufio.NewReader(bytes.NewReader(response)), &transport.Request{})
		if err != nil {
			panic(err)
		}
		m.responses[i] = res
	}
}

func (m *mockHTTPSTransport) CallCount() int {
	return m.callCount
}

func (m *mockHTTPSTransport) Request(index int) *transport.Request {
	if len(m.requests) < index {
		return nil
	}
	return m.requests[index]
}

func (m *mockHTTPSTransport) Requests() []*transport.Request {
	if m.requests == nil {
		return make([]*transport.Request, 0)
	}
	return m.requests
}

func (m *mockHTTPSTransport) Error(index int) error {
	if len(m.errors) < index {
		return nil
	}
	return m.errors[index]
}

func (m *mockHTTPSTransport) Errors() []error {
	if m.errors == nil {
		return make([]error, 0)
	}
	return m.errors
}

func (m *mockHTTPSTransport) Reset() {
	m.requests = make([]*transport.Request, 0)
	m.responses = make([]*transport.Response, 0)
	m.errors = make([]error, 0)
	m.callCount = 0
}

func (m *mockHTTPSTransport) init() {
	if m.requests == nil {
		m.requests = make([]*transport.Request, 0)
	}
	if m.responses == nil {
		m.responses = make([]*transport.Response, 0)
	}
	if m.errors == nil {
		m.errors = make([]error, 0)
	}
	m.callCount = 0
}

func (m *mockHTTPSTransport) checkAndAdaptBoundaries(req *transport.Request) {
	if m.requests == nil {
		m.requests = make([]*transport.Request, m.callCount)
	}
	if len(m.responses) <= m.callCount {
		if m.responses == nil {
			m.responses = make([]*transport.Response, m.callCount)
		}
		m.responses = append(m.responses, nil)
	}
	if len(m.errors) <= m.callCount {
		if m.errors == nil {
			m.errors = make([]error, m.callCount)
		}
		reqBytes, err := io.ReadAll(req.Body)
		if err != nil {
			m.errors = append(m.errors, fmt.Errorf("Unexpected request: %+#v\nBody: %v", req, err))
		} else {
			m.errors = append(m.errors, fmt.Errorf("Unexpected request: %+#v\nBody: %q", req, reqBytes))
		}
	}
}
