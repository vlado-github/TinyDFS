package main

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
	RegisterTimoutHandler(consensus.EventHandlerFunc)
}

type host struct {
	node           messaging.Node
	stateMachine   consensus.StateMachine
	timeoutHandler consensus.TimeoutHandler
	isLead         bool
	connParams     messaging.ConnParams
}

// NewHost creates a new instance of host
func NewHost(connParams messaging.ConnParams, isLead bool) Host {
	var node = messaging.NewNode(connParams, isLead)
	var stateMachine = consensus.NewStateMachine()
	var timeoutHandler = consensus.NewTimeoutHandler()
	return &host{node, stateMachine, timeoutHandler, isLead, connParams}
}

func (h *host) Start() {
	h.node.Run()
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

func (h *host) RegisterTimoutHandler(handler consensus.EventHandlerFunc) {
	h.timeoutHandler.RegisterHandler(handler)
}
