package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
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
}

// NewTCPService create tcpService struct
func NewTCPService(host string, port int) *Service {
	return &Service{
		host:          host,
		port:          port,
		reconnectTime: time.Second,
		Error:         make(chan error),
		WriteChan:     make(chan string),
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
	s.writeStream()
}

func (s *Service) writeStream() {
	for {
		// get command
		command := <-s.WriteChan
		s.connection.Write([]byte(command + "\r"))
	}
}

func (s *Service) readStream() {
	for {
		buffCommand := make([]byte, 2)
		_, err := s.connection.Read(buffCommand)
		if err != nil {
			s.Error <- err
			return
		}

		// TODO: use buffCommand

		msg, err := bufio.NewReader(s.connection).ReadString('\r')
		if err != nil {
			s.Error <- err
			continue
		}

		// Unmarshal
		var data interface{}
		data = nil
		err = json.Unmarshal([]byte(msg), &data)
		if err != nil {
			fmt.Println(msg)
		} else if reflect.TypeOf(data).String() == "[]interface {}" {
			list := data.([]interface{})
			for _, elem := range list {
				fmt.Println(elem)
			}
		}

	}
}
