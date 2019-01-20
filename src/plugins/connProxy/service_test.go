package connProxy

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"sync"
	"testing"
)

type testService struct{}

func (ts *testService) Echo(msg string, r *string) error { *r = msg; return nil }

func Test_Disabled(t *testing.T) {
	s := &Service{}
	ok, err := s.Init(&ConnConfig{Enable: false})

	assert.NoError(t, err)
	assert.False(t, ok)
}

func Test_RegisterNotConfigured(t *testing.T) {
	s := &Service{}
	assert.Error(t, s.Serve())
}

func Test_Enabled(t *testing.T) {
	s := &Service{}
	ok, err := s.Init(&ConnConfig{Enable: true, Conn: map[string]struct {
		Host       string
		Port       uint
		Protocol   string
		PoolNum    uint
		ListenAddr string
	}{
		"test": {
			Host:       "127.0.0.1",
			Port:       9001,
			Protocol:   "tcp",
			PoolNum:    30,
			ListenAddr: "127.0.0.1:9001",
		},
	}})
	assert.NoError(t, err)
	assert.True(t, ok)
}

func Test_StopNonServing(t *testing.T) {
	s := &Service{}
	ok, err := s.Init(&ConnConfig{Enable: true, Conn: map[string]struct {
		Host       string
		Port       uint
		Protocol   string
		PoolNum    uint
		ListenAddr string
	}{
		"test": {
			Host:       "127.0.0.1",
			Port:       9001,
			Protocol:   "tcp",
			PoolNum:    30,
			ListenAddr: "127.0.0.1:9001",
		},
	}})
	assert.NoError(t, err)
	assert.True(t, ok)
	s.Stop()
}

func Test_Serve_Errors(t *testing.T) {
	s := &Service{}
	ok, err := s.Init(&ConnConfig{Enable: true, Conn: map[string]struct {
		Host       string
		Port       uint
		Protocol   string
		PoolNum    uint
		ListenAddr string
	}{
		"test": {
			Host:       "127.0.0.1",
			Port:       9001,
			Protocol:   "tcp",
			PoolNum:    30,
			ListenAddr: "127.0.0.1:9001",
		},
	}})
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Error(t, s.Serve())
	s.Stop()
}

func Test_Serve_Client(t *testing.T) {
	buf := []byte("hello word!")
	s := &Service{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ln, err := net.Listen("tcp", "127.0.0.1:9000")
		assert.NoError(t, err)
		for {
			select {
			case <-ctx.Done():
				wg.Done()
			}
			conn, err := ln.Accept()
			if err != nil {
				assert.NoError(t, err)
				return
			}
			b, err := ioutil.ReadAll(conn)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			assert.Equal(t, b, len(buf))
		}
	}()
	ok, err := s.Init(&ConnConfig{Enable: true, Conn: map[string]struct {
		Host       string
		Port       uint
		Protocol   string
		PoolNum    uint
		ListenAddr string
	}{
		"test": {
			Host:       "127.0.0.1",
			Port:       9000,
			Protocol:   "tcp",
			PoolNum:    30,
			ListenAddr: "127.0.0.1:9001",
		},
	}})
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.NoError(t, s.Serve())
	client, err := net.Dial("tcp", "127.0.0.1:9001")
	assert.NoError(t, err)
	n, err := client.Write(buf)
	assert.NoError(t, err)
	fmt.Println(len(buf))
	assert.Equal(t, n, len(buf))
	cancel()
	s.Stop()
	wg.Wait()
}
