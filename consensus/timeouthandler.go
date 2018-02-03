package consensus

import (
	"logging"
	"time"
)

// TimeoutHandler handles Raft timeouts.
type TimeoutHandler interface {
	StartElectionTime(stateMachine StateMachine)
	ResetElectionTime(stateMachine StateMachine)
	RegisterHandler(handlerFunc EventHandlerFunc)
}

type timeouthandler struct {
	sendVoteOnElectionTimeout EventHandlerFunc
	timer                     *time.Timer
}

// NewTimeoutHandler creates new instance of TimeoutHandler
func NewTimeoutHandler() TimeoutHandler {
	return &timeouthandler{
		sendVoteOnElectionTimeout: NewEventHandlerFunc(),
	}
}

func (th *timeouthandler) StartElectionTime(stateMachine StateMachine) {
	electionTimeout := GetRandomElectionTimeout()
	logging.AddInfo("[Consensus] StartElectionTime: ", electionTimeout)
	th.timer = time.NewTimer(time.Duration(electionTimeout) * time.Millisecond)
	go func() {
		<-th.timer.C
		th.onElectionTimeout(stateMachine)
	}()
}

func (th *timeouthandler) ResetElectionTime(stateMachine StateMachine) {
	th.timer.Stop()
	th.StartElectionTime(stateMachine)
	logging.AddInfo("[Consensus] ResetElectionTime")
}

func (th *timeouthandler) RegisterHandler(handlerFunc EventHandlerFunc) {
	th.sendVoteOnElectionTimeout = handlerFunc
}

func (th *timeouthandler) onElectionTimeout(stateMachine StateMachine) {
	// set to Candidate
	stateMachine.SetState(CANDIDATE)
	logging.AddInfo("[Consensus] Changes state: ", stateMachine.GetCurrentState())

	// send Vote command
	th.sendVoteOnElectionTimeout()

	// reset election time
	th.ResetElectionTime(stateMachine)
}
