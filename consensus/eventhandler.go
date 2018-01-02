package consensus

// EventHandlerFunc represent a type of callback function.
type EventHandlerFunc func()

// NewEventHandlerFunc returns empty callback function.
func NewEventHandlerFunc() EventHandlerFunc {
	return func() {}
}
