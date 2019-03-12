package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"
)

// IsErrorAddressAlreadyInUse ...
func IsErrorAddressAlreadyInUse(err error) bool {
	errOpError, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	errSyscallError, ok := errOpError.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	errErrno, ok := errSyscallError.Err.(syscall.Errno)
	if !ok {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}
	return false
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func newServer(dir string, port int, logging bool) *http.Server {
	handler := fileServerHandlerWithLogging(dir)
	if !logging {
		handler = http.FileServer(http.Dir(dir))
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}

// Listener ...
func Listener(port int) (net.Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}

// Run ...
func Run(dir string, ln net.Listener) {
	srv := newServer(dir, ln.Addr().(*net.TCPAddr).Port, true)
	log.Printf("Serving [%s] on HTTP [%s]\n", dir, srv.Addr)

	defer ln.Close()
	log.Fatal(srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}))
}
