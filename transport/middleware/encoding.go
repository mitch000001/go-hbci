package transport

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/mitch000001/go-hbci/transport"

	"golang.org/x/text/encoding"
)

func UTF8Encoding(encoding encoding.Encoding) transport.Middleware {
	return func(t transport.Transport) transport.Transport {
		return transport.TransportFunc(func(req *transport.Request) (*transport.Response, error) {
			var buf bytes.Buffer
			encodingWriter := encoding.NewEncoder().Writer(&buf)
			_, err := io.Copy(encodingWriter, req.Body)
			if err != nil {
				return nil, err
			}
			encodedRequest := req
			encodedRequest.Body = ioutil.NopCloser(&buf)
			response, err := t.Do(encodedRequest)
			if err != nil {
				return nil, err
			}
			decodingReader := encoding.NewDecoder().Reader(response.Body)
			decodedResponse := response
			decodedResponse.Body = ioutil.NopCloser(decodingReader)
			return decodedResponse, nil
		})
	}
}

func Base64Encoding(encoding *base64.Encoding) transport.Middleware {
	return func(t transport.Transport) transport.Transport {
		return transport.TransportFunc(func(req *transport.Request) (*transport.Response, error) {
			var buf bytes.Buffer
			encodingWriter := base64.NewEncoder(encoding, &buf)
			_, err := io.Copy(encodingWriter, req.Body)
			if err != nil {
				return nil, err
			}
			encodingWriter.Close()
			encodedRequest := req
			encodedRequest.Body = ioutil.NopCloser(&buf)
			response, err := t.Do(encodedRequest)
			if err != nil {
				return nil, err
			}
			decodingReader := base64.NewDecoder(encoding, response.Body)
			decodedResponse := response
			decodedResponse.Body = ioutil.NopCloser(decodingReader)
			return decodedResponse, nil
		})
	}
}
