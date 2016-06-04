package config

import (
	"time"
)

const (
	CleanupInterval = 5 * time.Minute
	TCPPort         = 15000
	NoOfBuckets     = 10
	Timeout         = 5 * time.Minute
)

var (
	DBFileName = ""
)
