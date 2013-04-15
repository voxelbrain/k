package knsq

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"

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
	read   chan string
	write  chan *nsq.Command
	manage chan *manageRequest
}

// NewTCPClient creates an implementation of the Client interface
// using the TCP API of NSQd.
func NewTCPClient(addr string) (Client, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	_, err = c.Write(nsq.MagicV2)
	client := &rwcClient{ReadWriteCloser: c, read: make(chan string), write: make(chan *nsq.Command), manage: make(chan *manageRequest)}

	go client.readLoop()
	go client.writeLoop()
	go client.manageLoop()

	return client, err
}

func (c *rwcClient) readLoop() {
	for {
		resp, err := nsq.ReadResponse(c)
		if err != nil {
			log.Printf("ReadResponse returned: %v", err)
			break
		}
		_, resp, err = nsq.UnpackResponse(resp)
		if bytes.Equal(resp, []byte("_heartbeat_")) {
			c.write <- nsq.Nop()
			continue
		}
		c.read <- string(resp)
	}
}

func (c *rwcClient) writeLoop() {
	for {
		cmd := <-c.write
		cmd.Write(c)
	}
}

type manageRequest struct {
	Cmd      *nsq.Command
	Response chan string
}

func (c *rwcClient) manageLoop() {
	for {
		req := <-c.manage
		c.write <- req.Cmd
		req.Response <- <-c.read
	}
}

func (c *rwcClient) CreateTopic(topic string) error {
	// FIXME: This is awful. NSQd has no command to create a topic
	// except via the HTTP API. Needs patching.
	cmd := nsq.Publish(topic, []byte("dummymsg"))

	respchan := make(chan string)
	c.manage <- &manageRequest{Cmd: cmd, Response: respchan}
	resp := <-respchan

	if resp != "OK" {
		return errors.New(resp)
	}

	return nil
}

func (c *rwcClient) Send(topic string, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	cmd := nsq.Publish(topic, buf)
	respchan := make(chan string)
	c.manage <- &manageRequest{Cmd: cmd, Response: respchan}
	resp := <-respchan

	if resp != "OK" {
		return errors.New(resp)
	}

	return nil
}
