package sdk

type Balance interface {
	AddNode(addr string, weight int) error
	GetNode(token string) (string, error)
}
