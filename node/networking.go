package node

import (
	"fmt"
	"net"
)

const (
	protocol = "udp"
	port     = ":8080"
)

func (node *LocalNode) Server() error {
	conn, err := net.ListenPacket(protocol, port)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Println("Listening on port", port)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			continue
		}
		fmt.Printf("Received %s from address %s", string(buffer[:n]), addr)
	}
}

func (node *LocalNode) sendRequest(m message) error {
	if m.senderNode == nil || m.receiverNode == nil {
		return fmt.Errorf("Sender or receiver is nil")
	}
	if !m.request {
		return fmt.Errorf("Message is a reply")
	}
	raddr := &net.UDPAddr{
		IP:   m.receiverNode.IP,
		Port: m.receiverNode.Port,
	}
	conn, err := net.DialUDP(protocol, nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	input := []byte("Testing connectivity")
	_, err = conn.Write(input)

}
