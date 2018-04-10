package internal

import (
	"flag"
	"io"
	"log"
	"os"
)

func init() {
	flag.BoolVar(&debugMode, "debug", false, "true")
}

var debugMode = false
var infoMode = true
var debugFlags = log.Lshortfile | log.LstdFlags

// DebugMode returns true if debugMode is set, false otherwise
func DebugMode() bool {
	return debugMode
}

// SetDebugMode sets the debug mode
func SetDebugMode(debug bool) {
	debugMode = debug
}

// InfoLog returns true if logging to info level is enabled, false otherwise
func InfoLog() bool {
	return infoMode
}

// SetInfoLog sets whether the info log should be enabled
func SetInfoLog(info bool) {
	infoMode = info
}

// WithDebugLogging enables debug logging for the given function fn
func WithDebugLogging(fn func()) {
	debugMode = true
	fn()
	debugMode = false
}

// Debug is a logger which logs when debugMode is enabled. If disabled, it executes a noop
var Debug = newConditionalLogger(os.Stderr, "go-hbci: ", debugFlags, &debugMode)

// Info is a logger which logs when infoMode is enabled. If disabled, it executes a noop
var Info = newConditionalLogger(os.Stderr, "go-hbci: ", log.LstdFlags, &infoMode)

func newConditionalLogger(w io.Writer, prefix string, flag int, condition *bool) *log.Logger {
	condWriter := newConditionalWriter(w, condition)
	return log.New(condWriter, prefix, flag)
}

func newConditionalWriter(w io.Writer, condition *bool) io.Writer {
	return &conditionalWriter{Writer: w, condition: condition}
}

type conditionalWriter struct {
	condition *bool
	io.Writer
}

func (c *conditionalWriter) Write(p []byte) (n int, err error) {
	if *c.condition {
		return c.Writer.Write(p)
	}
	return 0, nil
}
