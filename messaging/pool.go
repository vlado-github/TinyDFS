package messaging

import (
	"net"
)

// Pool is a register of all tcp/ip network connections.
type Pool struct {
	conns map[string]net.Conn
}
