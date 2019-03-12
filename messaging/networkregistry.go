package messaging

import (
	"encoding/json"
	"logging"
	"sort"
)

// NetworkRegistry is a collection of TinyDFS IPs and ports
type NetworkRegistry interface {
	ToByteArray() ([]byte, error)
	FromByteArray(data []byte) error
	ToString() (string, error)
	GetItems() []NetworkTuple
	GetItemById(id string) (NetworkTuple, int)
	GetItemByRemoteAddPort(port string) (NetworkTuple, int)
	AddItem(networkTuple NetworkTuple)
	RemoveItem(index int)
}

type networkregistry struct {
	NetworkTuples []NetworkTuple `json:"NetworkTuples"`
}

// NewNetworkRegistry creates a new instance of network registry
func NewNetworkRegistry() NetworkRegistry {
	return &networkregistry{
		NetworkTuples: []NetworkTuple{},
	}
}

func (nr *networkregistry) AddItem(networkTuple NetworkTuple) {
	if networkTuple != nil {
		nr.NetworkTuples = append(nr.NetworkTuples, networkTuple)
	}
}

func (nr *networkregistry) RemoveItem(index int) {
	nr.NetworkTuples = append(nr.NetworkTuples[:index], nr.NetworkTuples[index+1:]...)
}

func (nr *networkregistry) GetItems() []NetworkTuple {
	sort.Slice(nr.NetworkTuples[:], func(i, j int) bool {
		return nr.NetworkTuples[i].GetPort() < nr.NetworkTuples[j].GetPort()
	})
	return nr.NetworkTuples
}

func (nr *networkregistry) GetItemById(id string) (NetworkTuple, int) {
	var index int
	if len(nr.NetworkTuples) == 0 {
		return nil, index
	}
	var result NetworkTuple
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetId() == id {
			result = tuples[i]
			index = i
		}
	}
	return result, index
}

func (nr *networkregistry) GetItemByRemoteAddPort(port string) (NetworkTuple, int) {
	var index int
	if len(nr.NetworkTuples) == 0 {
		return nil, index
	}
	var result NetworkTuple
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetPort() == port {
			result = tuples[i]
			index = i
		}
	}
	return result, index
}

// ToByteArray converts to Json string
func (nr *networkregistry) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(nr)

	if err != nil {
		logging.AddInfo("NetworkRegistry ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// FromByteArray converts byte array to NetworkRegistry
func (nr *networkregistry) FromByteArray(data []byte) error {
	err := json.Unmarshal(data, &nr)
	if err != nil {
		logging.AddError("NetworkRegistry FromByteArray ", err.Error())
		return err
	}
	return nil
}

// ToString converts to string
func (nr *networkregistry) ToString() (string, error) {
	result, err := json.Marshal(nr.NetworkTuples)
	if err != nil {
		logging.AddInfo("NetworkRegistry ToString", err.Error())
		return "", err
	}
	return string(result), nil
}
