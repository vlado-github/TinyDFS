package consensus

import "logging"

// CommandHandler handles sent/received Raft commands.
type CommandHandler interface {
	SetNumOfTotalNodes(totalNodes int)
	OnSent(cmd Command)
	OnReceived(cmd Command, stateMachine StateMachine)
}

type commandhandler struct {
	numOfTotalNodes   int
	receivedAckCount  int
	receivedNAckCount int
}

// NewCommandHandler createss new instance of CommandHandler
func NewCommandHandler() CommandHandler {
	return &commandhandler{
		numOfTotalNodes:   0,
		receivedAckCount:  0,
		receivedNAckCount: 0,
	}
}

func (ch *commandhandler) SetNumOfTotalNodes(totalNodes int) {
	ch.numOfTotalNodes = totalNodes
}

func (ch *commandhandler) OnSent(cmd Command) {
	ch.receivedAckCount = 0
	ch.receivedNAckCount = 0
}

func (ch *commandhandler) OnReceived(cmd Command, stateMachine StateMachine) {
	err := cmd.Validate()
	if err != nil {
		logging.AddError("Error: CommandHandler - Invalid command received.", err.Error())
		return
	}
	if cmd.IsAck() {
		ch.receivedAckCount++
	} else {
		ch.receivedNAckCount++
	}
	switch cmd.GetType() {
	case REQUESTVOTE:
		{
			if ch.checkVoteCommandStatus() == MAJORITYACK {
				// become leader
				stateMachine.SetState(LEADER)
			}
			break
		}
	}
}

func (ch *commandhandler) checkVoteCommandStatus() CommandStatuses {
	if ch.receivedAckCount > int(ch.numOfTotalNodes/2) {
		return MAJORITYACK
	} else if ch.receivedNAckCount > int(ch.numOfTotalNodes/2) {
		return MAJORITYNACK
	} else {
		return ACKPENDING
	}
}
