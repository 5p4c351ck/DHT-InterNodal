package node

import (
	"errors"
	"math/big"
	"sync"
	"time"
)

const (
	a = 3   //Concurrency parameter
	b = 160 //Key size in bits
	k = 20  //Bucket length
)

type RoutingTable struct {
	kbuckets      [][]*Node
	bucketRefresh [b]time.Time
	mutex         *sync.Mutex
	owner         *Node
	ownerCID      *big.Int
}

func NewRoutingTable(node *Node) (*RoutingTable, error) {
	if node == nil {
		return nil, errors.New("node cannot be nil")
	}
	rt := &RoutingTable{
		kbuckets: make([][]*Node, b),
		mutex:    &sync.Mutex{},
		owner:    node,
		ownerCID: new(big.Int).SetBytes(node.Cid[:]),
	}
	return rt, nil
}

func (r *RoutingTable) Dinstance(node *Node) (*big.Int, error) {
	if node == nil {
		return nil, errors.New("node cannot be nil")
	}
	nodeCID := new(big.Int).SetBytes(node.Cid[:])
	dinstance := new(big.Int).Xor(r.ownerCID, nodeCID)
	return dinstance, nil
}
