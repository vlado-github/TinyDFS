package consensus

import (
	"encoding/json"
	"messaging"
	"net"
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
		if message.Topic == HEARTBEAT {
			if message.Key != th.lastHeartbeat {
				tinylogging.AddTrace("****Receive a Heartbeat****")
				// if leader reset heartbeat timeout
				th.ResetHeartbeatTime()
				// if any node receives reset election timeout
				th.ResetElectionTime()
				// save Heartbeat key
				th.lastHeartbeat = message.Key
				// reply
				th.sendMessage(message)
			}
		} else if message.Topic == messaging.CLIENT_CONN_OPENED || message.Topic == messaging.CLIENT_CONN_CLOSED {
			tinylogging.AddTrace(message.Topic)
			basePayload := messaging.EmptyPayload()
			err := basePayload.ToPayload(message.Payload)
			if err != nil {
				tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
			} else {
				numOfNodes := basePayload.GetNumOfNodes()
				ipAddresses := basePayload.GetIPs()
				th.SetNumOfNodes(numOfNodes)

				// on every new change in connections we update our
				// network discovery register
				var filtered []NetworkTuple
				for i := range ipAddresses {
					if ipAddresses[i] != "" {

						host, port, _ := net.SplitHostPort(ipAddresses[i])
						tuple := NewNetworkTuple(host, port)
						filtered = append(filtered, tuple)
					}
				}
				registry := NewNetworkRegistry(filtered)
				th.SetNetworkRegistry(registry)
			}
		} else if message.Topic == LEADER_INFO {
			tinylogging.AddTrace(message.Topic)
			leaderInfoPayload := EmptyLeaderInfo()
			err := leaderInfoPayload.ToPayload(message.Payload)
			if err != nil {
				tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
			} else {
				if leaderInfoPayload.GetElectionID() != th.GetElectionID() {
					th.SetLeaderInfo(leaderInfoPayload)
					th.onLeaderElected()
				}
				th.ResetElectionTime()
			}
		} else if message.Topic == LEADER_VOTE {
			tinylogging.AddTrace(message.Topic)
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
						newVotePayload, err := json.Marshal(newVote)
						if err != nil {
							tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
						} else {
							var voteMsg = messaging.Message{
								Key:     uuid.New(),
								Topic:   LEADER_VOTE,
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

							// notify others about new leader
							newLeaderInfo := NewLeaderInfo(th.GetHostIP(), th.GetHostPort(), term, th.GetElectionID(), th.GetHostID().String())
							newLeaderInfoPayload, err := json.Marshal(newLeaderInfo)
							if err != nil {
								tinylogging.AddError("[Host] onMessageReceivedCallback ", err.Error())
							} else {
								tinylogging.AddTrace(newLeaderInfo)
								message := messaging.Message{
									Topic:   LEADER_INFO,
									Key:     uuid.New(),
									Payload: newLeaderInfoPayload,
								}
								th.sendMessage(message)
							}

							// renew queue connections
							// todo: (unlikely) possibiltiy for race condition in case
							// some node connects to leader queue before leader renews connections
							th.SetLeaderInfo(newLeaderInfo)
							th.onLeaderElected()
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
		payload, err := json.Marshal(vote)
		if err != nil {
			tinylogging.AddError("[Host] sendVoteOnElectionTimeoutCallback ", err.Error())
		} else {
			var voteMsg = messaging.Message{
				Key:     uuid.New(),
				Topic:   LEADER_VOTE,
				Payload: payload,
			}
			th.sendMessage(voteMsg)
		}
	}

	sendOnHeartbeatTimeoutFunc := func() {
		tinylogging.AddTrace("****HEARTBEAT**** from nodeID:", th.GetHostID(), " ElectionID:", th.GetElectionID())
		var message = messaging.Message{
			Key:   uuid.New(),
			Topic: HEARTBEAT,
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
