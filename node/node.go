package node

import (
	"crypto/ed25519"
	"crypto/sha1"
	"errors"
	"math/big"
	"net"
	"strconv"
)

// This is the representation of a remote node on the network
type Node struct {
	Cid  [b / 8]byte
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

// This is the representation of the local node on the network
type LocalNode struct {
	*Node
	*RoutingTable
	Codec
	pubKey ed25519.PublicKey
}

func NewLocalNode(node *Node) (*LocalNode, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	_ = priv
	pubsum := sha1.Sum(pub)
	/* Comment out for testing without creating files
	filename := fmt.Sprintf("./%x.key", pubsum)
	err = os.WriteFile(filename, priv, 0600)
	if err != nil {
		return nil, err
	}
	*/
	n := &LocalNode{
		Node:   node,
		Codec:  NewCodec(),
		pubKey: pub,
	}
	n.Cid = pubsum

	rt, err := NewRoutingTable(n)
	if err != nil {
		return nil, err
	}
	n.RoutingTable = rt
	return n, nil
}

func Dinstance(node1, node2 *Node) (*big.Int, error) {
	if node1 == nil || node2 == nil {
		return nil, errors.New("node cannot be nil")
	}

	cidlength := len(node1.Cid)
	result := make([]byte, cidlength)

	for i := 0; i < cidlength; i++ {
		result[i] = node1.Cid[i] ^ node2.Cid[i]
	}

	dinstance := new(big.Int).SetBytes(result)
	return dinstance, nil
}

func (node *LocalNode) Ping(n *Node) bool {
	msg := &message{
		SenderNode:   node.Node,
		ReceiverNode: n,
		Request:      true,
		Data:         []byte("data"),
	}

	codec := &CodecImp{}
	err := node.Request(msg, codec)
	return err == nil
}

func (node *LocalNode) Store(cid [20]byte, value []byte) error {
	return nil

}

func (node *LocalNode) FindNode(cid [20]byte) ([]Node, error) {
	return nil, nil

}

func (node *LocalNode) FindValue(cid [20]byte) ([]byte, error) {
	return nil, nil

}
