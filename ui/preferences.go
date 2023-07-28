package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/entities"
	"image/color"
	"strconv"
	"time"
)

// PreferencesPage manages application settings
func (v *viewProvider) PreferencesPage() *fyne.Container {
	sDesc := canvas.NewText("Selected Hub", color.White)
	sDesc.Alignment = fyne.TextAlignLeading
	sDesc.TextStyle = fyne.TextStyle{Italic: true}
	sDesc.TextSize = 24
	desc := container.NewPadded(sDesc)

	tdesc := canvas.NewText("Hub", color.White)
	tdesc.Alignment = fyne.TextAlignLeading
	tdesc.TextStyle = fyne.TextStyle{Italic: true}
	tdesc.TextSize = 24
	tDesc := container.NewPadded(tdesc)

	table := widget.NewTable(
		func() (int, int) { // length
			return len(v.hosts) + 1, 5
		},
		func() fyne.CanvasObject { // created
			return widget.NewRichTextWithText("")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			if id.Row == 0 { // headers
				switch id.Col {
				case 0:
					object.(*widget.RichText).ParseMarkdown("## Enabled")
				case 1:
					object.(*widget.RichText).ParseMarkdown("## Name")
				case 2:
					object.(*widget.RichText).ParseMarkdown("## IpAddress")
				case 3:
					object.(*widget.RichText).ParseMarkdown("## Average")
				case 4:
					object.(*widget.RichText).ParseMarkdown("## Access Token")
				}
				return
			}
			host := v.hosts[id.Row-1]
			switch id.Col {
			case 0: // Enabled
				label := "disabled"
				if host.IsEnabled() {
					label = "enabled"
				}
				object.(*widget.RichText).ParseMarkdown(label)
			case 1: // Name
				object.(*widget.RichText).ParseMarkdown(host.Name)

			case 2: // IpAddress
				object.(*widget.RichText).ParseMarkdown(host.IpAddress)

			case 3: // Graph
				object.(*widget.RichText).ParseMarkdown(strconv.Itoa(int(host.GraphingSamplePeriod)))

			case 4: // Access
				object.(*widget.RichText).ParseMarkdown(host.AccessToken)

			default:
				object.(*widget.RichText).ParseMarkdown("## Default")
			}
		},
	)

	enable := widget.NewCheck("", nil)
	enable.SetChecked(v.host.Enabled)

	dName := widget.NewEntry()
	dName.SetText(v.host.Name)

	dIp := widget.NewEntry()
	dIp.SetText(v.host.IpAddress)

	tIp := widget.NewEntry()
	tIp.SetText(commons.DefaultIp())

	hid := widget.NewEntry()
	hid.Disable()
	hid.SetText(v.host.Id)

	gPeriod := widget.NewEntry()
	gPeriod.SetText(strconv.Itoa(int(v.host.GraphingSamplePeriod)))

	dAccess := widget.NewEntry()
	dAccess.SetText(v.host.AccessToken)

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "is enabled", Widget: enable},
			{Text: "name", Widget: dName},
			{Text: "ip address", Widget: dIp},
			{Text: "This ip address", Widget: tIp, HintText: "defaults if blank"},
			{Text: "graph averaging count", Widget: gPeriod},
			{Text: "access token", Widget: dAccess, HintText: "defaults if blank"},
			{Text: "ID", Widget: hid},
		},
		SubmitText: "Apply edits",
	}
	form.OnSubmit = func() { // Apply optional, handle form submission
		var defaultListenerIp string
		if tIp.Text == "" {
			defaultListenerIp = commons.DefaultIp()
		} else {
			defaultListenerIp = tIp.Text
		}
		var defaultAccesstoken string
		if dAccess.Text == "" {
			defaultAccesstoken = commons.HubitatAccessToken()
		} else {
			defaultAccesstoken = dAccess.Text
		}
		if v.host.Id != hid.Text {
			gx, _ := strconv.Atoi(gPeriod.Text)
			v.host = entities.NewHubHost( // new
				dName.Text,
				dIp.Text,
				defaultAccesstoken,
				defaultListenerIp,
				time.Duration(gx),
				enable.Checked,
			)
		} else {
			x, _ := strconv.Atoi(gPeriod.Text)
			v.host.Update( // existing
				dName.Text,
				dIp.Text,
				dAccess.Text,
				defaultListenerIp,
				time.Duration(x),
				enable.Checked,
			)
		}
		v.cfg.Apply(v.host).Save()
		v.hosts = v.cfg.Hosts()
		table.Refresh()
		v.prfStatusLine.SetText(fmt.Sprintf("Form submitted for host:%s, restart for effect", v.host.Name))
	}

	table.OnSelected = func(id widget.TableCellID) {
		row := id.Row - 1
		if row > len(v.hosts) || id.Row == 0 {
			commons.DebugLog("Preferences::Table row invalid: ", row)
			return
		}
		v.host = v.hosts[row] // locally known

		dName.SetText(v.host.Name)
		dIp.SetText(v.host.IpAddress)
		tIp.SetText(v.host.ThisIpAddress)
		hid.SetText(v.host.Id)
		dAccess.SetText(v.host.AccessToken)
		z := strconv.Itoa(int(v.host.GraphingSamplePeriod))
		gPeriod.SetText(z)
		enable.Checked = v.host.Enabled

		form.Refresh()
		v.prfStatusLine.SetText(fmt.Sprintf("Selected row:%d:%d, col:%d, for host:%s", id.Row, row, id.Col, v.cfg.Hosts()[row].Name))
	}
	table.SetColumnWidth(0, 84)  // enable
	table.SetColumnWidth(1, 112) // name
	table.SetColumnWidth(2, 112) // ipAddress
	table.SetColumnWidth(3, 84)  // graph
	table.SetColumnWidth(4, 300) // access

	table.Select(widget.TableCellID{
		Row: 1,
		Col: 0,
	})

	page := container.NewGridWithColumns(1,
		settings.NewSettings().LoadAppearanceScreen(v.mainWindow),
		container.NewBorder(
			desc,
			nil,
			nil,
			nil,
			form,
		),

		container.NewBorder(
			tDesc,
			container.NewHBox(
				container.NewHBox(
					widget.NewButtonWithIcon("add selected", theme.ContentAddIcon(), func() {
						hid.SetText("new") // enable OnSubmit to treat this as a new record
						form.OnSubmit()
						v.prefsAddAction()
					}),
					widget.NewButtonWithIcon("del selected", theme.ContentRemoveIcon(), func() {
						v.prefsDelAction()
					}),
				),
				v.prfStatusLine,
			),
			nil,
			nil,
			table,
		),
	)
	return page
}
