package log

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)
import "testing"

func TestRotatingLog(t *testing.T) {
	// Create temp dir for unit test
	tmpdir, err := ioutil.TempDir(".", "ut-rotatelog-")
	if err != nil {
		t.Fatalf("Temp dir creation failed: %q", err)
	}

	rl := NewCustomRotatingLog(tmpdir, "unittest", func() time.Duration { return 4 * time.Second })
	defer rl.Close()

	// Message 1 goes to one file and Message 2 to the next
	msg1 := "Message 1"
	msg2 := "Message 2"
	rl.Write([]byte(msg1))
	time.Sleep(5 * time.Second)
	rl.Write([]byte(msg2))
	time.Sleep(1 * time.Second)

	matches, err := filepath.Glob(tmpdir + "/unittest-*.log")
	if err != nil || len(matches) != 2 {
		t.Fatalf("Expected 2 log files. Found %q. Err: %q", matches, err)
	}

	// Sorted list contains unittest-YYYY-mm-dd-<ts>.log and unittest-YYYY-mm-dd.log
	sort.Strings(matches)
	fp1, err1 := os.Open(matches[1])
	fp2, err2 := os.Open(matches[0])
	if fp1 == nil || fp2 == nil {
		t.Fatalf("Log file open failed. Err: %q, %q", err1, err2)
	}

	contents1, _ := ioutil.ReadAll(fp1)
	contents2, _ := ioutil.ReadAll(fp2)
	if string(contents1) != msg1 || string(contents2) != msg2 {
		t.Fatalf("Log file contents mismatch. File 1: %s. File 2: %s", string(contents1), string(contents2))
	}

	os.RemoveAll(tmpdir)
}

func TestNextDay(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(nextDay())

	if (tomorrow.YearDay() - now.YearDay()) != 1 {
		t.Fatalf("Next day calculation error. Now=%s Now+nextDay=%s", now, tomorrow)
	}
}
