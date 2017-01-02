package transport

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"

	"github.com/mitch000001/go-hbci/segment"
)

// Transport defines an interface over the various ways data got exchanged with
// HBCI servers. It is used by higher level implementations to perform requests
// against the bank servers.
type Transport interface {
	Do(*Request) (*Response, error)
}

// The TransportFunc type is an adapter to allow the use of ordinary functions
// as Transport handlers. If f is a function with the appropriate signature,
// TransportFunc(f) is a Transport that calls f.
type TransportFunc func(*Request) (*Response, error)

// Do calls fn(req).
func (fn TransportFunc) Do(req *Request) (*Response, error) {
	return fn(req)
}

// Middleware defines the interface for writing middleware for transports
type Middleware func(Transport) Transport

// A Request represents a client request to a HBCI server
type Request struct {
	// URL specifies the URI being requested
	//
	// For the sake of simplicity this is a string instead a *url.URL right now
	URL string
	// Body is the request's body.
	//
	// Body has always to be non-nil
	Body io.ReadCloser
}

// ReadResponse reads and returns a Response from r. It populates the embedded
// SegmentExtractor	to have it ready to use.
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

// A Response represents a server response from a HBCI server
type Response struct {
	// The SegmentExtractor can be used to conveniently query segments from the
	// Response.
	//
	// Right now the Response embeds the SegmentExtractor, which violates the SRP.
	// Future Response types may get rid of it.
	*segment.SegmentExtractor
	// Request is the request that was sent to obtain this Response.
	// Request's Body is nil (having already been consumed).
	Request *Request
	// MarshaledResponse represents the body payload as byte slice. It's strongly
	// discouraged to use that as it may or may not be populated by the Transport.
	//
	// It's only there for legacy reasons.
	MarshaledResponse []byte
	// Body represents the response body.
	Body io.ReadCloser
}

func (h *Response) IsEncrypted() bool {
	return h.SegmentExtractor.FindSegment(segment.EncryptionHeaderSegmentID) != nil
}
