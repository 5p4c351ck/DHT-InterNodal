package node

const (
	messagePing = iota
	messageStore
	messageFindNode
	messageFindValue
)

type message struct {
	messageType   int
	transactionID int64
	senderNode    *Node
	receiverNode  *Node
	request       bool
	data          interface{}
}
