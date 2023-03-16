package main

import (
	"fmt"
	cmd "ggcache/client"

	"time"
)

func main() {
	client, err := cmd.NewClient(":7091", cmd.Options{})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err = client.Set(key, value, 0)
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
