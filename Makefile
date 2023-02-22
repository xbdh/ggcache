build:
	go build -o bin/ggcache

run: build
	./bin/ggcache --listenAddr=":7070"
runf: build
	./bin/ggcache --listenAddr=":7080" --leaderAddr=":7070"

