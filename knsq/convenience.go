package knsq

import (
	"io"
	"log"
	"net/http"
)

// FIXME: Use socket API instead of HTTP

// CreateTopic makes sure a topic exists on the given nsqd.
func CreateTopic(nsqd string, topic string) error {
	resp, err := http.Get("http://" + nsqd + "/create_topic?topic=" + topic)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err
}

// Send sends a message on a topic to the specified nsqd.
func Send(nsqd string, topic string, body io.Reader) error {
	reqURL := nsqd + "/put?topic=" + topic
	log.Printf("Sending %s message to: %s", topic, reqURL)
	resp, err := http.DefaultClient.Post(reqURL, "application/json", body)
	defer resp.Body.Close()
	return err
}
