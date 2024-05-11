package load_balance

import (
	"math/rand"
	"sync"

	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
	"go.uber.org/zap"
)

type RandomBalance struct {
	nodes  []string
	length int
	mu     sync.RWMutex
}

// 随机负载均衡
func NewRandomBalance() *RandomBalance {
	return &RandomBalance{
		nodes:  make([]string, 0),
		length: 0,
		mu:     sync.RWMutex{},
	}
}

func (r *RandomBalance) AddNode(addr string, weight int) error {
	if addr == "" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	zlog.Zlog().Info("add node", zap.String("addr", addr), zap.Int("weight", weight))
	r.nodes = append(r.nodes, addr)
	r.length++
	return nil
}

func (r *RandomBalance) GetNode(token string) (string, error) {
	if len(r.nodes) == 0 {
		return "", sdk.ErrorNotFoundNode
	}
	index := rand.Intn(r.length)
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes[index], nil
}

func (r *RandomBalance) RmNode(addr string) {
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
