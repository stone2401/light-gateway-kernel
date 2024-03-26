package balance

import (
	"sync"

	"github.com/stone2401/light-gateway-kernel/pkg/monitor"
)

type RobinBalance struct {
	nodes    []string
	length   int
	mu       sync.RWMutex
	monitor  monitor.Monitor
	curIndex int
}

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
		return "", ErrorNotFoundNode
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
