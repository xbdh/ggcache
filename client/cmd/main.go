package main

import (
	"fmt"
	"ggcache/client"
	"time"
)

func main() {
	client, err := client.NewClient(":7070", client.Options{})
	if err != nil {
		panic(err)
	}
	err = client.Set("FOO", "BAR", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("client set complete")
	time.Sleep(1 * time.Second)
	val, err := client.Get("FOO")
	if err != nil {
		fmt.Println("client get error: ", err)
		//panic(err)
		return
	}
	fmt.Println("client get value: ===", val)
}
