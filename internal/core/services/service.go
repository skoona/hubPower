package services

import (
	"context"
	"github.com/skoona/hubPower/internal/adapters/repository"
	"github.com/skoona/hubPower/internal/commons"
	"github.com/skoona/hubPower/internal/core/entities"
	"github.com/skoona/hubPower/internal/core/ports"
)

type service struct {
	ctx          context.Context
	cfg          ports.Configuration
	hubProviders []ports.HubRepository
}

var _ ports.Service = (*service)(nil)

func NewService(ctx context.Context, cfg ports.Configuration) (ports.Service, error) {
	s := &service{
		ctx:          ctx,
		cfg:          cfg,
		hubProviders: []ports.HubRepository{},
	}

	err := s.begin()
	return s, err
}

func (s *service) begin() error {
	var err error

	// initialize Hubs
	for _, host := range s.cfg.Hosts() {
		if host.IsEnabled() {
			hpv := repository.NewHubitatRepository(s.ctx, host)
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
