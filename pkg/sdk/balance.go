package sdk

// 负载均衡器
type Balance interface {
	GetNode(token string) (string, error)
}
