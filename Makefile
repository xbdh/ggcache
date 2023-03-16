build:
	go build -o bin/ggcache

run: build
	./bin/ggcache --listenAddr=":7070"
runf: build
	./bin/ggcache --listenAddr=":7080" --leaderAddr=":7070"

runff: build
	./bin/ggcache --listenAddr=":7090" --leaderAddr=":7070"


run1: build
	./bin/ggcache --listenAddr=":7091" --nodename="node1" --port=8081 --store="./kvstore/s1"

run2: build
	./bin/ggcache --listenAddr=":7092" --nodename="node2" --port=8082 --store="./kvstore/s2"

run3: build
	./bin/ggcache --listenAddr=":7093" --nodename="node3" --port=8083 --store="./kvstore/s3"
