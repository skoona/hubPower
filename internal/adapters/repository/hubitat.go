package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/hubPower/internal/commons"
	"github.com/skoona/hubPower/internal/core/entities"
	"github.com/skoona/hubPower/internal/core/ports"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DeviceList                = "DeviceList"
	DeviceDetailsList         = "DeviceDetailsList"
	DeviceDetailById          = "DeviceDetailById"
	DeviceCapabilitiesById    = "DeviceCapabilitiesById"
	DeviceEventHistoryById    = "DeviceEventHistoryById"
	CreateDeviceEventListener = "CreateDeviceEventListener"
)

type HubError struct {
	Error   bool   `json:"error"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type hubitat struct {
	ctx               context.Context
	host              *entities.HubHost
	online            bool
	templates         map[string]string
	eventListenerChan chan entities.DeviceEventStream
	listener          http.Server
	shutdownGroup     sync.WaitGroup
}

var _ (ports.HubRepository) = (*hubitat)(nil)
var _ (ports.Provider) = (*hubitat)(nil)

// NewHubitatRepository creates a new hub provider to a configured hub host
// calls DeviceList to initialize the devices in provided hubHost
func NewHubitatRepository(ctx context.Context, hubHost *entities.HubHost) ports.HubRepository {
	provider := &hubitat{
		ctx:               ctx,
		eventListenerChan: make(chan entities.DeviceEventStream, 48),
		host:              hubHost,
		templates: map[string]string{
			DeviceList:                "http://IPADDRESS/apps/api/56/devices?access_token=TOKEN",
			DeviceDetailsList:         "http://IPADDRESS/apps/api/56/devices/all?access_token=TOKEN",
			DeviceDetailById:          "http://IPADDRESS/apps/api/56/devices/DEVICEID?access_token=TOKEN",
			DeviceCapabilitiesById:    "http://IPADDRESS/apps/api/56/devices/DEVICEID/capabilities?access_token=TOKEN",
			DeviceEventHistoryById:    "http://IPADDRESS/apps/api/56/devices/DEVICEID/events?access_token=TOKEN",
			CreateDeviceEventListener: "http://IPADDRESS/apps/api/56/postURLURI?access_token=TOKEN",
		},
		shutdownGroup: sync.WaitGroup{},
	}

	return provider
}

// prepareUri internal use only, applies current values to url templates
func (h *hubitat) prepareUri(template, value string) string {
	// replace IPADDRESS and TOKEN
	uri, ok := h.templates[template]
	if !ok {
		return ""
	}
	uri = strings.Replace(uri, "IPADDRESS", h.host.IpAddress, 1)
	uri = strings.Replace(uri, "TOKEN", h.host.AccessToken, 1)

	switch template {
	case DeviceDetailById, DeviceCapabilitiesById, DeviceEventHistoryById:
		uri = strings.Replace(uri, "DEVICEID", value, 1)

	case CreateDeviceEventListener: // URL and urlEncoded
		uri = strings.Replace(uri, "URI", value, 1)
	}
	commons.DebugLog("HubitatProvider::prepareUri() formatted uri ==> ", uri)
	return uri
}

// apiRequest internal use only, invokes the http get/post request
func (h *hubitat) apiRequest(uri string) ([]byte, error) {
	var err error
	var req *http.Request

	hub := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err = http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		commons.DebugLog(err.Error())
		return []byte{}, err
	}

	res, err := hub.Do(req)
	if err != nil {
		commons.DebugLog(err.Error(), " Code: ", res.Status)
		return []byte{}, err
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		commons.DebugLog(err.Error())
		return []byte{}, err
	}

	return body, nil
}

// DeviceList returns the list of device ids and names available
func (h *hubitat) DeviceList() []*entities.DeviceList {
	devices := []*entities.DeviceList{}
	body, err := h.apiRequest(h.prepareUri(DeviceList, ""))
	if err != nil {
		commons.DebugLog(err.Error())
		return devices
	}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		commons.DebugLog(err.Error())
		commons.DebugLog(string(body))
	}
	return devices
}

// DeviceDetailsList returns device details for all devices
func (h *hubitat) DeviceDetailsList() []*entities.DeviceDetails {
	devices := []*entities.DeviceDetails{}

	body, err := h.apiRequest(h.prepareUri(DeviceDetailsList, ""))
	if err != nil {
		commons.DebugLog(err.Error())
		return devices
	}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		commons.DebugLog(err.Error())
		commons.DebugLog(string(body))
	}
	for _, device := range devices {
		if device.BWattValue == nil {
			x, _ := strconv.ParseFloat(device.AttrByKey("Power").(string), 32)
			device.BWattValue = binding.NewFloat()
			_ = device.BWattValue.Set(x)
			z, _ := strconv.Atoi(device.AttrByKey("Voltage").(string))
			device.BVoltageValue = binding.NewInt()
			_ = device.BVoltageValue.Set(z)
		}
	}
	commons.DebugLog("DeviceDetailsList: devices ", devices)

	return devices
}

// DeviceDetailById returns a complete set of device details
func (h *hubitat) DeviceDetailById(id string) *entities.Device {
	device := &entities.Device{}
	body, err := h.apiRequest(h.prepareUri(DeviceDetailById, id))
	if err != nil {
		commons.DebugLog(err.Error())
		return device
	}
	err = json.Unmarshal(body, &device)
	if err != nil {
		commons.DebugLog(err.Error())
		commons.DebugLog(string(body))
	}

	return device
}

// DeviceCapabilitiesById returns a device's capabilities
func (h *hubitat) DeviceCapabilitiesById(id string) []*entities.DeviceCapabilities {
	caps := []*entities.DeviceCapabilities{}
	body, err := h.apiRequest(h.prepareUri(DeviceCapabilitiesById, id))
	if err != nil {
		commons.DebugLog(err.Error())
		return caps
	}
	err = json.Unmarshal(body, &caps)
	if err != nil {
		commons.DebugLog(err.Error())
		commons.DebugLog(string(body))
	}

	return caps
}

// DeviceEventHistoryById returns list of device events
func (h *hubitat) DeviceEventHistoryById(id string) []*entities.DeviceEvent {
	history := []*entities.DeviceEvent{}
	body, err := h.apiRequest(h.prepareUri(DeviceEventHistoryById, id))
	if err != nil {
		commons.DebugLog(err.Error())
		return history
	}
	err = json.Unmarshal(body, &history)
	if err != nil {
		commons.DebugLog(err.Error())
		commons.DebugLog(string(body))
	}

	return history
}

func logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Host, r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

// CreateDeviceEventListener instructs hub to post device event when they occur
func (h *hubitat) CreateDeviceEventListener() bool {
	var hubError HubError
	commons.DebugLog("HubitatProvider::CreateDeviceEventListener() BEGIN")

	body, err := h.apiRequest(h.prepareUri(CreateDeviceEventListener, "/"+url.QueryEscape(h.host.ListenerUri)))
	if err != nil {
		commons.DebugLog("HubitatProvider::CreateDeviceEventListener() Establish listener failed: ", err.Error())
		return false
	}
	_ = json.Unmarshal(body, &hubError)
	if hubError.Error {
		commons.DebugLog("HubitatProvider::CreateDeviceEventListener() Request listener failed: ", hubError)
		return false
	}

	go func(p *hubitat) { // listener
		// start a server to listen
		commons.DebugLog("HubitatProvider::CreateDeviceEventListener() Listener Starting")
		event := entities.DeviceEventStream{}
		mux := http.NewServeMux()
		p.shutdownGroup.Add(1)

		mux.HandleFunc("/hubEvents", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			err = json.Unmarshal(body, &event)
			if err != nil {
				commons.DebugLog("HubitatProvider::CreateDeviceEventListener() Listener Error: ", err.Error())
			} else {
				// on receive push to listener channel
				p.eventListenerChan <- event
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(202)
		})

		mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("received '/status': %s \n", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(202)
			err = json.NewEncoder(w).Encode(p.host)
		})

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("received '/': %s \n", r.Method)
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&HubError{Error: true, Type: "InvalidRequest", Message: "not supported"})
		})

		commons.DebugLog("HubitatProvider::CreateDeviceEventListener() Servers Listening on IpAddress: 0.0.0.0:2600")
		p.listener = http.Server{
			Addr:    "0.0.0.0:2600",
			Handler: logRequests(mux),
		}
		p.online = true
		if err := p.listener.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("error running http server: %s\n", err)
			}
		}
		p.shutdownGroup.Done()
		commons.DebugLog("HubitatProvider::CreateDeviceEventListener() Listener Ended")
	}(h)

	commons.DebugLog("HubitatProvider::CreateDeviceEventListener() END command response: ", string(body), ", hubError: ", hubError)

	return true
}

// CancelDeviceEventListener stops hub from sending events
func (h *hubitat) CancelDeviceEventListener() bool {
	commons.DebugLog("HubitatProvider::CancelDeviceEventListener() BEGIN")

	// expect: {"url":"cleared"}
	body, err := h.apiRequest(h.prepareUri(CreateDeviceEventListener, ""))
	if err != nil {
		commons.DebugLog("HubitatProvider::CancelDeviceEventListener() ERROR response: ", string(body), ", err: ", err.Error())
		return false
	}
	commons.DebugLog("HubitatProvider::CancelDeviceEventListener() END response: ", string(body))
	return true
}

// GetEventListenerChannel returns the channel that events are added to from Hub
func (h *hubitat) GetEventListenerChannel() chan entities.DeviceEventStream {
	return h.eventListenerChan
}

// Shutdown graceful shutdown of server listener and closes events channel
func (h *hubitat) Shutdown() {
	commons.DebugLog("HubitatProvider::Shutdown() BEGIN")
	if h.online {
		h.CancelDeviceEventListener()
		err := h.listener.Shutdown(h.ctx)
		if err != nil {
			commons.DebugLog(err.Error())
		}

		h.shutdownGroup.Wait() // give/wait server time to exit
		close(h.eventListenerChan)
	}
	commons.DebugLog("HubitatProvider::Shutdown() END")
}
