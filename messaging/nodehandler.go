package messaging

// NodeHandlerFunc represent a type of callback function.
type NodeHandlerFunc func(n Node)

// NewHandlerFunc returns empty callback function.
func NewHandlerFunc() NodeHandlerFunc {
	return func(n Node) {}
}
