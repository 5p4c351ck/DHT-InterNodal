package node

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	protocolUDP = "udp"
)

func (node *LocalNode) Server(ErrChan chan error, RpcChan chan int) {
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
		select {
		case RPC := <-RpcChan:
			go func() {
				msg, err := node.GenerateRequest(RPC)
				if err != nil {
					log.Printf("error: generating request %v", err)
					return
				}
				//send request
			}()
		default:
			conn.SetReadDeadline(time.Now().Add(time.Second)) // Set a timeout
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() { //Avoid logging timeouts
					continue
				}
				log.Printf("error: reading from connection %v", err)
			}
			//Delay new connections until the semaphore's buffer is not full
			semaphore <- struct{}{}

			stream := make([]byte, n)
			copy(stream, buffer[:])

			go func(stream []byte, address net.Addr, connection net.PacketConn) {
				defer func() { <-semaphore }()
				replyMsg, err := node.GenerateReply(stream)
				if err != nil {
					log.Printf("error: generating reply %v", err)
					return //Add logging
				}
				err = node.Send(replyMsg, address, connection)
				if err != nil {
					log.Printf("error: sending %v", err)
					return //Add logging
				}
			}(stream, addr, conn)
		}
	}
}

func (node *LocalNode) Send(msg *message, address net.Addr, connection net.PacketConn) error {
	msg.SenderNode = node.Node
	msg.R
	stream, err := node.Serialize(msg)
	if err != nil {
		return err
	}
	_, err = connection.WriteTo(stream, address)
	return err
}

func (node *LocalNode) GenerateRequest(RCP int) (*message, error) {
	var request interface{}
	switch RCP {
	case messagePing:
		request = "PING"
	case messageStore:
		request = "STORE"
	case messageFindNode:
		request = "FINDNODE"
	case messageFindValue:
		request = "FINDVALUE"
	default:
		return nil, fmt.Errorf("invalid RPC type %d", RCP)
	}
}

func (node *LocalNode) GenerateReply(stream []byte) (*message, error) {
	if len(stream) > 0 {
		msg, err := node.Deserialize(stream)
		var reply interface{}

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
			return nil, fmt.Errorf("invalid type in message with transaction ID %d", msg.TransactionID)
		}
		if err != nil {
			return nil, err
		}
		replyMessage := &message{
			MessageType:   msg.MessageType,
			TransactionID: msg.TransactionID,
			SenderNode:    node.Node,
			ReceiverNode:  msg.SenderNode,
			Request:       false,
			Data:          reply,
		}
		return replyMessage, nil
	}
	return nil, fmt.Errorf("stream length is 0")
}
