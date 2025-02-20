package connector

import (
	"sync"

	"github.com/google/uuid"
)

type ConnectorInterface interface {
	Connect(userID uuid.UUID, wg *sync.WaitGroup, connected chan<- struct{})
}
