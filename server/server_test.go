package server

import (
	"bufio"
	"fmt"
	"github.com/madhusudhancs/redis/config"
	"net"
	"os"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	config.DBFileName = "/tmp/dump.json"
	if _, err := NewServer(); err != nil {
		t.Fatalf("failed to initiate server. err: %v", err)
	}
}

func TestCleanup(t *testing.T) {
	config.DBFileName = "/tmp/dump.json"
	s, err := NewServer()
	if err != nil {
		t.Fatalf("failed to initiate server. err: %v", err)
	}

	s.cleanup()
}

func TestHandleSigInt(t *testing.T) {
	config.DBFileName = "/tmp/dump.json"
	s, err := NewServer()
	if err != nil {
		t.Fatalf("failed to initiate server. err: %v", err)
	}

	sigC = make(chan os.Signal, 1)
	cleanupC = make(chan bool, 1)
	sigC <- *new(os.Signal)
	s.handleIntSig()
	<-cleanupC
}

func TestStart(t *testing.T) {
	config.DBFileName = "/tmp/dump.json"
	s, err := NewServer()
	if err != nil {
		t.Fatalf("failed to initiate server. err: %v", err)
	}

	go s.Start()

	time.Sleep(10 * time.Second)

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", config.TCPPort))
	if err != nil {
		t.Fatalf("failed to connect to server. err: %v", err)
	}

	defer conn.Close()
	time.Sleep(10 * time.Second)
	writer := bufio.NewWriter(conn)

	//issue set command
	if _, err := writer.WriteString("set madhu sudhan"); err != nil {
		t.Fatalf("failed to issue set command. err: %v", err)
	}

	//issue get command
	if _, err := writer.WriteString("get madhu"); err != nil {
		t.Fatalf("failed to issue set command. err: %v", err)
	}
}
