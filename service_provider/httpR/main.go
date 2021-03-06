package main

import (
	"github.com/gin-gonic/gin"
	"net"
	"fmt"
	"os"
	"github.com/golang/glog"
	"zookeeper-microservice-go/service_provider"
	"flag"
	"strings"
)

var serviceName string
var servicePort int
var zkHost string
var timeout int

func main() {

	flag.StringVar(&serviceName, "name", "red",
		"the service name registered to zk")
	flag.IntVar(&servicePort, "port", 8080,
		"the port service listened on")
	flag.StringVar(&zkHost, "zkHost", "127.0.0.1:2181;",
		"the zk host list to registered")
	flag.IntVar(&timeout, "timeout", 1,
		"timeout to connect zk cluster")

	flag.Parse()
	// 获取本机IP地址
	IpAddress, err := getIp()
	if err != nil {
		glog.Error("cannot get ip address")
		panic(err)
	}

	glog.Infof("get local ip is %s", IpAddress)

	// 创建服务对象
	service := &service_provider.ServiceNode{
		Name: serviceName,
		Host: IpAddress,
		Port: servicePort,
	}

	// 创建client
	zkHostSlice := strings.Split(zkHost, ";")
	var zkServers []string
	for _, host := range zkHostSlice {
		zkServers = append(zkServers, host)
	}
	client, err := service_provider.NewClient(zkServers, "/api", timeout)
	if err != nil {
		panic(err)
	}
	glog.Infof("get connection to zk:%s", zkHost)

	// 注册服务
	if err := client.Register(service); err != nil {
		defer client.Close()
		panic(err)
	}
	glog.Infof("registry service:%s-%s:%d to zk", service.Name, service.Host, service.Port)
	defer client.Close()

	u := gin.Default()
	u.GET("/", service.IndexController)

	glog.Infof("listen and serve %s:%d", IpAddress, service_provider.PORT)
	u.Run(fmt.Sprintf("%s:%d",IpAddress, service_provider.PORT))
}

func getIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	err = service_provider.UnknowErr{
		Detail: "cannot get ip",
	}
	return "", err
}
