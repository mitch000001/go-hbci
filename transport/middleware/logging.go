package transport

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/mitch000001/go-hbci/internal"
	"github.com/mitch000001/go-hbci/transport"
)

// Logging creates a middleware that logs every request and response sent over
// the transport
func Logging(logger *log.Logger) transport.Middleware {
	if logger == nil {
		logger = internal.Debug
	}
	var count int
	return func(t transport.Transport) transport.Transport {
		return transport.Func(func(req *transport.Request) (*transport.Response, error) {
			count++
			var buf bytes.Buffer
			marshaledRequest, err := ioutil.ReadAll(io.TeeReader(req.Body, &buf))
			req.Body = ioutil.NopCloser(&buf)
			logger.Println("Request:")
			logger.Printf("%s\n", marshaledRequest)

			res, err := t.Do(req)
			if err != nil {
				logger.Printf("Error executing request:\n%v", err)
				return nil, err
			}
			var responseBuf bytes.Buffer
			marshaledResponse, err := ioutil.ReadAll(io.TeeReader(res.Body, &responseBuf))
			res.Body = ioutil.NopCloser(&responseBuf)
			logger.Println("Response:")
			logger.Printf("%s\n", marshaledResponse)
			return res, nil
		})
	}
}
