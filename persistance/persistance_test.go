package persistance

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var pathToDir string
var topic = "test_topic"
var fm FileManager

func TestMain(m *testing.M) {
	//setup
	setUp()

	retCode := m.Run()

	//teardown
	cleanUp()

	os.Exit(retCode)
}

func TestFileManager_Write(t *testing.T) {
	cmd := Command{
		Key:   uuid.New(),
		Text:  "This is testing message for persistance: writing.",
		Topic: topic,
	}
	err := fm.Write(cmd)
	if err != nil {
		t.Fail()
	}
}

func TestFileManager_UpdateShort(t *testing.T) {
	keyForUpdate := uuid.New()
	cmd01 := Command{
		Key:   uuid.New(),
		Text:  "This is testing message for persistance: updatingShort01.",
		Topic: topic,
	}
	cmd02 := Command{
		Key:   keyForUpdate,
		Text:  "This is testing message for persistance: updatingShort02.",
		Topic: topic,
	}
	cmd03 := Command{
		Key:   uuid.New(),
		Text:  "This is testing message for persistance: updatingShort03.",
		Topic: topic,
	}
	err01 := fm.Write(cmd01)
	err02 := fm.Write(cmd02)
	err03 := fm.Write(cmd03)
	if err01 != nil || err02 != nil || err03 != nil {
		t.Fail()
	}
	updateCmd := Command{
		Key:   keyForUpdate,
		Text:  "Text is changed.",
		Topic: topic,
	}
	err := fm.Update(updateCmd)
	if err != nil {
		t.Fail()
	}

	query := Query{
		Key:   keyForUpdate,
		Topic: topic,
	}
	data, err := fm.Read(query)
	fmt.Println(data)
	fmt.Println(updateCmd.Text)
	if err != nil {
		t.Fail()
	}
	if !strings.Contains(data, updateCmd.Text) {
		t.Fail()
	}
}

func TestFileManager_UpdateLong(t *testing.T) {
	keyForUpdate := uuid.New()
	cmd01 := Command{
		Key:   uuid.New(),
		Text:  "This is testing message for persistance: updating01.",
		Topic: topic,
	}
	cmd02 := Command{
		Key:   keyForUpdate,
		Text:  "This is testing message for persistance: updating02.",
		Topic: topic,
	}
	cmd03 := Command{
		Key:   uuid.New(),
		Text:  "This is testing message for persistance: updating03.",
		Topic: topic,
	}
	err01 := fm.Write(cmd01)
	err02 := fm.Write(cmd02)
	err03 := fm.Write(cmd03)
	if err01 != nil || err02 != nil || err03 != nil {
		t.Fail()
	}
	updateCmd := Command{
		Key:   keyForUpdate,
		Text:  "Message is changed with this new content. Blah blah blah blah blah blah blah blah blah blah blah blah",
		Topic: topic,
	}
	err := fm.Update(updateCmd)
	if err != nil {
		t.Fail()
	}

	query := Query{
		Key:   keyForUpdate,
		Topic: topic,
	}
	data, err := fm.Read(query)
	if err != nil {
		t.Fail()
	}
	if data != updateCmd.Text {
		t.Fail()
	}
}

func TestFileManager_Read(t *testing.T) {
	key := uuid.New()
	cmd := Command{
		Key:   key,
		Text:  "This is testing message for persistance: read line.",
		Topic: topic,
	}
	err := fm.Write(cmd)
	if err != nil {
		t.Fail()
	}
	query := Query{
		Key:   key,
		Topic: cmd.Topic,
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
		Key:   key,
		Text:  "This is testing message for persistance: read file.",
		Topic: topic,
	}
	err := fm.Write(cmd)
	if err != nil {
		t.Fail()
	}
	data, err := fm.ReadFile(cmd.Topic)
	if len(data) <= 0 {
		t.Fail()
	}
}

func setUp() {
	pathToDir = "C://go_testing//"
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Path of working directory not found.")
	} else {
		pathToDir = path
	}
	fmt.Println(path)
	fm = NewFileManager(pathToDir)
}

func cleanUp() {
	pathForDelete := path.Clean(path.Join(pathToDir, topic))
	err := os.Remove(pathForDelete)
	if err != nil {
		fmt.Println(err.Error())
	}
}
