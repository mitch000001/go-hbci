package dialog

import (
	"fmt"

	"github.com/mitch000001/go-hbci/transport"
)

type MockHttpsTransport struct {
	requests  []*transport.Request
	responses []*transport.Response
	errors    []error
	callCount int
}

func (m *MockHttpsTransport) init() {
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

func (m *MockHttpsTransport) checkAndAdaptBoundaries(req *transport.Request) {
	if m.requests == nil {
		m.requests = make([]*transport.Request, m.callCount)
	}
	if len(m.responses) <= m.callCount {
		if m.responses == nil {
			m.responses = make([]*transport.Response, m.callCount)
		} else {
			m.responses = append(m.responses, nil)
		}
	}
	if len(m.errors) <= m.callCount {
		if m.errors == nil {
			m.errors = make([]error, m.callCount)
		} else {
			m.errors = append(m.errors, fmt.Errorf("Unexpected request: %+#v\nBody: %q", req, req.MarshaledMessage))
		}
	}
}

func (m *MockHttpsTransport) Do(request *transport.Request) (*transport.Response, error) {
	m.checkAndAdaptBoundaries(request)
	m.requests = append(m.requests, request)
	response, err := m.responses[m.callCount], m.errors[m.callCount]
	m.callCount += 1
	return response, err
}

func (m *MockHttpsTransport) SetResponseMessage(message []byte) {
	response, err := transport.ReadResponse(message, &transport.Request{})
	if err != nil {
		panic(err)
	}
	m.init()
	m.responses = append(m.responses, response)
}

func (m *MockHttpsTransport) SetResponseMessages(responses [][]byte) {
	m.init()
	m.responses = make([]*transport.Response, len(responses))
	m.errors = make([]error, len(responses))
	for i, response := range responses {
		res, err := transport.ReadResponse(response, &transport.Request{})
		if err != nil {
			panic(err)
		}
		m.responses[i] = res
	}
}

func (m *MockHttpsTransport) CallCount() int {
	return m.callCount
}

func (m *MockHttpsTransport) Request(index int) *transport.Request {
	if len(m.requests) < index {
		return nil
	} else {
		return m.requests[index]
	}
}

func (m *MockHttpsTransport) Requests() []*transport.Request {
	if m.requests == nil {
		return make([]*transport.Request, 0)
	} else {
		return m.requests
	}
}

func (m *MockHttpsTransport) Error(index int) error {
	if len(m.errors) < index {
		return nil
	} else {
		return m.errors[index]
	}
}

func (m *MockHttpsTransport) Errors() []error {
	if m.errors == nil {
		return make([]error, 0)
	} else {
		return m.errors
	}
}

func (m *MockHttpsTransport) Reset() {
	m.requests = make([]*transport.Request, 0)
	m.responses = make([]*transport.Response, 0)
	m.errors = make([]error, 0)
	m.callCount = 0
}
