// Package stringutil contains utility functions for working with strings.
package qclient

// - Basic POST query
// - Basic GET query
// - Consistent hashing implementation equivalient to the python version

import ("encoding/json"
	"fmt"
	"net/url"
	"net/http")

type QClient struct {
	nodes []string
	currentNode string
}

type Query struct {
	Select   []interface{}    `json:"select,omitempty"`
	Where    []interface{}    `json:"where,omitempty"`
	OrderBy  []string       `json:"order_by,omitempty"`
	GroupBy  []string       `json:"group_by,omitempty"`
	Distinct []string       `json:"distinct,omitempty"`
	Offset   int            `json:"offset,omitempty"`
	Limit    int            `json:"limit,omitempty"`
	From     *Query         `json:"from,omitempty"`
}

func (c *QClient) Query(key string, q string) string {
	return "hello"
}


func (c *QClient) Get(key string, q Query) string {
	jq, _ := json.Marshal(q)
	fmt.Printf(string(jq[:]) + "\n")
	ujq := url.QueryEscape(string(jq[:]))
	fmt.Printf("Result: %s\n", ujq)
	http.Get(c.nodes[0] + "/qcache/dataset/" + key + "?q=" + ujq)
	return "hello"
}


func New(nodes []string) *QClient {
	return &QClient{nodes: nodes}
}
