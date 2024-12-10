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
			err = node.Reply(msg, addr, conn)
			if err != nil {
				return //Add logging
			}
		}(stream, addr, conn)
	}
}

func (node *LocalNode) Request(m *message) error {
	if m.SenderNode == nil || m.ReceiverNode == nil {
		return fmt.Errorf("sender or receiver is nil")
	}
	if !m.Request {
		return fmt.Errorf("message is a reply")
	}
	raddr := &net.UDPAddr{
		IP:   m.ReceiverNode.IP,
		Port: m.ReceiverNode.Port,
	}
	conn, err := net.DialUDP(protocolUDP, nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	codec := NewCodec()
	input, err := codec.Serialize(m)
	if err != nil {
		return err
	}
	_, err = conn.Write(input)
	if err != nil {
		return err
	}

	return nil
}

func (node *LocalNode) Reply(m *message, address net.Addr, connection net.PacketConn) error {
	replyMessage := &message{}
	switch m.MessageType {
	case messagePing:
		//Set message for Ping reply
	case messageStore:
		//Set message for Store reply
	case messageFindNode:
		//Set message for FindNode reply
	case messageFindValue:
		//Set message for FindValue reply
	default:
		//Set message for default case
	}
	codec := NewCodec()
	input, err := codec.Serialize(replyMessage)
	if err != nil {
		return err
	}
	connection.WriteTo(input, address)
	return nil
}
