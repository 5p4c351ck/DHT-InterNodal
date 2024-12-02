package node

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
)

const (
	messagePing = iota
	messageStore
	messageFindNode
	messageFindValue
)

type Codec interface {
	Serialize()
	Deserialize()
}

type CodecImp struct {
	msg    *message
	stream []byte
}

func (codec *CodecImp) Serialize() error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(codec.msg)
	if err != nil {
		return err
	}
	lenght := buffer.Len()

	var lengthSerialized [8]byte
	binary.PutUvarint(lengthSerialized[:], uint64(lenght))

	var SerializedMsg []byte
	SerializedMsg = append(SerializedMsg, lengthSerialized[:]...)
	SerializedMsg = append(SerializedMsg, buffer.Bytes()...)

	codec.stream = SerializedMsg
	return nil
}

func (codec *CodecImp) Deserialize() error {
	buffer := bytes.NewBuffer(codec.stream)
	msg := &message{}
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(msg)
	if err != nil {
		return err
	}
	codec.msg = msg
	return nil
}

func ReadIntoStream(conn io.Reader) ([]byte, error) {
	lengthSerialized := make([]byte, 8)
	_, err := conn.Read(lengthSerialized)
	if err != nil {
		return nil, err
	}
	length, err := binary.ReadUvarint(bytes.NewBuffer(lengthSerialized))
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, length)
	_, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

type message struct {
	messageType   int
	transactionID int64
	senderNode    *Node
	receiverNode  *Node
	request       bool
	data          interface{}
}

type storeRequest struct {
	Source bool
	Data   []byte
}

type storeReply struct {
	Success bool
}

type findNodeRequest struct {
	TargetNodeID []byte
}

type findNodeReply struct {
	ClosestNodes [][]*Node
}

type findValueRequest struct {
	TargetNodeID []byte
}

type findValueReply struct {
	Value        []byte
	ClosestNodes [][]*Node
}

func initMessages() {
	gob.Register(&storeRequest{})
	gob.Register(&storeReply{})
	gob.Register(&findNodeRequest{})
	gob.Register(&findNodeReply{})
	gob.Register(&findValueRequest{})
	gob.Register(&findValueReply{})
}
