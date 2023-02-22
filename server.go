package main

import (
	"fmt"
	"ggcache/cache"
	"ggcache/proto"
	"io"
	"net"
	"time"
)

type Server struct {
	ServerOpt
	c cache.Cacher
}
type ServerOpt struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

func NewServer(opt ServerOpt, c cache.Cacher) *Server {
	return &Server{
		ServerOpt: opt,
		c:         c,
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	fmt.Println("server start at  ", s.ListenAddr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Print("accept error: ", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	//buf := make([]byte, 2048)
	fmt.Println("server accept connection: ", conn.RemoteAddr().String())
	for {
		cmd, err := proto.ParseCommand(conn)

		if err != nil {
			if err == io.EOF {
				fmt.Println("server connection closed: ", conn.RemoteAddr().String())
				break
			}
			fmt.Println("server parse command error: ", err)
			break
		}
		//fmt.Println("server receive command: ", cmd)
		go s.handleCommand(conn, cmd)
	}
}
func (s *Server) handleCommand(conn net.Conn, b interface{}) {
	switch m := b.(type) {
	case *proto.CommandSet:
		if err := s.handleSetCmd(conn, m); err != nil {
			fmt.Println("server set command error: ", err)
			return
		}
	case *proto.CommandGet:
		if err := s.handleGetCmd(conn, m); err != nil {
			fmt.Println("server get command error: ", err)
			return
		}
	}
	//fmt.Println("server handle command: ", b)
	return
}

func (s *Server) handleSetCmd(conn net.Conn, m *proto.CommandSet) error {
	if err := s.c.Set([]byte(m.Key), []byte(m.Value), time.Duration(m.TTL)); err != nil {
		return err
	}

	//go s.sendToFollowers(context.TODO(), m)
	return nil

}
func (s *Server) handleGetCmd(conn net.Conn, m *proto.CommandGet) error {
	val, err := s.c.Get([]byte(m.Key))
	if err != nil {
		return err
	}
	//fmt.Println("server get value: ", string(val))
	_, err = conn.Write(val)
	return err
}
