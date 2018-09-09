package domain

import (
	"net"

	"github.com/satori/go.uuid"
)

// Boss struct
type Boss struct {
	ID              uuid.UUID
	Addr            net.Addr
	Conn            net.Conn
	IsAuthenticated bool
}
