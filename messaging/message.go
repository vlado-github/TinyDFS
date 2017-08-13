package messaging

import "github.com/google/uuid"

type Message struct {
	Key uuid.UUID
	Topic string
	Text  string
}
