package sdk

// 负载均衡器
type Balance interface {
	GetNode(token string) (string, error)
	AddNode(addr string, weight int) error
	RmNode(addr string)
}
