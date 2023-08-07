package ports

import (
	"github.com/skoona/hubPower/internal/core/entities"
)

type Service interface {
	HubEventsMessageChannel(hubId string) chan entities.DeviceEventStream
	Shutdown()
}
