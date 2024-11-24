package node

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strconv"
)

// This is the representation of a remote node on the network
type Node struct {
	Cid  [32]byte
	IP   net.IP
	Port int
}

func CreateNode(ip string, port string) (*Node, error) {
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
	pubKey ed25519.PublicKey
}

func CreateLocalNode(node *Node) (*LocalNode, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	pubsum := sha256.Sum256(pub)
	filename := fmt.Sprintf("./%x.key", pubsum)
	err = os.WriteFile(filename, priv, 0600)
	if err != nil {
		return nil, err
	}

	n := &LocalNode{
		Node:   node,
		pubKey: pub,
	}
	n.Cid = pubsum
	return n, nil
}

func (node *LocalNode) Ping(cid [32]byte) error {
	return nil
}

func (node *LocalNode) Store(cid [32]byte) error {
	return nil

}

func (node *LocalNode) FindNode(cid [32]byte) *Node {
	return nil

}

func (node *LocalNode) FindValue(cid [32]byte) error {
	return nil

}
