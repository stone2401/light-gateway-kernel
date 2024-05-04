package load_balance

import (
	"sync"

	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

type RobinBalance struct {
	nodes    []string
	length   int
	mu       sync.RWMutex
	curIndex int
}

// 　创建一个轮询负载均衡
func NewRobinBalance() *RobinBalance {
	return &RobinBalance{
		nodes:  make([]string, 0),
		length: 0,
		mu:     sync.RWMutex{},
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

func (r *RobinBalance) RmNode(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, v := range r.nodes {
		if v == addr {
			r.nodes = append(r.nodes[:i], r.nodes[i+1:]...)
			r.length--
			return
		}
	}
}
