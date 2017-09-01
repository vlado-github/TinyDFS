package messaging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func decodeMessage(message *Message, dec *json.Decoder) error {
	err := dec.Decode(&message)
	if err != nil {
		fmt.Println("Error: Decoding message.", err.Error())
	}
	return err
}

func encodeMessage(message *Message, enc *json.Encoder) error {
	err := enc.Encode(message)
	if err != nil {
		fmt.Println("Error: Encoding message.", err.Error())
	}
	return err
}

func getCurrentDirectory() string {
	var pathToDir = "C://tinydfs_storage//"
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Path of working directory not found.")
	} else {
		pathToDir = path
	}
	return pathToDir
}
