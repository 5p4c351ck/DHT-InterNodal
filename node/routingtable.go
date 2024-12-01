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
	kBuckets      [][]*Node
	bucketRefresh [b]time.Time
	mutex         *sync.Mutex
	owner         *LocalNode
}

func NewRoutingTable(localnode *LocalNode) (*RoutingTable, error) {
	if localnode == nil {
		return nil, errors.New("node cannot be nil")
	}
	rt := &RoutingTable{
		kBuckets: make([][]*Node, b),
		mutex:    &sync.Mutex{},
		owner:    localnode,
	}
	return rt, nil
}

func (r *RoutingTable) kBucketIndex(node *Node) (int, error) {
	dinstance, err := r.Dinstance(node)
	if err != nil {
		return 0, err
	}
	bitLen := dinstance.BitLen()
	return bitLen - 1, nil
}

func (r *RoutingTable) Dinstance(node *Node) (*big.Int, error) {
	if node == nil {
		return nil, errors.New("node cannot be nil")
	}

	cidlength := len(r.owner.Cid)
	result := make([]byte, cidlength)

	for i := 0; i < cidlength; i++ {
		result[i] = r.owner.Cid[i] ^ node.Cid[i]
	}

	dinstance := new(big.Int).SetBytes(result)
	return dinstance, nil
}

func (r *RoutingTable) InsertNode(node *Node) (int, error) {
	index, err := r.kBucketIndex(node)
	if err != nil {
		return -1, err
	}

	if index >= len(r.kBuckets) {
		return -1, errors.New("index out of range")
	}
	kBucket := r.kBuckets[index]

	//Check if node is already in the k-bucket
	for i, n := range kBucket {
		dinstance, err := Dinstance(n, node)
		if err != nil {
			return -1, err
		}
		if dinstance.Sign() == 0 {
			kBucket = append(kBucket[:i], kBucket[i+1:]...) //Remove the node from the k-bucket
			kBucket = append(kBucket, node)                 //Append the node to the tail of the k-bucket
			r.kBuckets[index] = kBucket
			return index, nil
		}
	}
	//Check if k-bucket is not full
	if len(kBucket) < k {
		r.kBuckets[index] = append(kBucket, node) //Append the node to the tail of the k-bucket
		return index, nil
	}
	//if the k-bucket is full we ping the least recently seen node
	lrsNode := kBucket[0]
	online := r.owner.Ping(lrsNode)
	if online {
		return -1, nil
	} else {
		r.kBuckets[index] = append(kBucket[1:], node)
		return index, nil
	}
}
