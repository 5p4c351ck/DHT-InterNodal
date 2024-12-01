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

func (node *LocalNode) sendRequest() error {

}
