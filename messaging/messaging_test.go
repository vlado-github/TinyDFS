package messaging

import (
	"testing"
	"os"
	"runtime"
	"github.com/google/uuid"
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
	var message = Message{Key:uuid.New(), Topic:"Test", Text:"Hello world!"}
	node.SendMessage(message)
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
