package node

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
)

type Peer interface {
	Ping(cid [32]byte) error
	Store()
	FindNode(cid [32]byte) ([][32]byte, []net.UDPAddr)
	FindValue()
}

type Node struct {
	cid     [32]byte
	pubKey  ed25519.PublicKey
	udpaddr net.UDPAddr
}

func (node *Node) Ping(cid [32]byte) error {
	return nil
}

func (node *Node) Store(cid [32]byte) error {
	return nil

}

func (node *Node) FindNode(cid [32]byte) ([][32]byte, []net.UDPAddr) {
	return nil, nil

}

func (node *Node) FindValue(cid [32]byte) error {
	return nil

}

func (node *Node) Cid() [32]byte {
	return node.cid

}

func CreateNode() (*Node, error) {
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

	n := &Node{cid: pubsum, pubKey: pub}
	return n, nil
}
