package fsm

import (
	"bytes"
	"fmt"
	"ggcache/bitcask"
	"ggcache/proto"
	"github.com/hashicorp/raft"
	"io"
)

type MyFSM struct {
	L *bitcask.Log
}

func (m *MyFSM) Apply(log *raft.Log) interface{} {

	fmt.Println("apply log: ", log.Data)
	r := bytes.NewReader(log.Data)

	cmdSetI, err := proto.ParseCommand(r)
	if err != nil {
		fmt.Println("ParseCommand error: ", err)
		return err
	}
	cmd := cmdSetI.(*proto.CommandSet)
	m.L.Append(cmd.Key, cmd.Value)

	return nil
}

func (m *MyFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (m *MyFSM) Restore(snapshot io.ReadCloser) error {
	return nil
}
