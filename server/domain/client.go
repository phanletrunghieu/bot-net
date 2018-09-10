package domain

import (
	"net"

	"github.com/satori/go.uuid"
)

// Client struct
type Client struct {
	ID   uuid.UUID `json:"id"`
	Addr net.Addr  `json:"address"`
	Conn net.Conn  `json:"-"`
}
