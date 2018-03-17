package tinydfs

import (
	"consensus"
	"math/rand"
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
	GetNumOfNodes() int
	SetNumOfNodes(numOfNodes int)
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
	lastVotes     map[string]int
	lastHeartbeat uuid.UUID
	numOfNodes    int
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
	lastHeartbeat := uuid.New()
	numOfNodes := 0
	return &host{
		node,
		electionID,
		stateMachine,
		timeoutHandler,
		isQueue,
		connParams,
		term,
		voteCount,
		lastVotes,
		lastHeartbeat,
		numOfNodes}
}

func (h *host) Start() {
	h.registerHandlers()
	h.node.Run()
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

func (h *host) GetNumOfNodes() int {
	if h.isQueue {
		h.numOfNodes = h.node.GetNumOfNodes()
	}
	return h.numOfNodes
}

func (h *host) SetNumOfNodes(numOfNodes int) {
	h.numOfNodes = numOfNodes
}

func (h *host) SendMessage(message messaging.Message) {
	h.node.SendMessage(message)
}

func (h *host) RegisterNodeHandler(handlerType messaging.HandlerType, handler messaging.NodeHandlerFunc) {
	h.node.RegisterHandler(handlerType, handler)
}

func (h *host) registerHandlers() {
	handlersRegistry := NewHandlersRegistry(h)
	h.node.RegisterMessageHandler(messaging.MESSAGERECEIVED, handlersRegistry.GetMessagingHandler(messaging.MESSAGERECEIVED))
	h.timeoutHandler.RegisterHandler(consensus.ELECTIONTIMEOUT, handlersRegistry.GetTimeoutHandler(consensus.ELECTIONTIMEOUT))
	h.timeoutHandler.RegisterHandler(consensus.HEARTBEATTIMEOUT, handlersRegistry.GetTimeoutHandler(consensus.HEARTBEATTIMEOUT))
}
