package consensus

import (
	"math/rand"
	"time"
)

const (
	// ELLECTIONMIN type
	ELLECTIONMIN int = 150
	// ELLECTIONMAX type
	ELLECTIONMAX int = 300
	// HEARTBEATMAX type
	HEARTBEATMAX int = 50
)

func getRandomElelctionTimeout() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(ELLECTIONMAX-ELLECTIONMIN) + ELLECTIONMIN
}
