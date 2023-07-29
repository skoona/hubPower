package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/entities"
	"github.com/skoona/hubPower/interfaces"
	"github.com/skoona/sknlinechart"
	"strconv"
	"syscall"
	"time"
)

// MonitorPage primary application page
func (v *viewProvider) MonitorPage() *fyne.Container {
	hostTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Hub Assets", theme.ComputerIcon(), v.OverviewPage()),
	)

	//v.handleUpdatesForMonitorPage(host, v.service, status, events, chart, knowledge)
	var chart sknlinechart.LineChart
	var err error
	for _, hub := range v.hosts {
		if hub.IsEnabled() {
			for _, device := range hub.DeviceDetails {
				// GraphSamplingPeriod for Charts
				v.chartPageData[hub.Name] = map[string]interfaces.GraphPointSmoothing{}
				for _, key := range v.chartKeys {
					intf := entities.NewGraphAverage(hub.Name, key, hub.GraphingSamplePeriod)
					v.chartPageData[hub.Name][key] = intf
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

				tab := container.NewTabItemWithIcon(device.Label, commons.SknSelectThemedResource("sensorOn"),
					container.NewAppTabs(
						container.NewTabItemWithIcon("History", theme.HistoryIcon(), chart),
						container.NewTabItemWithIcon("Detailed", theme.VisibilityIcon(), v.DeviceCard(device)),
					),
				)
				hostTabs.Append(tab)
			}

			go func(ctxx context.Context, vv *viewProvider, host *entities.HubHost, lc sknlinechart.LineChart) { // shutdown monitor
				eventChannel := vv.service.HubEventsMessageChannel(host.Id)
				commons.DebugLog("ViewProvider::MonitorPage() listener BEGIN")
			Gone:
				for {
					select {
					case <-ctxx.Done():
						commons.DebugLog("ViewProvider::MonitorPage() shutdown listener: ", ctxx.Err().Error())
						break Gone

					case ev := <-eventChannel:
						commons.DebugLog("ViewProvider::MonitorPage(", host.Name, ") listener received: ", ev)
					found:
						for _, device := range host.DeviceDetails {
							if device.Id == ev.Content.DeviceId {
								switch ev.Content.Name {
								case "power":
									z, _ := strconv.ParseFloat(ev.Content.Value, 32)
									err := device.BWattValue.Set(z)
									if err != nil {
										commons.DebugLog("ViewProvider::MonitorPage() listener(", ev.Content.Name, ") float parsing error: ", err.Error())
									}
									d64 := vv.chartPageData[host.Name]["Watts"].AddValue(z)
									point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorYellow, time.Now().Format(time.RFC1123))
									lc.ApplyDataPoint("Watts", &point)
									break found
								case "voltage":
									z, _ := strconv.ParseInt(ev.Content.Value, 10, 32)
									err := device.BVoltageValue.Set(int(z))
									if err != nil {
										commons.DebugLog("ViewProvider::MonitorPage() listener(", ev.Content.Name, ") int parsing error: ", err.Error())
									}
									d64 := vv.chartPageData[host.Name]["Voltage"].AddValue(float64(z))
									point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorGreen, time.Now().Format(time.RFC1123))
									lc.ApplyDataPoint("Voltage", &point)
									break found
								}
							}
						}

					default:
						time.Sleep(100 * time.Millisecond)
					}
				}
				commons.DebugLog("HubitatProvider::CreateDeviceEventListener() publisher END")
			}(v.ctx, v, hub, chart)
		}
	}

	return container.NewPadded(hostTabs)
}

/*
// handleUpdatesForMonitorPage does exactly that
func (v *viewProvider) handleUpdatesForMonitorPage(host *entities.ApcHost, service interfaces.Service, status *widget.Entry, events *widget.Entry, chart sknlinechart.LineChart, kChan chan map[string]string) {
	//hubEvents := v.service. .HubEventsMessageChannel()
	go func(h *entities.ApcHost, svc interfaces.Service, st *widget.Entry, ev *widget.Entry, chart sknlinechart.LineChart, knowledge chan map[string]string) {
		commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] BEGIN")
		rcvTuple := svc.MessageChannelByName(h.Name)
		var stBuild strings.Builder
		var evBuild strings.Builder
		var currentSt []string
		var currentEv []string
	pageExit:
		for {
			select {
			case <-v.ctx.Done():
				close(knowledge) // detail pages
				v.bondedUpsStatus[h.Name].UnbindUpsData()
				commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] fired:", v.ctx.Err().Error())
				break pageExit

			case msg := <-rcvTuple.Status:
				currentSt = msg
				stBuild.Reset()
				for idx, line := range currentSt {
					stBuild.WriteString(fmt.Sprintf("[%02d] %s\n", idx, line))
				}
				st.SetText(stBuild.String())
				st.Refresh()

			case msg := <-rcvTuple.Events:
				currentEv = msg
				evBuild.Reset()
				for idx, line := range currentEv {
					evBuild.WriteString(fmt.Sprintf("[%02d] %s\n", idx, line))
				}
				ev.SetText(evBuild.String())
				ev.Refresh()

			case "place":
			found:
				for _, hub := range v.cfgHubHosts {
					for _, dv := range hub.DeviceDetails {
						if dv.Id == ev.Content.DeviceId {
							z, _ := strconv.ParseFloat(ev.Content.Value, 32)
							_ = dv.BWattValue.Set(z)
							break found
						}
					}
				}

			default:
				var params map[string]string

				if len(currentSt) > 0 {
					params = svc.ParseStatus(currentSt)

					for k, vv := range params {
						floatStr := strings.Split(vv, " ")
						floatStr[0] = strings.TrimSpace(floatStr[0])
						// gapcmon charted: LINEV, LOADPCT, BCHARGE, CUMONBATT, TIMELEFT
						switch k {
						case "LINEV":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorYellow, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("LINEV", &point)
						case "LOADPCT":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorBlue, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("LOADPCT", &point)
						case "BCHARGE":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorGreen, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("BCHARGE", &point)
						case "TIMELEFT":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorPurple, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("TIMELEFT", &point)
						case "CUMONBATT":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorOrange, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("CUMONBATT", &point)
						case "STATUS":
							if strings.Contains(vv, "ONLINE") {
								h.State = commons.HostStatusOnline
							} else if strings.Contains(vv, "CHARG") {
								h.State = commons.HostStatusCharging
							} else if strings.Contains(vv, "TEST") {
								h.State = commons.PreferencesIcon
							} else if strings.Contains(vv, "ONBATT") {
								h.State = commons.HostStatusOnBattery
							}
						}

					}
					// details page updates
					knowledge <- params

					// ready next cycle
					currentSt = currentSt[:0]
				}
			}
		}
		// cleanup data syncs
		commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] ENDED")
	}(host, service, status, events, chart, kChan)
}
*/
