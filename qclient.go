// Package stringutil contains utility functions for working with strings.
package qclient

// - Basic POST query
// - Basic GET query
// - Consistent hashing implementation equivalient to the python version

import ("encoding/json"
	"net/url"
	"net/http"
	"io/ioutil"
	"io"
	"log")

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


func (c *QClient) Get(key string, q Query) ([]byte, error) {
	jq, _ := json.Marshal(q)
	ujq := url.QueryEscape(string(jq[:]))
	response, err := http.Get(c.nodes[0] + "/qcache/dataset/" + key + "?q=" + ujq)

	if err != nil {
		log.Fatal("Error getting data: ", err)
	}

	if response.StatusCode != 200 {
		return nil, nil
	}

	defer response.Body.Close()
        contents, _ := ioutil.ReadAll(response.Body)
	
	return contents, nil
}


func (c *QClient) Post(key string, bodyType string, body io.Reader) error {
	_, err := http.Post(c.nodes[0] + "/qcache/dataset/" + key, bodyType, body)
	return err
}


func NewClient(nodes []string) *QClient {
	return &QClient{nodes: nodes}
}
