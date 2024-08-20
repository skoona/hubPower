package repository

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/hubPower/internal/commons"
	"github.com/skoona/hubPower/internal/core/entities"
	"github.com/skoona/hubPower/internal/core/ports"
	"strings"
)

const (
	HubHostsPrefs = "HubHost"
)

type config struct {
	hubs  []*entities.HubHost
	prefs fyne.Preferences
}

var _ ports.Configuration = (*config)(nil)
var _ ports.Provider = (*config)(nil)

func NewConfigRepository(prefs fyne.Preferences) (ports.Configuration, error) {
	var err error
	var hubHosts []*entities.HubHost

	// create a default
	defaultHubHosts := []*entities.HubHost{
		entities.NewHubHost("Skoona's", commons.HubitatIP(), commons.HubitatAccessToken(), commons.DefaultIp(), 5, false),
	}

	// retrieve existing
	hubHostString := prefs.String(HubHostsPrefs)
	if hubHostString != "" {
		commons.DebugLog("NewConfigRepository() load HubHost preferences succeeded ")
		err = json.Unmarshal([]byte(hubHostString), &hubHosts)
		if err != nil {
			commons.DebugLog("NewConfigRepository() Unmarshal HubHosts failed: ", err.Error())
		}
	}
	if len(hubHosts) == 0 {
		commons.DebugLog("NewConfigRepository() load HubHost preferences failed using defaults ")
		save, err := json.Marshal(defaultHubHosts)
		if err != nil {
			commons.DebugLog("NewConfigRepository() Marshal saving HubHost prefs failed: ", err.Error())
		}
		prefs.SetString(HubHostsPrefs, string(save))
		hubHosts = defaultHubHosts
	}
	for _, h := range hubHosts {
		h.ListenerUri = strings.Replace("http://IPADDR:2600/hubEvents", "IPADDR", commons.DefaultIp(), 1)
		for _, dv := range h.DeviceDetails {
			dv.BWattValue = binding.NewFloat()
			dv.BVoltageValue = binding.NewFloat()
		}
	}

	cfg := &config{
		hubs:  hubHosts,
		prefs: prefs,
	}

	return cfg, err
}
func (c *config) ResetConfig() {
	c.prefs.SetString(HubHostsPrefs, "")
}
func (c *config) HostById(id string) *entities.HubHost {
	var host *entities.HubHost

	for _, h := range c.hubs {
		if h.Id == id {
			host = h
			break
		}
	}

	return host
}
func (c *config) Hosts() []*entities.HubHost {
	return c.hubs
}
func (c *config) Apply(h *entities.HubHost) ports.Configuration {
	index := -1
	for idx, hub := range c.hubs {
		if h.Id == hub.Id {
			index = idx
			break
		}
	}
	if index != -1 {
		c.hubs[index] = h
	} else {
		c.hubs = append(c.hubs, h)
	}
	return c
}
func (c *config) AddHost(h *entities.HubHost) {
	c.Apply(h).Save()
	commons.DebugLog("Config::AddHubHost() saved: .", h)
}
func (c *config) Save() {
	save, err := json.Marshal(c.hubs)
	if err != nil {
		commons.DebugLog("Configuration::Save() marshal hubHosts failed: ", err.Error())
	} else {
		c.prefs.SetString(HubHostsPrefs, string(save))
	}
}
func (c *config) Remove(id string) {
	if id == "" {
		return
	}
	index := -1
	for idx, h := range c.hubs {
		if h.Id == id {
			index = idx
			break
		}
	}
	if index != -1 {
		c.hubs = commons.RemoveIndexFromSlice(index, c.hubs)
		c.Save()
	}
}
func (c *config) HostIds() []string {
	ids := []string{}
	for _, h := range c.hubs {
		ids = append(ids, h.Id)
	}
	return ids
}

// Shutdown compliance with Provider Interface
func (c *config) Shutdown() {
	commons.DebugLog("Config::Shutdown() called.")
}
