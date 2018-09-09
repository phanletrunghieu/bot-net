package tcp

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"time"
)

// Service struct
type Service struct {
	host          string
	port          int
	connection    net.Conn
	reconnectTime time.Duration
	Error         chan error
	WriteChan     chan string
	ReadChan      chan string
}

// NewTCPService create tcpService struct
func NewTCPService(host string, port int) *Service {
	return &Service{
		host:          host,
		port:          port,
		reconnectTime: time.Second,
		Error:         make(chan error),
		WriteChan:     make(chan string),
		ReadChan:      make(chan string),
	}
}

// Run wait for a connection
func (s *Service) Run() {
	conn, err := net.Dial("tcp", s.host+":"+strconv.Itoa(s.port))
	if err != nil {
		s.Error <- err
		return
	}

	s.connection = conn
	defer s.connection.Close()

	go s.readStream()

	for {
		// get command
		command := <-s.WriteChan
		log.Println("send cmd:", command)
		s.connection.Write([]byte(command + "\r"))
	}
}

func (s *Service) readStream() {
	for {
		msg, err := bufio.NewReader(s.connection).ReadString('\r')
		if err != nil {
			s.Error <- err
			continue
		}
		s.ReadChan <- msg
	}
}
