package k

import (
	"net/url"
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
