package providers

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/entities"
	"github.com/skoona/hubPower/interfaces"
	"strconv"
)

const (
	HubHostsPrefs = "HubHost"
)

type config struct {
	hubs  []*entities.HubHost
	prefs fyne.Preferences
}

var _ interfaces.Configuration = (*config)(nil)
var _ interfaces.Provider = (*config)(nil)

func NewConfig(prefs fyne.Preferences) (interfaces.Configuration, error) {
	var err error
	var hubHosts []*entities.HubHost

	defaultHubHosts := []*entities.HubHost{
		entities.NewHubHost("Scotts", commons.HubitatIP(), "a79c07db-9178-4976-bd10-428aa0d3d159", commons.DefaultIp(), 5),
	}

	commons.DebugLog("This IpAddress: ", commons.DefaultIp())
	commons.DebugLog("Hubitat IpAddress: ", commons.HubitatIP())
	// retrieve existing
	hubHostString := prefs.String(HubHostsPrefs)
	if hubHostString != "" {
		commons.DebugLog("NewConfig() load HubHost preferences succeeded ")
		err = json.Unmarshal([]byte(hubHostString), &hubHosts)
		if err != nil {
			commons.DebugLog("NewConfig() Unmarshal HubHosts failed: ", err.Error())
		}
	}
	if len(hubHosts) == 0 {
		commons.DebugLog("NewConfig() load HubHost preferences failed using defaults ")
		save, err := json.Marshal(defaultHubHosts)
		if err != nil {
			commons.DebugLog("NewConfig() Marshal saving HubHost prefs failed: ", err.Error())
		}
		prefs.SetString(HubHostsPrefs, string(save))
		hubHosts = defaultHubHosts
	}
	for _, h := range hubHosts {
		h.ThisIpAddress = commons.DefaultIp() + ":2600"
		for _, dv := range h.DeviceDetails {
			dv.BWattValue = binding.NewFloat()
			dv.BVoltageValue = binding.NewInt()
			z, _ := strconv.Atoi(dv.AttrByKey("Voltage").(string))
			dv.BVoltageValue.Set(z)
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
func (c *config) Apply(h *entities.HubHost) interfaces.Configuration {
	c.hubs = append(c.hubs, h)
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
	index := 0
	for idx, h := range c.hubs {
		if h.Id == id {
			index = idx
			break
		}
	}
	commons.RemoveIndexFromSlice(index, c.hubs)
	c.Save()
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
