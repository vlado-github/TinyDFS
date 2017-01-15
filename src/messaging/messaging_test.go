package messaging

import (
	"testing"
	"os"
	"runtime"
)

var queueConnParams = connParams{
	"localhost","3333","tcp",
}

var nodeConnParams = connParams{
	"localhost","3333","tcp",
}

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	InitQueue(queueConnParams)
}

func TestParallelStart(t *testing.T) {
	InitNode(nodeConnParams)
}
