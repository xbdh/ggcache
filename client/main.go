package main

import (
	"fmt"
	"ggcache/proto"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}
type options struct {
}

func NewClient(endpoints string, opt options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoints)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Set(key string, value string, ttl uint32) error {
	cmds := &proto.CommandSet{
		Key:   []byte(key),
		Value: []byte(value),
		TTL:   ttl,
	}
	c.conn.Write(cmds.Bytes())
	return nil
}
func (c Client) Get(key string) (string, error) {
	cmdg := &proto.CommandGet{
		Key: []byte(key),
	}
	c.conn.Write(cmdg.Bytes())
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (c *Client) Close() error {
	return c.Close()
}
func main() {
	client, err := NewClient(":7070", options{})
	if err != nil {
		panic(err)
	}
	_ = client.Set("FOO", "BAR", 0)
	time.Sleep(1 * time.Second)
	val, err := client.Get("FOO")
	if err != nil {
		panic(err)
	}
	fmt.Println("client get value: ", val)
}
