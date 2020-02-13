package core

import (
	"github.com/gofrs/uuid"
)

type (
	// ID is a uniq identifier
	ID string
)

func generateID() ID {
	return ID(uuid.Must(uuid.NewV4()).String())
}
