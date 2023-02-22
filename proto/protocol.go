package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Command byte

const (
	CmdNone Command = iota
	CmdSet
	CmdGet
	CmdDel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   uint32
}
type CommandGet struct {
	Key []byte
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)
	binary.Write(buf, binary.LittleEndian, uint32(c.TTL))
	return buf.Bytes()

}
func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}

func ParseCommand(r io.Reader) (interface{}, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}
	switch cmd {
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	default:
		return nil, errors.New("unknown command")
	}
	//return nil, nil
}
func parseSetCommand(r io.Reader) *CommandSet {
	var keyLen uint32
	binary.Read(r, binary.LittleEndian, &keyLen)
	key := make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &key)
	var valueLen uint32
	binary.Read(r, binary.LittleEndian, &valueLen)
	value := make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &value)
	var ttl uint32
	binary.Read(r, binary.LittleEndian, &ttl)

	//fmt.Sprintf("parse setcmd-> key: %s, value: %s, ttl: %d", key, value, ttl)

	return &CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
}
func parseGetCommand(r io.Reader) *CommandGet {
	var keyLen uint32
	binary.Read(r, binary.LittleEndian, &keyLen)
	key := make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &key)

	return &CommandGet{
		Key: key,
	}
}
