package messaging

import (
	"encoding/json"
	"tinylogging"
)

type BaseMessagePayload interface {
	ToByteArray() ([]byte, error)
	ToPayload(data []byte) error
	GetNumOfNodes() int
	GetIPs() []string
}

type basemessagepayload struct {
	NumOfNodes int
	IpAddresses []string
}

func NewBaseMessagePayload(numOfNodes int, ipAddresses []string) BaseMessagePayload {
	return &basemessagepayload{
		NumOfNodes: numOfNodes,
		IpAddresses: ipAddresses,
	}
}

// EmptyPayload returns new instance of the payload with temporary data set
func EmptyPayload() BaseMessagePayload {
	return &basemessagepayload{
		NumOfNodes: -1,
		IpAddresses: nil,
	}
}

func (payload *basemessagepayload) GetNumOfNodes() int {
	return payload.NumOfNodes
}

func (payload *basemessagepayload) GetIPs() []string {
	return payload.IpAddresses
}

// ToByteArray converts to Json string
func (payload *basemessagepayload) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(payload)

	if err != nil {
		tinylogging.AddInfo("[Host] BaseMessagePayload ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// ToPayload converts byte array to BaseMessagePayload
func (payload *basemessagepayload) ToPayload(data []byte) error {
	err := json.Unmarshal(data, &payload)
	if err != nil {
		tinylogging.AddError("[Host] BaseMessagePayload ToPayload ", err.Error())
		return err
	}
	return nil
}
