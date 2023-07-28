package entities

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type HubHost struct {
	Id                   string
	Name                 string
	IpAddress            string
	AccessToken          string
	ListenerUri          string // "http://IPADDR:2600/hubEvents"
	GraphingSamplePeriod time.Duration
	DeviceDetails        []*DeviceDetails `json:"-"`
	Enabled              bool
}

func NewHubHost(name, ipaddress, accessToken, listenOnIp string, graphPeriod time.Duration, enabled bool) *HubHost {
	id, _ := uuid.NewUUID()
	return &HubHost{
		Id:          id.String(),
		Name:        name,
		IpAddress:   ipaddress,
		AccessToken: accessToken,
		// see config.go:56ish
		ListenerUri:          strings.Replace("http://IPADDR:2600/hubEvents", "IPADDR", listenOnIp, 1),
		GraphingSamplePeriod: graphPeriod,
		Enabled:              enabled,
	}
}

func (h *HubHost) IsEnabled() bool {
	return h.Enabled
}
