package messaging

import (
	"encoding/json"
	"os"
	"path/filepath"
	"tinylogging"
)

func decodeMessage(message *Message, dec *json.Decoder) error {
	err := dec.Decode(&message)
	if err != nil {
		tinylogging.AddError("Error: Decoding message.", err.Error())
	}
	return err
}

func encodeMessage(message *Message, enc *json.Encoder) error {
	err := enc.Encode(message)
	if err != nil {
		tinylogging.AddError("Error: Encoding message.", err.Error())
	}
	return err
}

func getCurrentDirectory() string {
	var pathToDir = "C://tinydfs_storage//"
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		tinylogging.AddError("Path of working directory not found.", err.Error())
	} else {
		pathToDir = path
	}
	return pathToDir
}
