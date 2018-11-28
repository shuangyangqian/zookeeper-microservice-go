package main

import (
	"flag"
	"zookeeper-microservice-go/service_provider"
	"strings"
	"github.com/golang/glog"
	"fmt"
	"net/http"
	"io/ioutil"
)

var serviceName string
var zkHost string

func main() {

	flag.StringVar(&serviceName, "name", "red",
		"the service name registered to zk")
	flag.StringVar(&zkHost, "zkHost", "127.0.0.1:2181;",
		"the zk host list to registered")

    flag.Parse()

	// 创建client
	zkHostSlice := strings.Split(zkHost, ";")
	var zkServers []string
	for _, host := range zkHostSlice {
		zkServers = append(zkServers, host)
	}
	client, err := service_provider.NewClient(zkServers, "/api", 10)
	if err != nil {
		panic(err)
	}
	glog.Infof("get connection to zk:%s", zkHost)

	nodes, err := client.GetNodes(serviceName)
	if err != nil {
		panic(err)
	}
	glog.Infof("get %d nodes with service name %s", len(nodes), serviceName)
	for _, node := range nodes {
		fmt.Println(contentInPage(node.Host, node.Port))
	}

}

func contentInPage(host string, port int) string {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d",
		host, port), nil)
	if err != nil {
		glog.Error(err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return ""
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return string(respByte[:])


}
