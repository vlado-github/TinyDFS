package consensus

// StateMachine manages the current state of node.
type StateMachine interface {
	GetCurrentState() State
}

type statemachine struct {
	currentState State
}

// NewStateMachine creates new instance of StateMachine.
func NewStateMachine() StateMachine {
	return &statemachine{
		currentState: FOLLOWER,
	}
}

func (sm *statemachine) GetCurrentState() State {
	return sm.currentState
}
