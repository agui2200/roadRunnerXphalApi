package connProxy

import (
	"context"
	"net"
)

type listener struct {
	net.Listener
	cancel context.CancelFunc
	closed chan struct{}
}
