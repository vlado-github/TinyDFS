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
)

type FileManager interface {
	Write(command Command) error
	Read(query Query) (string, error)
	ReadFile() ([]byte, error)
	Close() error
}

type fileManager struct{
	file *os.File
}

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
	w := bufio.NewWriter(fm.file)
	size, err := w.WriteString(command.Key.String() + " : " +command.Text+"\n")
	if err != nil {
		fmt.Println("Persistance: Write to file failed.", size, err.Error())
		log.Fatal(err)
	}
	w.Flush()
	return err
}

func (fm *fileManager) Read(query Query) (string, error) {
	scanner := bufio.NewScanner(fm.file)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		if strings.Contains(text, query.Key.String()) {
			return text, nil
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