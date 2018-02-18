package tinydfs

// MessagePayload represents additional message info for
// leader election
type MessagePayload struct {
	Term       int
	ElectionID string
}
