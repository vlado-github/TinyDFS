package tinydfs

import (
	"consensus"
	"fmt"
	"math/rand"
	"messaging"
	"strconv"
	"tinylogging"

	"github.com/google/uuid"
)

// Host is a single unit of distributed storage.
// It's composed of Node instance that allows messaging and persistance
// and uses Raft implementation for leader election.
type Host interface {
	Start()
	GetID() uuid.UUID
	GetElectionID() int
	GetIP() (string, error)
	SendMessage(message messaging.Message)
	RegisterNodeHandler(messaging.HandlerType, messaging.NodeHandlerFunc)
}

type host struct {
	node           messaging.Node
	electionID     int
	stateMachine   consensus.StateMachine
	timeoutHandler consensus.TimeoutHandler
	isQueue        bool
	connParams     messaging.ConnParams

	term      int
	voteCount int
	lastVotes map[string]int
}

// NewHost creates a new instance of host
func NewHost(connParams messaging.ConnParams, isQueue bool) Host {
	node := messaging.NewNode(connParams, isQueue)
	stateMachine := consensus.NewStateMachine()
	timeoutHandler := consensus.NewTimeoutHandler()
	term := 0
	voteCount := 0
	electionID := rand.Int()
	lastVotes := make(map[string]int)
	return &host{
		node,
		electionID,
		stateMachine,
		timeoutHandler,
		isQueue,
		connParams,
		term,
		voteCount,
		lastVotes}
}

func (h *host) Start() {
	h.node.Run()
	onMessageReceivedCallback := func(message messaging.Message) {
		if message.Topic == "LEADER_VOTE" {
			votePayload := EmptyVote()
			err := votePayload.ToPayload(message.Payload)
			if err != nil {
				tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
			} else {
				term := votePayload.GetTerm()
				electionID := votePayload.GetElectionID()
				nodeID := votePayload.GetNodeID()
				if electionID != strconv.Itoa(h.GetElectionID()) { // not me
					// give a vote
					lastVotedTerm := h.lastVotes[electionID]
					if lastVotedTerm != term {
						newVote := NewVote(term, electionID, h.GetID().String())
						newVotePayload, err := newVote.ToByteArray()
						if err != nil {
							tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
						} else {
							var voteMsg = messaging.Message{
								Key:     uuid.New(),
								Topic:   "LEADER_VOTE",
								Payload: newVotePayload,
							}
							h.lastVotes[electionID] = term
							fmt.Println("****Give a vote: TERM: ", term, " ElectionID: ", electionID)
							h.SendMessage(voteMsg)
						}
					}
				} else {
					if nodeID != h.GetID().String() { // not from me
						// receive a vote
						h.voteCount++
						fmt.Println("****Receive a vote: TERM: ", term, " ElectionID: ", electionID, " count:", h.voteCount)
					}
				}
			}
		}
	}
	h.node.RegisterMessageHandler(messaging.MESSAGERECEIVED, onMessageReceivedCallback)
	sendVoteOnElectionTimeoutCallback := func() {
		h.term++
		h.voteCount = 0
		fmt.Println("****Request a vote: TERM: ", h.term, " ElectionID: ", h.electionID, " count:", h.voteCount)
		vote := NewVote(h.term, strconv.Itoa(h.GetElectionID()), h.GetID().String())
		payload, err := vote.ToByteArray()
		if err != nil {
			tinylogging.AddError("[Host] sendVoteOnElectionTimeoutCallback ", err.Error())
		} else {
			var voteMsg = messaging.Message{
				Key:     uuid.New(),
				Topic:   "LEADER_VOTE",
				Payload: payload,
			}
			h.SendMessage(voteMsg)
		}
	}
	h.timeoutHandler.RegisterHandler(sendVoteOnElectionTimeoutCallback)
	h.timeoutHandler.StartElectionTime(h.stateMachine)
}

func (h *host) GetID() uuid.UUID {
	return h.node.GetID()
}

func (h *host) GetElectionID() int {
	return h.electionID
}

func (h *host) GetIP() (string, error) {
	return h.node.GetIP()
}

func (h *host) SendMessage(message messaging.Message) {
	h.node.SendMessage(message)
}

func (h *host) RegisterNodeHandler(handlerType messaging.HandlerType, handler messaging.NodeHandlerFunc) {
	h.node.RegisterHandler(handlerType, handler)
}
