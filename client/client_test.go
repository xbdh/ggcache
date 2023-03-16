package cmd

import (
	"fmt"
	_ "fmt"
	"testing"
)

func BenchmarkClient_Set(b *testing.B) {
	client, _ := NewClient(":7091", Options{})
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err := client.Set(key, value, 0)
		if err != nil {
			panic(err)
		}

	}
}
