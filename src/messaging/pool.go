package messaging

import (
	"net"
)

type Pool struct {
	conns map[string]net.Conn
}
