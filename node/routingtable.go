package node

import (
	"sync"
	"time"
)

const (
	a = 3   //Concurrency parameter
	b = 160 //Key size in bits
	k = 20  //Bucket length
)

type routingTableState struct {
	routingTable  [][]*Node
	bucketRefresh [b]time.Time
	mutex         *sync.Mutex
}

func (r *routingTableState) Init() error {
	return nil
}
