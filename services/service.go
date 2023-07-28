package services

import (
	"context"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/entities"
	"github.com/skoona/hubPower/interfaces"
	"github.com/skoona/hubPower/providers"
)

type service struct {
	ctx          context.Context
	cfg          interfaces.Configuration
	hubProviders []interfaces.HubProvider
}

var _ interfaces.Service = (*service)(nil)

func NewService(ctx context.Context, cfg interfaces.Configuration) (interfaces.Service, error) {
	s := &service{
		ctx:          ctx,
		cfg:          cfg,
		hubProviders: []interfaces.HubProvider{},
	}

	err := s.begin()
	return s, err
}

func (s *service) begin() error {
	var err error

	// initialize Hubs
	for _, host := range s.cfg.Hosts() {
		if host.IsEnabled() {
			hpv := providers.NewHubitatProvider(s.ctx, host)
			s.hubProviders = append(s.hubProviders, hpv)
			host.DeviceDetails = append(host.DeviceDetails, hpv.DeviceDetailsList()...)
			hpv.CreateDeviceEventListener()
		}
	}

	return err
}
func (s *service) Shutdown() {
	commons.DebugLog("Service::Shutdown() called.")
	for _, hp := range s.hubProviders {
		hp.Shutdown()
	}
}
func (s *service) HubEventsMessageChannel(hubId string) chan entities.DeviceEventStream {
	var index int

	for idx, device := range s.cfg.Hosts() {
		if device.Id == hubId {
			index = idx
			break
		}
	}

	return s.hubProviders[index].GetEventListenerChannel()
}
