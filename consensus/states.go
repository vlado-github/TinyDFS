package consensus

// NodeState enum
type State int

const (
	// FOLLOWER type
	FOLLOWER State = iota
	// CANDIDATE type
	CANDIDATE State = iota
	// LEADER type
	LEADER State = iota
)
