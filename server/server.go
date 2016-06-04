package server

import (
	"bufio"
	"fmt"
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/cmd"
	"github.com/madhusudhancs/redis/config"
	"github.com/madhusudhancs/redis/utils/log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	cleanupC chan bool
	sigC     chan os.Signal
)

type Server struct {
	cache *cache.Cache
}

//Initialize Server with cache
func NewServer() (*Server, error) {
	c := cache.NewCache(config.NoOfBuckets)

	t0 := time.Now()
	resp, err := c.Load(config.DBFileName)
	if err != nil {
		log.Errorf("Server: failed to load cache. err: %v", err)
		return nil, fmt.Errorf("Server: failed to load cache. err: %v", err)
	}

	if resp == "!OK" {
		log.Printf("Server: file '%s' does not exist", config.DBFileName)
	}

	log.Printf("Server: time taken to load db items to cache: %v", time.Since(t0))

	server := &Server{
		cache: c,
	}

	return server, nil
}

// cleanup of resources during SIGINT CTRL+C
func (s *Server) cleanup() {
	log.Printf("Server: performing cleanup.")

	log.Printf("Server: save the cache to db file: %s", config.DBFileName)
	if _, err := s.cache.Save(config.DBFileName); err != nil {
		log.Errorf("Server: failed to save cache to file: %s. err: %v", config.DBFileName, err)
	}
}

func (s *Server) handleIntSig() {
	<-sigC
	log.Printf("Server: interrupt signal received.")
	s.cleanup()
	cleanupC <- true
}

func (s *Server) registerIntSig() {
	sigC = make(chan os.Signal, 1)
	cleanupC = make(chan bool, 1)

	signal.Notify(sigC, os.Interrupt)
}

//waits for new conn request
//each conn request is served by goroutine
func (s *Server) waitForConn(l net.Listener) {
	log.Printf("Server: ready to accept connection at port %d", config.TCPPort)

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Server: failed to accept connection. err: %v", err)
		}

		log.Printf("Server: new connection request received from %s", conn.RemoteAddr())
		go s.handleConn(conn)
	}
}

//waits until the conn is killed
//serves all the request made by the clietn at this conn
func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)
	for {
		c.SetReadDeadline(time.Now().Add(config.Timeout))
		req, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("Server: unable to read request from %s. err: %v", c.RemoteAddr(), err)
			return
		}

		resp, exit := s.processReq(req)

		c.Write(resp)

		if exit {
			log.Printf("Server: received quit message from %s", c.RemoteAddr())
			return
		}
	}
}

//identifies the command to be executed and runs the command synchronously
func (s *Server) processReq(req string) ([]byte, bool) {
	req = strings.TrimSpace(req)

	command, err := cmd.NewCommand(req)
	if err != nil {
		errMsg := fmt.Sprintf("%v\r\n", err)
		log.Errorf("Server: invalid request. err: %v", err)
		return []byte(errMsg), false
	}

	return cmd.ExecuteCmd(command, s.cache)
}

//periodic cache cleanup to delete expired keyvalues
func (s *Server) cacheCleanup() {
	cleanupTickC := time.Tick(config.CleanupInterval)

	for now := range cleanupTickC {
		s.cache.Cleanup()
		log.Printf("Server: cache cleanup done at %v", now)
	}
}

func (s *Server) Start() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", config.TCPPort))
	if err != nil {
		log.Fatalf("Server: failed to start server. err: ", err)
	}

	defer l.Close()
	s.registerIntSig()

	go s.handleIntSig()
	go s.cacheCleanup()
	go s.waitForConn(l)

	//wait for cleanup to complete
	<-cleanupC
}
