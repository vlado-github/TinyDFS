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
	GetNextQueue() NetworkTuple
	SetQueueUnresponsive(ip string, queuePort string)
	GetQueueByIpAndPort(ip string, queuePort string) (NetworkTuple, int)
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

func (nr *networkregistry) SetQueueUnresponsive(ip string, queuePort string) {
	networkTuple, index := nr.GetQueueByIpAndPort(ip, queuePort)
	if networkTuple != nil {
		nr.NetworkTuples[index].SetIsAvailable(false)
	}
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
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetId() == id {
			return tuples[i], i
		}
	}
	return nil, index
}

func (nr *networkregistry) GetItemByRemoteAddPort(port string) (NetworkTuple, int) {
	var index int
	if len(nr.NetworkTuples) == 0 {
		return nil, index
	}
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetPort() == port {
			return tuples[i], i
		}
	}
	return nil, index
}

func (nr *networkregistry) GetQueueByIpAndPort(ip string, queuePort string) (NetworkTuple, int) {
	var index int
	if len(nr.NetworkTuples) == 0 {
		return nil, index
	}
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetQueuePort() == queuePort &&
			tuples[i].GetIP() == ip {
			return tuples[i], i
		}
	}
	return nil, index
}

func (nr *networkregistry) GetNextQueue() NetworkTuple {
	if len(nr.NetworkTuples) == 0 {
		return nil
	}
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetAvailableStatus() == true {
			return tuples[i]
		}
	}
	return nil
}

// ToByteArray converts to Json string
func (nr *networkregistry) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(nr.NetworkTuples)

	if err != nil {
		logging.AddInfo("NetworkRegistry ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// FromByteArray converts byte array to NetworkRegistry
func (nr *networkregistry) FromByteArray(data []byte) error {
	var tuples []networktuple
	err := json.Unmarshal(data, &tuples)
	if err != nil {
		logging.AddError("NetworkRegistry FromByteArray ", err.Error())
		return err
	}
	//todo: unmarshal doesn't work with interface,
	// check this out and remove this workaround code...
	if len(tuples) > 0 {
		nr.NetworkTuples = []NetworkTuple{}
	}
	for i := range tuples {
		nr.AddItem(NewNetworkTuple(tuples[i].GetId(),
			tuples[i].GetIP(),
			tuples[i].GetPort(),
			tuples[i].GetQueuePort()))
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
