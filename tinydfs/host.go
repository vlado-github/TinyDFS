package tinydfs

import (
	"consensus"
	"fmt"
	"math/rand"
	"messaging"
	"strconv"

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

	term          int
	voteCount     int
	lastVotedTerm int
}

// NewHost creates a new instance of host
func NewHost(connParams messaging.ConnParams, isQueue bool) Host {
	node := messaging.NewNode(connParams, isQueue)
	stateMachine := consensus.NewStateMachine()
	timeoutHandler := consensus.NewTimeoutHandler()
	term := 0
	voteCount := 0
	electionID := rand.Int()
	lastVotedTerm := 0
	return &host{node, electionID, stateMachine, timeoutHandler, isQueue, connParams, term, voteCount, lastVotedTerm}
}

func (h *host) Start() {
	h.node.Run()
	onMessageReceivedCallback := func(message messaging.Message) {
		term, _ := message.Payload.(int)
		electionID := message.Payload.(string)

		if message.Topic == "LEADER_VOTE" {
			if electionID != strconv.Itoa(h.GetElectionID()) {
				// give a vote
				if h.lastVotedTerm != term {
					var vote = messaging.Message{
						Key:     uuid.New(),
						Topic:   "LEADER_VOTE",
						Text:    "LEADER_VOTE",
						Payload: message.Payload,
					}
					h.lastVotedTerm = term
					fmt.Println("****Give a vote: ", term, " ", electionID)
					h.SendMessage(vote)
				}
			} else {
				// receive a vote
				h.voteCount++
				fmt.Println("****Receive a vote: ", term, " ", electionID, " count:", h.voteCount)
			}
		}
	}
	h.node.RegisterMessageHandler(messaging.MESSAGERECEIVED, onMessageReceivedCallback)
	onElectionTimeoutCallback := func() {
		h.term++
		h.voteCount = 1
		payload := MessagePayload{
			Term:       h.term,
			ElectionID: strconv.Itoa(h.GetElectionID()),
		}
		var message = messaging.Message{
			Key:     uuid.New(),
			Topic:   "LEADER_VOTE",
			Text:    "LEADER_VOTE",
			Payload: payload,
		}
		h.SendMessage(message)
	}
	h.timeoutHandler.RegisterHandler(onElectionTimeoutCallback)
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
