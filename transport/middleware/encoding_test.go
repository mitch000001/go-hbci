package transport

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/mitch000001/go-hbci/transport"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type mockEncoding struct {
	encoder *encoding.Encoder
	decoder *encoding.Decoder
}

func (m mockEncoding) NewEncoder() *encoding.Encoder {
	return m.encoder
}

func (m mockEncoding) NewDecoder() *encoding.Decoder {
	return m.decoder
}

type mockTransform struct {
	transform.Transformer
	transformCallCount int
	resetCallCount     int
}

func (m *mockTransform) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	m.transformCallCount++
	buf := bytes.ToUpper(src)
	return m.Transformer.Transform(dst, buf, atEOF)
}

func (m *mockTransform) Reset() {
	m.resetCallCount++
	m.Transformer.Reset()
}

func TestUTF8Encoding(t *testing.T) {
	called := false
	transportResponse := transport.Response{Body: ioutil.NopCloser(strings.NewReader("qux"))}
	var transportRequest *transport.Request
	innerTransport := transport.Func(func(req *transport.Request) (*transport.Response, error) {
		called = true
		transportRequest = req
		response := transportResponse
		return &response, nil
	})
	encoderTransform := mockTransform{Transformer: encoding.Nop.NewEncoder().Transformer}
	decoderTransform := mockTransform{Transformer: encoding.Nop.NewDecoder().Transformer}
	encoding := mockEncoding{
		encoder: &encoding.Encoder{
			Transformer: &encoderTransform,
		},
		decoder: &encoding.Decoder{
			Transformer: &decoderTransform,
		},
	}

	middleware := UTF8Encoding(encoding)

	wrappedTransport := middleware(innerTransport)

	request := &transport.Request{
		URL:  "foo",
		Body: ioutil.NopCloser(strings.NewReader("bar")),
	}

	response, err := wrappedTransport.Do(request)

	{
		actual := called
		expected := true
		if expected != actual {
			t.Logf("Expected transform to be called, was not")
			t.Fail()
		}
	}
	{
		actual, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		expected := []byte("QUX")
		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected response Body to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
			t.Fail()
		}
	}
	{
		actual, err := ioutil.ReadAll(transportRequest.Body)
		if err != nil {
			panic(err)
		}
		expected := []byte("BAR")
		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected transport request body to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
			t.Fail()
		}
	}
	{
		if err != nil {
			t.Logf("Expected transport err to equal nil, got\n%+#v\n", err)
			t.Fail()
		}
	}
	{
		actual := encoderTransform.transformCallCount
		expected := 1
		if expected != actual {
			t.Logf("Expected encoder transform to be called once, was %d", actual)
			t.Fail()
		}
	}
	{
		actual := decoderTransform.transformCallCount
		expected := 2
		if expected != actual {
			t.Logf("Expected decoder transform to be called once, was %d", actual)
			t.Fail()
		}
	}
}

func mustDecode(res []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return res
}

func TestBase64Encoding(t *testing.T) {
	encoding := base64.StdEncoding

	called := false
	transportResponse := transport.Response{
		Body: ioutil.NopCloser(strings.NewReader(encoding.EncodeToString([]byte("qux")))),
	}
	var transportRequest *transport.Request
	innerTransport := transport.Func(func(req *transport.Request) (*transport.Response, error) {
		called = true
		transportRequest = req
		response := transportResponse
		return &response, nil
	})

	middleware := Base64Encoding(encoding)

	wrappedTransport := middleware(innerTransport)

	request := &transport.Request{
		URL:  "foo",
		Body: ioutil.NopCloser(strings.NewReader("bar")),
	}

	response, _ := wrappedTransport.Do(request)

	{
		actual := called
		expected := true
		if expected != actual {
			t.Logf("Expected transform to be called, was not")
			t.Fail()
		}
	}
	{
		actual, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		expected := []byte("qux")
		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected response Body to equal\n%s\n\tgot\n%s\n", expected, actual)
			t.Fail()
		}
	}
	{
		actual, err := ioutil.ReadAll(transportRequest.Body)
		if err != nil {
			panic(err)
		}
		expected := []byte(encoding.EncodeToString([]byte("bar")))
		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected transport request body to equal\n%s\n\tgot\n%s\n", expected, actual)
			t.Fail()
		}
	}
}
