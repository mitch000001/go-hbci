package hbci

import "github.com/mitch000001/go-hbci/internal"

// Version represents the current library version
const Version = "0.2.0"

// SetDebugMode enables or disables logging on the debug logger
func SetDebugMode(debug bool) {
	internal.SetDebugMode(debug)
}

// SetInfoLog enables or disables logging on the info logger
func SetInfoLog(info bool) {
	internal.SetInfoLog(info)
}
