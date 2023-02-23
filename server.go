package main

import (
	"encoding/binary"
	"fmt"
	"ggcache/cache"
	"ggcache/client"
	"ggcache/proto"
	"io"
	"net"
	"time"
)

type Server struct {
	ServerOpt
	c       cache.Cacher
	members map[*client.Client]struct{}
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
		members:   make(map[*client.Client]struct{}),
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	if !s.IsLeader && len(s.LeaderAddr) > 0 {
		go func() {
			err := s.dialLeader()
			if err != nil {
				fmt.Println("follower can not dial to leader error: ", err)
			}
		}()
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
func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		fmt.Println("dial leader error: ", err)
		return err
	}
	fmt.Println("follower dial leader success: ", s.LeaderAddr)

	// send join command
	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)

	s.handleConn(conn)
	return nil
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
	//fmt.Println("server connection closed: ", conn.RemoteAddr().String())

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
	case *proto.CommandJoin:
		if err := s.handleJoinCmd(conn, m); err != nil {
			fmt.Println("server join command error: ", err)
			return
		}
	}
	//fmt.Println("server handle command: ", b)
	return
}

func (s *Server) handleSetCmd(conn net.Conn, m *proto.CommandSet) error {
	var respSet proto.ResponseSet
	err := s.c.Set([]byte(m.Key), []byte(m.Value), time.Duration(m.TTL))

	go func() {
		for member := range s.members {
			err = member.Set(string(m.Key), string(m.Value), m.TTL)
			if err != nil {
				fmt.Println("member set kv data from leader error: ", err)
				continue
			}
		}
	}()

	if err != nil {
		respSet.Status = proto.StatusErr
		return err
	}

	respSet.Status = proto.StatusOK
	_, err = conn.Write(respSet.Bytes())
	return err

}
func (s *Server) handleGetCmd(conn net.Conn, m *proto.CommandGet) error {
	var respGet proto.ResponseGet
	val, err := s.c.Get([]byte(m.Key))
	if err != nil {
		respGet.Status = proto.StatusNotFound
		_, err = conn.Write(respGet.Bytes())
		return err
	}
	respGet.Status = proto.StatusOK
	respGet.Value = val
	_, err = conn.Write(respGet.Bytes())
	return err
}

func (s *Server) handleJoinCmd(conn net.Conn, m *proto.CommandJoin) error {
	fmt.Println("member  joined the cluster :", conn.RemoteAddr().String())
	c := client.NewFromConn(conn)
	s.members[c] = struct{}{}
	return nil
}
