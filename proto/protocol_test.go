package proto

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSetCommand(t *testing.T) {
	cmds := &CommandSet{
		Key:   []byte("FOO"),
		Value: []byte("BAR"),
		TTL:   2,
	}
	rs := bytes.NewReader(cmds.Bytes())
	res, err := ParseCommand(rs)
	assert.Nil(t, err)
	assert.Equal(t, cmds, res)

}

func TestParseGetCommand(t *testing.T) {
	cmdg := &CommandGet{
		Key: []byte("FOO"),
	}
	rg := bytes.NewReader(cmdg.Bytes())
	resg, err := ParseCommand(rg)
	assert.Nil(t, err)

	assert.Equal(t, cmdg, resg)
}
