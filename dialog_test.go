package hbci

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
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

func (m *MockHttpTransport) checkAndAdaptBoundaries() {
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
			m.errors = append(m.errors, nil)
		}
	}
}

func (m *MockHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.checkAndAdaptBoundaries()
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

func TestDialogSyncClientID(t *testing.T) {
	transport := &MockHttpTransport{}
	httpClient := &http.Client{Transport: transport}

	url := "http://localhost"
	clientID := "12345"
	bankID := domain.BankId{280, "10000000"}
	dialog := NewPinTanDialog(bankID, url, clientID)
	dialog.SetPin("abcde")
	dialog.httpClient = httpClient

	syncResponseMessage := "HNHBK:1:3++220+abcde+1+'HNVSK:998:2:+998+1+1::0+1:20150713:173634+2:2:13:@8@\x00\x00\x00\x00\x00\x00\x00\x00:5:1:+280:10000000:12345:V:0:0+0+'HNVSD:999:1:+@30@HISYN:3:3:8+newClientSystemID''HNHBS:4:1:+1'"

	transport.SetResponsePayloads([][]byte{
		[]byte(syncResponseMessage),
		[]byte(""),
	})

	res, err := dialog.SyncClientSystemID()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedClientSystemID := "newClientSystemID"

	if dialog.ClientSystemID != expectedClientSystemID {
		t.Logf("Expected ClientSystemID to equal %q, got %q\n", expectedClientSystemID, dialog.ClientSystemID)
		t.Fail()
	}

	expected := ""

	if res != expected {
		t.Logf("Expected response to equal\n%q\n\tgot\n%q\n", expected, res)
		t.Fail()
	}

}
