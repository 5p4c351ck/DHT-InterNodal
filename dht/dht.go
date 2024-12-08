package dht

import (
	"github.com/5p4c351ck/DHT-InterNodal/node"
)

type DHT interface {
	Ping(n *node.Node) bool
	Store(cid [20]byte, value []byte) error
	FindNode(cid [20]byte) ([]node.Node, error)
	FindValue(cid [20]byte) ([]byte, error)
}

type DHTimpl struct {
	*node.LocalNode
}

func NewDHT(ln *node.LocalNode) (DHT, error) {
	dht := &DHTimpl{
		LocalNode: ln,
	}
	return dht, nil
}
