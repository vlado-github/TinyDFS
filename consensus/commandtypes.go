package consensus

// CommandType enum
type CommandType int

const (
	// REQUESTVOTE type
	REQUESTVOTE CommandType = iota
	// APPENDENTRY type
	APPENDENTRY CommandType = iota
	// HEARTBEAT type
	HEARTBEAT CommandType = iota
)
