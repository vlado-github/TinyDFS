package consensus

import (
	"encoding/json"
	"tinylogging"
)

// VoteMessagePayload represents additional message info for
// leader election
type VoteMessagePayload interface {
	ToByteArray() ([]byte, error)
	ToPayload(data []byte) error
	GetTerm() int
	GetElectionID() string
	GetNodeID() string
}

type votemessagepayload struct {
	Term       int    `json:"term"`
	ElectionID string `json:"election_id"`
	NodeID     string `json:"node_id"`
}

// NewVote creates new instance of the message payload for vote
func NewVote(term int, electionID string, nodeID string) VoteMessagePayload {
	return &votemessagepayload{
		Term:       term,
		ElectionID: electionID,
		NodeID:     nodeID,
	}
}

// EmptyVote returns new instance of the message with temporary data set
func EmptyVote() VoteMessagePayload {
	return &votemessagepayload{
		Term:       -1,
		ElectionID: "",
		NodeID:     "",
	}
}

func (payload *votemessagepayload) GetTerm() int {
	return payload.Term
}

func (payload *votemessagepayload) GetElectionID() string {
	return payload.ElectionID
}

func (payload *votemessagepayload) GetNodeID() string {
	return payload.NodeID
}

// ToByteArray converts to Json string
func (payload *votemessagepayload) ToByteArray() ([]byte, error) {
	result, err := json.Marshal(payload)

	if err != nil {
		tinylogging.AddInfo("[Host] VoteMessagePayload ToByteArray", err.Error())
		return nil, err
	}
	return result, nil
}

// ToPayload converts byte array to VoteMessagePayload
func (payload *votemessagepayload) ToPayload(data []byte) error {
	err := json.Unmarshal(data, &payload)
	if err != nil {
		tinylogging.AddError("[Host] VoteMessagePayload ToPayload ", err.Error())
		return err
	}
	return nil
}
