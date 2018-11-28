package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net"
	"fmt"
	"os"
	"github.com/golang/glog"
	"github.com/shuangyangqian/zookeeper-microservice-go/service_provider"
	"flag"
	"strings"
)

var serviceName string
var servicePort int
var zkHost string

func main() {

	flag.StringVar(&serviceName, "--name", "color",
		"the service name registered to zk")
	flag.IntVar(&servicePort, "--port", 8080,
		"the port service listened on")
	flag.StringVar(&zkHost, "--zkHost", "127.0.0.1:2181;",
		"the zk host list to registered")

	flag.Parse()
	// 获取本机IP地址
	IpAddress := getIp()
	if IpAddress == "" {
		glog.Error("cannot get ip address")
	}

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
	client, err := service_provider.NewClient(zkServers, "/api", 10)
	if err != nil {
		panic(err)
	}

	// 注册服务
	if err := client.Register(service); err != nil {
		client.Close()
		panic(err)
	}
	client.Close()

	u := gin.Default()
	u.GET("/", controller)

	u.Run(fmt.Sprintf("%s:%d",IpAddress, service_provider.PORT))
}

func controller(ctx *gin.Context)  {
	ctx.JSON(http.StatusOK, gin.H{"message": "Hi, this is a blue page"})
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
