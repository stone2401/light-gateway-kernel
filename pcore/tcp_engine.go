package pcore

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
)

type TcpEngine struct {
	listener net.Listener
	balance  sdk.Balance
	network  string
	addr     string

	readTimeout   time.Duration
	writeTimeout  time.Duration
	keepAliveTime time.Duration
}
type conn struct {
	net.Conn
	network string
	addr    string
}

type WithHandler func(time.Duration) func(*conn)

func NewConn(network, addr string) *conn {
	return &conn{
		network: network,
		addr:    addr,
	}
}

func WithWriteTimeout(d time.Duration) func(*conn) {
	return func(c *conn) {
		if d == 0 {
			return
		}
		c.SetWriteDeadline(time.Now().Add(d))
	}
}

func WithReadTimeout(d time.Duration) func(*conn) {
	return func(c *conn) {
		if d == 0 {
			return
		}
		c.SetReadDeadline(time.Now().Add(d))
	}
}

func WithKeepAlivePeriod(d time.Duration) func(*conn) {
	return func(c *conn) {
		if d == 0 {
			return
		}
		if tcpConn, ok := c.Conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
}

func (t *conn) Serve(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Zlog().Error(err.(error).Error())
		}
		conn.Close()
	}()
	nodeConn, err := net.Dial(t.network, t.addr)
	if err != nil {
		zlog.Zlog().Error(err.Error())
		return
	}
	defer func() {
		if err := recover(); err != nil {
			zlog.Zlog().Error(err.(error).Error())
		}
		nodeConn.Close()
	}()
	errch := make(chan error, 1)
	go t.Copy(errch, nodeConn, conn)
	go t.Copy(errch, conn, nodeConn)
	<-errch
}

func (t *conn) Copy(errch chan error, nodeConn, conn net.Conn) {
	_, err := io.Copy(nodeConn, conn)
	errch <- err
}

func NewTcpEngine(b sdk.Balance, readTimeout, writeTimeout, keepAliveTime time.Duration) *TcpEngine {
	return &TcpEngine{
		balance: b,
	}
}

func (t *TcpEngine) ListenAndServe(network, addr string) error {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	t.listener = lis
	t.network = network
	t.addr = addr
	return t.ServeTcp(lis)
}

func (t *TcpEngine) ServeTcp(lis net.Listener) error {
	for {
		rw, err := lis.Accept()
		if err != nil {
			return err
		}
		node, err := t.balance.GetNode(rw.RemoteAddr().String())
		if err != nil {
			zlog.Zlog().Error(err.Error())
			continue
		}
		conn := NewConn(t.network, node)
		WithWriteTimeout(t.writeTimeout)(conn)
		WithReadTimeout(t.readTimeout)(conn)
		WithKeepAlivePeriod(t.keepAliveTime)(conn)
		zlog.Zlog().Info(rw.RemoteAddr().String())
		zlog.Zlog().Info(node)
		go conn.Serve(context.Background(), rw)
	}
}

func (t *TcpEngine) Close() error {
	return t.listener.Close()
}
