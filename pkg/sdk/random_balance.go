package sdk

import (
	"math/rand"
	"sync"

	"github.com/stone2401/light-gateway-kernel/pkg/monitor"
)

type RandomBalance struct {
	nodes  []string
	length int
	mu     sync.RWMutex
	// 监听者
	monitor monitor.Monitor
}

func NewRandomBalance() *RandomBalance {
	return &RandomBalance{
		nodes:   make([]string, 0),
		length:  0,
		mu:      sync.RWMutex{},
		monitor: nil,
	}
}

func (r *RandomBalance) AddNode(addr string, weight int) error {
	if addr == "" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes = append(r.nodes, addr)
	r.length++
	return nil
}

func (r *RandomBalance) GetNode(token string) (string, error) {
	if len(r.nodes) == 0 {
		return "", ErrorNotFoundNode
	}
	index := rand.Intn(r.length)
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes[index], nil
}
