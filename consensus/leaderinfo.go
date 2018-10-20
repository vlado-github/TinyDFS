package consensus

import (
	"encoding/json"
	"tinylogging"
)

// LeaderInfoMessagePayload represents additional message info for
// leader election
type LeaderInfo interface {
	ToByteArray() ([]byte, error)
	ToPayload(data []byte) error
	GetIP() string
	GetPort() string
	GetTerm() int
	GetElectionID() int
	GetNodeID() string
}

type leaderinfomessagepayload struct {
	IP         string `json:"ip"`
	Port       string `json:"port"`
	Term       int    `json:"term"`
	ElectionID int    `json:"election_id"`
	NodeID     string `json:"node_id"`
}

// NewLeaderInfo creates new instance of the message payload for vote
func NewLeaderInfo(ip string, port string, term int, electionID int, nodeID string) LeaderInfo {
	return &leaderinfomessagepayload{
		IP:         ip,
		Port:       port,
		Term:       term,
		ElectionID: electionID,
		NodeID:     nodeID,
	}
}

// EmptyLeaderInfo returns new instance of the message with temporary data set
func EmptyLeaderInfo() LeaderInfo {
	return &leaderinfomessagepayload{
		IP:         "",
		Port:       "",
		Term:       -1,
		ElectionID: -1,
		NodeID:     "",
	}
}

func (payload *leaderinfomessagepayload) GetIP() string {
	return payload.IP
}

func (payload *leaderinfomessagepayload) GetPort() string {
	return payload.Port
}

func (payload *leaderinfomessagepayload) GetTerm() int {
	return payload.Term
}

func (payload *leaderinfomessagepayload) GetElectionID() int {
	return payload.ElectionID
}

func (payload *leaderinfomessagepayload) GetNodeID() string {
	return payload.NodeID
}

// ToByteArray converts to Json string
func (payload *leaderinfomessagepayload) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(payload)

	if err != nil {
		tinylogging.AddInfo("[Host] LeaderInfoMessagePayload ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// ToPayload converts byte array to LeaderInfoMessagePayload
func (payload *leaderinfomessagepayload) ToPayload(data []byte) error {
	err := json.Unmarshal(data, &payload)
	if err != nil {
		tinylogging.AddError("[Host] LeaderInfoMessagePayload ToPayload ", err.Error())
		return err
	}
	return nil
}
