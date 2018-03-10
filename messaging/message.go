package messaging

import "github.com/google/uuid"

// Message used in protocol with unique identifier
type Message struct {
	Key     uuid.UUID
	Topic   string
	Payload []byte
}
