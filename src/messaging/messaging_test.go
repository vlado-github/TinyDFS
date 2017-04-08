package messaging

import (
	"testing"
	"os"
	"runtime"
)

var queueConnParams = ConnParams{
	"localhost","3333","tcp",
}

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	var queue = NewQueue(queueConnParams)
	queue.Run()
}

func TestConnectingToQueue(t *testing.T) {
	var node = NewNode(queueConnParams)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
}

func TestSendingToQueue(t *testing.T) {
	var node = NewNode(queueConnParams)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
	node.SendMessage("Hello world!")
}

func TestCloseNode(t *testing.T) {
	var node = NewNode(queueConnParams)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
	close_err := node.CloseConn()
	if close_err != nil {
		t.Fail()
	}
}
