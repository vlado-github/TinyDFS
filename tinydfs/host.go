package tinydfs

import (
	"consensus"
	"messaging"

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
	stateMachine   consensus.StateMachine
	timeoutHandler consensus.TimeoutHandler
	isQueue        bool
	connParams     messaging.ConnParams
}

// NewHost creates a new instance of host
func NewHost(connParams messaging.ConnParams, isQueue bool) Host {
	var node = messaging.NewNode(connParams, isQueue)
	var stateMachine = consensus.NewStateMachine()
	var timeoutHandler = consensus.NewTimeoutHandler()
	return &host{node, stateMachine, timeoutHandler, isQueue, connParams}
}

func (h *host) Start() {
	h.node.Run()
	onElectionTimeoutCallback := func() {
		var message = messaging.Message{Key: uuid.New(), Topic: "LEADER_VOTE", Text: "LEADER_VOTE"}
		h.SendMessage(message)
	}
	h.timeoutHandler.RegisterHandler(onElectionTimeoutCallback)
	h.timeoutHandler.StartElectionTime(h.stateMachine)
}

func (h *host) GetID() uuid.UUID {
	return h.node.GetID()
}

func (h *host) GetElectionID() int {
	return h.node.GetElectionID()
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
