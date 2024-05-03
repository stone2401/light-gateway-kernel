package pcore

import (
	"sync"

	"github.com/stone2401/light-gateway-kernel/pkg/monitor"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

type RobinBalance struct {
	nodes    []string
	length   int
	mu       sync.RWMutex
	monitor  monitor.Monitor
	curIndex int
}

// 　创建一个轮询负载均衡
func NewRobinBalance() *RobinBalance {
	return &RobinBalance{
		nodes:   make([]string, 0),
		length:  0,
		mu:      sync.RWMutex{},
		monitor: nil,
	}
}

func (r *RobinBalance) AddNode(addr string, weight int) error {
	if addr == "" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes = append(r.nodes, addr)
	r.length++
	return nil
}

func (r *RobinBalance) GetNode(token string) (string, error) {
	if len(r.nodes) == 0 {
		return "", sdk.ErrorNotFoundNode
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	defer func() {
		if r.curIndex >= r.length {
			r.curIndex = 0
		}
		r.curIndex = (r.curIndex + 1) % r.length
	}()
	return r.nodes[r.curIndex], nil
}