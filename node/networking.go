package node

import (
	"fmt"
	"net"
)

const (
	protocolUDP = "udp"
	port        = ":8080"
)

func (node *LocalNode) Server() error {
	conn, err := net.ListenPacket(protocolUDP, port)
	if err != nil {
		return err
	}
	defer conn.Close()
	bufferSize := 1024
	maxConns := 10

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

		go func(stream []byte, addr net.Addr) {
			defer func() { <-semaphore }()
			codec := &CodecImp{}
			msg, err := codec.Deserialize(stream)
			if err != nil {
				return //Add logging
			}
			err = node.Reply(msg)
			if err != nil {
				return //Add logging
			}
		}(stream, addr)
	}
}

func (node *LocalNode) Request(m *message, codec Codec) error {
	if m.senderNode == nil || m.receiverNode == nil {
		return fmt.Errorf("sender or receiver is nil")
	}
	if !m.request {
		return fmt.Errorf("message is a reply")
	}
	raddr := &net.UDPAddr{
		IP:   m.receiverNode.IP,
		Port: m.receiverNode.Port,
	}
	conn, err := net.DialUDP(protocolUDP, nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()
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

func (node *LocalNode) Reply(m *message) error {
	return nil
}
