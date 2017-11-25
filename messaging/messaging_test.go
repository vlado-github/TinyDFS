package messaging

import (
	"os"
	"runtime"
	"testing"

	"github.com/google/uuid"
)

var queueConnParams = ConnParams{
	"localhost", "3333", "tcp",
}

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	// setup
	var masterNode = NewNode(queueConnParams, true)
	go masterNode.Run()
	go func() {
		retCode := m.Run()
		//cleanup
		masterNode.CloseConn()
		os.Exit(retCode)
	}()
}

func TestConnectingToQueue(t *testing.T) {
	var node = NewNode(queueConnParams, false)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
}

func TestSendingToQueue(t *testing.T) {
	var node = NewNode(queueConnParams, false)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
	var message = Message{Key: uuid.New(), Topic: "Test", Text: "Hello world!"}
	node.SendMessage(message)
}

func TestCloseNode(t *testing.T) {
	var node = NewNode(queueConnParams, false)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
	closeerr := node.CloseConn()
	if closeerr != nil {
		t.Fail()
	}
}
