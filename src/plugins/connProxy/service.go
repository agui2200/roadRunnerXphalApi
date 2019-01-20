package connProxy

import (
	"context"
	"fmt"
	"github.com/agui2200/roadrunner/cmd/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"roadRunnerXPhalApi/src/plugins/connProxy/genericPool"
	"time"
)

const ID = "connProxy"

type Service struct {
	pool      map[string]genericPool.Pool
	listeners map[string]*listener
}

func (s *Service) Init(config *ConnConfig) (ok bool, err error) {
	s.pool = make(map[string]genericPool.Pool)
	s.listeners = make(map[string]*listener)
	for n, c := range config.Conn {
		p, err := genericPool.NewGenericPool(0, int(c.PoolNum), 10*time.Second, func() (poolable genericPool.Poolable, e error) {
			c, e := newConn(c.Protocol, fmt.Sprintf("%s:%d", c.Host, c.Port))
			return genericPool.Poolable(c), nil
		})
		if err != nil {
			return false, err
		}
		s.pool[n] = p
		ln, err := net.Listen(c.Protocol, c.ListenAddr)
		if err != nil {
			return false, err
		}
		s.listeners[n] = &listener{
			Listener: ln,
			closed:   make(chan struct{}),
		}
	}
	return true, nil
}

func (s *Service) Serve() error {
	if s.pool == nil {
		return errors.New("connect pool not configured")
	}
	ctx := context.Background()
	for n, ln := range s.listeners {
		ctx, cancel := context.WithCancel(ctx)
		ln.cancel = cancel
		go func() {
			util.Printf("<green>on %s start</reset>\n", n)
			for {
				select {
				case <-ctx.Done():
					ln.closed <- struct{}{}
					return
				default:
					c, err := ln.Accept()
					if err != nil {
						util.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
						break
					}
					buf, err := ioutil.ReadAll(c)
					if err != nil {
						util.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
						break
					}
					oc, err := s.pool[n].Acquire()
					if err != nil {
						util.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
						break
					}
					_, err = oc.Write(buf)
					if err != nil {
						util.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
						break
					}
				}

			}
		}()
	}

	return nil
}

func (s *Service) Stop() {
	for n, c := range s.pool {
		err := c.Shutdown()
		if err != nil {
			util.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
		}
		ln := s.listeners[n]
		ln.cancel()
		select {
		case <-time.After(3 * time.Second):
			util.Printf("<red+hb>Error:</reset> <red>%s timeout!</reset>\n", n)
		case <-ln.closed:
			util.Printf("<green>on %s stop</reset>\n", n)
		}
	}
}
