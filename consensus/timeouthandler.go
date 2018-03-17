package consensus

import (
	"time"
	"tinylogging"
)

// TimeoutHandler handles Raft timeouts.
type TimeoutHandler interface {
	StartElectionTime(stateMachine StateMachine)
	ResetElectionTime(stateMachine StateMachine)
	ChangeStateToLeader(stateMachine StateMachine)
	RegisterHandler(electionType ElectionTimeoutType, handlerFunc EventHandlerFunc)

	StartHeartbeatTime(stateMachine StateMachine)
	ResetHeartbeatTime(stateMachine StateMachine)
}

type timeouthandler struct {
	sendVoteOnElectionTimeout EventHandlerFunc
	sendOnHeartbeatTimeout    EventHandlerFunc
	timer                     *time.Timer
	heartbeat                 *time.Timer
}

// NewTimeoutHandler creates new instance of TimeoutHandler
func NewTimeoutHandler() TimeoutHandler {
	return &timeouthandler{
		sendVoteOnElectionTimeout: NewEventHandlerFunc(),
		sendOnHeartbeatTimeout:    NewEventHandlerFunc(),
	}
}

func (th *timeouthandler) StartElectionTime(stateMachine StateMachine) {
	electionTimeout := GetRandomElectionTimeout()
	tinylogging.AddInfo("[Consensus] StartElectionTime: ", electionTimeout)
	th.timer = time.NewTimer(time.Duration(electionTimeout) * time.Millisecond)
	go func() {
		<-th.timer.C
		th.onElectionTimeout(stateMachine)
	}()
}

func (th *timeouthandler) ResetElectionTime(stateMachine StateMachine) {
	th.timer.Stop()
	th.StartElectionTime(stateMachine)
	tinylogging.AddInfo("[Consensus] ResetElectionTime")
}

func (th *timeouthandler) StartHeartbeatTime(stateMachine StateMachine) {
	if stateMachine.GetCurrentState() == LEADER {
		tinylogging.AddInfo("[Consensus] StartHeartbeatTime: ", HEARTBEATMAX)
		th.heartbeat = time.NewTimer(time.Duration(HEARTBEATMAX) * time.Millisecond)
		go func() {
			<-th.heartbeat.C
			th.onHeartbeatTimeout(stateMachine)
		}()
	}
}

func (th *timeouthandler) ResetHeartbeatTime(stateMachine StateMachine) {
	if stateMachine.GetCurrentState() == LEADER {
		th.heartbeat.Stop()
		th.StartHeartbeatTime(stateMachine)
		tinylogging.AddTrace("[Consensus] ResetHeartbeatTime")
	}
}

func (th *timeouthandler) ChangeStateToLeader(stateMachine StateMachine) {
	stateMachine.SetState(LEADER)
	tinylogging.AddInfo("[Consensus] Changes state: ", stateMachine.GetCurrentState())
	th.StartHeartbeatTime(stateMachine)
}

func (th *timeouthandler) RegisterHandler(electionType ElectionTimeoutType, handlerFunc EventHandlerFunc) {
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

func (th *timeouthandler) onElectionTimeout(stateMachine StateMachine) {
	// set to Candidate
	stateMachine.SetState(CANDIDATE)
	tinylogging.AddInfo("[Consensus] Changes state: ", stateMachine.GetCurrentState())

	// send Vote command
	th.sendVoteOnElectionTimeout()

	// reset election time
	th.ResetElectionTime(stateMachine)
}

func (th *timeouthandler) onHeartbeatTimeout(stateMachine StateMachine) {
	if stateMachine.GetCurrentState() == LEADER {
		// send Heartbeat command
		th.sendOnHeartbeatTimeout()

		// reset election time
		th.ResetHeartbeatTime(stateMachine)
	}
}
