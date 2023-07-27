package interfaces

import "github.com/skoona/hubPower/entities"

type HubProvider interface {
	DeviceList() []*entities.DeviceList
	DeviceDetailsList() []*entities.DeviceDetails
	DeviceDetailById(id string) *entities.Device
	DeviceCapabilitiesById(id string) []*entities.DeviceCapabilities
	DeviceEventHistoryById(id string) []*entities.DeviceEvent
	CreateDeviceEventListener() bool
	CancelDeviceEventListener() bool
	GetEventListenerChannel() chan entities.DeviceEventStream
	Shutdown()
}
