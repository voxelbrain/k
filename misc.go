package k

import (
	"net/url"
	"os"
)

// MustURL calls net/url.Parse() and panics if it returns a non-nil
// error. Useful for URL constants.
func MustURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return u
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
