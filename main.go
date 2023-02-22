package main

import (
	"flag"
	"ggcache/cache"
)

func main() {
	listenAddr := flag.String("listenAddr", ":7070", "server listen addr")
	leaderAddr := flag.String("leaderAddr", "", "leader addr")
	flag.Parse()

	ops := &ServerOpt{
		ListenAddr: *listenAddr,
		LeaderAddr: *leaderAddr,
		IsLeader:   len(*leaderAddr) == 0,
	}

	s := NewServer(*ops, cache.NewCache())
	err := s.Start()

	if err != nil {
		panic(err)
	}

}

// command
