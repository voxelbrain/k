package k

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPLogger(t *testing.T) {
	testData := []struct {
		StatusCode int
		TestData []byte
		Method string
		URI string
	}{
		{ 200, []byte("hello world"), "GET", "/foobar" },
		{ 304, []byte("this is just some test data, ignore"), "GET", "/q/x?foo=bar" },
		{ 404, []byte("Not found"), "GET", "/favicon.ico" },
		{ 500, []byte("Server Error"), "POST", "/cgi-bin/vintage.pl" },
	}

	for _, data := range testData {
		logLine := ""
		httpLoggerHandler := HTTPLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(data.StatusCode)
			w.Write(data.TestData)
		}), func(l string) {
			logLine = l
		})
		req, _ := http.NewRequest(data.Method, data.URI, nil)
		resp := httptest.NewRecorder()
		httpLoggerHandler.ServeHTTP(resp, req)
		got := strings.Join(strings.Split(logLine, " ")[1:6], " ")
		expected := fmt.Sprintf("\"%s %s %s\" %d %d", req.Method, req.RequestURI, req.Proto, data.StatusCode, len(data.TestData))
		if got != expected {
			t.Errorf("expected %s, got %s", expected, got)
		}
	}
}
