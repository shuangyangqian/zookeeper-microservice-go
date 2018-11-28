package service_provider

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"encoding/json"
	"fmt"
)

const PORT = 8080

type ServiceNode struct {
	Name string `json:"name"` // 服务名称，这里是user
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SdClient struct {
	zkServers []string // 多个节点地址
	zkRoot    string // 服务根节点，这里是/api
	conn      *zk.Conn // zk的客户端连接
}

func NewClient(zkServers []string, zkRoot string, timeout int) (*SdClient, error) {
	client := new(SdClient)
	client.zkServers = zkServers
	client.zkRoot = zkRoot
	// 连接服务器
	conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	client.conn = conn
	// 创建服务根节点
	if err := client.ensureRoot(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

// 关闭连接，释放临时节点
func (s *SdClient) Close() {
	s.conn.Close()
}

func (s *SdClient) ensureRoot() error {
	exists, _, err := s.conn.Exists(s.zkRoot)
	if err != nil {
		return err
	}
	if !exists {
		_, err := s.conn.Create(s.zkRoot, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *SdClient) Register(node *ServiceNode) error {
	if err := s.ensureName(node.Name); err != nil {
		return err
	}
	path := s.zkRoot + "/" + node.Name + "/n"
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	_, err = s.conn.CreateProtectedEphemeralSequential(path, data, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

func (s *SdClient) ensureName(name string) error {
	path := s.zkRoot + "/" + name
	exists, _, err := s.conn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		_, err := s.conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *SdClient) GetNodes(name string) ([]*ServiceNode, error) {
	path := s.zkRoot + "/" + name
	// 获取字节点名称
	childs, _, err := s.conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return []*ServiceNode{}, nil
		}
		return nil, err
	}
	nodes := []*ServiceNode{}
	for _, child := range childs {
		fullPath := path + "/" + child
		data, _, err := s.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return nil, err
		}
		node := new(ServiceNode)
		err = json.Unmarshal(data, node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

type UnknowErr struct {
	Detail string
}

func (i UnknowErr) Error() string {
	return fmt.Sprintf("get unknow Err with detail:%s", i.Detail)
}