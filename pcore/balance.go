package pcore

import (
	"github.com/stone2401/light-gateway-kernel/pcore/load_balance"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

type LoadBalance string

const (
	// 轮询
	LoadBalanceRoundRobin LoadBalance = "round_robin"
	// 随机
	LoadBalanceRandom LoadBalance = "random"
	// 一致性哈希
	LoadBalanceConsistentHash LoadBalance = "consistent_hash"
	// 权重
	LoadBalanceWeight LoadBalance = "weight"
)

func NewLoadBalance(balance LoadBalance) sdk.Balance {
	switch balance {
	case LoadBalanceRoundRobin:
		return load_balance.NewRobinBalance()
	case LoadBalanceRandom:
		return load_balance.NewRandomBalance()
	case LoadBalanceConsistentHash:
		return load_balance.NewConsistentHashBanlance()
	case LoadBalanceWeight:
		return load_balance.NewWeightBalance()
	default:
		return load_balance.NewRandomBalance()
	}
}

func GetLoadBalance(balance int) LoadBalance {
	switch balance {
	case 0:
		return LoadBalanceRoundRobin
	case 1:
		return LoadBalanceRandom
	case 2:
		return LoadBalanceConsistentHash
	case 3:
		return LoadBalanceWeight
	default:
		return LoadBalanceRandom
	}
}
