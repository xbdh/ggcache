package main

import (
	"encoding/binary"
	"fmt"
	"ggcache/bitcask"
	"ggcache/fsm"
	"ggcache/proto"
	"github.com/hashicorp/raft"
	"io"
	"net"
	"os"
	"time"
)

type Server struct {
	ServerOpt
	l *bitcask.Log
	r *raft.Raft
	//members map[*client.Client]struct{}
}
type ServerOpt struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
	// for test
	Port     int
	NodeName string
}

func NewServer(opt ServerOpt, l *bitcask.Log) *Server {
	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID(opt.NodeName) // test

	fsm := fsm.MyFSM{L: l}
	log := raft.NewInmemStore()
	stable := raft.NewInmemStore()
	snap := raft.NewInmemSnapshotStore()
	//transport:=raft.NewInmemTransport(raft.ServerAddress("1"))
	ip := "127.0.0.1"
	binaddr := fmt.Sprintf("%s:%d", ip, opt.Port)
	addr, _ := net.ResolveTCPAddr("tcp", binaddr)
	transport, err := raft.NewTCPTransport(addr.String(), addr, 5, 10*time.Second, os.Stderr)
	if err != nil {
		fmt.Println("transport error: ", err)
		return nil
	}
	server1 := raft.Server{
		//Suffrage: raft.Voter,
		ID:      raft.ServerID("node1"),
		Address: raft.ServerAddress("127.0.0.1:8081"),
	}
	server2 := raft.Server{
		//Suffrage: raft.Voter,
		ID:      raft.ServerID("node2"),
		Address: raft.ServerAddress("127.0.0.1:8082"),
	}
	server3 := raft.Server{
		//Suffrage: raft.Voter,
		ID:      raft.ServerID("node3"),
		Address: raft.ServerAddress("127.0.0.1:8083"),
	}
	config := raft.Configuration{
		Servers: []raft.Server{server1, server2, server3},
	}
	r, err := raft.NewRaft(cfg, &fsm, log, stable, snap, transport)
	fmt.Printf("raft info --------: %v", r)

	future := r.BootstrapCluster(config)
	if err := future.Error(); err != nil {
		fmt.Println("bootstrap error: +++++", err)
		return nil
	}
	fmt.Println("bootstrap success: ", future)

	//go func() {
	//	<-r.LeaderCh()
	//	id, serverID := r.LeaderWithID()
	//	fmt.Println("leader id: ", id, "serverID: ", serverID)
	//}()

	return &Server{
		ServerOpt: opt,
		l:         l,
		r:         r,
		//members:   make(map[*client.Client]struct{}),
	}

}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	//if !s.IsLeader && len(s.LeaderAddr) > 0 {
	//	go func() {
	//		err := s.dialLeader()
	//		if err != nil {
	//			fmt.Println("follower can not dial to leader error: ", err)
	//		}
	//	}()
	//}

	// todo：还需要判断是否是leader:获取leader的id之后写入server的结构体中，
	//然后在每次请求的时候判断是否是leader，只有leader才能处理请求

	go func() {
		for {
			time.Sleep(3 * time.Second)
			id, serverID := s.r.LeaderWithID()
			fmt.Println("leader id: ", id, "serverID: ", serverID)
		}
	}()
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
		//case *proto.CommandJoin:
		//	if err := s.handleJoinCmd(conn, m); err != nil {
		//		fmt.Println("server join command error: ", err)
		//		return
		//	}
	}
	//fmt.Println("server handle command: ", b)
	return
}

func (s *Server) handleSetCmd(conn net.Conn, m *proto.CommandSet) error {
	var respSet proto.ResponseSet
	var err error
	//err := s.c.Set([]byte(m.Key), []byte(m.Value), time.Duration(m.TTL))

	//go func() {
	//	for member := range s.members {
	//		err = member.Set(string(m.Key), string(m.Value), m.TTL)
	//		if err != nil {
	//			fmt.Println("member set kv data from leader error: ", err)
	//			continue
	//		}
	//	}
	//}()
	future := s.r.Apply(m.Bytes(), time.Duration(3*time.Second))
	if err := future.Error(); err != nil {
		fmt.Println("raft apply error: ", err)
		// need return?
	}

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
	val, err := s.l.Read([]byte(m.Key))
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

//func (s *Server) handleJoinCmd(conn net.Conn, m *proto.CommandJoin) error {
//	fmt.Println("member  joined the cluster :", conn.RemoteAddr().String())
//	c := client.NewFromConn(conn)
//	s.members[c] = struct{}{}
//	return nil
//}
