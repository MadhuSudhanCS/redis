package log

import (
	"fmt"
	"os"
	"time"
)

// RotatingLog is a log that generates one file per day in a specified log dir.
type RotatingLog struct {
	logDir         string
	prefix         string
	nextRotationFn func() time.Duration

	fp    *os.File
	logC  chan []byte
	quitC chan struct{}
}

const logBufferSize = 100

// NewRotatingLog initializes a log file that rotates daily.
// Creates one file per day (<prefix>-YYYY-mm-dd) in logDir.
func NewRotatingLog(logDir, prefix string) *RotatingLog {
	return NewCustomRotatingLog(logDir, prefix, nextDay)
}

// NewCustomRotatingLog initializes a log file that rotates at an interval determined by nextRotationFn.
func NewCustomRotatingLog(logDir, prefix string, nextRotationFn func() time.Duration) *RotatingLog {
	rl := &RotatingLog{
		logC:           make(chan []byte, logBufferSize),
		quitC:          make(chan struct{}),
		logDir:         logDir,
		prefix:         prefix,
		nextRotationFn: nextRotationFn,
	}
	go rl.start()
	return rl
}

// Write implements the io.Writer interface
func (rl *RotatingLog) Write(p []byte) (n int, err error) {
	rl.logC <- p
	return len(p), nil
}

// Close implements the io.Closer interface
func (rl *RotatingLog) Close() error {
	close(rl.quitC)
	return nil
}

func (rl *RotatingLog) start() {
	rl.initLog(true)
	rotateTimer := time.NewTimer(rl.nextRotationFn())

	for {
		select {
		case p := <-rl.logC:
			if rl.fp != nil {
				if n, err := rl.fp.Write(p); err != nil {
					fmt.Printf("ERROR %s: Write error. Wrote %d bytes of '%s' (%d bytes). Err: %s", rl.fp.Name(), n, p, len(p), err)
				}
			}

		case <-rotateTimer.C:
			rl.initLog(false)
			rotateTimer = time.NewTimer(rl.nextRotationFn())

		case <-rl.quitC:
			if rl.fp != nil {
				rl.fp.Close()
			}
			return
		}
	}
}

func (rl *RotatingLog) String() string {
	return fmt.Sprintf("RotatingLog {Dir: %s, prefix: %s, logFile: %s}", rl.logDir, rl.prefix, rl.fp.Name())
}

func (rl *RotatingLog) initLog(useExistingFile bool) {
	if rl.fp != nil {
		rl.fp.Sync()
		rl.fp.Close()
		rl.fp = nil
	}

	now := time.Now()

	fp, err := os.OpenFile(rl.logFile(now, useExistingFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("ERROR Initializing log file: %q", err)
		return
	}
	rl.fp = fp
}

func (rl *RotatingLog) logFile(t time.Time, useExisting bool) string {
	fname := fmt.Sprintf("%s/%s-%4d-%02d-%02d.log", rl.logDir, rl.prefix, t.Year(), t.Month(), t.Day())
	if fi, _ := os.Stat(fname); fi != nil && !useExisting {
		// File for this day exists. Best effort attempt to generate a new filename by appending Unix epoch.
		fname = fmt.Sprintf("%s/%s-%4d-%02d-%02d-%d.log", rl.logDir, rl.prefix, t.Year(), t.Month(), t.Day(), t.Unix())
	}
	return fname
}

func nextDay() time.Duration {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 1, 0, now.Location())
	return tomorrow.Sub(now)
}
