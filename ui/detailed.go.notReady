package ui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/entities"
	"image/color"
)

func (v *viewProvider) Performance(bond *entities.UpsStatusValueBindings) *fyne.Container {
	desc := canvas.NewText("Performance Summary", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	lbl := widget.NewLabel("Selftest running")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bselftest))

	lbl = widget.NewLabel("Number of transfers")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bnumxfers))

	lbl = widget.NewLabel("Reason last transfer")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Blastxfer))

	lbl = widget.NewLabel("Last transfer to battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bxonbatt))

	lbl = widget.NewLabel("Last transfer off battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bxoffbatt))

	lbl = widget.NewLabel("Time on battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Btonbatt))

	lbl = widget.NewLabel("Cummulative on battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bcumonbatt))

	return titleBorder
}
func (v *viewProvider) Metrics(bond *entities.UpsStatusValueBindings) *fyne.Container {
	var desc *canvas.Text
	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	desc = canvas.NewText("UPS Metrics", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	lbl := widget.NewLabel("Utility line")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Blinev))

	lbl = widget.NewLabel("Battery DC")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bbattv))

	lbl = widget.NewLabel("Percent battery charge")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bbcharge))

	lbl = widget.NewLabel("Percent load capacity")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bloadpct))

	lbl = widget.NewLabel("Minutes remaining")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Btimeleft))

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)
	return titleBorder
}
func (v *viewProvider) Software(bond *entities.UpsStatusValueBindings) *fyne.Container {
	desc := canvas.NewText("Software Information ", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	lbl := widget.NewLabel("APCUPSD version")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bversion))

	lbl = widget.NewLabel("This node")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bhostname))

	lbl = widget.NewLabel("Monitored UPS name")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bupsname))

	t, _ := bond.Bmaster.Get()
	if t != "" {
		lbl = widget.NewLabel("Parent Node ip")
		lbl.Alignment = fyne.TextAlignTrailing
		items.Add(lbl)
		items.Add(widget.NewLabelWithData(bond.Bmaster))
	}

	lbl = widget.NewLabel("Cable driver type")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bcable))

	lbl = widget.NewLabel("Driver interface")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bdriver))

	lbl = widget.NewLabel("Configuration mode")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bupsmode))

	lbl = widget.NewLabel("Last started")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bstarttime))

	lbl = widget.NewLabel("UPS state")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bstatus))

	return titleBorder
}
func (v *viewProvider) Product(bond *entities.UpsStatusValueBindings) *fyne.Container {
	desc := canvas.NewText("Product Information ", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	lbl := widget.NewLabel("DeviceList model")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bmodel))

	lbl = widget.NewLabel("Serial number")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bserialno))

	lbl = widget.NewLabel("Manufacture date")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bmandate))

	lbl = widget.NewLabel("Firmware")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bfirmware))

	lbl = widget.NewLabel("Battery date")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bbattdate))

	lbl = widget.NewLabel("Internal temp")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(bond.Bitemp))

	return titleBorder
}

// DetailPage manages the different cards in the Details view based on type of ups node
func (v *viewProvider) DetailPage(params chan map[string]string, bond *entities.UpsStatusValueBindings) *fyne.Container {
	page := container.NewGridWithColumns(2)

	go func(bondData *entities.UpsStatusValueBindings, status chan map[string]string, content *fyne.Container) {
		oneShot := true
		commons.DebugLog("ViewProvider::UpsStatusValueBindings[", bondData.Host.Name, "] BEGIN")
		// wait for a msg to determine what type of UPS,
		// len < 24 (or status key of 'MASTER') is a networked node without a local UPS,
		// while len > 24 assumes node with a local UPS and battery attached
		for msg := range status {
			bondData.Apply(msg)
			if oneShot {
				content.Add(v.Performance(bondData))
				if len(msg) > 24 {
					content.Add(v.Metrics(bondData))
				}
				content.Add(v.Software(bondData))
				if len(msg) > 24 {
					content.Add(v.Product(bondData))
				}
				v.overviewTable.Refresh()
				oneShot = false
			}
		}
		commons.DebugLog("ViewProvider::UpsStatusValueBindings[", bondData.Host.Name, "] END")
	}(bond, params, page)

	return page
}
