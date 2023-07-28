package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/hubPower/entities"
	"image/color"
)

func (v *viewProvider) DeviceCard(dv *entities.DeviceDetails) *fyne.Container {
	desc := canvas.NewText(dv.Label, theme.PrimaryColor())
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

	lbl := widget.NewLabel("Name")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Name))

	lbl = widget.NewLabel("Type")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Type))

	lbl = widget.NewLabel("Id")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Id))

	lbl = widget.NewLabel("Date")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(*dv.Date))

	lbl = widget.NewLabel("Switch")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Attributes.Switch))

	lbl = widget.NewLabel("Room")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Room))

	lbl = widget.NewLabel("Energy")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Attributes.Energy))

	lbl = widget.NewLabel("Amperage")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(dv.Attributes.Amperage))

	lbl = widget.NewLabel("Watts")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(
		binding.FloatToStringWithFormat(dv.BWattValue, "%4.1f")),
	)

	lbl = widget.NewLabel("Voltage")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabelWithData(
		binding.IntToString(dv.BVoltageValue)),
	)

	return titleBorder
}
