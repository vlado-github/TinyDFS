package messaging

import (
	"encoding/json"
	"tinylogging"
)

type BaseMessagePayload interface {
	ToByteArray() ([]byte, error)
	ToPayload(data []byte) error
	GetNumOfNodes() int
}

type basemessagepayload struct {
	NumOfNodes int
}

func NewBaseMessagePayload(numOfNodes int) BaseMessagePayload {
	return &basemessagepayload{
		NumOfNodes: numOfNodes,
	}
}

// EmptyPayload returns new instance of the payload with temporary data set
func EmptyPayload() BaseMessagePayload {
	return &basemessagepayload{
		NumOfNodes: -1,
	}
}

func (payload *basemessagepayload) GetNumOfNodes() int {
	return payload.NumOfNodes
}

// ToByteArray converts to Json string
func (payload *basemessagepayload) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(payload)

	if err != nil {
		tinylogging.AddInfo("[Host] VoteMessagePayload ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// ToPayload converts byte array to BaseMessagePayload
func (payload *basemessagepayload) ToPayload(data []byte) error {
	err := json.Unmarshal(data, &payload)
	if err != nil {
		tinylogging.AddError("[Host] VoteMessagePayload ToPayload ", err.Error())
		return err
	}
	return nil
}
