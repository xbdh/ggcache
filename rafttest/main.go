package main

import (
	"fmt"
	"github.com/hashicorp/raft"
	"os"
	"time"
)

func main() {
	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID("mynode1")
	fsm := raft.MockFSM{}
	log := raft.NewInmemStore()
	stable := raft.NewInmemStore()
	snap := raft.NewInmemSnapshotStore()
	//transport:=raft.NewInmemTransport(raft.ServerAddress("1"))
	transport, err := raft.NewTCPTransport("127.0.0.1:8080", nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		fmt.Println("transport error: ", err)
		panic(err)
	}
	server1 := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(cfg.LocalID),
		Address:  raft.ServerAddress("127.0.0.1:8080"),
	}
	//server2 := raft.Server{
	//	Suffrage: raft.Voter,
	//	ID:       raft.ServerID("mynode2"),
	//	Address:  raft.ServerAddress("127.0.0.1:8090"),
	//}
	config := raft.Configuration{
		Servers: []raft.Server{server1},
	}
	r, err := raft.NewRaft(cfg, &fsm, log, stable, snap, transport)
	fmt.Printf("raft: %v", r)

	future := r.BootstrapCluster(config)
	if err := future.Error(); err != nil {
		fmt.Println("bootstrap error: ", err)
		panic(err)
	}

}
