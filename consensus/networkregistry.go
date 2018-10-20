package consensus

import (
	"encoding/json"
	"sort"
	"tinylogging"
)

// NetworkRegistry is a collection of TinyDFS IPs and ports
type NetworkRegistry interface {
	ToByteArray() ([]byte, error)
	ToPayload(data []byte) error
	GetItems() []NetworkTuple
	GetItemByIP(ip string) []NetworkTuple
}

type networkregistry struct {
	networkTuples []NetworkTuple `json:"networkTuples"`
}

// NewNetworkRegistry creates a new instance of network registry
func NewNetworkRegistry(network []NetworkTuple) NetworkRegistry {
	for i := range network {
		println(">>>" + network[i].GetIP() + " " + network[i].GetPort())
	}
	return &networkregistry{
		networkTuples: network,
	}
}

// EmptyNetworkRegistry creates an empty instance of network registry
func EmptyNetworkRegistry() NetworkRegistry {
	return &networkregistry{
		networkTuples: []NetworkTuple{},
	}
}

func (nr *networkregistry) GetItems() []NetworkTuple {
	if len(nr.networkTuples) == 0 {
		return nr.networkTuples
	}
	sort.Slice(nr.networkTuples[:], func(i, j int) bool {
		return nr.networkTuples[i].GetPort() < nr.networkTuples[j].GetPort()
	})
	return nr.networkTuples
}

func (nr *networkregistry) GetItemByIP(ip string) []NetworkTuple {
	if len(nr.networkTuples) == 0 {
		return nr.networkTuples
	}
	var result []NetworkTuple
	tuples := nr.GetItems()
	for i := range tuples {
		if tuples[i].GetIP() == ip {
			result = append(result, tuples[i])
		}
	}
	return result
}

// ToByteArray converts to Json string
func (nr *networkregistry) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(nr)

	if err != nil {
		tinylogging.AddInfo("[Host] NetworkRegistry ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// ToPayload converts byte array to NetworkRegistry
func (nr *networkregistry) ToPayload(data []byte) error {
	err := json.Unmarshal(data, &nr)
	if err != nil {
		tinylogging.AddError("[Host] NetworkRegistry ToPayload ", err.Error())
		return err
	}
	return nil
}
