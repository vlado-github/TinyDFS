package messaging

// HandlerType enum
type HandlerType int

const (
	// NODECONNCLOSED type
	NODECONNCLOSED HandlerType = iota
	// NODECONNOPENED type
	NODECONNOPENED HandlerType = iota
	// NETWORKCHANGED type
	NETWORKCHANGED HandlerType = iota
)
