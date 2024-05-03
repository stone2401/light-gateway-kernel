package pcore

import (
	"sync"

	"github.com/stone2401/light-gateway-kernel/pkg/monitor"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

type WeightBalance struct {
	totalWeight int
	nodes       []*weightNode
	length      int
	mu          sync.RWMutex
	monitor     monitor.Monitor
}

type weightNode struct {
	node            string
	weight          int
	effectiveWeight int
}

// 权重轮巡
func NewWeightBalance() *WeightBalance {
	return &WeightBalance{
		totalWeight: 0,
		nodes:       make([]*weightNode, 0),
		length:      0,
		mu:          sync.RWMutex{},
		monitor:     nil,
	}
}

func (w *WeightBalance) AddNode(addr string, weight int) error {
	if addr == "" {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.nodes = append(w.nodes, &weightNode{
		node:            addr,
		weight:          weight,
		effectiveWeight: weight,
	})
	w.length++
	w.totalWeight += weight
	return nil
}

func (w *WeightBalance) GetNode(token string) (string, error) {
	if len(w.nodes) == 0 {
		return "", sdk.ErrorNotFoundNode
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	var base *weightNode
	for _, item := range w.nodes {
		item.effectiveWeight += item.weight

		if base == nil || item.effectiveWeight >= w.totalWeight {
			base = item
		}
	}
	base.effectiveWeight -= w.totalWeight
	return base.node, nil
}
