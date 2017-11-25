package messaging

// HandlerType enum
type HandlerType int

const (
	// NODECONNCLOSED type
	NODECONNCLOSED HandlerType = iota
	// NODECONNOPENED type
	NODECONNOPENED HandlerType = iota
)
