package persistance

import (
	"testing"
	"github.com/google/uuid"
	"os"
	"path"
	"path/filepath"
	"fmt"
)

var pathToDir string
var topic = "test_topic"
var fm FileManager

func TestMain(m *testing.M) {
	//setup
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Path of working directory not found.")
	}
	pathToDir = path
	fm = NewFileManager(pathToDir, topic)
	os.Exit(m.Run())
	//teardown
	cleanUp()
}

func TestFileManager_Write(t *testing.T) {
	cmd := Command{
		Key: uuid.New(),
		Text: "This is testing message for persistance: writing.",
	}
	err := fm.Write(cmd)
	if err != nil {
		t.Fail()
	}
}

func TestFileManager_Read(t *testing.T) {
	key := uuid.New()
	cmd := Command{
		Key: key,
		Text: "This is testing message for persistance: read line.",
	}
	err := fm.Write(cmd)
	if err != nil {
		t.Fail()
	}
	query := Query{
		Key: key,
	}
	data, err := fm.Read(query)
	if err != nil {
		t.Fail()
	}
	if data != cmd.Text {
		t.Fail()
	}
}

func TestFileManager_ReadFile(t *testing.T) {
	key := uuid.New()
	cmd := Command{
		Key: key,
		Text: "This is testing message for persistance: read file.",
	}
	err := fm.Write(cmd)
	if err != nil {
		t.Fail()
	}
	data, err := fm.ReadFile()
	if len(data) <= 0 {
		t.Fail()
	}
}

func cleanUp(){
	fm.Close()
	os.Remove(path.Join(pathToDir, topic))
}
