// Package log is a simple wrapper over golang std logger for minimal log level support.
//
// - Debug logs will be logged only if debug flag is true.
// - Error and Warn logs will have ERROR and WARN preprended.
//
package log

import (
	"fmt"
	"log"
	"os"
)

var debug = false

func init() {
	if os.Getenv("DEBUG") == "1" {
		SetLevel("DEBUG")
	}
}

// SetLevel - use "DEBUG" to enable debug logging
func SetLevel(l string) {
	debug = (l == "DEBUG")
	log.Printf("Set log level debug = %t", debug)
}

// IsDebug returns true if log level is DEBUG.
func IsDebug() bool {
	return debug
}

// Printf - These logs cannot be suppresed.
func Printf(format string, l ...interface{}) {
	log.Printf(format, l...)
}

// Debugf - These logs will be suppresed unless log level=DEBUG.
func Debugf(format string, l ...interface{}) {
	if debug {
		log.Printf("DEBUG %s", fmt.Sprintf(format, l...))
	}
}

// Errorf logs with ERROR preprended. These logs cannot be suppressed.
func Errorf(format string, l ...interface{}) {
	log.Printf("ERROR %s", fmt.Sprintf(format, l...))
}

// Warnf logs with WARN preprended. These logs cannot be suppressed.
func Warnf(format string, l ...interface{}) {
	log.Printf("WARN %s", fmt.Sprintf(format, l...))
}

// Fatalf logs with FATAL preprended and exit.
func Fatalf(format string, l ...interface{}) {
	log.Fatalf("FATAL %s", fmt.Sprintf(format, l...))
}
