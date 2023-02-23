package client

import (
	"errors"
	"fmt"
	"ggcache/proto"
	"net"
)

type Client struct {
	conn net.Conn
}
type Options struct {
}

func NewFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}

}
func NewClient(endpoints string, opt Options) (*Client, error) {
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

	_, err := c.conn.Write(cmds.Bytes())
	if err != nil {
		return err
	}
	res, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		fmt.Println("client parse set response error: ", err)
		return err
	}
	if res.Status == proto.StatusOK {
		return nil
	} else {
		return errors.New("client set error")
	}

	return nil
}
func (c *Client) Get(key string) (string, error) {
	cmdg := &proto.CommandGet{
		Key: []byte(key),
	}
	_, err := c.conn.Write(cmdg.Bytes())
	if err != nil {
		return "", err
	}

	res, err := proto.ParseGetResponse(c.conn)

	if err != nil {
		fmt.Println("client parse get response error: ", err)
	}
	if res.Status == proto.StatusOK {
		//fmt.Println("client get value:--- ", string(res.Value))
		return string(res.Value), nil
	} else if res.Status == proto.StatusNotFound {
		return "", errors.New("key not found")
	} else {
		return "", errors.New("client get error")
	}
	return "", nil
}

func (c *Client) Close() error {
	return c.Close()
}
