package main

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		panic(err)
	}
	has := sdk.NewMurmurHasher()
	node := &pcore.NodeInfo{
		Ip:     "127.0.0.1",
		Weight: 1,
	}
	b, _ := json.Marshal(node)
	enc := has.Encrypt(string(b))
	client.Put(context.Background(), "test"+strconv.Itoa(int(enc)), string(b))
}
