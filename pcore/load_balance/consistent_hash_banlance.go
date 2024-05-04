package load_balance

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

type ConsistentHashBanlance struct {
	sdk.Ring
	sdk.Encryptor
}

// hash 一致性哈希
func NewConsistentHashBanlance() *ConsistentHashBanlance {
	return &ConsistentHashBanlance{
		Ring:      sdk.NewHashRing(),
		Encryptor: sdk.NewMurmurHasher(),
	}
}

func (c *ConsistentHashBanlance) AddNode(addr string, weight int) error {
	if len(c.Ring.GetNodeReplicas(context.Background(), addr)) != 0 {
		return sdk.ErrorNodeExists
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

func (c *ConsistentHashBanlance) RmNode(addr string) {
	c.Lock(context.Background(), addr, 10)
	defer c.Unlock(context.Background(), addr)
	c.Ring.Rem(context.Background(), c.Encryptor.Encrypt(addr), addr)
}
