package node

import (
	"fmt"
	"net"
)

const (
	protocolUDP = "udp"
)

func (node *LocalNode) Server(ErrChan chan error) {
	address := fmt.Sprintf("%s:%d", node.IP.String(), node.Port)
	conn, err := net.ListenPacket(protocolUDP, address)
	if err != nil {
		ErrChan <- err
		return
	}
	defer conn.Close()
	ErrChan <- nil
	bufferSize := 1024
	maxConns := 10
	initMessages()
	//Implement a semaphore using a buffered channel
	var semaphore = make(chan struct{}, maxConns)

	buffer := make([]byte, bufferSize)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			continue
		}
		//Delay new connections until the semaphore's buffer is not full
		semaphore <- struct{}{}

		stream := make([]byte, n)
		copy(stream, buffer[:])

		go func(stream []byte, address net.Addr, connection net.PacketConn) {
			defer func() { <-semaphore }()
			codec := NewCodec()
			msg, err := codec.Deserialize(stream)
			if err != nil {
				return //Add logging
			}
			replyMsg, err := node.ConstructReply(msg)
			if err != nil {
				return //Add logging
			}
			replyStream, err := codec.Serialize(replyMsg)
			if err != nil {
				return //Add logging
			}
			err = node.Send(replyStream, address, connection)
			if err != nil {
				return //Add logging
			}
		}(stream, addr, conn)
	}
}

func (node *LocalNode) Send(stream []byte, address net.Addr, connection net.PacketConn) error {
	//Handle any extra logic
	connection.WriteTo(stream, address)
}

func (node *LocalNode) ConstructReply(msg *message) (*message, error) {
	replyMessage := &message{}
	var reply interface{}
	var err error

	switch msg.MessageType {
	case messagePing:
		reply = "PONG"
	case messageStore:
		reply, err = storeReplyMsg(msg)
	case messageFindNode:
		reply, err = findNodeReplyMsg(msg)
	case messageFindValue:
		reply, err = findValueReplyMsg(msg)
	default:
		return nil, fmt.Errorf("Invalid type in message with transaction ID %d", msg.TransactionID)
	}
	if err != nil {
		return nil, err
	}
	replyMessage.Data = reply
	return replyMessage, nil
}
