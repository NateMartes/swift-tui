package ui

import (
	"github.com/rivo/tview"
)

// Get main box for the TUI application
func GetMainBox() *tview.Box {

	val := tview.NewBox().
		SetBorder(true).
		SetTitle("OpenStack Swift TUI")

	return val
}

// Gets the main view for the TUI application
func GetMainTUI() *tview.Application {

	app := tview.NewApplication()
	app.SetRoot(GetMainBox(), true)
	return app
}

// Closes the TUI application
func EndMainTUI(app *tview.Application) {
	app.Stop()
}
