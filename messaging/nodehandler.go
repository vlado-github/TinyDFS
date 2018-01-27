package messaging

// NodeHandlerFunc represent a type of callback function.
type NodeHandlerFunc func()

// NewNodeHandlerFunc returns empty callback function.
func NewNodeHandlerFunc() NodeHandlerFunc {
	return func() {}
}
