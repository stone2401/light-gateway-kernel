package sdk

import (
	"context"
	"sync"
	"time"
)

// ring 接口
type Ring interface {
	// 锁
	Lock(ctx context.Context, owner string, expireTime int64) bool
	// 解锁
	Unlock(ctx context.Context, owner string)
	// 添加节点
	Add(ctx context.Context, score int32, node string) error
	// 删除节点
	Rem(ctx context.Context, score int32, node string) error
	// 获取最后一个节点
	Last(ctx context.Context) string
	// 获取第一个节点
	First(ctx context.Context) string
	// 获取上一个节点
	Floor(ctx context.Context, score int32) (string, error)
	// 获取下一个节点
	Ceil(ctx context.Context, score int32) (string, error)
	// 更新节点映射
	UpdateNodeReplicas(ctx context.Context, node string, replicas []string)
	// 获取节点映射
	GetNodeReplicas(ctx context.Context, node string) []string
}

type LockMutex struct {
	mu         sync.Mutex
	owner      string
	expireTime time.Time
	flag       bool
}

func NewLockMutex() *LockMutex {
	return &LockMutex{
		mu:         sync.Mutex{},
		owner:      "",
		expireTime: time.Now(),
		flag:       false,
	}
}

func (l *LockMutex) Lock(ctx context.Context, owner string, expireTime int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if !l.flag || now.Before(l.expireTime) {
		l.owner = owner
		l.expireTime = now.Add(time.Duration(expireTime) * time.Second)
		l.flag = true
		return true
	}
	return false
}

func (l *LockMutex) Unlock(ctx context.Context, owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag && l.owner == owner {
		l.owner = ""
		l.expireTime = time.Now()
		l.flag = false
	}
}

type HashRing struct {
	// 锁
	*LockMutex
	// 虚拟节点
	virtualNodes []*VirtualNode
	// 虚拟节点 与真实节点映射
	virtualNodeReplicas map[string][]string
}

// 虚拟节点
type VirtualNode struct {
	nodes []string
	score int32
}

// 创建hash环
func NewHashRing() *HashRing {
	return &HashRing{
		LockMutex:           NewLockMutex(),
		virtualNodes:        make([]*VirtualNode, 0),
		virtualNodeReplicas: make(map[string][]string),
	}
}

func (h *HashRing) Add(ctx context.Context, score int32, node string) error {
	// 获取节点，如果已经存在则返回
	targetNode, _ := h.getScore(ctx, score)
	if targetNode != nil {
		for _, item := range targetNode.nodes {
			if item == node {
				return nil
			}
		}
		// 如果score相同，直接插入
		targetNode.nodes = append(targetNode.nodes, node)
		return nil
	}
	// 先处理特殊情况：如果没有虚拟节点，直接插入
	if len(h.virtualNodes) == 0 {
		h.virtualNodes = append(h.virtualNodes, &VirtualNode{
			nodes: []string{node},
			score: score,
		})
		return nil
	}
	// 比 0 小，插入到第一个
	if score < h.virtualNodes[0].score {
		h.virtualNodes = append([]*VirtualNode{{
			nodes: []string{node},
			score: score,
		}}, h.virtualNodes...)
		return nil
	}
	// 比最后一个大，插入到最后
	if score > h.virtualNodes[len(h.virtualNodes)-1].score {
		h.virtualNodes = append(h.virtualNodes, &VirtualNode{
			nodes: []string{node},
			score: score,
		})
		return nil
	}
	// 不相同，寻找插入位置
	for i := 0; i < len(h.virtualNodes); i++ {
		if h.virtualNodes[i].score > score {
			h.virtualNodes = append(h.virtualNodes[:i], append([]*VirtualNode{{
				nodes: []string{node},
				score: score,
			}}, h.virtualNodes[i:]...)...)
			break
		}
	}
	return nil
}

func (h *HashRing) Rem(ctx context.Context, score int32, node string) error {
	targetNode, index := h.getScore(ctx, score)
	if targetNode == nil {
		return ErrorNodeNotExists
	}
	// 1. 如果只有一个节点，删除
	if len(targetNode.nodes) == 1 {
		h.virtualNodes = append(h.virtualNodes[:index], h.virtualNodes[index+1:]...)
		return nil
	}
	// 2. 如果不止一个节点，删除
	for i := 0; i < len(targetNode.nodes); i++ {
		if targetNode.nodes[i] == node {
			targetNode.nodes = append(targetNode.nodes[:i], targetNode.nodes[i+1:]...)
			break
		}
	}
	return nil
}

