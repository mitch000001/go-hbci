package hbci

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

type mockHttpTransport struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (m *mockHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.request = req
	return m.response, m.err
}

func (m *mockHttpTransport) setResponsePayload(payload []byte) {
	encodedMessage := base64.StdEncoding.EncodeToString(payload)
	reader := strings.NewReader(encodedMessage)
	m.response = &http.Response{
		Body:          ioutil.NopCloser(reader),
		ContentLength: int64(len(payload)),
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
	}
}

func TestDialogSyncClientID(t *testing.T) {
	transport := &mockHttpTransport{}
	httpClient := &http.Client{Transport: transport}

	url := "http://localhost"
	clientID := "12345"
	bankID := domain.BankId{280, "10000000"}
	dialog := NewPinTanDialog(bankID, url, clientID)
	dialog.SetPin("abcde")
	dialog.httpClient = httpClient

	transport.setResponsePayload([]byte("HNHBK:1:3+++abcde'"))

	res, err := dialog.SyncClientSystemID()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := ""

	if res != expected {
		t.Logf("Expected response to equal\n%q\n\tgot\n%q\n", expected, res)
		t.Fail()
	}

}
