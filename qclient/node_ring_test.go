package qclient

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func keyDistribution(ring *nodeRing, keyCount int) map[string]map[string]bool {
	nodeToKeys := make(map[string]map[string]bool)
	for i := 0; i < keyCount; i++ {
		iStr := strconv.Itoa(i)
		node, _ := ring.getNode(iStr)
		if nodeToKeys[node] == nil {
			nodeToKeys[node] = make(map[string]bool)
		}
		nodeToKeys[node][iStr] = true
	}

	return nodeToKeys
}

func isSubset(a, b map[string]bool) bool {
	result := true
	for k := range a {
		if !b[k] {
			result = false
		}
	}

	return result
}

func TestDistribution(t *testing.T) {
	nodes := []string{"aaa", "bbb", "ccc", "ddd", "eee"}
	ring, _ := newNodeRing(nodes)

	// Initial distribution should be fairly even
	initialDistribution := keyDistribution(ring, 60000)
	for _, node := range nodes {
		length := len(initialDistribution[node])
		assert.True(t, 10500 < length && length < 13000)
	}

	// Removing one node should redistribute keys from the failed node
	// but all keys on the nodes that remain should stay still
	ring.removeNode("aaa")
	reducedDistribution := keyDistribution(ring, 60000)
	for _, node := range []string{"bbb", "ccc", "ddd", "eee"} {
		length := len(reducedDistribution[node])
		assert.True(t, 13000 < length && length < 16000)
		assert.True(t, isSubset(initialDistribution[node], reducedDistribution[node]))
	}

	// Re-adding the node should result in the initial distribution
	ring.addNode("aaa")
	revampedDistribution := keyDistribution(ring, 60000)
	assert.Equal(t, revampedDistribution, initialDistribution)
}

func TestNotPossibleToCreateRingWithoutNodes(t *testing.T) {
	ring, err := newNodeRing([]string{})
	assert.Nil(t, ring)
	assert.NotNil(t, err)
}

func TestNilNodeReturnedWhenNoNodesExist(t *testing.T) {
	ring, _ := newNodeRing([]string{"aaa"})
	ring.removeNode("aaa")
	node, err := ring.getNode("xyz")
	assert.Equal(t, node, "")
	assert.NotNil(t, err)
}
