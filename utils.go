package openc

import (
	"github.com/gofrs/uuid"
)

func generateID() ID {
	return ID(uuid.Must(uuid.NewV4()).String())
}
