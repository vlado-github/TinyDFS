package consensus

import (
	"math/rand"
	"time"
)

// ElectionTimeoutType enum
type ElectionTimeoutType int

const (
	// ELECTIONTIMEOUT type
	ELECTIONTIMEOUT ElectionTimeoutType = iota
	// HEARTBEATTIMEOUT type
	HEARTBEATTIMEOUT ElectionTimeoutType = iota
)

const (
	/* production

	// ELLECTIONMIN type (in ms)
	ELLECTIONMIN int = 150
	// ELLECTIONMAX type (in ms)
	ELLECTIONMAX int = 300
	// HEARTBEATMAX type (in ms)
	HEARTBEATMAX int = 50

	*/

	/* testing */

	// ELLECTIONMIN type (in ms)
	ELLECTIONMIN int = 45000
	// ELLECTIONMAX type (in ms)
	ELLECTIONMAX int = 60000
	// HEARTBEATMAX type (in ms)
	HEARTBEATMAX int = 5000
)

// GetRandomElectionTimeout returns random value for time span in ms.
func GetRandomElectionTimeout() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(ELLECTIONMAX-ELLECTIONMIN) + ELLECTIONMIN
}
