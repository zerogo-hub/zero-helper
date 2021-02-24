package main

import (
	"fmt"
	"net/http"
	"os"

	graceful "github.com/zerogo-hub/zero-helper/graceful/http"
	"github.com/zerogo-hub/zero-helper/logger"
)

// defaultServer 一个简单的 http 服务器
type defaultServer struct {
	// 使用 graceful.Server 替代 httpServer
	server graceful.Server
}

// ServeHTTP 实现 http.Handler 接口
func (ds *defaultServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// 这里进行逻辑处理，比如按照路由进行处理

	pid := os.Getpid()
	message := fmt.Sprintf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
	res.Write([]byte(message))
}

func main() {
	ds := new(defaultServer)
	logger := logger.NewSampleLogger()

	ds.server = graceful.NewServer(ds, logger)

	addr := "127.0.0.1:8877"
	logger.Infof("listen on: http://%s, pid: %d", addr, os.Getpid())

	// 监听信号
	ds.server.ListenSignal()

	// 启动服务，接受连接
	err := ds.server.ListenAndServe(addr)
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("server closed")
		} else {
			logger.Error(err.Error())
		}
	}
}
