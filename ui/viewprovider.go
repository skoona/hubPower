package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/entities"
	"github.com/skoona/hubPower/interfaces"
)

// viewProvider control structure for view management
type viewProvider struct {
	ctx           context.Context
	cfg           interfaces.Configuration
	service       interfaces.Service
	mainWindow    fyne.Window
	prefsWindow   fyne.Window
	overviewTable *widget.Table
	prfStatusLine *widget.Label
	chartPageData map[string]map[string]interfaces.GraphPointSmoothing
	chartKeys     []string
	hosts         []*entities.HubHost
	host          *entities.HubHost
}

// compiler helpers to insure interfaces requirements are meet
var (
	_ interfaces.ViewProvider = (*viewProvider)(nil)
	_ interfaces.Provider     = (*viewProvider)(nil)
)

// NewViewProvider manages all UI views and implements the ViewProvider Interface
func NewViewProvider(ctx context.Context, cfg interfaces.Configuration, service interfaces.Service) interfaces.ViewProvider {
	hh := cfg.Hosts()
	stLine := widget.NewLabel("click entry in table to edit, or click add to add.")
	//stLine.Wrapping = fyne.TextWrapWord -- causes pref page to break
	view := &viewProvider{
		ctx:           ctx,
		cfg:           cfg,
		service:       service,
		mainWindow:    fyne.CurrentApp().NewWindow("Hubitat Power Monitor"),
		prefsWindow:   fyne.CurrentApp().NewWindow("Preferences"),
		prfStatusLine: stLine,
		chartPageData: map[string]map[string]interfaces.GraphPointSmoothing{}, // [host][chartkey]struct
		hosts:         hh,
		host:          hh[0],
		chartKeys:     []string{"Watts", "Voltage"},
	}
	view.mainWindow.Resize(fyne.NewSize(726, 448))
	view.mainWindow.SetCloseIntercept(func() { view.mainWindow.Hide() })
	view.mainWindow.SetMaster()
	view.mainWindow.SetIcon(commons.SknSelectThemedResource(commons.AppIcon))

	view.prefsWindow.Resize(fyne.NewSize(632, 572))
	view.prefsWindow.SetCloseIntercept(func() { view.prefsWindow.Hide() })
	view.mainWindow.SetIcon(commons.SknSelectThemedResource(commons.PreferencesIcon))

	view.SknTrayMenu()
	view.SknMenus()

	return view
}

// ShowMainPage display the primary application page
func (v *viewProvider) ShowMainPage() {
	v.mainWindow.SetContent(v.MonitorPage())
	v.mainWindow.Show()
}

// ShowPrefsPage displays teh settings por preferences page
func (v *viewProvider) ShowPrefsPage() {
	v.hosts = v.cfg.Hosts()
	v.host = v.hosts[0]
	v.prefsWindow.SetContent(v.PreferencesPage())
	v.prefsWindow.Show()
}

// Shutdown closes all go routine
func (v *viewProvider) Shutdown() {
	commons.DebugLog("ViewProvider::Shutdown() called.")
}

// prefsAddAction adds or replaces the host in the form
func (v *viewProvider) prefsAddAction() {
	v.cfg.AddHost(v.host)
	v.ShowPrefsPage()
	v.prfStatusLine.SetText("Host " + v.host.Name + " was added")
}

// prefsDelAction deletes the select host
func (v *viewProvider) prefsDelAction() {
	n := v.host.Name
	v.cfg.Remove(v.host.Id)
	v.ShowPrefsPage()
	v.prfStatusLine.SetText("Host " + n + " was removed")
}
