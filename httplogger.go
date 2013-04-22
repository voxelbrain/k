package k

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

type HTTPLoggerFunc func(logLine string)

type httpLoggingHandler struct {
	h       http.Handler
	logfunc HTTPLoggerFunc
}

type logResponseWriter struct {
	http.ResponseWriter
	RespCode int
	Size     int
}

// HTTPLogger wraps an http.Handler and calls a function logFunc with 
// a line that can be written to an HTTP access log when a request 
// to the handler has been completed.
func HTTPLogger(h http.Handler, logFunc HTTPLoggerFunc) http.Handler {
	return &httpLoggingHandler{h: h, logfunc: logFunc}
}

// ServeHTTP implements the http.Handler interface.
func (h *httpLoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := &logResponseWriter{ResponseWriter: w}
	t := time.Now()
	h.h.ServeHTTP(lrw, r)
	duration := time.Since(t).String()
	if lrw.RespCode == 0 {
		lrw.RespCode = 200
	}
	if h.logfunc != nil {
		h.logfunc(fmt.Sprintf("%s \"%s %s %s\" %d %d (%s)", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, lrw.RespCode, lrw.Size, duration))
	}
}

func (w *logResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *logResponseWriter) Write(data []byte) (s int, err error) {
	s, err = w.ResponseWriter.Write(data)
	w.Size += s
	return
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		panic("w.ResponseWriter is not a http.Hijacker")
	}
	return hj.Hijack()
}

func (w *logResponseWriter) WriteHeader(r int) {
	w.ResponseWriter.WriteHeader(r)
	w.RespCode = r
}
