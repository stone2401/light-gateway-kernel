package balance

import "errors"

var (
	// 无可用节点
	ErrorNotFoundNode = errors.New("not found node")
	// 节点存在
	ErrorNodeExists = errors.New("node exists")
	// 节点不存在
	ErrorNodeNotExists = errors.New("node not exists")
	// 节点不可用
	ErrorNodeNotAvailable = errors.New("node not available")
)

type Balance interface {
	AddNode(addr string, weight int) error
	GetNode(token string) (string, error)
}
