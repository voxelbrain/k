package k

import (
	"encoding/json"
	"log"
	"net"
	"net/url"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// FIXME: Make Must* functions consistent

// MustURL calls net/url.Parse() and panics if it returns a non-nil
// error.
func MustURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return u
}

// MustBytes is a function which usually wraps functions which return
// a byte slice and an error. It panics, if the given error is not nil.
func MustBytes(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

// DefaultEnv returns the value of an environment variable, provided
// it has been set. If it is unset (i.e. empty), the specified default
// value is returned.
func DefaultEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

// MustTCPAddr calls net.ResolveTCPAddr and panics if it returns
// a non-nil error.
func MustTCPAddr(rawaddr string) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", rawaddr)
	if err != nil {
		panic(err)
	}

	return addr
}

// Hostname returns def if no hostname could be determined
// for this machine, the machine's hostname otherwise.
func Hostname(def string) string {
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return def
}

// JsonRemarshal takes old, marshals it into json and unmarshals it
// into new. If an error occurs along the way, it is returned.
func JsonRemarshal(new interface{}, old interface{}) error {
	data, err := json.Marshal(old)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, new)
}
