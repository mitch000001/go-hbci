package transport

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"

	"github.com/mitch000001/go-hbci/segment"
)

type Transport interface {
	Do(*Request) (*Response, error)
}

type TransportFunc func(*Request) (*Response, error)

func (fn TransportFunc) Do(req *Request) (*Response, error) {
	return fn(req)
}

type Middleware func(Transport) Transport

type Request struct {
	URL  string
	Body io.ReadCloser
}

func ReadResponse(r *bufio.Reader, req *Request) (*Response, error) {
	var buf bytes.Buffer
	marshaledMessage, err := ioutil.ReadAll(io.TeeReader(r, &buf))
	if err != nil {
		return nil, err
	}
	extractor := segment.NewSegmentExtractor(marshaledMessage)
	_, err = extractor.Extract()
	if err != nil {
		return nil, err
	}
	response := &Response{
		Request:           req,
		MarshaledResponse: marshaledMessage,
		SegmentExtractor:  extractor,
		Body:              ioutil.NopCloser(&buf),
	}
	return response, nil
}

type Response struct {
	*segment.SegmentExtractor
	Request           *Request
	MarshaledResponse []byte
	Body              io.ReadCloser
}

func (h *Response) IsEncrypted() bool {
	return h.SegmentExtractor.FindSegment(segment.EncryptionHeaderSegmentID) != nil
}
