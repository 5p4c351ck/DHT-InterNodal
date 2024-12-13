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

func (node *LocalNode) Server(ErrChan chan error, RpcChan chan *message) {
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
			go func(msg *message, connection net.PacketConn) {
				err = node.Send(msg, connection)
				if err != nil {
					log.Printf("error: sending RPC %v", err)
					return
				}
				//send request
			}(RPC, conn)
		default:
			conn.SetReadDeadline(time.Now().Add(time.Second)) // Set a timeout
			n, _, err := conn.ReadFrom(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() { //Avoid logging timeouts
					continue
				}
				log.Printf("error: reading from connection %v", err)
				continue
			}
			//Delay new connections until the semaphore's buffer is not full
			semaphore <- struct{}{}

			stream := make([]byte, n)
			copy(stream, buffer[:])

			go func(stream []byte, connection net.PacketConn) {
				defer func() { <-semaphore }()
				replyMsg, err := node.GenerateRpcReply(stream)
				if err != nil {
					log.Printf("error: generating reply %v", err)
					return
				}
				err = node.Send(replyMsg, connection)
				if err != nil {
					log.Printf("error: replying %v", err)
					return
				}
			}(stream, conn)
		}
	}
}

func GenerateAddress(msg *message) (net.Addr, error) {
	IP := msg.ReceiverNode.IP
	PORT := msg.ReceiverNode.Port
	if !msg.Request {
		IP = msg.SenderNode.IP
		PORT = msg.SenderNode.Port
	}
	if PORT < 0 || PORT > 65535 {
		return nil, fmt.Errorf("invalid Port number %d", PORT)
	}
	ParsedIP := net.ParseIP(IP.String())
	if ParsedIP == nil {
		return nil, fmt.Errorf("invalid IP address %s", IP.String())
	}
	udpAddr := &net.UDPAddr{
		IP:   ParsedIP,
		Port: PORT,
	}
	return udpAddr, nil
}

func (node *LocalNode) Send(msg *message, connection net.PacketConn) error {
	address, err := GenerateAddress(msg)
	if err != nil {
		return err
	}
	msg.SenderNode = node.Node
	stream, err := node.Serialize(msg)
	if err != nil {
		return err
	}
	_, err = connection.WriteTo(stream, address)
	return err
}

func (node *LocalNode) GenerateRpcRequest(RPC *message) (*message, error) {
	var requestPayload interface{}
	switch RPC.MessageType {
	case messagePing:
		requestPayload = "PING"
	case messageStore:
		requestPayload = "STORE"
	case messageFindNode:
		requestPayload = "FINDNODE"
	case messageFindValue:
		requestPayload = "FINDVALUE"
	default:
		return nil, fmt.Errorf("invalid RPC type %d", RPC.MessageType)
	}
	RPC.Request = true
	RPC.Payload = requestPayload
	return RPC, nil
}

func (node *LocalNode) GenerateRpcReply(stream []byte) (*message, error) {
	if len(stream) > 0 {
		msg, err := node.Deserialize(stream)
		var replyPayload interface{}

		switch msg.MessageType {
		case messagePing:
			replyPayload = "PONG"
		case messageStore:
			replyPayload, err = storeReplyMsg(msg)
		case messageFindNode:
			replyPayload, err = findNodeReplyMsg(msg)
		case messageFindValue:
			replyPayload, err = findValueReplyMsg(msg)
		default:
			return nil, fmt.Errorf("invalid type in message with transaction ID %d", msg.TransactionID)
		}
		if err != nil {
			return nil, err
		}
		replyMessage := &message{
			MessageType:   msg.MessageType,
			TransactionID: msg.TransactionID,
			SenderNode:    nil,
			ReceiverNode:  msg.SenderNode,
			Request:       false,
			Payload:       replyPayload,
		}
		return replyMessage, nil
	}
	return nil, fmt.Errorf("stream length is 0")
}
