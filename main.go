package main

import (
	"flag"
	"fmt"
	"ggcache/bitcask"
)

func main() {
	listenAddr := flag.String("listenAddr", ":7070", "server listen addr")
	leaderAddr := flag.String("leaderAddr", "", "leader addr")
	port := flag.Int("port", 8080, "server listen addr")
	nodename := flag.String("nodename", "node1", "node name")
	store := flag.String("store", "./store", "store path")
	fmt.Println("listenAddr: ", *listenAddr)
	fmt.Println("leaderAddr: ", *leaderAddr)
	fmt.Println("port: ", *port)
	fmt.Println("nodename: ", *nodename)
	fmt.Println("store: ", *store)
	flag.Parse()

	ops := &ServerOpt{
		ListenAddr: *listenAddr,
		LeaderAddr: *leaderAddr,
		IsLeader:   len(*leaderAddr) == 0,
		Port:       *port,
		NodeName:   *nodename,
	}

	lops := bitcask.Options{
		Dir:          *store,
		MaxStoreSize: 1024 * 1024 * 1024,
		SyncWrite:    true,
	}
	l, err := bitcask.NewLog(lops)
	if err != nil {
		fmt.Println("log init error: ", err)
		panic(err)
	}
	s := NewServer(*ops, l)
	err = s.Start()

	if err != nil {
		fmt.Println("server start error: ", err)
		panic(err)
	}

}

// command
