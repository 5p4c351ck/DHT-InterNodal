package node

import (
	"crypto/ed25519"
	"crypto/sha1"
	"fmt"
	"net"
	"os"
	"strconv"
)

// This is the representation of a remote node on the network
type Node struct {
	Cid  [20]byte
	IP   net.IP
	Port int
}

func NewNode(ip string, port string) (*Node, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	n := &Node{
		IP:   net.ParseIP(ip),
		Port: p,
	}
	return n, nil
}

// This is the representation of a local node on the network
type LocalNode struct {
	*Node
	*RoutingTable
	pubKey ed25519.PublicKey
}

func NewLocalNode(node *Node) (*LocalNode, error) {
	rt, err := NewRoutingTable(node)
	if err != nil {
		return nil, err
	}

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	pubsum := sha1.Sum(pub)
	filename := fmt.Sprintf("./%x.key", pubsum)
	err = os.WriteFile(filename, priv, 0600)
	if err != nil {
		return nil, err
	}

	n := &LocalNode{
		Node:         node,
		RoutingTable: rt,
		pubKey:       pub,
	}
	n.Cid = pubsum
	return n, nil
}

func (node *LocalNode) Ping(cid [20]byte) error {
	return nil
}

func (node *LocalNode) Store(cid [20]byte) error {
	return nil

}

func (node *LocalNode) FindNode(cid [20]byte) *Node {
	return nil

}

func (node *LocalNode) FindValue(cid [20]byte) error {
	return nil

}
