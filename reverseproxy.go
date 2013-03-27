package k

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewSingleHostReverseProxy is an augmentation of
// `net/http/httputil.NewSingleHostReverseProxy` which additionally sets
// the `Host` header.
func NewSingleHostReverseProxy(url *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(url)
	oldDirector := rp.Director
	rp.Director = func(r *http.Request) {
		oldDirector(r)
		r.Host = url.Host
	}
	return rp
}
