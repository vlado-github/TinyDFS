package tinydfs

import (
	"consensus"
	"messaging"
	"strconv"
	"tinylogging"

	"github.com/google/uuid"
)

// HandlersRegistry collection of Host's handlers.
type HandlersRegistry interface {
	GetMessagingHandler(messagingType messaging.HandlerType) func(message messaging.Message)
	GetTimeoutHandler(timeoutType consensus.ElectionTimeoutType) func()
}

type handlersregistry struct {
	onMessageReceivedHandler         func(message messaging.Message)
	sendVoteOnElectionTimeoutHandler func()
	sendOnHeartbeatTimeoutHandler    func()
}

// NewHandlersRegistry creates new handlers registry for host.
func NewHandlersRegistry(h *host) HandlersRegistry {
	onMessageReceivedFunc := func(message messaging.Message) {
		if message.Topic == "HEARTBEAT" {
			if message.Key != h.lastHeartbeat {
				tinylogging.AddTrace("****Receive a Heartbeat****")
				// if leader reset heartbeat timeout
				h.timeoutHandler.ResetHeartbeatTime(h.stateMachine)
				// if any node receives reset election timeout
				h.timeoutHandler.ResetElectionTime(h.stateMachine)
				// reply
				h.SendMessage(message)
				// save Heartbeat ket
				h.lastHeartbeat = message.Key
			}
		} else if message.Topic == "CLIENT_CONN_OPENED" || message.Topic == "CLIENT_CONN_CLOSED" {
			tinylogging.AddTrace(message.Topic)
			basePayload := messaging.EmptyPayload()
			err := basePayload.ToPayload(message.Payload)
			if err != nil {
				tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
			} else {
				numOfNodes := basePayload.GetNumOfNodes()
				h.SetNumOfNodes(numOfNodes)
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
				if electionID != strconv.Itoa(h.GetElectionID()) { // not me, vote from other nodes
					// give a vote
					lastVotedTerm := h.lastVotes[electionID]
					if lastVotedTerm != term {
						newVote := NewVote(term, electionID, h.GetID().String())
						newVotePayload, err := newVote.ToByteArray()
						if err != nil {
							tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
						} else {
							var voteMsg = messaging.Message{
								Key:     uuid.New(),
								Topic:   "LEADER_VOTE",
								Payload: newVotePayload,
							}
							h.lastVotes[electionID] = term
							// if votes, host will reset the election timeout
							h.timeoutHandler.ResetElectionTime(h.stateMachine)
							// send vote
							tinylogging.AddTrace("****Give a vote: TERM: ", term, " ElectionID: ", electionID)
							h.SendMessage(voteMsg)
						}
					}
				} else {
					if nodeID != h.GetID().String() { // not from me
						// receive a vote
						h.voteCount++
						tinylogging.AddTrace("****Receive a vote: TERM: ", term, " ElectionID: ", electionID, " count:", h.voteCount, " totalNodes:", h.GetNumOfNodes())
						if h.voteCount > int(h.GetNumOfNodes()/2) {
							h.timeoutHandler.ChangeStateToLeader(h.stateMachine)
							tinylogging.AddTrace("****BECAME A LEADER****")
						}
					}
				}
			}
		}
	}

	sendVoteOnElectionTimeoutFunc := func() {
		h.term++
		h.voteCount = 0
		tinylogging.AddTrace("****Request a vote: TERM: ", h.term, " ElectionID: ", h.electionID, " count:", h.voteCount)
		vote := NewVote(h.term, strconv.Itoa(h.GetElectionID()), h.GetID().String())
		payload, err := vote.ToByteArray()
		if err != nil {
			tinylogging.AddError("[Host] sendVoteOnElectionTimeoutCallback ", err.Error())
		} else {
			var voteMsg = messaging.Message{
				Key:     uuid.New(),
				Topic:   "LEADER_VOTE",
				Payload: payload,
			}
			h.SendMessage(voteMsg)
		}
	}

	sendOnHeartbeatTimeoutFunc := func() {
		tinylogging.AddTrace("****HEARTBEAT**** from nodeID:", h.GetID(), " ElectionID:", h.GetElectionID())
		var message = messaging.Message{
			Key:   uuid.New(),
			Topic: "HEARTBEAT",
		}
		h.SendMessage(message)
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

func (r *handlersregistry) GetTimeoutHandler(timoutType consensus.ElectionTimeoutType) func() {
	switch timoutType {
	case consensus.ELECTIONTIMEOUT:
		{
			return r.sendVoteOnElectionTimeoutHandler
		}
	case consensus.HEARTBEATTIMEOUT:
		{
			return r.sendOnHeartbeatTimeoutHandler
		}
	}
	return nil
}