// 寻找下一个节点
func (h *HashRing) Ceil(ctx context.Context, score int32) (string, error) {
	targetNode := h.ceiling(ctx, score)
	if targetNode == nil {
		return "", ErrorNodeNotAvailable
	}
	return targetNode.nodes[0], nil
}

// 寻找上一个节点
func (h *HashRing) Floor(ctx context.Context, score int32) (string, error) {
	targetNode := h.floor(ctx, score)
	if targetNode == nil {
		return "", ErrorNodeNotExists
	}
	return targetNode.nodes[0], nil
}

// 第一个节点
func (h *HashRing) First(ctx context.Context) string {
	if len(h.virtualNodes) == 0 {
		return ""
	}
	return h.virtualNodes[0].nodes[0]
}

// 最后一个节点
func (h *HashRing) Last(ctx context.Context) string {
	if len(h.virtualNodes) == 0 {
		return ""
	}
	return h.virtualNodes[len(h.virtualNodes)-1].nodes[0]
}

// 修改节点映射
func (h *HashRing) UpdateNodeReplicas(ctx context.Context, node string, replicas []string) {
	if h.virtualNodeReplicas == nil {
		h.virtualNodeReplicas = make(map[string][]string)
	}
	if _, ok := h.virtualNodeReplicas[node]; !ok {
		h.virtualNodeReplicas[node] = replicas
	} else {
		h.virtualNodeReplicas[node] = append(h.virtualNodeReplicas[node], replicas...)
	}
}

// 获取节点映射
func (h *HashRing) GetNodeReplicas(ctx context.Context, node string) []string {
	return h.virtualNodeReplicas[node]
}
func (h *HashRing) getScore(_ context.Context, score int32) (*VirtualNode, int) {
	if len(h.virtualNodes) == 0 {
		return nil, 0
	}
	// score的节点
	pre, end := 0, len(h.virtualNodes)-1
	for pre <= end {
		mid := (pre + end) / 2
		if h.virtualNodes[mid].score == score {
			return h.virtualNodes[mid], mid
		} else if h.virtualNodes[mid].score < score {
			pre = mid + 1
		} else if h.virtualNodes[mid].score > score {
			end = mid - 1
		}
	}
	return nil, 0
}

// 获取最接近 score的节点
func (h *HashRing) ceiling(_ context.Context, score int32) *VirtualNode {
	if len(h.virtualNodes) == 0 {
		return nil
	}
	pre, end := 0, len(h.virtualNodes)-1
	for pre <= end {
		mid := (pre + end) / 2
		if mid == 0 && h.virtualNodes[mid].score >= score {
			return h.virtualNodes[mid]
		} else if mid == len(h.virtualNodes)-1 && h.virtualNodes[mid].score < score {
			return h.virtualNodes[0]
		} else if h.virtualNodes[mid].score < score && h.virtualNodes[mid+1].score >= score {
			return h.virtualNodes[mid+1]
		} else if h.virtualNodes[mid].score > score {
			end = mid - 1
		} else if h.virtualNodes[mid+1].score < score {
			pre = mid + 1
		}
	}
	return nil
}

func (h *HashRing) floor(_ context.Context, score int32) *VirtualNode {
	if len(h.virtualNodes) == 0 {
		return nil
	}
	pre, end := 0, len(h.virtualNodes)-1
	for pre <= end {
		mid := (pre + end) / 2
		if mid == 0 && h.virtualNodes[mid].score > score {
			return h.virtualNodes[len(h.virtualNodes)-1]
		} else if mid == len(h.virtualNodes)-1 && h.virtualNodes[mid].score < score {
			return h.virtualNodes[len(h.virtualNodes)-1]
		} else if h.virtualNodes[mid].score < score && h.virtualNodes[mid+1].score >= score {
			return h.virtualNodes[mid]
		} else if h.virtualNodes[mid].score > score {
			end = mid - 1
		} else if h.virtualNodes[mid+1].score < score {
			pre = mid + 1
		}
	}
	return nil
}
