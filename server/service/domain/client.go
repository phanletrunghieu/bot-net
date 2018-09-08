package domain

import (
	"net"

	"github.com/satori/go.uuid"
)

// Client struct
type Client struct {
	ID   uuid.UUID
	Addr net.Addr
	Conn net.Conn
}
