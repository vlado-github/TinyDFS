package persistance

import "github.com/google/uuid"

type Command struct {
	Key   uuid.UUID
	Topic string
	Text  string
}
