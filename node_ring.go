package qclient

import (
	"sort"
	"fmt"
	"math"
	"crypto/md5"
"errors"
)

type nodeRing struct {
	nodeMap map[int]string
	sortedKeys []int
	virtualCount int
}

func generateKey(key string) int {
	byteKey := md5.Sum([]byte(key))
	return (int(byteKey[3]) << 24) | (int(byteKey[2]) << 16) | (int(byteKey[1]) << 8) | int(byteKey[0])
}

func keysForNode(node string, virtualCount int) []int {
	keys := make([]int, virtualCount)
	for i := 0; i < virtualCount; i++ {
		keys = append(keys, generateKey(fmt.Sprintf("%s-%s", node, i)))
	}

	return keys
}

func newNodeRing(nodes []string) nodeRing {
	ring := nodeRing{nodeMap: make(map[int]string),
			 sortedKeys: make([]int, 0),
			 virtualCount: int(math.Ceil(1000.0/float64(len(nodes))))}
	ring.addNodes(nodes)
	return ring
}

func (ring *nodeRing) addNodes(nodes []string) {
	for _, node := range nodes {
		for _, key := range keysForNode(node, ring.virtualCount) {
			ring.nodeMap[key] = node
			ring.sortedKeys = append(ring.sortedKeys, key)
		}
	}

	sort.Ints(ring.sortedKeys)
}

func (ring *nodeRing) addNode(node string) {
	ring.addNodes([]string{node})
}

func intSet(ints []int) map[int]bool {
	set := make(map[int]bool)
	for _, i := range ints {
		set[i] = true
	}
	return set
}

func (ring *nodeRing) removeNode(node string) {
	nodeKeys := intSet(keysForNode(node, ring.virtualCount))
	newSortedKeys := make([]int, 0)
	for _, key := range ring.sortedKeys {
		if !nodeKeys[key] {
			newSortedKeys = append(newSortedKeys, key)
		} else {
			delete(ring.nodeMap, key)
		}
	}

	ring.sortedKeys = newSortedKeys
}

func (ring *nodeRing)getNode(tableKey string) (string, error) {
	if len(ring.sortedKeys) == 0 {
		return "", errors.New("No nodes available")
	}

	key := generateKey(tableKey)
	pos := sort.Search(len(ring.sortedKeys), func(i int) bool { return ring.sortedKeys[i] >= key })
	pos %= len(ring.sortedKeys)
	return ring.nodeMap[ring.sortedKeys[pos]], nil
}
