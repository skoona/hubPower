package ports

import (
	"github.com/skoona/hubPower/internal/core/entities"
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
