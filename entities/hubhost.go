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
	ThisIpAddress        string
	GraphingSamplePeriod time.Duration
	DeviceDetails        []*DeviceDetails `json:"-"`
}

func NewHubHost(name, ipaddress, accessToken, listenOnIp string, graphPeriod time.Duration) *HubHost {
	id, _ := uuid.NewUUID()
	return &HubHost{
		Id:                   id.String(),
		Name:                 name,
		IpAddress:            ipaddress,
		AccessToken:          accessToken,
		ListenerUri:          strings.Replace("http://IPADDR:2600/hubEvents", "IPADDR", listenOnIp, 1),
		ThisIpAddress:        listenOnIp + ":2600",
		GraphingSamplePeriod: graphPeriod,
	}
}
