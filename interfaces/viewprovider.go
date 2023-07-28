package interfaces

import "fyne.io/fyne/v2"

type ViewProvider interface {
	ShowPrefsPage()
	ShowMainPage() fyne.Window
	Shutdown()
}
