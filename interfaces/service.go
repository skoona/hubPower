package interfaces

import "github.com/skoona/hubPower/entities"

type Service interface {
	HubEventsMessageChannel(hubId string) chan entities.DeviceEventStream
	Shutdown()
}
