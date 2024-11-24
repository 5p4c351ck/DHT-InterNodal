package routingtable

import (
	"sync"
	"time"

	"github.com/5p4c351ck/DHT-InterNodal/node"
)

const (
	a = 3   //Concurrency parameter
	b = 160 //Key size in bits
	k = 20  //Bucket length
)

type routingTableState struct {
	routingTable  [][]*node.Node
	bucketRefresh [b]time.Time
	mutex         *sync.Mutex
}
