package messaging

// MsgQueueHandlerFunc represent a type of callback function.
type MsgQueueHandlerFunc func(queue MessageQueue)

// NewMsgQueueHandlerFunc returns empty callback function.
func NewMsgQueueHandlerFunc() MsgQueueHandlerFunc {
	return func(queue MessageQueue) {}
}
