package persistance

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"logging"
	"os"
	"path"
	"strings"
	"sync"
)

type FileManager interface {
	Write(command Command) error
	Update(command Command) error
	Read(query Query) (string, error)
	ReadFile(topic string) ([]byte, error)
	//Close() error
}

type fileManager struct {
	pathToDir string
}

var mutex = &sync.Mutex{}
var pos int64

// Creates instance of FileManager
func NewFileManager(pathDir string) FileManager {
	pathToDir := path.Clean(path.Join(pathDir))
	err := os.MkdirAll(pathToDir, os.ModePerm)
	if err != nil {
		logging.AddError("Persistance: Can not create a directory.", err.Error())
	}

	return &fileManager{
		pathToDir: pathToDir,
	}
}

func (fm *fileManager) Write(command Command) error {
	mutex.Lock()
	defer mutex.Unlock()
	pathToFile := path.Clean(path.Join(fm.pathToDir, command.Topic))
	f, err := createOrAppendFile(pathToFile)
	defer f.Close()
	w := bufio.NewWriter(f)
	size, err := fmt.Fprint(w, command.Key.String()+":"+command.Text)
	if err != nil {
		logging.AddError("Persistance: Write to file failed.", size, err.Error())
	}
	w.Flush()
	return err
}

func (fm *fileManager) Update(command Command) error {
	mutex.Lock()
	defer mutex.Unlock()
	pathToFile := path.Clean(path.Join(fm.pathToDir, command.Topic))
	fileHandle, _ := os.OpenFile(pathToFile, os.O_RDWR, 0777)
	defer fileHandle.Close()
	scanner := bufio.NewScanner(fileHandle)
	splitFunc := newSplitFunc()
	scanner.Split(splitFunc)
	for scanner.Scan() {
		oldBytesText := scanner.Bytes()
		text := string(oldBytesText)
		if strings.Contains(text, command.Key.String()) {
			newText := command.Key.String() + ":" + command.Text + "\n"
			newBytesText := []byte(newText)
			diff := len(oldBytesText) - len(newBytesText)
			// fits the message size - performs replacement
			if diff >= 0 {
				var additional string
				for i := 0; i <= diff; i++ {
					additional += " "
				}
				newBytesText = []byte(command.Key.String() + ":" + command.Text + additional + "\n")
				newBytesText = bytes.Replace(oldBytesText, oldBytesText, newBytesText, -1)
				fileHandle.Seek(0, 0)
				n, err := fileHandle.WriteAt(newBytesText, pos-int64(len(oldBytesText)))
				if err != nil {
					logging.AddError("Update failed. Writen bytes: ", n, err.Error())
					return err
				}
			} else {
				// new message larger than old one - delete and append
				var emptyLine string
				for i := 0; i < len(oldBytesText); i++ {
					emptyLine += " "
				}
				emptyLineBytes := []byte(emptyLine)
				fileHandle.Seek(0, 0)
				n, errDelete := fileHandle.WriteAt(emptyLineBytes, pos-int64(len(oldBytesText)))
				if errDelete != nil {
					logging.AddError("Update failed. Writen bytes: ", n, errDelete.Error())
					return errDelete
				}
				fileHandle.Seek(0, 2) //EOF
				n, errAppend := fileHandle.WriteString(newText)
				if errAppend != nil {
					logging.AddError("Update failed. Writen bytes: ", n, errAppend.Error())
					return errAppend
				}
			}

			return nil
		}
	}
	err := errors.New("Item not found")
	logging.AddError(err.Error())
	return err
}

func (fm *fileManager) Read(query Query) (string, error) {
	pathToFile := path.Clean(path.Join(fm.pathToDir, query.Topic))
	fileHandle, _ := os.Open(pathToFile)
	defer fileHandle.Close()
	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, query.Key.String()) {
			result := strings.SplitN(text, ":", 2)
			if len(result) == 2 {
				return result[1], nil
			}
		}
	}
	err := errors.New("Item not found")
	return "", err
}

func (fm *fileManager) ReadFile(topic string) ([]byte, error) {
	pathToFile := path.Clean(path.Join(fm.pathToDir, topic))
	byteArray, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		if err != nil {
			logging.AddError("Persistance: Can not create a file.", err.Error())
		}
	}
	return byteArray, err
}

func newSplitFunc() bufio.SplitFunc {
	var n int64
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if err == nil && (atEOF || advance > len(token)) {
			// We found the end of a line
			pos = n + int64(len(token))
			token = append(([]byte)(nil), token...)
		}
		n += int64(advance)
		return
	}
}

func createOrAppendFile(pathToFile string) (*os.File, error) {
	_, err := os.Stat(pathToFile)
	if os.IsNotExist(err) {
		f, err := os.Create(pathToFile)
		if err != nil {
			logging.AddError("Persistance: Can not create a file.", err.Error())
		}
		return f, err
	}
	f, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		logging.AddError("Persistance: Can not create a file.", err.Error())
	}

	return f, err
}
