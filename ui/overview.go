package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/hubPower/commons"
	"image/color"
	"strings"
)

func (v *viewProvider) OverviewPage() *fyne.Container {
	table := widget.NewTable(
		func() (int, int) { // length, columns
			return len(v.hosts) + 1, 3
		},
		func() fyne.CanvasObject { // created
			i := widget.NewIcon(theme.StorageIcon())
			i.Hide()

			l := widget.NewRichTextFromMarkdown("")

			return container.NewHBox(i, l) // issue container minSize is 0
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			// ICON - NAME, Device Listing,
			//              Device Listing,
			// Row, Col
			if id.Row == 0 { // headers
				object.(*fyne.Container).Objects[0].Hide()
				switch id.Col {
				case 0:
					object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown("")
					object.(*fyne.Container).Objects[1].Show()
				case 1:
					object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown("## Hub")
					object.(*fyne.Container).Objects[1].Show()
				case 2:
					object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown("## Device Information")
					object.(*fyne.Container).Objects[1].Show()
				}
				return
			}
			host := v.hosts[id.Row-1]
			devices := host.DeviceDetails
			switch id.Col {
			case 0: // State
				icon := "ThumbsUp"
				if !host.IsEnabled() {
					icon = "ThumbsDown"
				}
				object.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(commons.SknSelectThemedResource(icon))
				object.(*fyne.Container).Objects[0].(*widget.Icon).Resize(fyne.NewSize(40, 40))
				object.(*fyne.Container).Objects[0].Show()
				object.(*fyne.Container).Objects[1].Hide()

			case 1: // Enabled
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown("## " + strings.ToUpper(host.Name))
				object.(*fyne.Container).Objects[1].Show()

			case 2: // descriptions
				z := ""
				if !host.IsEnabled() {
					z = "## no data available until Hub is enabled."
				} else {
					for _, dv := range devices {
						vac, _ := dv.BVoltageValue.Get()
						watts, _ := dv.BWattValue.Get()
						str := fmt.Sprintf("### Id:%s %11s %s VAC: %3v Watts: %4.1f\n\n", dv.Id, dv.Label, dv.Name, vac, watts)
						z += str
					}
				}
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown(z)
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()
				object.(*fyne.Container).Refresh()

			default:
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Default")
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()
			}
		},
	)

	table.SetColumnWidth(0, 42)  // icon
	table.SetColumnWidth(1, 96)  // status
	table.SetColumnWidth(2, 432) // description
	for idx := range v.hosts {
		table.SetRowHeight(idx, 72)
	}

	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeWidth = 4
	rect.StrokeColor = theme.PrimaryColor()

	v.overviewTable = table // allow external refresh

	return container.NewPadded(table)
}
