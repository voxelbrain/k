package knsq

import (
	"github.com/bitly/nsq/nsq"
)

// MustReader calls nsq.NewReader and panics if nsq.NewReader
// returned an error. nsq.NewReader only fails on invalid
// topic and channel names.
func MustReader(topic, channel string) *nsq.Reader {
	r, err := nsq.NewReader(topic, channel)
	if err != nil {
		panic(err)
	}
	return r
}

// HandlerFunc is a type that makes a single function
// implement to the nsq.Handler interface which consists
// of only one function (just like http.HandlerFunc and http.Handler).
type HandlerFunc func(message *nsq.Message) error

func (f HandlerFunc) HandleMessage(message *nsq.Message) error {
	return f(message)
}

// AttachHandler creates a new nsq.Reader for a topic and attaches the
// given handler to it.
func AttachHandler(topic, channel string, lookupd string, handler nsq.Handler) error {
	mountReader, err := nsq.NewReader(topic, channel)
	if err != nil {
		return err
	}
	mountReader.AddHandler(handler)
	return mountReader.ConnectToLookupd(lookupd)
}
