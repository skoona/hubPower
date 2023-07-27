package interfaces

type ViewProvider interface {
	ShowPrefsPage()
	ShowMainPage()
	Shutdown()
}
