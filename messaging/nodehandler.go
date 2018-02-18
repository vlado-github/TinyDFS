package messaging

// NodeHandlerFunc represent a type of callback function.
type NodeHandlerFunc func()

// MessageHandlerFunc represent a type of callback function.
type MessageHandlerFunc func(message Message)

// NewNodeHandlerFunc returns empty callback function.
func NewNodeHandlerFunc() NodeHandlerFunc {
	return func() {}
}

// NewMessageHandlerFunc returns empty callback function.
func NewMessageHandlerFunc() MessageHandlerFunc {
	return func(message Message) {}
}
