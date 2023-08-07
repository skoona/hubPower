package ports

import "fyne.io/fyne/v2"

type ViewHandler interface {
	ShowPrefsPage()
	ShowMainPage() fyne.Window
	Shutdown()
}
