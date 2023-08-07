package handler

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/hubPower/internal/commons"
	"github.com/skoona/hubPower/internal/core/entities"
	"github.com/skoona/hubPower/internal/core/ports"
)

// viewHandler control structure for view management
type viewHandler struct {
	ctx           context.Context
	cfg           ports.Configuration
	service       ports.Service
	mainWindow    fyne.Window
	prefsWindow   fyne.Window
	overviewTable *widget.Table
	prfStatusLine *widget.Label
	chartPageData map[string]map[string]ports.GraphPointSmoothing
	chartKeys     []string
	hosts         []*entities.HubHost
	host          *entities.HubHost
}

// compiler helpers to insure ports requirements are meet
var (
	_ ports.ViewHandler = (*viewHandler)(nil)
	_ ports.Provider    = (*viewHandler)(nil)
)

// NewViewHandler manages all UI views and implements the ViewHandler Interface
func NewViewHandler(ctx context.Context, cfg ports.Configuration, service ports.Service) ports.ViewHandler {
	hh := cfg.Hosts()
	stLine := widget.NewLabel("click entry in table to edit, or click add to add.")
	//stLine.Wrapping = fyne.TextWrapWord -- causes pref page to break
	view := &viewHandler{
		ctx:           ctx,
		cfg:           cfg,
		service:       service,
		mainWindow:    fyne.CurrentApp().NewWindow("Hubitat Power Monitor"),
		prefsWindow:   fyne.CurrentApp().NewWindow("Preferences"),
		prfStatusLine: stLine,
		chartPageData: map[string]map[string]ports.GraphPointSmoothing{}, // [host][chartkey]struct
		hosts:         hh,
		host:          hh[0],
		chartKeys:     []string{"Watts", "Voltage"},
	}
	view.mainWindow.Resize(fyne.NewSize(726, 448))
	view.mainWindow.SetCloseIntercept(func() { view.mainWindow.Hide() })
	view.mainWindow.SetMaster()
	view.mainWindow.CenterOnScreen()
	view.mainWindow.SetIcon(commons.SknSelectThemedResource(commons.AppIcon))

	view.prefsWindow.Resize(fyne.NewSize(632, 572))
	view.prefsWindow.SetCloseIntercept(func() { view.prefsWindow.Hide() })
	view.mainWindow.SetIcon(commons.SknSelectThemedResource(commons.PreferencesIcon))

	view.SknTrayMenu()
	view.SknMenus()

	return view
}

// ShowMainPage display the primary application page
func (v *viewHandler) ShowMainPage() fyne.Window {
	v.mainWindow.SetContent(v.MonitorPage())
	v.mainWindow.Show()
	return v.mainWindow
}

// ShowPrefsPage displays teh settings por preferences page
func (v *viewHandler) ShowPrefsPage() {
	v.hosts = v.cfg.Hosts()
	v.host = v.hosts[0]
	v.prefsWindow.SetContent(v.PreferencesPage())
	v.prefsWindow.Show()
}

// Shutdown closes all go routine
func (v *viewHandler) Shutdown() {
	commons.DebugLog("ViewHandler::Shutdown() called.")
}

// prefsAddAction adds or replaces the host in the form
func (v *viewHandler) prefsAddAction() {
	v.cfg.AddHost(v.host)
	v.ShowPrefsPage()
	v.prfStatusLine.SetText("Host " + v.host.Name + " was added")
}

// prefsDelAction deletes the select host
func (v *viewHandler) prefsDelAction() {
	n := v.host.Name
	v.cfg.Remove(v.host.Id)
	v.ShowPrefsPage()
	v.prfStatusLine.SetText("Host " + n + " was removed")
}
