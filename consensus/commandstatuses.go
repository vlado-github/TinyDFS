package consensus

// CommandType enum
type CommandStatuses int

const (
	// MAJORITYACK type
	MAJORITYACK CommandStatuses = iota
	// MAJORITYNACK type
	MAJORITYNACK CommandStatuses = iota
	// ACKPENDING type
	ACKPENDING CommandStatuses = iota
)
