package main

import (
	"context"
	"net"

	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
)

type TcpSever struct {
	listener net.Listener
}

func (*TcpSever) ServeTcp(ctx context.Context, src net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Zlog().Error(err.(error).Error())
		}
		src.Close()
	}()
	zlog.Zlog().Info(src.LocalAddr().String())
	zlog.Zlog().Info(src.LocalAddr().Network())
	zlog.Zlog().Info(src.RemoteAddr().String())
	zlog.Zlog().Info(src.RemoteAddr().Network())
	for {
		buf := make([]byte, 1024)
		n, err := src.Read(buf)
		if err != nil {
			zlog.Zlog().Error(err.Error())
			return
		}
		zlog.Zlog().Info(string(buf[:n]))
		src.Write(buf[:n])
	}
}

func (t *TcpSever) ListenAndServe() error {
	for {
		rw, err := t.listener.Accept()
		if err != nil {
			return err
		}
		go t.ServeTcp(context.Background(), rw)
	}
}

func main() {
	zlog.Zlog().Info("start tcp server")
	s := &TcpSever{}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		zlog.Zlog().Error(err.Error())
		return
	}
	s.listener = ln
	err = s.ListenAndServe()
	if err != nil {
		zlog.Zlog().Error(err.Error())
		return
	}
}
