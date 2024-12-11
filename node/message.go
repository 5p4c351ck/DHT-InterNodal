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
	Serialize(msg *message) ([]byte, error)
	Deserialize(bytestream []byte) (*message, error)
}

type CodecImp struct{} //Possibly add some state later

func NewCodec() Codec {
	codec := &CodecImp{}
	return codec
}

func (codec *CodecImp) Serialize(msg *message) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(msg)
	if err != nil {
		return nil, err
	}
	lenght := buffer.Len()

	var lengthSerialized [8]byte
	binary.PutUvarint(lengthSerialized[:], uint64(lenght))

	var SerializedMsg []byte
	SerializedMsg = append(SerializedMsg, lengthSerialized[:]...)
	SerializedMsg = append(SerializedMsg, buffer.Bytes()...)

	return SerializedMsg, nil
}

func (codec *CodecImp) Deserialize(bytestream []byte) (*message, error) {
	buffer := bytes.NewBuffer(bytestream[8:])
	msg := &message{}
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
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
	MessageType   int
	TransactionID int64
	SenderNode    *Node
	ReceiverNode  *Node
	Request       bool
	Data          interface{}
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

func storeReplyMsg(msg *message) (storeReply, error) {

}

func findNodeReplyMsg(msg *message) (findNodeReply, error) {

}

func findValueReplyMsg(msg *message) (findValueReply, error) {

}

func initMessages() {
	gob.Register(&storeRequest{})
	gob.Register(&storeReply{})
	gob.Register(&findNodeRequest{})
	gob.Register(&findNodeReply{})
	gob.Register(&findValueRequest{})
	gob.Register(&findValueReply{})
}
