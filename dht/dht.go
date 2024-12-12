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

func NewDHT(ip string, port string) (DHT, error) {
	n, err := node.NewNode(ip, port)
	if err != nil {
		return nil, err
	}
	ln, err := node.NewLocalNode(n)
	if err != nil {
		return nil, err
	}
	serverErrChan := make(chan error, 1)
	serverRpcChan := make(chan int, 1)
	go func() {
		ln.Server(serverErrChan, serverRpcChan)
	}()
	if err = <-serverErrChan; err != nil {
		return nil, err
	}
	dht := &DHTimpl{
		LocalNode: ln,
	}
	return dht, nil
}
