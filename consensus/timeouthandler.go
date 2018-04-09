package consensus

import (
	"messaging"
	"time"
	"tinylogging"
	"github.com/google/uuid"
)

// TimeoutHandler handles Raft timeouts.
type TimeoutHandler interface {
	StartElectionTime()
	ResetElectionTime()
	ChangeStateToLeader()
	RegisterSendCallback(callback func(message messaging.Message))
	RegisterOnLeaderElectedHandler(eventHandler EventHandlerFunc)

	StartHeartbeatTime()
	ResetHeartbeatTime()

	GetNumOfNodes() int
	SetNumOfNodes(numOfNodes int)
	GetElectionID() int
	GetNetworkRegistry() []string
	SetNetworkRegistry(ipAddresses []string)
	GetLeaderInfo() LeaderInfo
	SetLeaderInfo(info LeaderInfo)

	GetHandlersRegistry() HandlersRegistry
}

type timeouthandler struct {
	networkRegistry			  []string
	handlersRegistry		  HandlersRegistry
	sendMessage				  func(message messaging.Message)
	sendVoteOnElectionTimeout EventHandlerFunc
	sendOnHeartbeatTimeout    EventHandlerFunc
	onLeaderElected	          EventHandlerFunc
	timer                     *time.Timer
	heartbeat                 *time.Timer
	stateMachine              StateMachine
	voteCount     			  int
	lastVotes     			  map[string]int
	lastHeartbeat 			  uuid.UUID
	term 			          int
	numOfNodes                int
	electionID                int
	hostID 		              uuid.UUID
	hostIP					  string
	isQueue					  bool
	leaderInfo                LeaderInfo
}

// NewTimeoutHandler creates new instance of TimeoutHandler
func NewTimeoutHandler(electionID int, hostIP string, hostID uuid.UUID, isQueue bool) TimeoutHandler {
	voteCount := 0
	lastVotes := make(map[string]int)
	lastHeartbeat := uuid.New()
	numOfNodes := 0
	term := 0
	stateMachine := NewStateMachine()
	th := &timeouthandler{
		sendVoteOnElectionTimeout: NewEventHandlerFunc(),
		sendOnHeartbeatTimeout:    NewEventHandlerFunc(),
		onLeaderElected:		   NewEventHandlerFunc(),
		stateMachine: stateMachine,
		voteCount: voteCount,
		lastVotes: lastVotes,
		lastHeartbeat: lastHeartbeat,
		numOfNodes: numOfNodes,
		term: term,
		electionID: electionID,
		hostID: hostID,
		hostIP: hostIP,
		isQueue: isQueue,
	}
	th.handlersRegistry = NewHandlersRegistry(th)
	th.registerHandler(ELECTIONTIMEOUT, th.handlersRegistry.GetTimeoutHandler(ELECTIONTIMEOUT))
	th.registerHandler(HEARTBEATTIMEOUT, th.handlersRegistry.GetTimeoutHandler(HEARTBEATTIMEOUT))
	return th
}

func (th *timeouthandler) StartElectionTime() {
	electionTimeout := GetRandomElectionTimeout()
	tinylogging.AddInfo("[Consensus] StartElectionTime: ", electionTimeout)
	th.timer = time.NewTimer(time.Duration(electionTimeout) * time.Millisecond)
	go func() {
		<-th.timer.C
		th.onElectionTimeout()
	}()
}

func (th *timeouthandler) ResetElectionTime() {
	th.timer.Stop()
	th.StartElectionTime()
	tinylogging.AddInfo("[Consensus] ResetElectionTime")
}

func (th *timeouthandler) StartHeartbeatTime() {
	if th.stateMachine.GetCurrentState() == LEADER {
		tinylogging.AddInfo("[Consensus] StartHeartbeatTime: ", HEARTBEATMAX)
		th.heartbeat = time.NewTimer(time.Duration(HEARTBEATMAX) * time.Millisecond)
		go func() {
			<-th.heartbeat.C
			th.onHeartbeatTimeout()
		}()
	}
}

func (th *timeouthandler) ResetHeartbeatTime() {
	if th.stateMachine.GetCurrentState() == LEADER {
		th.heartbeat.Stop()
		th.StartHeartbeatTime()
		tinylogging.AddTrace("[Consensus] ResetHeartbeatTime")
	}
}

func (th *timeouthandler) ChangeStateToLeader() {
	th.stateMachine.SetState(LEADER)
	tinylogging.AddInfo("[Consensus] Changes state: ", th.stateMachine.GetCurrentState())
	th.StartHeartbeatTime()
}

func (th *timeouthandler) RegisterSendCallback(callback func(message messaging.Message)) {
	th.sendMessage = callback
}

func (th *timeouthandler) RegisterOnLeaderElectedHandler(eventHandler EventHandlerFunc) {
	th.onLeaderElected = eventHandler
}

func (th *timeouthandler) registerHandler(electionType ElectionTimeoutType, handlerFunc EventHandlerFunc) {
	switch electionType {
	case ELECTIONTIMEOUT:
		{
			th.sendVoteOnElectionTimeout = handlerFunc
			break
		}
	case HEARTBEATTIMEOUT:
		{
			th.sendOnHeartbeatTimeout = handlerFunc
		}
	}
}

func (th *timeouthandler) onElectionTimeout() {
	// set to Candidate
	th.stateMachine.SetState(CANDIDATE)
	tinylogging.AddInfo("[Consensus] Changes state: ", th.stateMachine.GetCurrentState())

	// send Vote command
	th.sendVoteOnElectionTimeout()

	// reset election time
	th.ResetElectionTime()
}

func (th *timeouthandler) onHeartbeatTimeout() {
	if th.stateMachine.GetCurrentState() == LEADER {
		// send Heartbeat command
		th.sendOnHeartbeatTimeout()

		// reset election time
		th.ResetHeartbeatTime()
	}
}

// getters and setters

func (th *timeouthandler) GetHandlersRegistry() HandlersRegistry {
	return th.handlersRegistry
}

func (th *timeouthandler) GetNumOfNodes() int {
	return th.numOfNodes
}

func (th *timeouthandler) SetNumOfNodes(numOfNodes int) {
	th.numOfNodes = numOfNodes
}

func (th *timeouthandler) GetElectionID() int {
	return th.electionID
}

func (th *timeouthandler) GetHostIP() string {
	return th.hostIP
}

func (th *timeouthandler) GetHostID() uuid.UUID {
	return th.hostID
}

func (th *timeouthandler) GetNetworkRegistry() []string {
	return th.networkRegistry
}

func (th *timeouthandler) SetNetworkRegistry(ipAddresses []string) {
	th.networkRegistry = ipAddresses
}

func (th *timeouthandler) GetLeaderInfo() LeaderInfo {
	return th.leaderInfo
}

func (th *timeouthandler) SetLeaderInfo(info LeaderInfo) {
	th.leaderInfo = info
}
