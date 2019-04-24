package bus

import (
	"github.com/asaskevich/EventBus"
)

// Event bus decorator
type Bus struct {
	EventBus.Bus
}

// Event bus constructor
func NewEventBus() EventBus.Bus {
	b := Bus{
		EventBus.New(),
	}

	return b
}
