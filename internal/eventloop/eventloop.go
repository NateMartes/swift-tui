package eventloop

import (
	swiftSdk "github.com/ncw/swift/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Adds event handlers to the application before starting the application event loop
func SetupEventLoop(client *swiftSdk.Connection, view *tview.Application) {
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Rune() == 'q' {
            view.Stop()
            return nil
        }
        return event
    })
}
