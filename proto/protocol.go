package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Command byte
type Status byte

func (s Status) String() string {
	switch s {
	case StatusNone:
		return "StatusNone"
	case StatusOK:
		return "StatusOK"
	case StatusErr:
		return "StatusErr"
	case StatusNotFound:
		return "StatusNotFound"
	default:
		return "StatusUnknown"
	}
}

const (
	CmdNone Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

const (
	StatusNone Status = iota
	StatusOK
	StatusErr
	StatusNotFound
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   uint32
}
type CommandGet struct {
	Key []byte
}
type CommandJoin struct{}

type ResponseGet struct {
	Status Status
	Value  []byte
}
type ResponseSet struct {
	Status Status
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

func (r *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	binary.Write(buf, binary.LittleEndian, uint32(len(r.Value)))
	binary.Write(buf, binary.LittleEndian, r.Value)
	return buf.Bytes()
}
func (r *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	return buf.Bytes()
}

func ParseCommand(r io.Reader) (interface{}, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}
	switch cmd {
	case CmdSet:
		return ParseSetCommand(r), nil
	case CmdGet:
		return ParseGetCommand(r), nil
	case CmdJoin:
		return &CommandJoin{}, nil
	default:
		return nil, errors.New("unknown command")
	}
	//return nil, nil
}

//func parseJoinCommand(r io.Reader) *CommandJoin {
//	var cmd Command
//	binary.Read(r, binary.LittleEndian, &cmd)
//
//	return &CommandJoin{}
//}
func ParseSetCommand(r io.Reader) *CommandSet {
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
func ParseGetCommand(r io.Reader) *CommandGet {
	var keyLen uint32
	binary.Read(r, binary.LittleEndian, &keyLen)
	key := make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &key)

	return &CommandGet{
		Key: key,
	}
}

func ParseGetResponse(r io.Reader) (*ResponseGet, error) {
	var resp ResponseGet
	var status Status
	err := binary.Read(r, binary.LittleEndian, &status)
	if err != nil {
		return nil, err
	}

	var valueLen uint32
	binary.Read(r, binary.LittleEndian, &valueLen)
	value := make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &value)
	if status == StatusErr {
		resp.Status = StatusErr
	}
	if status == StatusNotFound {
		resp.Status = StatusNotFound
	}
	if status == StatusOK {
		resp.Status = StatusOK
		resp.Value = value
	}
	return &resp, nil
}
func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	var status Status
	if err := binary.Read(r, binary.LittleEndian, &status); err != nil {
		return nil, err
	}

	return &ResponseSet{
		Status: status,
	}, nil
}
