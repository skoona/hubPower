package ui

import (
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/skoona/hubPower/commons"
	"net/url"
)

func (v *viewProvider) shortcutFocused(s fyne.Shortcut) {
	if focused, ok := v.mainWindow.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func (v *viewProvider) SknTrayMenu() {
	// Add SystemBar Menu
	if desk, ok := fyne.CurrentApp().(desktop.App); ok {
		m := fyne.NewMenu("gAPC Monitor",
			fyne.NewMenuItem("Show monitor...", func() {
				v.mainWindow.Show()
			}),
			fyne.NewMenuItem("Show preferences...", func() {
				if v.prefsWindow.Content().Size().Width <= 10 {
					v.ShowPrefsPage()
				} else {
					v.prefsWindow.Show()
				}
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(commons.SknSelectResource("apcupsd"))
	}
}
func (v *viewProvider) SknMenus() {

	settingsItem := fyne.NewMenuItem("Preferences", func() {
		if v.prefsWindow.Content().Size().Width <= 10 {
			v.ShowPrefsPage()
		} else {
			v.prefsWindow.Show()
		}
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		v.shortcutFocused(&fyne.ShortcutCut{
			Clipboard: v.mainWindow.Clipboard(),
		})
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		v.shortcutFocused(&fyne.ShortcutCopy{
			Clipboard: v.mainWindow.Clipboard(),
		})
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		v.shortcutFocused(&fyne.ShortcutPaste{
			Clipboard: v.mainWindow.Clipboard(),
		})
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
	)
	file := fyne.NewMenu("File")
	if !fyne.CurrentDevice().IsMobile() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem),
		helpMenu,
	)
	v.mainWindow.SetMainMenu(mainMenu)
}
