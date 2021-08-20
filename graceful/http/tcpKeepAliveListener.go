package gracefulhttp

import (
	"net"

	zerotime "github.com/zerogo-hub/zero-helper/time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err = tc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err = tc.SetKeepAlivePeriod(zerotime.Minute(3)); err != nil {
		return nil, err
	}

	return tc, nil
}
