// modified from https://github.com/WiiLink24/wfc-server/blob/main/nas/listener.go
// licensed under GNU AFFERO GENERAL PUBLIC LICENSE Version 3, see LICENSE

package main

import (
	"bufio"
	"io"
	"net"
)

type inConn struct {
	net.Conn
	in io.Reader
}

func (c inConn) Read(b []byte) (int, error) { return c.in.Read(b) }

func (c inConn) Close() error {
	err := c.Conn.Close()
	if closer, ok := c.in.(io.Closer); ok {
		_ = closer.Close()
	}
	return err
}

type sslListener struct {
	net.Listener
}

func (l *sslListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	sslConn := sslConn{
		Conn: conn,
	}

	pr, pw := io.Pipe()
	go func() {
		_, err := io.Copy(pw, bufio.NewReader(&sslConn))
		_ = pw.CloseWithError(err)
	}()
	return &inConn{
		Conn: &sslConn,
		in:   pr,
	}, nil
}
