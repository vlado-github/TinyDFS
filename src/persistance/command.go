package persistance

import "github.com/google/uuid"

type Command struct {
	Key uuid.UUID
	Text  string
}