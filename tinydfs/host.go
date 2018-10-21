package tinydfs

import (
	"consensus"
	"math/rand"
	"messaging"
	"tinylogging"

	"github.com/google/uuid"
)

// This is a flag to distinct between explicit call
// and application (leader) crash
var isConnCloseExplicit = false

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

	ConnectToLeaderQueue() error
	ConnectToNextAvailableQueue() error
	CloseConnToQueue() error
}

type host struct {
	node           messaging.Node
	electionID     int
	timeoutHandler consensus.TimeoutHandler
	connParams     messaging.ConnParams

	term          int
	voteCount     int
	lastVotes     map[string]int
	lastHeartbeat uuid.UUID
}

// NewHost creates a new instance of host
func NewHost(connParams messaging.ConnParams, broadcastConnParams messaging.ConnParams, port string) Host {
	node := messaging.NewNode(connParams, broadcastConnParams, port)
	hostIP, _ := node.GetIP()
	hostPort := node.GetPort()
	term := 0
	voteCount := 0
	electionID := rand.Int()
	lastVotes := make(map[string]int)
	lastHeartbeat := uuid.New()
	timeoutHandler := consensus.NewTimeoutHandler(electionID, hostIP, hostPort, node.GetID())
	return &host{
		node,
		electionID,
		timeoutHandler,
		connParams,
		term,
		voteCount,
		lastVotes,
		lastHeartbeat}
}

func (h *host) Start() {
	h.registerHandlers()
	h.node.Run()
	h.timeoutHandler.RegisterSendCallback(h.SendMessage)
	h.timeoutHandler.StartElectionTime()
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

// For case a leader is unaccessable next queue for messaging is selected.
func (h *host) ConnectToNextAvailableQueue() error {
	if !isConnCloseExplicit {
		isConnCloseExplicit = false
		tinylogging.AddTrace("next queue")
		leaderInfo := h.timeoutHandler.GetLeaderInfo()
		networkRegistry := h.timeoutHandler.GetNetworkRegistry()
		if networkRegistry != nil {
			tuples := networkRegistry.GetItems()
			for i := range tuples {
				item := tuples[i]
				if item.GetIP() != leaderInfo.GetIP() || item.GetPort() != leaderInfo.GetPort() {
					tinylogging.AddTrace("NextQueue connect:", item.GetIP(), item.GetPort(), h.connParams.Protocol)
					return h.node.ConnectToQueue(h.connParams.Protocol, item.GetIP()+":"+item.GetPort())
				}
			}
		}
	}
	isConnCloseExplicit = false
	return nil
}

// Opens connection to elected leader's queue
func (h *host) ConnectToLeaderQueue() error {
	leaderInfo := h.timeoutHandler.GetLeaderInfo()
	err := h.node.ConnectToQueue(h.connParams.Protocol, leaderInfo.GetIP()+":"+leaderInfo.GetPort())
	if err == nil {
		tinylogging.AddTrace("***** CONNECTED TO LEADER *****")
	}
	return err
}

// Closes connection to current queue
func (h *host) CloseConnToQueue() error {
	isConnCloseExplicit = true
	return h.node.CloseConn()
}

func (h *host) RegisterNodeHandler(handlerType messaging.HandlerType, handler messaging.NodeHandlerFunc) {
	h.node.RegisterHandler(handlerType, handler)
}

func (h *host) registerHandlers() {
	onLeaderElectedCallback := func() {
		leaderInfo := h.timeoutHandler.GetLeaderInfo()
		if leaderInfo != nil {
			h.CloseConnToQueue()
			h.ConnectToLeaderQueue()
		}
	}
	h.timeoutHandler.RegisterOnLeaderElectedHandler(onLeaderElectedCallback)

	onQueueConnClosedCallback := func() {
		h.ConnectToNextAvailableQueue()
	}
	h.RegisterNodeHandler(messaging.QUEUECONNCLOSED, onQueueConnClosedCallback)

	h.node.RegisterMessageHandler(messaging.MESSAGERECEIVED,
		h.timeoutHandler.GetHandlersRegistry().GetMessagingHandler(messaging.MESSAGERECEIVED))
}
