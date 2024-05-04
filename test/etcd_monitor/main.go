package main

import (
	"sync"
	"time"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
)

var wg sync.WaitGroup

func main() {
	monitor := pcore.NewEtcdMonitor([]string{"127.0.0.1:2379"})
	balance := monitor.Register("test", pcore.LoadBalanceRandom)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			node, err := balance.GetNode("test")
			if err != nil {
				continue
			}
			zlog.Zlog().Info(node)
		}
	}()
	monitor.SyncWatch()
	// client, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	// if err != nil {
	// 	panic(err)
	// }
	// watch := client.Watch(context.Background(), "test", clientv3.WithPrefix())
	// for wresp := range watch {
	// 	for _, ev := range wresp.Events {
	// 		fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	// 	}
	// }
}
