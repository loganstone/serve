package server

import (
	"log"
	"net/http"
)

type hasStatusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *hasStatusCodeResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func fileServerHandlerWithLogging(dir string) http.Handler {
	fileServerHandler := http.FileServer(http.Dir(dir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-> %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

		hrw := &hasStatusCodeResponseWriter{w, http.StatusOK}
		fileServerHandler.ServeHTTP(hrw, r)

		log.Printf("<- %d %s\n",
			hrw.statusCode, http.StatusText(hrw.statusCode))
	})
}
