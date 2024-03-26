package sdk

import "errors"

var (
	// 节点已存在
	ErrorNodeExists = errors.New("node exists")
	// 节点不存在
	ErrorNodeNotExists = errors.New("node not exists")
	// 节点不可用
	ErrorNodeNotAvailable = errors.New("node not available")
)
