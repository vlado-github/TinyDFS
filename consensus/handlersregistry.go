package consensus

import (
	"messaging"
	"strconv"
	"tinylogging"

	"github.com/google/uuid"
)

// HandlersRegistry collection of Host's handlers.
type HandlersRegistry interface {
	GetMessagingHandler(messagingType messaging.HandlerType) func(message messaging.Message)
	GetTimeoutHandler(timeoutType ElectionTimeoutType) func()
}

type handlersregistry struct {
	onMessageReceivedHandler         func(message messaging.Message)
	sendVoteOnElectionTimeoutHandler func()
	sendOnHeartbeatTimeoutHandler    func()
}

// NewHandlersRegistry creates new handlers registry.
func NewHandlersRegistry(th *timeouthandler) HandlersRegistry {
	onMessageReceivedFunc := func(message messaging.Message) {
		if message.Topic == "HEARTBEAT" {
			if message.Key != th.lastHeartbeat {
				tinylogging.AddTrace("****Receive a Heartbeat****")
				// if leader reset heartbeat timeout
				th.ResetHeartbeatTime()
				// if any node receives reset election timeout
				th.ResetElectionTime()
				// reply
				th.sendMessage(message)
				// save Heartbeat ket
				th.lastHeartbeat = message.Key
			}
		} else if message.Topic == "CLIENT_CONN_OPENED" || message.Topic == "CLIENT_CONN_CLOSED" {
			tinylogging.AddTrace(message.Topic)
			basePayload := messaging.EmptyPayload()
			err := basePayload.ToPayload(message.Payload)
			if err != nil {
				tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
			} else {
				numOfNodes := basePayload.GetNumOfNodes()
				th.SetNumOfNodes(numOfNodes)
			}
		} else if message.Topic == "LEADER_VOTE" {
			votePayload := EmptyVote()
			err := votePayload.ToPayload(message.Payload)
			if err != nil {
				tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
			} else {
				// prepare vote data
				term := votePayload.GetTerm()
				electionID := votePayload.GetElectionID()
				nodeID := votePayload.GetNodeID()
				if electionID != strconv.Itoa(th.GetElectionID()) { // not me, vote from other nodes
					// give a vote
					lastVotedTerm := th.lastVotes[electionID]
					if lastVotedTerm != term {
						newVote := NewVote(term, electionID, th.GetHostID().String())
						newVotePayload, err := newVote.ToByteArray()
						if err != nil {
							tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
						} else {
							var voteMsg = messaging.Message{
								Key:     uuid.New(),
								Topic:   "LEADER_VOTE",
								Payload: newVotePayload,
							}
							th.lastVotes[electionID] = term
							// if votes, host will reset the election timeout
							th.ResetElectionTime()
							// send vote
							tinylogging.AddTrace("****Give a vote: TERM: ", term, " ElectionID: ", electionID)
							th.sendMessage(voteMsg)
						}
					}
				} else {
					if nodeID != th.GetHostID().String() { // not from me
						// receive a vote
						th.voteCount++
						tinylogging.AddTrace("****Receive a vote: TERM: ", term, " ElectionID: ", electionID, " count:", th.voteCount, " totalNodes:", th.GetNumOfNodes())
						if th.voteCount > int(th.GetNumOfNodes()/2) {
							th.ChangeStateToLeader()
							tinylogging.AddTrace("****BECAME A LEADER****")
						}
					}
				}
			}
		}
	}

	sendVoteOnElectionTimeoutFunc := func() {
		th.term++
		th.voteCount = 0
		tinylogging.AddTrace("****Request a vote: TERM: ", th.term, " ElectionID: ", th.electionID, " count:", th.voteCount)
		vote := NewVote(th.term, strconv.Itoa(th.GetElectionID()), th.GetHostID().String())
		payload, err := vote.ToByteArray()
		if err != nil {
			tinylogging.AddError("[Host] sendVoteOnElectionTimeoutCallback ", err.Error())
		} else {
			var voteMsg = messaging.Message{
				Key:     uuid.New(),
				Topic:   "LEADER_VOTE",
				Payload: payload,
			}
			th.sendMessage(voteMsg)
		}
	}

	sendOnHeartbeatTimeoutFunc := func() {
		tinylogging.AddTrace("****HEARTBEAT**** from nodeID:", th.GetHostID(), " ElectionID:", th.GetElectionID())
		var message = messaging.Message{
			Key:   uuid.New(),
			Topic: "HEARTBEAT",
		}
		th.sendMessage(message)
	}

	return &handlersregistry{
		onMessageReceivedHandler:         onMessageReceivedFunc,
		sendVoteOnElectionTimeoutHandler: sendVoteOnElectionTimeoutFunc,
		sendOnHeartbeatTimeoutHandler:    sendOnHeartbeatTimeoutFunc,
	}
}

func (r *handlersregistry) GetMessagingHandler(messagingType messaging.HandlerType) func(message messaging.Message) {
	switch messagingType {
	case messaging.MESSAGERECEIVED:
		{
			return r.onMessageReceivedHandler
		}
	}
	return nil
}

func (r *handlersregistry) GetTimeoutHandler(timoutType ElectionTimeoutType) func() {
	switch timoutType {
	case ELECTIONTIMEOUT:
		{
			return r.sendVoteOnElectionTimeoutHandler
		}
	case HEARTBEATTIMEOUT:
		{
			return r.sendOnHeartbeatTimeoutHandler
		}
	}
	return nil
}
