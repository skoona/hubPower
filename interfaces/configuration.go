package interfaces

import (
	"github.com/skoona/hubPower/entities"
)

type Configuration interface {
	Hosts() []*entities.HubHost
	HostIds() []string
	HostById(id string) *entities.HubHost
	AddHost(h *entities.HubHost)
	Apply(h *entities.HubHost) Configuration
	Save()
	Remove(id string)
	ResetConfig()
	Shutdown()
}
