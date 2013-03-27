package k

import (
	"net/http"
	"testing"
)

func TestSingleHostReverseProxy(t *testing.T) {
	r, err := http.NewRequest("GET", "http://my.domain/my/path", nil)
	if err != nil {
		t.Fatalf("Could not parse request URL: %s", err)
	}

	u := MustURL("http://their.domain/their/path")
	rp := NewSingleHostReverseProxy(u)
	rp.Director(r)
	got_uri := r.URL.String()
	expected_uri := "http://their.domain/their/path/my/path"
	if r.URL.String() != expected_uri {
		t.Fatalf("Unexpected request URI. Expected %s, got %s", expected_uri, got_uri)
	}

	if r.Host != "their.domain" {
		t.Fatalf("Unexpected hostname. Expected %s, got %s", "their.domain", r.Host)
	}
}
