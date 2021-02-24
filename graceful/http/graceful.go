package gracefulhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/zerogo-hub/zero-helper/logger"
	"github.com/zerogo-hub/zero-helper/time"
)

// Server 用来替代 http.Server
type Server interface {

	// ListenAndServe 用于替代 `http.Server.ListenAndServe`
	ListenAndServe(addr string) error

	// ListenAndServeTLS 用于替代 `http.Server.ListenAndServeTLS`
	ListenAndServeTLS(addr, certFile, keyFile string) error

	// Close 直接关闭服务器
	Close()

	// Shutdown 优雅关闭服务器
	// 关闭监听
	// 执行之前注册的关闭函数(RegisterShutdownHandler)，可以用于清理资源等
	// 关闭空闲连接，等待激活的连接变为空闲，再关闭它
	Shutdown()

	// Restart 重启服务
	Restart()

	// SetShutdownTimeout 设置优雅退出超时时间
	// 服务器会每隔500毫秒检查一次连接是否都断开处理完毕
	// 如果超过超时时间，就不再检查，直接退出
	// 如果要单独给指定的服务器设置 超时时间，可以使用 WithTimeout
	//
	// ms: 单位：毫秒，当 <= 0 时无效，直接退出
	SetShutdownTimeout(ms int)

	// RegisterShutdownHandler 注册关闭函数
	// 按照注册的顺序调用这些函数
	// 所有已经添加的服务器都会响应这个函数
	RegisterShutdownHandler(f func())

	// ListenSignal 监听信号
	ListenSignal()
}

type server struct {
	httpServer *http.Server

	// tc 用于获取监控套接字文件
	tc tcpKeepAliveListener

	// shutdownTimeout 退出时的超时时间，单位: 秒
	shutdownTimeout int

	// log 日志
	logger logger.Logger
}

var (
	// defaultShutdownTimeout 默认关闭等待时间
	defaultShutdownTimeout = 5

	envNewKey = "ZERO_HELPER_GRACEFUL"
)

// NewServer 生成服务器，用来替代 http.Server
func NewServer(handler http.Handler, logger logger.Logger) Server {
	return &server{
		httpServer:      &http.Server{Handler: handler},
		shutdownTimeout: defaultShutdownTimeout,
		logger:          logger,
	}
}

// listener 创建监听套接字
func (s *server) listener(addr string) (ln net.Listener, err error) {
	if s.isChild() {
		fp := os.NewFile(3, "")
		defer fp.Close()
		ln, err = net.FileListener(fp)
	} else {
		ln, err = net.Listen("tcp", addr)
	}
	return
}

func (s *server) isChild() bool {
	_, ok := syscall.Getenv(envNewKey)
	return ok
}

// listenFile 拷贝当前的监听套接字文件
func (s *server) listenFile() (*os.File, error) {
	file, err := s.tc.File()
	if err != nil {
		return nil, nil
	}

	return file, nil
}

// ListenAndServe 用于替代 `http.Server.ListenAndServe`
func (s *server) ListenAndServe(addr string) error {
	if addr == "" {
		addr = ":http"
	}

	ln, err := s.listener(addr)
	if err != nil {
		return err
	}

	tc := tcpKeepAliveListener{ln.(*net.TCPListener)}

	s.tc = tc
	s.httpServer.Addr = addr

	return s.httpServer.Serve(s.tc)
}

// ListenAndServeTLS 用于替代 `http.Server.ListenAndServeTLS`
func (s *server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	if addr == "" {
		addr = ":https"
	}

	ln, err := s.listener(addr)
	if err != nil {
		return err
	}
	tc := tcpKeepAliveListener{ln.(*net.TCPListener)}

	s.tc = tc
	s.httpServer.Addr = addr

	defer s.tc.Close()

	return s.httpServer.ServeTLS(s.tc, certFile, keyFile)
}

// Close 直接关闭服务器
func (s *server) Close() {
	err := s.httpServer.Close()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Errorf("server close, err: %s", err.Error())
	} else {
		s.logger.Info("server exiting")
	}
}

// Shutdown 优雅关闭服务器
// 关闭监听
// 执行之前注册的关闭函数(RegisterShutdownHandler)，可以用于清理资源等
// 关闭空闲连接，等待激活的连接变为空闲，再关闭它
func (s *server) Shutdown() {
	logger := s.logger
	timeout := s.shutdownTimeout

	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second(timeout))
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("server shutdown, err: %s", err.Error())
		} else {
			logger.Info("server shutdown")
		}

		select {
		case <-ctx.Done():
			logger.Infof("server timeout of %d seconds", timeout)
		}
	} else {
		ctx := context.TODO()
		err := s.httpServer.Shutdown(ctx)

		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("server shutdown, err: %s", err.Error())
		} else {
			logger.Info("server shutdown")
		}
	}
}

// Restart 重启服务
func (s *server) Restart() {
	logger := s.logger

	dir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("get dir failed: %s", err.Error())
	}

	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	listenFile, err := s.listenFile()
	if err != nil {
		logger.Fatalf("get listenFile failed: %s", err.Error())
	}

	// listenFile 是复制出来的
	defer listenFile.Close()

	files = append(files, listenFile)

	env := []string{}
	for _, s := range os.Environ() {
		if !strings.HasPrefix(s, envNewKey) {
			env = append(env, s)
		}
	}
	env = append(env, fmt.Sprintf("%s=1", envNewKey))

	// 获取可执行文件路径
	name, err := exec.LookPath(os.Args[0])
	if err != nil {
		logger.Fatalf("%s look path failed: %s", os.Args[0], err.Error())
	}

	s.logger.Infof("bin file: %s", name)

	s.httpServer.SetKeepAlivesEnabled(false)

	process, err := os.StartProcess(name, os.Args, &os.ProcAttr{
		Dir: dir,
		Env: env,
		// 新的进程拥有拷贝了当前的监听套接字
		Files: files,
	})

	if err != nil {
		logger.Fatalf("start new process failed: %s", err.Error())
		return
	}

	logger.Infof("restart success, new pid: %d, ppid: %d", process.Pid)

	// 父进程退出
	s.httpServer.Close()
}

// SetShutdownTimeout 设置优雅退出超时时间
// 服务器会每隔500毫秒检查一次连接是否都断开处理完毕
// 如果超过超时时间，就不再检查，直接退出
// 如果要单独给指定的服务器设置 超时时间，可以使用 WithTimeout
//
// timeout: 单位：毫秒，当 <= 0 时无效，直接退出
func (s *server) SetShutdownTimeout(ms int) {
	s.shutdownTimeout = ms
}

// RegisterShutdownHandler 注册关闭函数
// 按照注册的顺序调用这些函数
// 所有已经添加的服务器都会响应这个函数
func (s *server) RegisterShutdownHandler(f func()) {
	s.httpServer.RegisterOnShutdown(f)
}

// ListenSignal 监听信号
func (s *server) ListenSignal() {
	go s.waitSignal()
}

func (s *server) waitSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	sig := <-ch
	signal.Stop(ch)

	s.logger.Infof("received signal, sig: %+v", sig)

	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		s.logger.Info("close signal .. shutdown server ..")
		s.Shutdown()
	case syscall.SIGUSR2:
		s.logger.Info("restart signal .. restart server ..")
		s.Restart()
	default:
		s.logger.Errorf("unsupport signal: %s", sig.String())
	}
}
