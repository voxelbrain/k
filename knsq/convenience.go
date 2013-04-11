package knsq

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/bitly/nsq/nsq"
)

type Client interface {
	// CreateTopic makes sure a topic exists on the given nsqd.
	CreateTopic(topic string) error
	// Send sends a message on a topic to the specified nsqd.
	Send(topic string, body io.Reader) error
}

// HttpClient implements the Client interface using HTTP requests.
type HttpClient struct {
	*url.URL
}

func (h HttpClient) CreateTopic(topic string) error {
	resp, err := http.Get(h.String() + "/create_topic?topic=" + topic)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err
}

func (h HttpClient) Send(topic string, body io.Reader) error {
	reqURL := h.String() + "/put?topic=" + topic
	log.Printf("Sending %s message to: %s", topic, reqURL)
	resp, err := http.DefaultClient.Post(reqURL, "application/json", body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	return err
}

// TODO: Does this thing need to be thread-safe?
type rwcClient struct {
	io.ReadWriteCloser
}

// NewTCPClient creates an implementation of the Client interface
// using the TCP API of NSQd.
func NewTCPClient(addr string) (Client, error) {
	if err != nil {
		return nil, err
	}

	_, err = c.Write(nsq.MagicV2)
	return &rwcClient{c}, err
}

func (c *rwcClient) CreateTopic(topic string) error {
	// FIXME: This is awful. NSQd has no command to create a topic
	// except via the HTTP API. Needs patching.
	cmd := nsq.Publish(topic, []byte("dummymsg"))
	return cmd.Write(c)
}

func (c *rwcClient) Send(topic string, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	cmd := nsq.Publish(topic, buf)
	return cmd.Write(c)
}
