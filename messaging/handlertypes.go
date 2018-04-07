package messaging

// HandlerType enum
type HandlerType int

const (
	// NODECONNCLOSED type
	NODECONNCLOSED HandlerType = iota
	// NODECONNOPENED type
	NODECONNOPENED HandlerType = iota
	// MESSAGERECEIVED type
	MESSAGERECEIVED HandlerType = iota
	// QUEUECONNCLOSED type
	QUEUECONNCLOSED HandlerType = iota
)
