package discovery

import (
"context"
"go.etcd.io/etcd/clientv3"
"log"
	"math/rand"
	"strings"
	"sync"
	"time"
"encoding/json"
"fmt"
etcd "github.com/micro/go-micro/v2/registry"
)

//the detail of service 定义服务结构，唯一id加ip地址
type ServiceInfo struct{
	Name   string
	ID string
	IP   string
}

type Master struct {
	Path         string
	Nodes         map[string] *Node
	Client         *clientv3.Client
	mutex sync.Mutex
}

//node is a client
type Node struct {
	State    bool
	Key        string
	Info    map[string]*ServiceInfo
}

func NewMaster(endpoints []string, watchPath string) (*Master,error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:    endpoints,
		DialTimeout: 3*time.Second,
	})

	if err != nil {
		log.Fatal(err)
		return nil,err
	}

	master := &Master {
		Path:    watchPath,
		Nodes:    make(map[string]*Node),
		Client: cli,
	}
	go master.getService(cli)
	go master.WatchNodes()
	return master,err
}

func (m *Master)getService(etcdclient *clientv3.Client)  {
	fmt.Println("get start")
	kv :=clientv3.NewKV(etcdclient)
	fmt.Println("get start1")
	rangeResp,err:=kv.Get(context.Background(),m.Path,clientv3.WithPrefix())
	fmt.Println("get start2")
	if err !=nil {
		panic(err)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _,kv :=range rangeResp.Kvs{
		info := &ServiceInfo{}
		fmt.Println("json start")



		info = GetServiceInfo(&kv.Value)
		if info ==nil{
			continue
		}
		m.AddNode(info.ID, info)
	}
}

func (m *Master) AddNode(key string, info *ServiceInfo) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	node, ok := m.Nodes[key]
	if !ok {
		node = &Node{
			State:    true,
			Key:    key,
			Info: map[string]*ServiceInfo{info.ID:info},
		}
		m.Nodes[node.Key] = node
	} else {
		node.Info[info.ID] = info
	}
}

func (m *Master) DeleteNode(Id string, info *ServiceInfo) {
	node, ok := m.Nodes[Id]
	if !ok {
		return
	} else {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		delete(node.Info, info.ID)
	}
}


func GetServiceInfo(values *[]byte) (*ServiceInfo) {
	info := &ServiceInfo{}
	var serverinfo etcd.Service
	err :=json.Unmarshal(*values,&serverinfo)
	fmt.Println("nodes num:",len(serverinfo.Nodes))
	//err := json.Unmarshal([]byte(kv.Value), info)
	if (err != nil) || (len(serverinfo.Nodes)<1){
		log.Println(err)
		//panic(err)
		return nil
	}
	info.Name=serverinfo.Name
	info.IP=serverinfo.Nodes[0].Address
	info.ID=serverinfo.Nodes[0].Id
	return info
}

func (m *Master) WatchNodes()  {
	fmt.Println(m.Path)
	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				//fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(&ev.Kv.Value)
				m.AddNode(string(ev.Kv.Key), info)
				fmt.Println(*info)
			case clientv3.EventTypeDelete:
				//fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(&ev.Kv.Value)
				m.DeleteNode(string(ev.Kv.Key), info)
			}
		}
	}
	fmt.Println("path over")
}

func (m *Master)GetNodeInfo(serviceName string)string  {
	m.mutex.Lock()
	m.mutex.Unlock()
	ipList :=make([]string,1)
	for _,node := range m.Nodes{
		for _,service := range node.Info{
			if 0 == strings.Compare(serviceName,service.Name) {
				ipList = append(ipList,service.IP)
			}
		}
	}
	if  0== len(ipList) {
		return ""
	}
	index := rand.Intn(len(ipList))
	return ipList[index]
}