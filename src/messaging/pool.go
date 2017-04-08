package messaging

import "net"

type Pool struct {
	conns []net.Conn
}
