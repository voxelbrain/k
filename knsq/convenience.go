package knsq

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// TODO: Add socket implementation

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
