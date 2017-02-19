package messaging

import (
	"testing"
	"os"
	"runtime"
)

var queueConnParams = connParams{
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

func TestConnectingToQueue(t *testing.T) {
	err := InitNode(queueConnParams)
	if err != nil {
		t.Fail()
	}
}

func TestSendingToQueue(t *testing.T) {
	err := InitNode(queueConnParams)
	if err == nil {
		SendEvent("Hello world!")
	} else {
		t.Fail()
	}
}
