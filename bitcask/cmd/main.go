// package main
//
// import (
//
//	"fmt"
//	"gihub.com/xbdh/bitcask"
//
// )
//
//	func main() {
//		//os.Mkdir("./store", 0755)
//		ops := bitcask.Options{
//			Dir:          "./store",
//			MaxStoreSize: 100,
//			SyncWrite:    true,
//		}
//
//		log, err := bitcask.NewLog(ops)
//		if err != nil {
//			fmt.Printf("创建log失败:=%v", err)
//			panic(err)
//		}
//		//for i := 0; i < 20; i++ {
//		//	key := fmt.Sprintf("key%d", i)
//		//	value := fmt.Sprintf("value%d", i)
//		//	err = log.Append([]byte(key), []byte(value))
//		//	if err != nil {
//		//		panic(err)
//		//	}
//		//}
//		for i := 0; i < 20; i++ {
//			key := fmt.Sprintf("key%d", i)
//			value, err := log.Read([]byte(key))
//			if err != nil {
//				panic(err)
//			}
//			fmt.Println(string(value))
//		}
//
// }
package cmd
