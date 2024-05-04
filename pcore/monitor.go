package pcore

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Monitor interface {
	// 监听打开
	Watch()
	// 监听关闭
	UnWatch()
	// 同步监听
	SyncWatch()
	// 注册
	Register(name string, balance LoadBalance) sdk.Balance
	// 注销
	UnRegister(name string)
}

type EtcdMonitor struct {
	endpoints   []string
	client      *clientv3.Client
	mu          sync.Mutex
	watchGroups map[string]sdk.Balance
	// cases 初始化第一位是stop信道
	cases     []reflect.SelectCase
	chanInfos []ChanInfo
	stopChan  chan bool
	// node map, 如果是修改，则应该先删除再添加
	nodeMap map[string]*NodeInfo
}

type ChanInfo struct {
	name string
	ch   any
}

type NodeInfo struct {
	Ip     string `json:"ip"`
	Weight int    `json:"weight"`
}

func (node *NodeInfo) Unmarshal(data []byte) error {
	return json.Unmarshal(data, node)
}

func NewEtcdMonitor(endpoints []string) Monitor {
	client, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		panic(err)
	}
	chanStop := make(chan bool)
	cases := []reflect.SelectCase{{Chan: reflect.ValueOf(chanStop), Dir: reflect.SelectRecv}}
	chanInfo := []ChanInfo{{name: "stop", ch: chanStop}}
	return &EtcdMonitor{
		endpoints:   endpoints,
		client:      client,
		mu:          sync.Mutex{},
		watchGroups: map[string]sdk.Balance{},
		cases:       cases,
		chanInfos:   chanInfo,
		stopChan:    chanStop,
		nodeMap:     map[string]*NodeInfo{},
	}
}

func (m *EtcdMonitor) SyncWatch() {
	for {
		// 上锁，线程安全
		m.mu.Lock()
		// 监听全部channel，等待响应
		chosen, recv, ok := reflect.Select(m.cases)
		zlog.Zlog().Info("select", zap.Any("chosen", chosen), zap.Any("recv", recv), zap.Any("ok", ok))
		// 如果是0 也就是stop channel，那么直接返回，如果ok是false，说明channel关闭了，直接退出
		if chosen == 0 && !ok {
			return
		}
		if chosen == 0 {
			// 等待 1 毫秒，防止锁争夺导致的不安全问题
			m.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			continue
		}
		// 响应，处理监听结果
		response, ok := recv.Interface().(clientv3.WatchResponse)
		if !ok {
			m.mu.Unlock()
			continue
		}
		balance := m.watchGroups[m.chanInfos[chosen].name]
		for _, v := range response.Events {
			// 如果是PUT事件，则可能是新增或者修改
			// 通过nodeMap判断, key 是否存在，如果存在则是修改，否则是新增
			// 修改需要先删除老的node，再添加新的node
			// 如果是DELETE事件，一定是删除
			if v.Type == clientv3.EventTypePut {
				// 解析，如果解析失败直接忽略
				node := &NodeInfo{}
				if err := node.Unmarshal(v.Kv.Value); err != nil {
					continue
				}
				if tmp, ok := m.nodeMap[string(v.Kv.Key)]; ok {
					balance.RmNode(tmp.Ip)
				}
				balance.AddNode(node.Ip, node.Weight)
				m.nodeMap[string(v.Kv.Key)] = node
				zlog.Zlog().Info("watch", zap.Any("type", v.Type), zap.Any("key", string(v.Kv.Key)), zap.Any("node", node))
			} else if v.Type == clientv3.EventTypeDelete {
				if tmp, ok := m.nodeMap[string(v.Kv.Key)]; ok {
					balance.RmNode(tmp.Ip)
				}
				delete(m.nodeMap, string(v.Kv.Key))
				fmt.Printf("delete %#v\n", v.Kv.String())
				zlog.Zlog().Info("watch", zap.Any("type", v.Type), zap.Any("key", string(v.Kv.Key)))
			}
		}
		m.mu.Unlock()
	}
}

func (m *EtcdMonitor) Watch() {
	go func() {
		m.SyncWatch()
	}()
}

func (m *EtcdMonitor) UnWatch() {
	close(m.stopChan)
}

// 注册
func (m *EtcdMonitor) Register(name string, balance LoadBalance) sdk.Balance {
	zlog.Zlog().Info("register", zap.String("name", name))
	// 1. 先停止上一次监听，之后上锁
	fmt.Println("try lock")
	if !m.mu.TryLock() {
		// 加锁失败
		m.stopChan <- true
		m.mu.Lock()
	}
	defer m.mu.Unlock()
	// 2. 如果此路由以及注册，则直接返回，反之则初始化新的负载均衡器
	if balance, ok := m.watchGroups[name]; ok {
		return balance
	}
	loadBalance := NewLoadBalance(balance)
	// 3. 开始监听，并注册到监听chan中
	watchChan := m.client.Watch(context.Background(), name, clientv3.WithPrefix())
	m.cases = append(m.cases, reflect.SelectCase{Chan: reflect.ValueOf(watchChan), Dir: reflect.SelectRecv})
	m.chanInfos = append(m.chanInfos, ChanInfo{name: name, ch: watchChan})
	m.watchGroups[name] = loadBalance
	// 4. 获取现有节点信息
	resp, err := m.client.Get(context.Background(), name, clientv3.WithPrefix())
	zlog.Zlog().Info("get", zap.Any("resp", resp), zap.Any("err", err))
	if err != nil {
		return nil
	}
	for _, v := range resp.Kvs {
		node := &NodeInfo{}
		if err := node.Unmarshal(v.Value); err != nil {
			continue
		}
		m.nodeMap[string(v.Key)] = node
		zlog.Zlog().Info("watch", zap.Any("key", string(v.Key)), zap.Any("node", node))
		loadBalance.AddNode(node.Ip, node.Weight)
	}
	return loadBalance
}

func (m *EtcdMonitor) UnRegister(name string) {
	m.stopChan <- true
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.watchGroups, name)
	for i, v := range m.chanInfos {
		if v.name == name {
			m.cases = append(m.cases[:i], m.cases[i+1:]...)
			m.chanInfos = append(m.chanInfos[:i], m.chanInfos[i+1:]...)
			break
		}
	}
}
