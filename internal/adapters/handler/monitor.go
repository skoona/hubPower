package handler

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/hubPower/internal/commons"
	"github.com/skoona/hubPower/internal/core/entities"
	"github.com/skoona/hubPower/internal/core/ports"
	"github.com/skoona/sknlinechart"
	"strconv"
	"syscall"
	"time"
)

// MonitorPage primary application page
func (v *viewHandler) MonitorPage() *fyne.Container {
	hostTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Hub Assets", theme.ComputerIcon(), v.OverviewPage()),
	)

	//v.handleUpdatesForMonitorPage(host, v.service, status, events, chart, knowledge)
	var chart sknlinechart.LineChart
	var err error
	for _, hub := range v.hosts {
		v.chartPages[hub.Id] = make(map[string]sknlinechart.LineChart)
		if hub.IsEnabled() {
			for _, device := range hub.DeviceDetails {
				// GraphSamplingPeriod for Charts
				v.chartPageData[device.Id] = map[string]ports.GraphPointSmoothing{}

				for _, key := range v.chartKeys {
					intf := entities.NewGraphAverage(hub.Name, key, hub.GraphingSamplePeriod)
					v.chartPageData[device.Id][key] = intf
				}
				// for chart page updates
				data := map[string][]*sknlinechart.ChartDatapoint{}
				chart, err = sknlinechart.NewLineChart(
					hub.Name,
					"",
					65,
					&data,
				)
				if err != nil {
					dialog.ShowError(err, v.mainWindow)
					commons.ShutdownSignals <- syscall.SIGINT
				}
				chart.SetBottomLeftLabel(hub.Name + "@" + hub.IpAddress + " has " + strconv.Itoa(len(hub.DeviceDetails)) + " devices")
				v.chartPages[hub.Id][device.Id] = chart

				tab := container.NewTabItemWithIcon(device.Label, commons.SknSelectThemedResource("sensorOn"),
					container.NewAppTabs(
						container.NewTabItemWithIcon("History", theme.HistoryIcon(), chart),
						container.NewTabItemWithIcon("Detailed", theme.VisibilityIcon(), v.DeviceCard(device)),
					),
				)
				hostTabs.Append(tab)
			}

			go func(ctxx context.Context, vv *viewHandler, host *entities.HubHost) { // shutdown monitor
				eventChannel := vv.service.HubEventsMessageChannel(host.Id)
				commons.DebugLog("ViewHandler::MonitorPage() listener BEGIN")
				var lc sknlinechart.LineChart

			Gone:
				for {
					select {
					case <-ctxx.Done():
						commons.DebugLog("ViewHandler::MonitorPage() shutdown listener: ", ctxx.Err().Error())
						break Gone

					case ev, ok := <-eventChannel:
						commons.DebugLog("ViewHandler::MonitorPage(", host.Name, ") listener received: ", ev)
						if !ok {
							break Gone
						}
					found:
						for _, device := range host.DeviceDetails {
							if device.Id == ev.Content.DeviceId {
								switch ev.Content.Name {
								case "power":
									z, _ := strconv.ParseFloat(ev.Content.Value, 32)
									err := device.BWattValue.Set(z)
									if err != nil {
										commons.DebugLog("ViewHandler::MonitorPage() listener(", ev.Content.Name, ") float parsing error: ", err.Error())
									}
									d64 := vv.chartPageData[device.Id]["Watts"].AddValue(z)
									point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorYellow, time.Now().Format(time.RFC1123))
									lc = vv.chartPages[host.Id][device.Id]
									lc.ApplyDataPoint("Watts", &point)
									break found
								case "voltage":
									z, _ := strconv.ParseInt(ev.Content.Value, 10, 32)
									err := device.BVoltageValue.Set(int(z))
									if err != nil {
										commons.DebugLog("ViewHandler::MonitorPage() listener(", ev.Content.Name, ") int parsing error: ", err.Error())
									}
									d64 := vv.chartPageData[device.Id]["Voltage"].AddValue(float64(z))
									point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorGreen, time.Now().Format(time.RFC1123))
									lc = vv.chartPages[host.Id][device.Id]
									lc.ApplyDataPoint("Voltage", &point)
									break found
								}
							}
						}
					}
				}
				commons.DebugLog("HubitatProvider::CreateDeviceEventListener() publisher END")
			}(v.ctx, v, hub)
		}
	}

	return container.NewPadded(hostTabs)
}
