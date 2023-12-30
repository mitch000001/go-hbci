package transport

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mitch000001/go-hbci/transport"
	middleware "github.com/mitch000001/go-hbci/transport/middleware"
)

func TestHttpsBase64Transport(t *testing.T) {
	response := []byte("HNHBK:1:3+abc'HNDGC:2:3+def'")
	roundtripper := &MockHTTPTransport{}
	roundtripper.SetResponsePayloads([][]byte{
		response,
		response,
	})
	httpClient := &http.Client{Transport: roundtripper}

	httpsTransport := middleware.Base64Encoding(base64.StdEncoding)(&HTTPSTransport{httpClient})
	httpsBase64Transport := &HTTPSBase64Transport{httpClient}

	httpReq := &transport.Request{
		URL:  "foo",
		Body: ioutil.NopCloser(strings.NewReader("bar")),
	}

	httpBase64Req := &transport.Request{
		URL:  "foo",
		Body: ioutil.NopCloser(strings.NewReader("bar")),
	}

	httpResponse, httpError := httpsTransport.Do(httpReq)
	httpBase64Respose, httpBase64Error := httpsBase64Transport.Do(httpBase64Req)

	if httpError != nil {
		t.Logf("HTTP: Expected no error, got %v\n", httpError)
		t.Fail()
	}

	if httpBase64Error != nil {
		t.Logf("HTTPBase64: Expected no error, got %v\n", httpBase64Error)
		t.Fail()
	}

	if httpResponse == nil {
		t.Logf("HTTP: Expected response not to be nil\n")
		t.Fail()
	}

	if httpBase64Respose == nil {
		t.Logf("HTTBase64: Expected response not to be nil\n")
		t.Fail()
	}

	httpResponseBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		panic(err)
	}
	httpBase64ResponseBytes, err := io.ReadAll(httpBase64Respose.Body)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(httpBase64ResponseBytes, httpResponseBytes) {
		t.Logf("Expected response body to equal\n%q\n\tgot\n%q\n", httpBase64ResponseBytes, httpResponseBytes)
		t.Fail()
	}

	requests := roundtripper.Requests()

	if len(requests) != 2 {
		t.Logf("Expected 2 requests, got %d\n", len(requests))
		t.FailNow()
	}

	if requests[0] == nil {
		t.Logf("HTTP: Expected request not to be nil\n")
		t.Fail()
	}

	if requests[1] == nil {
		t.Logf("HTTBase64: Expected request not to be nil\n")
		t.Fail()
	}

	httpRequest, err := io.ReadAll(requests[0].Body)
	if err != nil {
		panic(err)
	}
	httpBase64Request, err := io.ReadAll(requests[1].Body)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(httpBase64Request, httpRequest) {
		t.Logf("Expected request to equal\n%q\n\tgot\n%q\n", httpBase64Request, httpRequest)
		t.Fail()
	}
}
