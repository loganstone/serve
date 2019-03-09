package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

const defaultDir = "."
const defaultPort = 9000

var dirToServe string
var portToListen int

type hasStatusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *hasStatusCodeResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func isErrorAddressAlreadyInUse(err error) bool {
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

func fileServerHandlerWithLogging() http.Handler {
	fileServerHandler := http.FileServer(http.Dir(dirToServe))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-> %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

		hrw := &hasStatusCodeResponseWriter{w, http.StatusOK}
		fileServerHandler.ServeHTTP(hrw, r)

		log.Printf("<- %d %s\n", hrw.statusCode, http.StatusText(hrw.statusCode))
	})
}

func newServer() *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", portToListen),
		Handler: fileServerHandlerWithLogging(),
	}
}

func abs(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		return dir, nil
	}
	return filepath.Abs(dir)
}

func run() {
	server := newServer()
	log.Printf("Serving [%s] on HTTP [%s]\n", dirToServe, server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		if isErrorAddressAlreadyInUse(err) {
			log.Println(err)
			portToListen++
			log.Printf("Change port: [%d]\n", portToListen)
			run()
		}
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&dirToServe, "d", defaultDir, "directory to serve")
	flag.IntVar(&portToListen, "p", defaultPort, "port to listen on")
}

func main() {
	flag.Parse()

	absPath, err := abs(dirToServe)
	if err != nil {
		log.Fatal(err)
	}

	dirInfo, err := os.Stat(absPath)
	if err != nil {
		log.Fatal(err)
	}

	if !dirInfo.IsDir() {
		log.Fatal("-d option value must be directory")
	}

	dirToServe = absPath
	run()
}
