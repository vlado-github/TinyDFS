package persistance

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"io/ioutil"
	"strings"
	"errors"
	"path"
	"sync"
	"bytes"
)

type FileManager interface {
	Write(command Command) error
	Update(command Command) error
	Read(query Query) (string, error)
	ReadFile() ([]byte, error)
	Close() error
}

type fileManager struct{
	file *os.File
}

var mutex = &sync.Mutex{}
var pos int64

func NewFileManager(pathToDir string, topic string) FileManager{
	pathToFile := path.Clean(path.Join( pathToDir, topic))
	f, err := os.Create(pathToFile)
	if err != nil {
		fmt.Println("Persistance: Can not create a file.", err.Error())
		log.Fatal(err)
	}
	return &fileManager{
		file : f,
	}
}

func (fm *fileManager) Write(command Command) error {
	mutex.Lock()
	defer mutex.Unlock()
	w := bufio.NewWriter(fm.file)
	size, err := w.WriteString(command.Key.String() + ":" +command.Text+"\n")
	if err != nil {
		fmt.Println("Persistance: Write to file failed.", size, err.Error())
		log.Fatal(err)
	}
	w.Flush()
	return err
}

func (fm *fileManager) Update(command Command) error {
	mutex.Lock()
	defer mutex.Unlock()
	fileHandle, _ := os.OpenFile(fm.file.Name(), os.O_RDWR, 0777)
	defer fileHandle.Close()
	scanner := bufio.NewScanner(fileHandle)
	splitFunc := newSplitFunc()
	scanner.Split(splitFunc)
	for scanner.Scan() {
		oldBytesText := scanner.Bytes()
		text := string(oldBytesText)
		if strings.Contains(text, command.Key.String()) {
			newText := command.Key.String() + ":" + command.Text+"\n"
			newBytesText := []byte(newText)
			diff := len(oldBytesText) - len(newBytesText)
			// fits the message size - performs replacement
			if diff >= 0 {
				var additional string
				for i:=0; i<=diff; i++ {
					additional += " "
				}
				newBytesText = []byte(command.Key.String() + ":" + command.Text+additional+"\n")
				newBytesText = bytes.Replace(oldBytesText, oldBytesText, newBytesText, -1)
				fileHandle.Seek(0, 0)
				n, err := fileHandle.WriteAt(newBytesText, pos-int64(len(oldBytesText)))
				if err != nil{
					fmt.Println("Update failed. Writen bytes: ",n)
					fmt.Println(err)
					return err
				}
			} else {
				// new message larger than old one - delete and append
				var emptyLine string
				for i:=0; i<len(oldBytesText); i++ {
					emptyLine += " "
				}
				emptyLineBytes := []byte(emptyLine)
				fileHandle.Seek(0, 0)
				n, errDelete := fileHandle.WriteAt(emptyLineBytes, pos-int64(len(oldBytesText)))
				if errDelete != nil{
					fmt.Println("Update failed. Writen bytes: ",n)
					fmt.Println(errDelete)
					return errDelete
				}
				fileHandle.Seek(0, 2) //EOF
				n, errAppend := fileHandle.WriteString(newText)
				if errAppend != nil{
					fmt.Println("Update failed. Writen bytes: ",n)
					fmt.Println(errAppend)
					return errAppend
				}
			}


			return nil
		}
	}
	err := errors.New("Item not found")
	return err
}

func (fm *fileManager) Read(query Query) (string, error) {
	fileHandle, _ := os.Open(fm.file.Name())
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

func (fm *fileManager) ReadFile() ([]byte, error) {
	byteArray, err := ioutil.ReadFile(fm.file.Name())
	if err != nil {
		if err != nil {
			fmt.Println("Persistance: Can not create a file.", err.Error())
			log.Fatal(err)
		}
	}
	return byteArray, err
}

// Close file stream
func (fm *fileManager) Close() error{
	err := fm.file.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
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