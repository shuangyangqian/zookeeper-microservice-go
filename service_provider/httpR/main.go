package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net"
	"fmt"
	"os"
	"github.com/golang/glog"
	"zookeeper-go-demo/service_provider"
	"flag"
	"strings"
)

var serviceName string
var servicePort int
var zkHost string

func main() {

	flag.StringVar(&serviceName, "--name", "red",
		"the service name registered to zk")
	flag.IntVar(&servicePort, "--port", 8080,
		"the port service listened on")
	flag.StringVar(&zkHost, "--zkHost", "127.0.0.1:2181;",
		"the zk host list to registered")

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
	client, err := service_provider.NewClient(zkServers, "/api", 10)
	if err != nil {
		panic(err)
	}
	glog.Infof("get connection to zk:%s", zkHost)

	// 注册服务
	if err := client.Register(service); err != nil {
		client.Close()
		panic(err)
	}
	glog.Infof("registry service:%s-%s:%s to zk", service.Name, service.Host, service.Port)
	client.Close()

	u := gin.Default()
	u.GET("/", controller)

	u.Run(fmt.Sprintf("%s:%d",IpAddress, service_provider.PORT))
}

func controller(ctx *gin.Context)  {
	ctx.JSON(http.StatusOK, gin.H{"message":
		fmt.Sprintf("Hi, this is a %s page", serviceName)})
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
	err = new(error)
	return "", err
}
