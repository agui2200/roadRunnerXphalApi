package connProxy

import (
	"net"
	"time"
)

func newConn(network string, addr string) (*conn, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &conn{
		Conn:           c,
		lastActiveTime: time.Now(),
	}, nil
}

type conn struct {
	net.Conn
	lastActiveTime time.Time
}

func (c *conn) GetActiveTime() time.Time {
	return c.lastActiveTime
}
