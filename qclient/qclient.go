package qclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type QClient struct {
	ring *nodeRing
}

type Clause []interface{}

type Query struct {
	Select   Clause   `json:"select,omitempty"`
	Where    Clause   `json:"where,omitempty"`
	OrderBy  []string `json:"order_by,omitempty"`
	GroupBy  []string `json:"group_by,omitempty"`
	Distinct []string `json:"distinct,omitempty"`
	Offset   int      `json:"offset,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	From     *Query   `json:"from,omitempty"`
}

type QueryResult struct {
	Content           []byte
	UnslicedResultLen int
}

func And(clauses ...Clause) Clause {
	return clause("&", clauses...)
}

func clause(op string, clauses ...Clause) Clause {
	result := Clause{op}
	for _, clause := range clauses {
		result = append(result, clause)
	}

	return result
}

func Op(op string, operands ...interface{}) Clause {
	result := Clause{op}
	for _, operand := range operands {
		result = append(result, operand)
	}

	return result
}

func Eq(x, y interface{}) Clause {
	return Op("==", x, y)
}

func (c *QClient) Query(key string, q string) string {
	return "hello"
}

func (c *QClient) Get(key string, q Query) (*QueryResult, error) {
	jq, _ := json.Marshal(q)
	ujq := url.QueryEscape(string(jq[:]))
	node, err := c.ring.getNode(key)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(node + "/qcache/dataset/" + key + "?q=" + ujq)
	if err != nil {
		log.Fatal("Error getting data: ", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		// Not found in cache
		return nil, nil
	}

	content, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		// Set error but also include content, it may provide valuable information about the error
		return &QueryResult{Content: content}, errors.New(fmt.Sprintf("Unexpected HTTP response code %s", response.StatusCode))
	}

	unslicedLength, err := strconv.Atoi(response.Header.Get("X-QCache-unsliced-length"))
	return &QueryResult{Content: content, UnslicedResultLen: unslicedLength}, nil
}

func (c *QClient) Post(key string, bodyType string, body io.Reader) error {
	node, err := c.ring.getNode(key)
	if err != nil {
		return err
	}
	_, err = http.Post(node+"/qcache/dataset/"+key, bodyType, body)
	return err
}

func NewClient(nodes []string) *QClient {
	ring, _ := newNodeRing(nodes)
	return &QClient{ring: ring}
}
