package balance

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestRandomBalance(t *testing.T) {

	b := NewRandomBalance()
	b.AddNode("127.0.0.1:8080", 1)
	b.AddNode("127.0.0.1:8081", 2)
	b.AddNode("127.0.0.1:8082", 3)
	b.AddNode("127.0.0.1:8083", 4)
	b.AddNode("127.0.0.1:8084", 5)
	b.AddNode("127.0.0.1:8085", 6)
	b.AddNode("127.0.0.1:8086", 7)
	b.AddNode("127.0.0.1:8087", 8)
	b.AddNode("127.0.0.1:8088", 9)
	b.AddNode("127.0.0.1:8089", 10)

	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
}

func TestRobinBalance(t *testing.T) {

	b := NewRobinBalance()
	b.AddNode("127.0.0.1:8080", 1)
	b.AddNode("127.0.0.1:8081", 2)
	b.AddNode("127.0.0.1:8082", 3)
	b.AddNode("127.0.0.1:8083", 4)
	b.AddNode("127.0.0.1:8084", 5)
	b.AddNode("127.0.0.1:8085", 6)
	b.AddNode("127.0.0.1:8086", 7)
	b.AddNode("127.0.0.1:8087", 8)
	b.AddNode("127.0.0.1:8088", 9)
	b.AddNode("127.0.0.1:8089", 10)

	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
}

func TestWeightBalance(t *testing.T) {

	b := NewWeightBalance()
	b.AddNode("127.0.0.1:8080", 1)
	b.AddNode("127.0.0.1:8081", 1)
	b.AddNode("127.0.0.1:8082", 3)

	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
	fmt.Println(b.GetNode(""))
}

func TestConsistentHashBanlance(t *testing.T) {

	b := NewConsistentHashBanlance()
	b.AddNode("127.0.0.1:20000", 20)
	b.AddNode("127.0.0.1:10000", 20)
	b.AddNode("127.0.0.1:15000", 20)
	fmt.Println(b.GetNodeReplicas(context.Background(), "127.0.0.1:5000"))
	fmt.Println(b.GetNodeReplicas(context.Background(), "127.0.0.1:10000"))
	fmt.Println(b.GetNodeReplicas(context.Background(), "127.0.0.1:15000"))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
	fmt.Println(b.GetNode(uuid.NewString()))
}

func TestUUID(t *testing.T) {
	uid, _ := uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
	uid, _ = uuid.NewRandom()
	fmt.Println(uid.String())
}
