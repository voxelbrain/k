package knsq

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// Producer represents a running NSQd instance registered with a
// nsqlookupd.
type Producer struct {
	Hostname         string   `json:"hostname"`
	BroadcastAddress string   `json:"broadcast_address"`
	TcpPort          int      `json:"tcp_port"`
	HttpPort         int      `json:"http_port"`
	Version          string   `json:"version"`
	Topics           []string `json:"topics"`
}

// Returns the TCP address of the producer in the form of hostname:port.
func (p Producer) TCPAddr() string {
	return fmt.Sprintf("%s:%d", p.Hostname, p.TcpPort)
}

// Returns the HTTP address of the producer in the form of hostname:port.
func (p Producer) HTTPAddr() string {
	return fmt.Sprintf("%s:%d", p.Hostname, p.HttpPort)
}

type lookupdNodes struct {
	StatusCode int    `json:"status_code"`
	StatusTxt  string `json:"status_txt"`
	Data       struct {
		Producers []Producer `json:"producers"`
	} `json:"data"`
}

// Lookupd represents a running nsqlookupd instance. The underlying
// string should be the HTTP address in the form of hostname:port.
type Lookupd string

// NSQdInstnaces returns all producers registered with the nsqlookupd.
func (lookupd Lookupd) NSQdInstances() ([]Producer, error) {
	resp, err := http.Get("http://" + string(lookupd) + "/nodes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	nodes := &lookupdNodes{}
	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&nodes)
	if err != nil {
		return nil, err
	}

	if nodes.StatusCode != 200 {
		return nil, fmt.Errorf("lookupd returned %d (%s)", nodes.StatusCode, nodes.StatusTxt)
	}

	nsqds := nodes.Data.Producers
	if len(nsqds) == 0 {
		return nil, fmt.Errorf("no nsqd instances found at lookupd %s", lookupd)
	}

	return nsqds, nil
}

// NSQdInstnaces returns a random producer registered with the nsqlookupd.
func (lookupd Lookupd) NSQdInstance() (Producer, error) {
	nsqds, err := lookupd.NSQdInstances()
	if err != nil {
		return Producer{}, err
	}
	perm := rand.Perm(len(nsqds))
	return nsqds[perm[0]], nil
}
