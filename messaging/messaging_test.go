package messaging

import (
	"os"
	"persistance"
	"runtime"
	"testing"

	"github.com/google/uuid"
)

var defaultTestPath = "C://go_testing//"

var queueConnParams = ConnParams{
	"localhost", "3333", "tcp",
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
	var guid = uuid.New().String()
	var fm = persistance.NewFileManager(defaultTestPath+guid, "test")
	var node = NewNode(queueConnParams, fm)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
}

func TestSendingToQueue(t *testing.T) {
	var guid = uuid.New().String()
	var fm = persistance.NewFileManager(defaultTestPath+guid, "test")
	var node = NewNode(queueConnParams, fm)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
	var message = Message{Key: uuid.New(), Topic: "Test", Text: "Hello world!"}
	node.SendMessage(message)
}

func TestCloseNode(t *testing.T) {
	var guid = uuid.New().String()
	var fm = persistance.NewFileManager(defaultTestPath+guid, "test")
	var node = NewNode(queueConnParams, fm)
	err := node.Run()
	if err != nil {
		t.Fail()
	}
	close_err := node.CloseConn()
	if close_err != nil {
		t.Fail()
	}
}
