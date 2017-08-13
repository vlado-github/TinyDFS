package messaging

import (
	"encoding/json"
	"fmt"
)

func decodeMessage(message *Message, dec *json.Decoder)(error){
	err := dec.Decode(&message)
	if err != nil {
		fmt.Println("Error: Decoding message.", err.Error())
	}
	return err
}

func encodeMessage(message *Message, enc *json.Encoder) (error){
	err := enc.Encode(message);
	if err != nil {
		fmt.Println("Error: Encoding message.", err.Error())
	}
	return err
}
