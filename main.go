package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"encoding/json"
	"fmt"
)




func main() {
	// 服务器地址列表
	servers := []string{"10.104.109.48:2181", "10.104.109.5:2181", "10.104.109.6:2181"}
	client, err := NewClient(servers, "/api", 10)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	node1 := &ServiceNode{"user", "127.0.0.1", 4000}
	node2 := &ServiceNode{"user", "127.0.0.1", 4001}
	node3 := &ServiceNode{"user", "127.0.0.1", 4002}
	if err := client.Register(node1); err != nil {
		panic(err)
	}
	if err := client.Register(node2); err != nil {
		panic(err)
	}
	if err := client.Register(node3); err != nil {
		panic(err)
	}
	nodes, err := client.GetNodes("user")
	if err != nil {
		panic(err)
	}
	for _, node := range nodes {
		fmt.Println(node.Host, node.Port)
	}
}

