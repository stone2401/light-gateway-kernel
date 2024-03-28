package sdk

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone2401/light-gateway-kernel/pkg/monitor"
)

type ConsistentHashBanlance struct {
	Ring
	Encryptor
	monitor monitor.Monitor
}

func NewConsistentHashBanlance() *ConsistentHashBanlance {
	return &ConsistentHashBanlance{
		Ring:      NewHashRing(),
		Encryptor: NewMurmurHasher(),
		monitor:   nil,
	}
}

func (c *ConsistentHashBanlance) AddNode(addr string, weight int) error {
	if len(c.Ring.GetNodeReplicas(context.Background(), addr)) != 0 {
		return ErrorNodeExists
	}
	if weight == 0 {
		weight = 1
	} else if weight > 20 {
		weight = 20
	}
	c.Lock(context.Background(), addr, 10)
	defer c.Unlock(context.Background(), addr)
	replicas := make([]string, weight)
	for i := 0; i < weight; i++ {
		tmp := addr + "-" + uuid.NewString()
		replicas[i] = tmp
		score := c.Encryptor.Encrypt(tmp)
		c.Ring.Add(context.Background(), score, addr)
	}
	c.Ring.UpdateNodeReplicas(context.Background(), addr, replicas)
	return nil
}

func (c *ConsistentHashBanlance) GetNode(token string) (string, error) {
	score := c.Encryptor.Encrypt(token)

	addr, err := c.Ring.Ceil(context.Background(), score)
	if err != nil {
		return "", err
	}
	return addr, nil
}
